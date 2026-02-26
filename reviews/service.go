package reviews

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

type ReviewBase struct {
	RDB *utils.Base
}

func (rdb ReviewBase) ReviewAccessCheckHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	switch r.Method {
	case "POST":
		db := rdb.RDB.DB
		var id_bouquet int
		err := json.NewDecoder(r.Body).Decode(&id_bouquet)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Некорректный формат данных",
				Code:    http.StatusBadRequest,
			}
		}

		id_client, _, app_err := utils.GetUserData(w, r, db)

		if app_err != nil {
			return app_err
		}

		context := r.Context()
		query := "SELECT id_status FROM orders o JOIN orderitems i ON o.id_order = i.id_order WHERE o.id_client = ? AND i.id_bouquet = ? AND o.id_status = 4 LIMIT 1;"
		var statusID int
		err = db.QueryRowContext(context, query, id_client, id_bouquet).Scan(&statusID)

		if err != nil {
			if err == sql.ErrNoRows {
				// Человек не покупал этот букет
				return &utils.AppError{Message: "Сначала купите этот букет", Code: http.StatusForbidden}
			} else {
				return &utils.AppError{
					Err:     err,
					Message: "Ошибка сканирования результата запроса",
					Code:    http.StatusInternalServerError,
				}
			}
		}
		if statusID != 4 {
			return &utils.AppError{Message: "Вы сможете оставить отзыв, когда получите заказ", Code: http.StatusForbidden}
		}

		return utils.Form_response(w, "Напишите отзыв!", http.StatusOK)

	default:
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}
}

func (rdb ReviewBase) ReviewsHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {

	switch r.Method {
	case "GET":
		db := rdb.RDB.DB

		var id_bouquet int //из тела

		idStr := strings.TrimPrefix(r.URL.Path, "/reviews/")

		id_bouquet, err := strconv.Atoi(idStr)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Неверный ID",
				Code:    http.StatusBadRequest,
			}
		}

		context := r.Context()
		query := "SELECT c.last_name, c.first_name, r.message, r.grade, r.created_at FROM reviews r JOIN clients c ON r.id_client = c.id_client WHERE id_bouquet = ?"

		rows, err := db.QueryContext(context, query, id_bouquet)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка SQL запроса",
				Code:    http.StatusInternalServerError,
			}
		}

		defer rows.Close()
		reviews := []models.Reviews{}
		for rows.Next() {
			review := models.Reviews{}
			err := rows.Scan(&review.Last_name, &review.First_name, &review.Message, &review.Grade, &review.CreatedAt)
			if err != nil {
				return &utils.AppError{
					Err:     err,
					Message: "Не удалось распаковать данные",
					Code:    http.StatusInternalServerError,
				}
			}
			reviews = append(reviews, review)
		}
		return utils.Form_response(w, reviews, http.StatusOK)

	case "POST":
		db := rdb.RDB.DB

		type reviewPost struct {
			IDBouquet int    `json:"id_bouquet"`
			Message   string `json:"message"`
			Grade     int    `json:"grade"`
		}
		review := reviewPost{}
		err := json.NewDecoder(r.Body).Decode(&review)

		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Некорректный формат данных",
				Code:    http.StatusBadRequest,
			}
		}

		id_client, _, app_err := utils.GetUserData(w, r, db)
		if app_err != nil {
			return app_err
		}

		context := r.Context()

		created_at := time.Now()
		query := "INSERT INTO reviews(id_client, id_bouquet, message, grade, created_at)" +
			" VALUES (?, ?, ?, ?, ?)"
		result, err := db.ExecContext(context, query, id_client, review.IDBouquet, review.Message, review.Grade, created_at)

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

		return utils.Form_response(w, "Спасибо за отзыв!", http.StatusOK)

	default:
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}

}
