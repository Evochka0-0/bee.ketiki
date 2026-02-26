package orderitems

import (
	"errors"
	"fmt"
	"myproject/models"
	"myproject/utils"
	"net/http"
	"strconv"
)

type ItemBase struct {
	IDB *utils.Base
}

func (idb ItemBase) OrderIdHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	switch r.Method {
	case "GET":
		db := idb.IDB.DB
		//----ПРИНИМАЕТ ID ЗАКАЗА
		id_str := r.URL.Query().Get("id_order")
		id, err := strconv.Atoi(id_str)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Неверный ID заказа",
				Code:    http.StatusBadRequest,
			}
		}

		context := r.Context()
		query := "SELECT quantity, image_url, reserve_image_url, type FROM " +
			"orderitems INNER JOIN orders ON orderitems.id_order = orders.id_order" +
			" INNER JOIN bouquets ON bouquets.id_bouquet = orderitems.id_bouquet WHERE orders.id_order = ?"
		//query := "SELECT id_bouquet, quantity FROM orderitems INNER JOIN orders ON orderitems.order_i` = orders.id WHERE orders.id = ?"

		result, err := db.QueryContext(context, query, id)
		if err != nil {
			message := fmt.Sprintf("Ошибка SQL: %s", err)
			return &utils.AppError{
				Err:     err,
				Message: message,
				Code:    http.StatusInternalServerError,
			}
		}

		defer result.Close()
		var items []models.OrderItemWithImage
		for result.Next() {
			item := models.OrderItemWithImage{}
			err := result.Scan(&item.Quantity, &item.ImageUrl, &item.ReserveImageUrl, &item.Type)

			if err != nil {
				return &utils.AppError{
					Err:     err,
					Message: "Не удалось распаковать данные item",
					Code:    http.StatusInternalServerError,
				}
			}
			items = append(items, item)
		}
		return utils.Form_response(w, items, http.StatusOK)

	default:
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}
}
