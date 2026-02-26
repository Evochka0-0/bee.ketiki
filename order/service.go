package order

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"myproject/models"
	"myproject/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type OrderBase struct {
	ODB *utils.Base
}

type OrderClientInfo struct {
	ClientID int `json:"id_client"`

	StatusID        int     `json:"id_status"`
	TotalCost       float64 `json:"total_cost"`
	PaymentStatus   string  `json:"payment_status"`
	Deadline        string  `json:"deadline"`
	PickUp_DateTime string  `json:"pickup_datetime"`
	CreatedAt       string  `json:"created_at"`

	Last_name  string `json:"last_name"`
	First_name string `json:"first_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

type BouquetInOrder struct {
	IDBouquet       int     `json:"id_bouquet"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	ImageUrl        string  `json:"image_url"`
	ReserveImageUrl string  `json:"reserve_image_url"`
	Type            string  `json:"type"`
	Quantity        int     `json:"quantity"`
}

func (odb OrderBase) OrderHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	switch r.Method {
	case "GET":
		db := odb.ODB.DB
		id_str := r.URL.Query().Get("user_id")
		id, err := strconv.Atoi(id_str)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Неверный ID",
				Code:    http.StatusBadRequest,
			}
		}

		context := r.Context()
		query := "SELECT id_order, id_client, id_status, total_cost, payment_status, payment_ref, deadline, pickup_datetime, created_at FROM orders WHERE id_client = ? ORDER BY deadline ASC"
		rows, err := db.QueryContext(context, query, id)

		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка SQL запроса",
				Code:    http.StatusInternalServerError,
			}
		}
		defer rows.Close()
		var orders []models.Order
		for rows.Next() {
			order := models.Order{}

			var paymentStatus sql.NullString
			var paymentRef sql.NullString
			var pickup sql.NullString
			err := rows.Scan(
				&order.IDOrder,
				&order.ClientID,
				&order.StatusID,
				&order.TotalCost,
				&paymentStatus,
				&paymentRef,
				&order.Deadline,
				&pickup,
				&order.CreatedAt,
			)
			if paymentStatus.Valid {
				order.PaymentStatus = paymentStatus.String
			} else {
				order.PaymentStatus = ""
			}
			if paymentRef.Valid {
				order.PaymentRef = paymentRef.String
			} else {
				order.PaymentRef = ""
			}
			if pickup.Valid {
				order.PickUp_DateTime = pickup.String
			} else {
				order.PickUp_DateTime = ""
			}

			if err != nil {
				return &utils.AppError{
					Err:     err,
					Message: "Не удалось распаковать данные",
					Code:    http.StatusInternalServerError,
				}
			}
			orders = append(orders, order)
		}

		return utils.Form_response(w, orders, http.StatusOK)

	case "POST":
		context := r.Context()
		db := odb.ODB.DB

		status := 1
		created_at := time.Now()

		type PostOrder struct {
			ClientID  int     `json:"client_id"`
			TotalCost float64 `json:"total_cost"`
			Deadline  string  `json:"deadline"`
			Items     []struct {
				IDBouquet int `json:"id_bouquet"`
				Quantity  int `json:"count"`
			} `json:"items"`
		}

		var request_data PostOrder

		err := json.NewDecoder(r.Body).Decode(&request_data)

		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Некорректный формат данных",
				Code:    http.StatusBadRequest,
			}
		}

		// Валидация
		if len(request_data.Items) == 0 {
			return &utils.AppError{
				Err:     errors.New("пустой заказ"),
				Message: "Заказ должен содержать хотя бы один товар",
				Code:    http.StatusBadRequest,
			}
		}

		for _, item := range request_data.Items {
			if item.Quantity <= 0 {
				return &utils.AppError{
					Err:     errors.New("некорректное количество"),
					Message: "Количество товара должно быть больше нуля",
					Code:    http.StatusBadRequest,
				}
			}
		}

		if len(request_data.Deadline) == 0 {
			return &utils.AppError{
				Err:     errors.New("Null value od datetime"),
				Message: "Укажите дату и время выдачи заказа",
				Code:    http.StatusBadRequest,
			}
		}

		// Начинаем транзакцию
		tx, err := db.BeginTx(context, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка начала транзакции",
				Code:    http.StatusInternalServerError,
			}
		}

		query := "INSERT INTO orders(id_client, id_status, total_cost, deadline, created_at) VALUES (?,?,?,?,?)"
		result, err := tx.ExecContext(context, query, request_data.ClientID, status, request_data.TotalCost, request_data.Deadline, created_at)

		if err != nil {
			tx.Rollback()
			message := fmt.Sprintf("Ошибка создания заказа: %s", err)
			return &utils.AppError{
				Err:     err,
				Message: message,
				Code:    http.StatusInternalServerError,
			}
		}

		orderID, err := result.LastInsertId()
		if err != nil {
			tx.Rollback()
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка получения ID заказа",
				Code:    http.StatusInternalServerError,
			}
		}
		stmt, err := tx.PrepareContext(context, "INSERT INTO orderitems(id_order, id_bouquet, quantity) VALUES (?,?,?)")
		if err != nil {
			tx.Rollback()
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка подготовки запроса для позиций заказа",
				Code:    http.StatusInternalServerError,
			}
		}
		defer stmt.Close()

		for _, item := range request_data.Items {
			_, err = stmt.ExecContext(context, orderID, item.IDBouquet, item.Quantity)
			if err != nil {
				tx.Rollback()
				message := fmt.Sprintf("Ошибка добавления позиции в заказ: %s", err)
				return &utils.AppError{
					Err:     err,
					Message: message,
					Code:    http.StatusInternalServerError,
				}
			}
		}

		if err := tx.Commit(); err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка сохранения заказа",
				Code:    http.StatusInternalServerError,
			}
		}

		response := map[string]interface{}{
			"message":  "Заказ оформлен! Спасибо за покупку",
			"id_order": orderID,
		}
		return utils.Form_response(w, response, http.StatusOK)

	default:
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}
}

func (adb OrderBase) OrderIdHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	// данные одного заказа для страницы заказа
	db := adb.ODB.DB

	strId := strings.TrimPrefix(r.URL.Path, "/order_data/")
	id_order, err := strconv.Atoi(strId)
	if err != nil {
		return &utils.AppError{
			Err:     err,
			Message: "Неверный ID",
			Code:    http.StatusBadRequest,
		}
	}

	context := r.Context()
	// заказ и клиент
	query := "SELECT o.id_client, o.id_status, o.total_cost, o.payment_status, o.deadline, o.pickup_datetime," +
		" o.created_at, c.last_name, c.first_name, c.phone, c.email FROM orders o JOIN clients c " +
		"ON o.id_client = c.id_client WHERE id_order = ?;"

	var order_client OrderClientInfo
	var paymentStatus sql.NullString
	var pickup sql.NullString
	err = db.QueryRowContext(context, query, id_order).Scan(
		&order_client.ClientID,
		&order_client.StatusID,
		&order_client.TotalCost,
		&paymentStatus,
		&order_client.Deadline,
		&pickup,
		&order_client.CreatedAt,
		&order_client.Last_name,
		&order_client.First_name,
		&order_client.Phone,
		&order_client.Email,
	)

	if paymentStatus.Valid {
		order_client.PaymentStatus = paymentStatus.String
	} else {
		order_client.PaymentStatus = ""
	}
	if pickup.Valid {
		order_client.PickUp_DateTime = pickup.String
	} else {
		order_client.PickUp_DateTime = ""
	}

	if err != nil {
		mess := fmt.Sprintf("Ошибка SQL запроса: %s", err)
		return &utils.AppError{
			Err:     err,
			Message: mess,
			Code:    http.StatusInternalServerError,
		}
	}

	// детали заказа и букеты
	query = "SELECT b.id_bouquet, b.name, b.description, b.price, b.image_url, b.reserve_image_url, b.type, i.quantity FROM orderitems i JOIN bouquets b " +
		"ON i.id_bouquet = b.id_bouquet WHERE i.id_order = ?;"

	rows, err := db.QueryContext(context, query, id_order)
	if err != nil {
		mess := fmt.Sprintf("Ошибка выполнения SQL запроса: %s", err)
		return &utils.AppError{
			Err:     err,
			Message: mess,
			Code:    http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	var bouquets_data []BouquetInOrder
	for rows.Next() {
		var bouquet BouquetInOrder
		err := rows.Scan(
			&bouquet.IDBouquet,
			&bouquet.Name,
			&bouquet.Description,
			&bouquet.Price,
			&bouquet.ImageUrl,
			&bouquet.ReserveImageUrl,
			&bouquet.Type,
			&bouquet.Quantity,
		)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Не удалось распаковать данные",
				Code:    http.StatusInternalServerError,
			}
		}
		bouquets_data = append(bouquets_data, bouquet)
	}
	if err = rows.Err(); err != nil {
		return &utils.AppError{
			Err:     errors.New("Ошибка итерации"),
			Message: "rows.Next() закончился с ошибками",
			Code:    http.StatusInternalServerError,
		}
	}

	type OrderFullDetailsResponse struct {
		OrderClient   OrderClientInfo  `json:"order_client"`
		OrderBouquets []BouquetInOrder `json:"order_bouquets"`
	}

	var order_info OrderFullDetailsResponse
	order_info.OrderClient = order_client
	order_info.OrderBouquets = bouquets_data

	return utils.Form_response(w, order_info, http.StatusOK)
}
