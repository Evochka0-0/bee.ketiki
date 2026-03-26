package occasions

import (
	"errors"
	"fmt"
	"myproject/models"
	"myproject/utils"
	"net/http"
)

type OccasionBase struct {
	ODB *utils.Base
}

func (odb OccasionBase) OccasionsHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	if r.Method != "GET" {
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}

	db := odb.ODB.DB
	context := r.Context()
	query := "SELECT * FROM occasion WHERE 1"

	rows, err := db.QueryContext(context, query)
	if err != nil {
		message := fmt.Sprintf("Ошибка SQL: %s", err)
		return &utils.AppError{
			Err:     err,
			Message: message,
			Code:    http.StatusInternalServerError,
		}
	}
	occasions := []models.Occasion{}
	defer rows.Close()
	for rows.Next() {
		occasion := models.Occasion{}
		err := rows.Scan(&occasion.IDOccasion, &occasion.OccasionName)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Не удалось распаковать данные",
				Code:    http.StatusInternalServerError,
			}
		}
		occasions = append(occasions, occasion)
	}
	if err = rows.Err(); err != nil {
		return &utils.AppError{
			Err:     errors.New("Ошибка итерации"),
			Message: "rows.Next() закончился с ошибками",
			Code:    http.StatusInternalServerError,
		}
	}
	return utils.Form_response(w, occasions, http.StatusOK)
}
