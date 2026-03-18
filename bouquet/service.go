package bouquet

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

	_ "github.com/go-sql-driver/mysql"
)

type BouquetBase struct {
	BDB *utils.Base // Embedding встраивание
}

type ListIdBouquets struct {
	IDs []int `json:"ids"`
}

func (b BouquetBase) ListBouquetsHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	switch r.Method {
	case "POST":
		db := b.BDB.DB
		context := r.Context()

		var list_id ListIdBouquets

		err := json.NewDecoder(r.Body).Decode(&list_id)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Некорректный формат данных",
				Code:    http.StatusBadRequest,
			}
		}

		var query_elements []string
		var query_questions []string
		query_begin := "SELECT id_bouquet, name, description, price, image_url, reserve_image_url, id_base_color, type FROM bouquets WHERE id_bouquet IN ("
		query_elements = append(query_elements, query_begin)

		for i := 0; i < len(list_id.IDs); i++ {
			query_questions = append(query_questions, "?")
		}

		questions := strings.Join(query_questions, ",")

		query_elements = append(query_elements, questions)
		query_elements = append(query_elements, ")")

		query := strings.Join(query_elements, "")

		args := make([]any, len(list_id.IDs))

		for i := 0; i < len(list_id.IDs); i++ {
			args[i] = list_id.IDs[i]
		}

		rows, err := db.QueryContext(context, query, args...)

		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка SQL запроса",
				Code:    http.StatusInternalServerError,
			}
		}

		defer rows.Close()
		bouquets := []models.Bouquet{}

		for rows.Next() {
			bouquet := models.Bouquet{}
			err := rows.Scan(&bouquet.IDBouquet, &bouquet.Name, &bouquet.Description, &bouquet.Price, &bouquet.ImageUrl, &bouquet.ReserveImageUrl, &bouquet.IDBaseColor, &bouquet.Type)
			if err != nil {
				return &utils.AppError{
					Err:     err,
					Message: "Не удалось распаковать данные",
					Code:    http.StatusInternalServerError,
				}
			}
			bouquets = append(bouquets, bouquet)
		}

		return utils.Form_response(w, bouquets, http.StatusOK)
	default:
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}
}

func (b BouquetBase) BouquetIdHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	db := b.BDB.DB
	idStr := strings.TrimPrefix(r.URL.Path, "/bouquets/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return &utils.AppError{
			Err:     err,
			Message: "Неверный ID",
			Code:    http.StatusBadRequest,
		}
	}
	switch r.Method {
	case "GET":
		data := models.Bouquet{}
		context := r.Context()
		query := "SELECT id_bouquet, name, description, price, image_url, reserve_image_url, id_base_color, type FROM bouquets WHERE id_bouquet = ?"
		err := db.QueryRowContext(context, query, id).Scan(&data.IDBouquet, &data.Name, &data.Description, &data.Price, &data.ImageUrl, &data.ReserveImageUrl, &data.IDBaseColor, &data.Type)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка отправки запроса SELECT WHERE id",
				Code:    http.StatusInternalServerError,
			}
		}
		return utils.Form_response(w, data, http.StatusOK)

	default:
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}
}

func (b BouquetBase) BouquetsHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	db := b.BDB.DB
	switch r.Method {
	case "GET":
		requ_type := r.URL.Query().Get("type")

		if requ_type != "usual" && requ_type != "special" && requ_type != "all" {
			return &utils.AppError{
				Err:     errors.New("unsupported type"),
				Message: "Несуществующий тип букета",
				Code:    http.StatusInternalServerError,
			}
		}

		var rows *sql.Rows
		var err error
		context := r.Context()
		if requ_type == "all" {
			query := "SELECT id_bouquet, name, description, price, image_url, reserve_image_url, id_base_color, type FROM bouquets"
			rows, err = db.QueryContext(context, query)
		}

		if requ_type == "usual" || requ_type == "special" {
			query := "SELECT id_bouquet, name, description, price, image_url, reserve_image_url, id_base_color, type FROM bouquets WHERE type = ?"
			rows, err = db.QueryContext(context, query, requ_type)
		}

		if err != nil {
			message := fmt.Sprintf("Ошибка SQL: %s", err)
			return &utils.AppError{
				Err:     err,
				Message: message,
				Code:    http.StatusInternalServerError,
			}
		}

		defer rows.Close()
		var bouquets []models.Bouquet
		for rows.Next() {
			bouquet := models.Bouquet{}
			if err := rows.Scan(&bouquet.IDBouquet, &bouquet.Name, &bouquet.Description, &bouquet.Price, &bouquet.ImageUrl,
				&bouquet.ReserveImageUrl, &bouquet.IDBaseColor, &bouquet.Type); err != nil {
				return &utils.AppError{
					Err:     err,
					Message: "Ошибка чтения данных о товарах",
					Code:    http.StatusInternalServerError,
				}
			}
			bouquets = append(bouquets, bouquet)
		}
		if err = rows.Err(); err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка итерации",
				Code:    http.StatusInternalServerError,
			}
		}

		return utils.Form_response(w, bouquets, http.StatusOK)

	case "POST":
		context := r.Context()

		bouquet_data := models.Bouquet{}

		err := json.NewDecoder(r.Body).Decode(&bouquet_data)

		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Некорректный формат данных",
				Code:    http.StatusBadRequest,
			}
		}
		query := "INSERT INTO bouquets(name, description, price, image_url) VALUES (?, ?, ?, ?)"

		result, err := db.ExecContext(context, query, bouquet_data.Name, bouquet_data.Description, bouquet_data.Price, bouquet_data.ImageUrl)

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

		return utils.Form_response(w, "Новый товар добавлен", http.StatusOK)

	default:
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}

	}
}
