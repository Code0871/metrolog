// ============================================================
// Файл: connect.go
// Назначение: Подключение и управление соединением с PostgreSQL
// Библиотеки: sqlx (работа с БД), lib/pq (драйвер PostgreSQL)
// ============================================================
//
// Что делает этот файл:
//   1. Устанавливает соединение с базой данных PostgreSQL
//   2. Настраивает пул соединений (максимум соединений, время жизни)
//   3. Предоставляет глобальный доступ к БД через переменную DB
//   4. Закрывает соединение при завершении работы
// Пример использования:
//   cfg := config.LoadConfig()
//   err := database.ConnectDB(cfg)
//   defer database.CloseDB()

package db

import (
	"fmt"
	"service_park/config"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB(cfg *config.Config) error {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	fmt.Printf("Connection string: %s\n", connStr)

	var err error

	DB, err = sqlx.Connect("postgres", connStr)

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// SetMaxOpenConns устанавливает максимальное количество одновременных соединений с БД
	// 25 - при превышении запросы ждут в очереди
	DB.SetMaxOpenConns(25)

	// SetMaxIdleConns устанавливает количество простаивающих соединений в пуле
	// 5 - всегда держать 5 готовых соединений для быстрых запросов
	DB.SetMaxIdleConns(5)

	// SetConnMaxLifetime устанавливает максимальное время жизни соединения
	// 5 минут - после этого соединение закрывается и создаётся новое
	// Это предотвращает использование "залипших" или разорванных соединений
	DB.SetConnMaxLifetime(5 * time.Minute)

	return nil

}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
