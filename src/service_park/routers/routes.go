package routers

import (
	"net/http"
	"service_park/handler"
	"service_park/repository"
	"service_park/service"

	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := c.Writer.Header()
		headers.Set("Access-Control-Allow-Origin", "*")
		headers.Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		headers.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(corsMiddleware())

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
			miinstance.GET("", miInstanceHandler.GetAll)
			miinstance.GET("/", miInstanceHandler.GetAll)
			miinstance.GET("/passport", miInstanceHandler.GetByPassport)
			miinstance.POST("/", miInstanceHandler.Create)
			miinstance.DELETE("/passport/:passport", miInstanceHandler.Delete)
			miinstance.DELETE("/passport", miInstanceHandler.DeleteMulti)
			miinstance.HEAD("/passport/:passport", miInstanceHandler.Head)
			miinstance.PUT("/passport/:passport", miInstanceHandler.Update)
		}
	}

	// ============================================
	// HEALTH CHECK
	// ============================================
	router.GET("/health", handler.HealthCheck)
	router.GET("/ping", handler.Ping)

	return router

}
