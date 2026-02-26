package payment

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"myproject/utils"
	"net/http"
	"strings"
)

type PaymentBase struct {
	PDB *utils.Base
}

type initiateRequest struct {
	OrderID int     `json:"id_order"`
	Amount  float64 `json:"amount"`
}

type confirmRequest struct {
	OrderID    int    `json:"id_order"`
	PaymentRef string `json:"payment_ref"`
}

type paymentResponse struct {
	Status     string `json:"status"`
	PaymentRef string `json:"payment_ref,omitempty"`
	OrderID    int    `json:"id_order,omitempty"`
	Message    string `json:"message,omitempty"`
}

// проверяет данные заказа и если все правильно выдает уникальный пропуск на оплату (payment_ref)
func (p PaymentBase) PaymentInitiateHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	if r.Method != http.MethodPost {
		return &utils.AppError{Err: errors.New("MethodNotAllowed"), Message: "Метод не поддерживается", Code: http.StatusMethodNotAllowed}
	}

	db := p.PDB.DB
	clientID, appErr := getAuthClientID(db, r)
	if appErr != nil {
		return appErr
	}

	var req initiateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &utils.AppError{Err: err, Message: "Некорректный формат данных", Code: http.StatusBadRequest}
	}
	if req.OrderID <= 0 || req.Amount <= 0 {
		return &utils.AppError{Err: errors.New("bad request"), Message: "order_id и amount должны быть > 0", Code: http.StatusBadRequest}
	}

	var ownerID int
	var total float64
	var payStatus sql.NullString // NULL в sql
	var payRef sql.NullString
	err := db.QueryRow(
		"SELECT id_client, total_cost, payment_status, payment_ref FROM orders WHERE id_order = ?",
		req.OrderID,
	).Scan(&ownerID, &total, &payStatus, &payRef)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &utils.AppError{Err: err, Message: "Заказ не найден", Code: http.StatusNotFound}
		}
		return &utils.AppError{Err: err, Message: "Ошибка SQL при получении заказа", Code: http.StatusInternalServerError}
	}

	if ownerID != clientID {
		return &utils.AppError{Err: errors.New("forbidden"), Message: "Чужой заказ", Code: http.StatusForbidden}
	}

	if payStatus.Valid && strings.EqualFold(payStatus.String, "paid") { // сравнение строк без учета регистра
		return utils.Form_response(w, paymentResponse{
			Status:     "paid",
			PaymentRef: payRef.String,
			OrderID:    req.OrderID,
			Message:    "Заказ уже оплачен",
		}, http.StatusOK)
	}

	if fmt.Sprintf("%.2f", total) != fmt.Sprintf("%.2f", req.Amount) {
		return &utils.AppError{
			Err:     errors.New("amount mismatch"),
			Message: "Сумма не совпадает с суммой заказа",
			Code:    http.StatusBadRequest,
		}
	}

	newRef := randomRef(12)
	_, err = db.Exec("UPDATE orders SET payment_status = 'pending', payment_ref = ? WHERE id_order = ?", newRef, req.OrderID)
	if err != nil {
		return &utils.AppError{Err: err, Message: "Не удалось сохранить платёж", Code: http.StatusInternalServerError}
	}

	return utils.Form_response(w, paymentResponse{
		Status:     "pending",
		PaymentRef: newRef,
		OrderID:    req.OrderID,
		Message:    "Платёж инициирован",
	}, http.StatusOK)
}

func (p PaymentBase) PaymentConfirmHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	if r.Method != http.MethodPost {
		return &utils.AppError{Err: errors.New("MethodNotAllowed"), Message: "Метод не поддерживается", Code: http.StatusMethodNotAllowed}
	}

	db := p.PDB.DB
	clientID, appErr := getAuthClientID(db, r)
	if appErr != nil {
		return appErr
	}

	var req confirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &utils.AppError{Err: err, Message: "Некорректный формат данных", Code: http.StatusBadRequest}
	}
	if req.OrderID <= 0 || strings.TrimSpace(req.PaymentRef) == "" {
		return &utils.AppError{Err: errors.New("bad request"), Message: "order_id > 0 и payment_ref обязателен", Code: http.StatusBadRequest}
	}

	var ownerID int
	var payStatus sql.NullString
	var payRef sql.NullString
	err := db.QueryRow(
		"SELECT id_client, payment_status, payment_ref FROM orders WHERE id_order = ?",
		req.OrderID,
	).Scan(&ownerID, &payStatus, &payRef)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &utils.AppError{Err: err, Message: "Заказ не найден", Code: http.StatusNotFound}
		}
		return &utils.AppError{Err: err, Message: "Ошибка SQL при получении заказа", Code: http.StatusInternalServerError}
	}

	if ownerID != clientID {
		return &utils.AppError{Err: errors.New("forbidden"), Message: "Чужой заказ", Code: http.StatusForbidden}
	}

	if payStatus.Valid && strings.EqualFold(payStatus.String, "paid") {
		return utils.Form_response(w, paymentResponse{
			Status:     "paid",
			PaymentRef: payRef.String,
			OrderID:    req.OrderID,
			Message:    "Уже оплачен",
		}, http.StatusOK)
	}

	if payRef.Valid && payRef.String != "" && payRef.String != req.PaymentRef {
		return &utils.AppError{
			Err:     errors.New("payment_ref mismatch"),
			Message: "Неверный payment_ref",
			Code:    http.StatusBadRequest,
		}
	}

	_, err = db.Exec("UPDATE orders SET payment_status = 'paid', payment_ref = ?, id_status = 1 WHERE id_order = ?", req.PaymentRef, req.OrderID)
	if err != nil {
		return &utils.AppError{Err: err, Message: "Не удалось обновить статус оплаты", Code: http.StatusInternalServerError}
	}

	return utils.Form_response(w, paymentResponse{
		Status:     "paid",
		PaymentRef: req.PaymentRef,
		OrderID:    req.OrderID,
		Message:    "Оплата подтверждена",
	}, http.StatusOK)
}

func getAuthClientID(db *sql.DB, r *http.Request) (int, *utils.AppError) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, &utils.AppError{Err: err, Message: "Пользователь не авторизирован", Code: http.StatusUnauthorized}
	}

	var clientID int
	query := "SELECT id_client FROM sessions WHERE token = ? AND expires_at > " + utils.NowSQL()
	if err := db.QueryRow(query, cookie.Value).Scan(&clientID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, &utils.AppError{Err: err, Message: "Сессия не найдена", Code: http.StatusUnauthorized}
		}
		return 0, &utils.AppError{Err: err, Message: "Ошибка проверки сессии", Code: http.StatusInternalServerError}
	}
	return clientID, nil
}

func randomRef(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "payref-fallback"
	}
	return hex.EncodeToString(buf) //превращает каждый байт в два символа шестнадцатеричной системы (0-9 и a-f)
}

/*Действия мошенника:
Умный школьник открывает консоль разработчика (F12) на твоем сайте.
Находит код отправки запроса.
Меняет сумму 5000 на 1 рубль.
Отправляет запрос.
Если сервер просто принимает деньги "одним махом", он спишет 1 рубль и закроет заказ на 5000.
Как спасает разделение (Твой код):
Initiate: Frontend говорит: "Хочу оплатить заказ №1". Сумму он не диктует.
Сервер сам лезет в БД, видит, что заказ стоит 5000, и создает в базе запись:
"Жду оплату 5000 рублей, вот секретный код подтверждения (payment_ref)".
Если мошенник потом попытается в этапе Confirm прислать другую сумму или другой код
— сервер сверит это с записью, созданной на этапе Initiate, увидит несовпадение и заблокирует операцию.*/
