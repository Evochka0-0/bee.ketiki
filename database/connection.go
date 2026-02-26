package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	//_ "modernc.org/sqlite"
)

type DBConfig struct {
	Driver string
	DSN    string
}

// читает конфигурацию из переменных окружения
func GetDBConfig() DBConfig {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "mysql"
	}

	var dsn string
	dsn = os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/bee.кетики1?parseTime=true&multiStatements=true"
	}

	return DBConfig{
		Driver: driver,
		DSN:    dsn,
	}
}

// устанавливаем соединение с БД
func Connect(config DBConfig) (*sql.DB, error) {
	db, err := sql.Open(config.Driver, config.DSN)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия БД: %w", err)
	}

	// Настройки пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	log.Printf("Подключение к БД (%s) успешно установлено!", config.Driver)

	return db, nil
}

func InitSchema(db *sql.DB) error {
	var schema string
	schema = mysqlSchema

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("ошибка инициализации схемы: %w", err)
	}

	log.Println("Схема БД успешно инициализирована")
	return nil
}

const mysqlSchema = `
CREATE TABLE IF NOT EXISTS clients (
    id_client INT AUTO_INCREMENT PRIMARY KEY,
    last_name VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(191) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS sessions (
    id_session INT AUTO_INCREMENT PRIMARY KEY,
    id_client INT NOT NULL,
    token VARCHAR(191) NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (id_client) REFERENCES clients(id_client) ON DELETE CASCADE,
    INDEX idx_token (token),
    INDEX idx_expires (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS bouquets (
    id_bouquet INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    image_url VARCHAR(500)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS orderstatuses (
    id_status INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


INSERT IGNORE INTO orderstatuses (id_status, name) VALUES
    (1, 'Новый'),
    (2, 'Собирается'),
    (3, 'Готов к выдаче'),
    (4, 'Выдан'),
    (5, 'Отменён');

CREATE TABLE IF NOT EXISTS orders (
    id_order INT AUTO_INCREMENT PRIMARY KEY,
    id_client INT NOT NULL,
    id_status INT NOT NULL DEFAULT 1, 
    total_cost DECIMAL(10,2) NOT NULL,
    payment_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    payment_ref VARCHAR(100),
    pickup_datetime DATETIME,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    /* [!] ИСПРАВЛЕНО: ссылки на новые ID */
    FOREIGN KEY (id_client) REFERENCES clients(id_client) ON DELETE CASCADE,
    FOREIGN KEY (id_status) REFERENCES orderstatuses(id_status) ON DELETE RESTRICT,
    INDEX idx_client (id_client),
    INDEX idx_status (id_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS reviews (
    id_review INT AUTO_INCREMENT PRIMARY KEY,
    id_client INT NOT NULL,
    id_bouquet INT NOT NULL,
    message TEXT,
    grade TINYINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_client) REFERENCES clients(id_client) ON DELETE CASCADE,
    FOREIGN KEY (id_bouquet) REFERENCES bouquets(id_bouquet) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS orderitems (
    id_items INT AUTO_INCREMENT PRIMARY KEY,
    id_order INT NOT NULL,
    id_bouquet INT NOT NULL,
    quantity INT NOT NULL,
    /* [!] ИСПРАВЛЕНО: ссылки на новые ID */
    FOREIGN KEY (id_order) REFERENCES orders(id_order) ON DELETE CASCADE,
    FOREIGN KEY (id_bouquet) REFERENCES bouquets(id_bouquet) ON DELETE RESTRICT,
    INDEX idx_order (id_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT IGNORE INTO bouquets (id_bouquet, name, description, price, image_url) VALUES
    (1, 'Букет "Весна"', 'Нежные тюльпаны и нарциссы', 2500.00, '/images/roses1.jpg'),
    (2, 'Букет "Романтика"', 'Красные розы для любимых', 3500.00, '/images/roses2.jpg'),
    (3, 'Букет "Солнечный"', 'Яркие герберы и хризантемы', 2000.00, '/images/piones1.jpg');
`
