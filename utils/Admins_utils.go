package utils

import (
	"database/sql"
	"net/http"
)

// проверка токена и роли для функций админа
func GetUserData(w http.ResponseWriter, r *http.Request, db *sql.DB) (int, string, *AppError) {

	//делаем запрос по Id клиента и проверяем роль
	// возвращаем ошибку если нет прав админа

	context := r.Context()

	cookie, err := r.Cookie("session_token")
	// проверяем есть ли у пользователя вообще пешенько
	if err != nil {
		return 0, "", &AppError{
			Err:     err,
			Message: "Пользователь не авторизирован",
			Code:    http.StatusUnauthorized,
		}
	}
	token := cookie.Value

	// проверяем есть ли токен пешенька в базе
	//query := "SELECT id_client FROM sessions WHERE token = ? AND expires_at > " + NowSQL() // проверяем не истек ли токен пешеньки
	query := "SELECT s.id_client, c.role FROM sessions s JOIN clients c ON s.id_client = c.id_client WHERE s.token = ? AND s.expires_at > " + NowSQL()

	var role string
	var id_client int

	err = db.QueryRowContext(context, query, token).Scan(&id_client, &role)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", &AppError{
				Err:     err,
				Message: "Сессия не найдена",
				Code:    http.StatusUnauthorized, //401
			}
		} else {
			return 0, "", &AppError{
				Err:     err,
				Message: "Ошибка SQL запроса SELECT ... FROM sessions WHERE token = ? ...",
				Code:    http.StatusInternalServerError,
			}
		}
	}

	return id_client, role, nil
}
