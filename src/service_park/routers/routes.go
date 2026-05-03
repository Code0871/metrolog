package routers

import (
	"service_park/handler"
	"service_park/repository"
	"service_park/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// ============================================
	// ИНИЦИАЛИЗАЦИЯ ХЕНДЛЕРОВ
	// ============================================

	miInstaceRepo := repository.NewMiInstanceRepository()               // Инициализация репозитория для средств измерения
	miInstaceService := service.NewMiInstanceService(miInstaceRepo)     // Инициализация сервиса для средств измерения
	miInstanceHandler := handler.NewMiInstanceHandler(miInstaceService) // Инициализация хендлера для средств измерения

	// ============================================
	// API МАРШРУТЫ
	// ============================================

	api := router.Group("/api")
	{
		miinstance := api.Group("/miinstance")
		{
			miinstance.GET("/", miInstanceHandler.GetAll)
			miinstance.GET("/passport", miInstanceHandler.GetByPassport)
		}
	}

	// ============================================
	// HEALTH CHECK
	// ============================================
	router.GET("/health", handler.HealthCheck)
	router.GET("/ping", handler.Ping)

	return router

}
