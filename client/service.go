package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"myproject/models"
	"myproject/utils"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type ClientBase struct {
	CDB *utils.Base
}

// зарегестрироваться как админ нельзя, поэтому роль по умолчанию user
func (cdb *ClientBase) ClientHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	switch r.Method {
	case "POST":
		db := cdb.CDB.DB
		var client_data models.Client
		err := json.NewDecoder(r.Body).Decode(&client_data)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Не удалось получить данные пользователя",
				Code:    http.StatusBadRequest,
			}
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(client_data.Password), 10)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Не удалось хешировать пароль",
				Code:    http.StatusInternalServerError,
			}
		}

		client_data.Password = string(hash)
		context := r.Context()
		//
		query := "INSERT INTO clients (last_name, first_name, phone, email, password) VALUES (?,?,?,?,?)"

		result, err := db.ExecContext(context, query, client_data.Last_name, client_data.First_name, client_data.Phone, client_data.Email, client_data.Password)
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

		// получить id нового польз

		query = "SELECT id_client FROM clients WHERE phone = ?"
		var id_client int
		err = db.QueryRowContext(context, query, client_data.Phone).Scan(&id_client)
		if err != nil {
			message := fmt.Sprintf("Ошибка SQL: %s", err)
			return &utils.AppError{
				Err:     err,
				Message: message,
				Code:    http.StatusInternalServerError,
			}
		}
		result_message := "Регистрация"
		return utils.AddCookie(result_message, db, id_client, w, r)

	case "GET":
		db := cdb.CDB.DB

		client_id_str := r.URL.Query().Get("id_client")
		client_id, err := strconv.Atoi(client_id_str)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Неверный id клиента",
				Code:    http.StatusInternalServerError,
			}
		}

		context := r.Context()
		query := "SELECT * FROM clients WHERE id_client =?"
		client_data := models.Client{}
		err = db.QueryRowContext(context, query, client_id).Scan(&client_data.IDClient, &client_data.Last_name, &client_data.First_name, &client_data.Phone, &client_data.Email, &client_data.Password, &client_data.Role)

		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка отправки запроса SELECT * FROM clients WHERE id",
				Code:    http.StatusInternalServerError,
			}
		}
		return utils.Form_response(w, client_data, http.StatusOK)

	case "PUT":
		db := cdb.CDB.DB

		client_data := models.Client{}
		err := json.NewDecoder(r.Body).Decode(&client_data)

		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Не удалось получить обновленные данные",
				Code:    http.StatusBadRequest,
			}
		}

		context := r.Context()
		query := "UPDATE clients SET last_name= ?,first_name= ?,phone=?,email=? WHERE id_client = ?"
		result, err := db.ExecContext(context, query, client_data.Last_name, client_data.First_name, client_data.Phone, client_data.Email, client_data.IDClient)

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
		} else {
			return utils.Form_response(w, "Данные обновлены!", http.StatusOK)
		}

	default:
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}
}
