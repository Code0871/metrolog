// ============================================================================
// Пакет handlers - health check endpoints
//
// Эндпоинты:
//   - GET /health - проверка работоспособности сервиса
//   - GET /ping   - простой ping
// ============================================================================

package handler

import (
	"net/http"
	"service_park/db"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck - проверка состояния сервиса
// Возвращает статус сервера и соединения с БД
func HealthCheck(c *gin.Context) {

	status := "ok"
	dbStatus := "ok"

	if db.DB != nil {
		if err := db.DB.Ping(); err != nil {
			dbStatus = "error: " + err.Error()
			status = "error"
		}
	} else {
		dbStatus = "error: database not initialized"
		status = "degraded"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    status,
		"service":   "service_park",
		"database":  dbStatus,
		"timestamp": time.Now().Unix(),
	})
}

// Ping - простой ping для проверки доступности
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
