// ============================================================================
// Пакет handlers обрабатывает HTTP запросы для средств измерения
//
// Назначение:
//   - Принимает запросы от фронтенда
//   - Вызывает соответствующие методы сервиса
//   - Формирует JSON ответ
//   - Обрабатывает ошибки
//
// Эндпоинты:
//   - GET  /api/miinstance           - список всех СИ (с пагинацией)
//   - GET  /api/miinstance/:passport - получение СИ по паспорту
// ============================================================================

package handler

import (
	"net/http"
	"service_park/models"
	"service_park/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type MIInstanceHandler struct {
	service *service.MiInstanceService
}

func NewMiInstanceHandler(service *service.MiInstanceService) *MIInstanceHandler {
	return &MIInstanceHandler{service: service}
}

// GetAll - GET /api/miinstance
// Возвращает список всех средств измерения с пагинацией
//
// Параметры запроса:
//   - limit: количество записей на странице (по умолчанию 20, максимум 200)
//   - offset: смещение для пагинации (по умолчанию 0)
//
// Ответ:
//
//	{
//	  "success": true,
//	  "data": [...],
//	  "total": 265336
//	}
func (h *MIInstanceHandler) GetAll(c *gin.Context) {
	const (
		defaultLimit = 20
		maxLimit     = 200
	)

	// парсинг параметров пагинации
	limit := defaultLimit
	offset := 0
	query := strings.TrimSpace(c.Query("query"))
	expiringRange := strings.TrimSpace(c.Query("expiring_range"))

	if expiringRange != "" && expiringRange != "all" && expiringRange != "week" && expiringRange != "month" && expiringRange != "year" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "expiring_range must be one of: all, week, month, year",
		})
		return
	}

	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			if val > maxLimit {
				limit = maxLimit
			} else {
				limit = val
			}
		}
	}

	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	// вызов сервиса для получения данных
	instances, total, err := h.service.GetAll(limit, offset, query, expiringRange)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{

			Success: false,
			Error:   "Error retrieving data: " + err.Error(),
		})
		return
	}

	// формирование ответа
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    instances,
		Total:   total,
	})

}

func (h *MIInstanceHandler) GetByPassport(c *gin.Context) {
	passportsStr := c.Query("idx")
	if passportsStr == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Passport parameter is required",
		})
		return
	}

	// Разделяем строку по запятой и очищаем пробелы
	rawPassports := strings.Split(passportsStr, ",")
	passports := make([]string, 0, len(rawPassports))
	for _, passport := range rawPassports {
		trimmedPassport := strings.TrimSpace(passport)
		if trimmedPassport != "" {
			passports = append(passports, trimmedPassport)
		}
	}

	if len(passports) == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Passport parameter is required",
		})
		return
	}

	// Вызываем метод для нескольких паспортов
	instances, err := h.service.GetByPassport(passports)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Error retrieving data: " + err.Error(),
		})
		return
	}

	if len(instances) == 0 {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "No records found for passports: " + passportsStr,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    instances,
		Total:   len(instances),
	})
}
