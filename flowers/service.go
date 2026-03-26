package flowers

import (
	"errors"
	"fmt"
	"myproject/models"
	"myproject/utils"
	"net/http"
)

type FlowerBase struct {
	FDB *utils.Base
}

func (fdb FlowerBase) FlowersHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	if r.Method != "GET" {
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}

	db := fdb.FDB.DB
	context := r.Context()
	query := "SELECT * FROM `flowers` WHERE 1"

	rows, err := db.QueryContext(context, query)
	if err != nil {
		message := fmt.Sprintf("Ошибка SQL: %s", err)
		return &utils.AppError{
			Err:     err,
			Message: message,
			Code:    http.StatusInternalServerError,
		}
	}
	flowers := []models.Flower{}
	defer rows.Close()
	for rows.Next() {
		flower := models.Flower{}
		err := rows.Scan(&flower.IDFlower, &flower.NameFlower)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Не удалось распаковать данные",
				Code:    http.StatusInternalServerError,
			}
		}
		flowers = append(flowers, flower)
	}
	if err = rows.Err(); err != nil {
		return &utils.AppError{
			Err:     errors.New("Ошибка итерации"),
			Message: "rows.Next() закончился с ошибками",
			Code:    http.StatusInternalServerError,
		}
	}
	return utils.Form_response(w, flowers, http.StatusOK)
}
