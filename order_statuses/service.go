package orderstatuses

import (
	"errors"
	"myproject/models"
	"myproject/utils"
	"net/http"
)

type OrderStatusesBase struct {
	SDB *utils.Base
}

func (sdb OrderStatusesBase) OrderStatusHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	switch r.Method {
	case "GET":
		db := sdb.SDB.DB

		context := r.Context()
		query := "SELECT id_status, name FROM orderstatuses WHERE 1"

		rows, err := db.QueryContext(context, query)

		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка SQL запроса SELECT * FROM orderstatuses WHERE 1",
				Code:    http.StatusInternalServerError,
			}
		}

		defer rows.Close()
		var Statuses []models.OrderStatus
		for rows.Next() {
			status := models.OrderStatus{}
			err := rows.Scan(&status.IDStatus, &status.NameStatus)

			if err != nil {
				return &utils.AppError{
					Err:     err,
					Message: "Не удалось распаковать данные",
					Code:    http.StatusInternalServerError,
				}
			}

			Statuses = append(Statuses, status)
		}

		return utils.Form_response(w, Statuses, http.StatusOK)
	default:
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}

}
