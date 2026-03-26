package main

import (
	"context"
	"database/sql"
	"log"
	"myproject/admin"
	basecolorpalette "myproject/base_color_palette"
	"myproject/bouquet"
	"myproject/client"
	"myproject/database"
	"myproject/flowers"
	"myproject/occasions"
	"myproject/order"
	orderitems "myproject/order_items"
	orderstatuses "myproject/order_statuses"
	"myproject/payment"
	"myproject/reviews"
	"myproject/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config := database.GetDBConfig()

	db, err := database.Connect(config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := database.InitSchema(db); err != nil {
		log.Fatal(err)
	}

	// фоновая очистка сессий
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			cleanupSessions(db)
		}
	}()

	base := utils.Base{
		DB: db,
	}

	bouquet_base := bouquet.BouquetBase{
		BDB: &base,
	}

	client_base := client.ClientBase{
		CDB: &base,
	}

	order_base := order.OrderBase{
		ODB: &base,
	}

	status_base := orderstatuses.OrderStatusesBase{
		SDB: &base,
	}

	item_base := orderitems.ItemBase{
		IDB: &base,
	}

	payment_base := payment.PaymentBase{
		PDB: &base,
	}

	admin_base := admin.AdminBase{
		ADB: &base,
	}

	review_base := reviews.ReviewBase{
		RDB: &base,
	}

	color_base := basecolorpalette.ColorBase{
		CDB: &base,
	}

	occasion_base := occasions.OccasionBase{
		ODB: &base,
	}

	flower_base := flowers.FlowerBase{
		FDB: &base,
	}

	// Настраиваем маршруты
	mux := http.NewServeMux()
	mux.HandleFunc("/bouquets", utils.ErrorHandler(bouquet_base.BouquetsHandler))
	mux.HandleFunc("/bouquets/", utils.ErrorHandler(bouquet_base.BouquetIdHandler))
	mux.HandleFunc("/clients", utils.ErrorHandler(client_base.ClientHandler))                   // регистрация клиента
	mux.HandleFunc("/login", utils.ErrorHandler(base.LogInHandler))                             // проверка пароля
	mux.HandleFunc("/orders", utils.ErrorHandler(order_base.OrderHandler))                      // "мои заказы"
	mux.HandleFunc("/order_statuses", utils.ErrorHandler(status_base.OrderStatusHandler))       // статусы заказов
	mux.HandleFunc("/order_items", utils.ErrorHandler(item_base.OrderIdHandler))                // GET содержимое заказа по id заказа с картинкой из bouquets
	mux.HandleFunc("/payments", utils.ErrorHandler(payment_base.PaymentInitiateHandler))        // инициация оплаты
	mux.HandleFunc("/payments/confirm", utils.ErrorHandler(payment_base.PaymentConfirmHandler)) // подтверждение оплаты (мок)
	mux.HandleFunc("/auth_access", utils.ErrorHandler(base.AuthHandler))
	mux.HandleFunc("/cart", utils.ErrorHandler(bouquet_base.ListBouquetsHandler))         // POST со списком id букетов для корзины
	mux.HandleFunc("/admin/orders", utils.ErrorHandler(admin_base.OrdersForAdminHandler)) // все заказы для панели админа
	mux.HandleFunc("/admin/status", utils.ErrorHandler(admin_base.OrderIdHandler))        //обновление статуса заказа
	mux.HandleFunc("/reviews/", utils.ErrorHandler(review_base.ReviewsHandler))           //отзывы
	mux.HandleFunc("/reviews/access", utils.ErrorHandler(review_base.ReviewAccessCheckHandler))
	mux.HandleFunc("/order_data/", utils.ErrorHandler(order_base.OrderIdHandler)) // страница заказа

	//----------ФИЛЬТРЫ------
	mux.HandleFunc("/colors", utils.ErrorHandler(color_base.ColorsHandler))                         // список colors для фильтров
	mux.HandleFunc("/occasions", utils.ErrorHandler(occasion_base.OccasionsHandler))                //список назначений букетов
	mux.HandleFunc("/flowers", utils.ErrorHandler(flower_base.FlowersHandler))                      // спписок цветков
	mux.HandleFunc("/price_extremes", utils.ErrorHandler(bouquet_base.MaxMinPricesBouquetsHandler)) // максимальная и минимальная цена букетов для фильтра

	// Создаем файловый сервер, который смотрит в папку "./static"
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fileServer)

	// Создаём HTTP сервер
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Println("Сервер запущен на http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ждём сигнал завершения (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении сервера: %v", err)
	}

	log.Println("Сервер остановлен")
}

// cleanupSessions удаляет истёкшие сессии из БД
func cleanupSessions(db *sql.DB) {
	query := "DELETE FROM sessions WHERE expires_at < " + utils.NowSQL()
	result, err := db.Exec(query)
	if err != nil {
		log.Printf("Ошибка очистки сессий: %v", err)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("Ошибка получения количества удалённых сессий: %v", err)
		return
	}

	if rows > 0 {
		log.Printf("Удалено %d истёкших сессий", rows)
	}
}
