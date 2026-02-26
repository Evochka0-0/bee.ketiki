package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"myproject/models"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func AddCookie(operation string, db *sql.DB, id int, w http.ResponseWriter, r *http.Request) *AppError {
	//создает сессию в бд и печеньку с токеном

	context := r.Context()

	expires_at := time.Now().Add(24 * time.Hour)
	token := uuid.NewString()

	query := "INSERT INTO sessions (id_client, token, expires_at) VALUES (?,?,?)"
	result, err := db.ExecContext(context, query, id, token, expires_at)

	if err != nil {
		message := fmt.Sprintf("Ошибка SQL: %s", err)
		return &AppError{
			Err:     err,
			Message: message,
			Code:    http.StatusInternalServerError,
		}
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return &AppError{
			Err:     err,
			Message: "Ошибка RowsAffected()",
			Code:    http.StatusInternalServerError,
		}
	}
	if rows == 0 {
		return &AppError{
			Err:     errors.New("ноль строк изменено"),
			Message: "Затронуто 0 строк",
			Code:    http.StatusInternalServerError,
		}
	}

	// создание печеньки
	cookie := http.Cookie{
		Name:  "session_token",
		Value: token,
		Path:  "/", //куки доступна на всех страницах сайта
		//HttpOnly: true,
		Secure:   false,
		MaxAge:   86400,                //время действия 24 часа
		SameSite: http.SameSiteLaxMode, // защита от подделки запросов
	}
	http.SetCookie(w, &cookie)

	result_message := fmt.Sprintf("Успешная %s", operation)

	return Form_response(w, result_message, http.StatusOK)
}

func (ldb *Base) LogInHandler(w http.ResponseWriter, r *http.Request) *AppError {
	switch r.Method {
	case "POST":
		db := ldb.DB
		var login_data = models.LogInData{} //здесь записан  телефон и введенный пароль string
		err := json.NewDecoder(r.Body).Decode(&login_data)

		if err != nil {
			return &AppError{
				Err:     err,
				Message: fmt.Sprintf("Не удалось получить данные, введенные пользователем, %s", err),
				Code:    http.StatusInternalServerError,
			}
		}
		context := r.Context()
		query := "SELECT id_client, password FROM `clients` WHERE phone = ? "
		var id int
		var hash string //находим hash правильного пароля
		err = db.QueryRowContext(context, query, login_data.Phone).Scan(&id, &hash)
		if err != nil {
			if err == sql.ErrNoRows {
				return &AppError{
					Err:     err,
					Message: "Телефон не верный",
					Code:    http.StatusUnauthorized, // 401 - Нет доступа
				}
			}
			return &AppError{
				Err:     err,
				Message: "Ошибка сервера при поиске",
				Code:    http.StatusInternalServerError, // 500
			}
		}

		check_password := bcrypt.CompareHashAndPassword([]byte(hash), []byte(login_data.Password))

		if check_password != nil {
			return &AppError{
				Err:     check_password,
				Message: "Пароль не верный",
				Code:    http.StatusUnauthorized,
			}
		}
		result_message := "Авторизация"
		return AddCookie(result_message, db, id, w, r)

	default:
		return &AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}
}
