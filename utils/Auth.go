package utils

import (
	//"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (base Base) AuthHandler(w http.ResponseWriter, r *http.Request) *AppError {

	switch r.Method {
	case "GET":
		db := base.DB
		user_id, role, err := GetUserData(w, r, db)

		if err != nil {
			return err
		}

		type Response struct {
			UserId int    `json:"user_id"`
			Role   string `json:"role"`
		}
		response := Response{
			UserId: user_id,
			Role:   role,
		}
		return Form_response(w, response, http.StatusOK)

	case "PUT": //выход из сессии
		db := base.DB

		cookie, err := r.Cookie("session_token")
		// проверяем есть ли у пользователя вообще пешенько
		if err != nil {
			return &AppError{
				Err:     err,
				Message: "Пользователь не авторизирован",
				Code:    http.StatusUnauthorized,
			}
		}
		token := cookie.Value
		context := r.Context()
		// обновляем запись сессии: срок действия меням на текущее время -> значит пользователь вышел
		query := "UPDATE sessions SET expires_at = ? WHERE token = ?"
		time_now := time.Now()
		result, err := db.ExecContext(context, query, time_now, token)

		if err != nil {
			message := fmt.Sprintf("Ошибка SQL: %s", err)
			return &AppError{
				Err:     err,
				Message: message,
				Code:    http.StatusInternalServerError,
			}
		}
		_, err = result.RowsAffected()
		if err != nil {
			return &AppError{
				Err:     err,
				Message: "Ошибка RowsAffected()",
				Code:    http.StatusInternalServerError,
			}
		}

		exitCookie := http.Cookie{
			Name:  "session_token",
			Value: "",
			Path:  "/",
			// HttpOnly: true, // [🤓]  прости ева...
			Secure:   false,
			MaxAge:   -1, //Если поставить 0 или отрицательное число — кука удалится сразу
			SameSite: http.SameSiteLaxMode,
		}
		//браузер удаляет печеньку
		http.SetCookie(w, &exitCookie)

		return Form_response(w, "Выход из учетной записи", http.StatusOK)

	default:
		return &AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}
}
