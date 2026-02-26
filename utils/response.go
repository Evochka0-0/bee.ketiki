package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"
	//"os"
)

type Base struct {
	DB *sql.DB
}

func NowSQL() string {
	return "NOW()"
}

type AppError struct {
	Err     error
	Message string
	Code    int
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

type AppHandler func(http.ResponseWriter, *http.Request) *AppError

func ErrorHandler(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.Code)

			jsonResponse := struct {
				Error string `json:"error"`
			}{
				Error: err.Message,
			}
			if json.NewEncoder(w).Encode(jsonResponse) != nil {
				// Если даже отправка JSON не удалась, отправляем простой текст.
				http.Error(w, "Критическая ошибка сервера", http.StatusInternalServerError)
			}
		}
	}
}

func Form_response(w http.ResponseWriter, objects interface{}, statusCode int) *AppError {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(objects)
	if err != nil {
		return &AppError{
			Err:     err,
			Message: "Ошибка при формировании JSON-ответа",
			Code:    http.StatusInternalServerError,
		}
	}
	return nil
}
