// ============================================================================
// Главный модуль приложения Service Park
//
// Назначение:
//   - Загрузка конфигурации
//   - Подключение к базе данных
//   - Запуск HTTP сервера
//
// ============================================================================
package main

import (
	"log"
	"service_park/config"
	"service_park/db"
	"service_park/routers"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Подключение к базе данных
	if err := db.ConnectDB(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.CloseDB()
	log.Println("Database connected successfully")

	// Настройка маршрутов и запуск HTTP сервера
	router := routers.SetupRouter()

	log.Printf("Server starting on port %s", cfg.ServerPort)

	// Запуск сервера
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
