package admin

/*
	АДМИН
	1) создавать, удалать, редактировать товары
	2) убрать корзину полностью

	ОБЩИЕ
	1) витрина
	2) фильтр каталога
	3) просмотр данных заказа


	3) сортировка заказов по дедлайну
	КЛИЕНТ
	1) выбор времени самовывоза

*/

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"myproject/utils"
	"net/http"
	"time"
)

type AdminBase struct {
	ADB *utils.Base
}

type OrderFullData struct {
	IDOrder         int     `json:"id_order"`
	ClientID        int     `json:"id_client"`
	StatusID        int     `json:"id_status"`
	TotalCost       float64 `json:"total_cost"`
	PaymentStatus   string  `json:"payment_status"`
	PaymentRef      string  `json:"payment_ref"`
	Deadline        string  `json:"deadline"`
	PickUp_DateTime string  `json:"pickup_datetime"`
	CreatedAt       string  `json:"created_at"`
	Last_name       string  `json:"last_name"`
	First_name      string  `json:"first_name"`
	Phone           string  `json:"phone"`
	Email           string  `json:"email"`
}

func (adb AdminBase) OrdersForAdminHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	//получить данные всех заказов
	switch r.Method {
	case "GET":
		db := adb.ADB.DB
		//проверяем права
		_, role, app_err := utils.GetUserData(w, r, db)

		if app_err != nil {
			return app_err
		}

		if role != "admin" {
			return &utils.AppError{
				Err:     errors.New("Access denied"),
				Message: "Доступ запрещен, у вас нет прав администратора",
				Code:    http.StatusForbidden,
			}
		}

		// делаем запрос к orders
		query := "SELECT o.id_order, o.id_client, o.id_status, o.total_cost, o.payment_status, o.payment_ref, o.deadline, o.pickup_datetime," +
			" o.created_at, c.last_name, c.first_name, c.phone, c.email" +
			" FROM orders o JOIN clients c ON o.id_client = c.id_client ORDER BY deadline ASC"

		context := r.Context()

		rows, err := db.QueryContext(context, query)

		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка SQL запроса",
				Code:    http.StatusInternalServerError,
			}
		}
		defer rows.Close()
		var orders_data []OrderFullData
		for rows.Next() {
			data := OrderFullData{}
			var paymentStatus sql.NullString
			var paymentRef sql.NullString
			var pickup sql.NullString

			err := rows.Scan(
				&data.IDOrder,
				&data.ClientID,
				&data.StatusID,
				&data.TotalCost,
				&paymentStatus,
				&paymentRef,
				&data.Deadline,
				&pickup,
				&data.CreatedAt,
				&data.Last_name,
				&data.First_name,
				&data.Phone,
				&data.Email,
			)

			if paymentStatus.Valid {
				data.PaymentStatus = paymentStatus.String
			} else {
				data.PaymentStatus = ""
			}
			if paymentRef.Valid {
				data.PaymentRef = paymentRef.String
			} else {
				data.PaymentRef = ""
			}
			if pickup.Valid {
				data.PickUp_DateTime = pickup.String
			} else {
				data.PickUp_DateTime = ""
			}

			if err != nil {
				return &utils.AppError{
					Err:     err,
					Message: "Не удалось распаковать данные",
					Code:    http.StatusInternalServerError,
				}
			}

			orders_data = append(orders_data, data)
		}
		if err = rows.Err(); err != nil {
			return &utils.AppError{
				Err:     errors.New("Ошибка итерации"),
				Message: "rows.Next() закончился с ошибками",
				Code:    http.StatusInternalServerError,
			}
		}

		return utils.Form_response(w, orders_data, http.StatusOK)

	default:
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}
}

func (adb AdminBase) OrderIdHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	db := adb.ADB.DB
	//проверяем права

	_, role, app_err := utils.GetUserData(w, r, db)

	if app_err != nil {
		return app_err
	}

	if role != "admin" {
		return &utils.AppError{
			Err:     errors.New("Access denied"),
			Message: "Доступ запрещен, у вас нет прав администратора",
			Code:    http.StatusForbidden,
		}
	}
	switch r.Method {
	case "PUT": // обновлияем статус заказа

		type StatusData struct {
			StatusID int `json:"id_status"`
			ID       int `json:"id_order"`
		}
		var status_data StatusData
		err := json.NewDecoder(r.Body).Decode(&status_data)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Не удалось получить данные статуса",
				Code:    http.StatusBadRequest,
			}
		}

		var pickupTime sql.NullTime

		// если устанавливае статус выдан, записываем дату самовывоза
		if status_data.StatusID == 4 {
			pickupTime.Time = time.Now()
			pickupTime.Valid = true // записываем время из первого поля
		} else {
			pickupTime.Valid = false
		}

		context := r.Context()
		query := "UPDATE orders SET id_status = ?, pickup_datetime = ? WHERE id_order = ?"

		result, err := db.ExecContext(context, query, status_data.StatusID, pickupTime, status_data.ID)

		if err != nil {
			message := fmt.Sprintf("Ошибка SQL: %s", err)
			return &utils.AppError{
				Err:     err,
				Message: message,
				Code:    http.StatusInternalServerError,
			}
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка RowsAffected()",
				Code:    http.StatusInternalServerError,
			}
		}
		if rows == 0 {
			return &utils.AppError{
				Err:     errors.New("ноль строк изменено"),
				Message: "Затронуто 0 строк",
				Code:    http.StatusInternalServerError,
			}
		}

		return utils.Form_response(w, "Статус обновлен", http.StatusOK)

	default:
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}
}
