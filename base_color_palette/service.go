package basecolorpalette

import (
	"errors"
	"fmt"
	"myproject/models"
	"myproject/utils"
	"net/http"
)

type ColorBase struct {
	CDB *utils.Base
}

func (cdb ColorBase) ColorsHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	if r.Method != "GET" {
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}

	db := cdb.CDB.DB
	context := r.Context()
	query := "SELECT * FROM base_color_palette WHERE 1"

	rows, err := db.QueryContext(context, query)
	if err != nil {
		message := fmt.Sprintf("Ошибка SQL: %s", err)
		return &utils.AppError{
			Err:     err,
			Message: message,
			Code:    http.StatusInternalServerError,
		}
	}
	colors := []models.BaseColorPalette{}
	defer rows.Close()
	for rows.Next() {
		color := models.BaseColorPalette{}
		err := rows.Scan(&color.IDBaseColor, &color.Hex, &color.Name)

		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Не удалось распаковать данные",
				Code:    http.StatusInternalServerError,
			}
		}

		colors = append(colors, color)
	}

	if err = rows.Err(); err != nil {
		return &utils.AppError{
			Err:     errors.New("Ошибка итерации"),
			Message: "rows.Next() закончился с ошибками",
			Code:    http.StatusInternalServerError,
		}
	}

	return utils.Form_response(w, colors, http.StatusOK)
}
