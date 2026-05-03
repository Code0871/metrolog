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
	// парсинг параметров пагинации
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val > 0 {
			offset = val
		}
	}

	// вызов сервиса для получения данных
	instances, total, err := h.service.GetAll(limit, offset)

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
	passport := c.Param("passport")

	if passport == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Passport parameter is required",
		})
		return
	}

	instance, err := h.service.GetByPassport(passport)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Error retrieving data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    instance,
	})

}
