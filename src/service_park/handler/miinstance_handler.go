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
	"fmt"
	"net/http"
	"service_park/models"
	"service_park/service"
	"service_park/utils"
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
	passports := utils.ParseCommaSeparatedString(passportsStr)

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

func (h *MIInstanceHandler) DeleteMulti(c *gin.Context) {

	passportsStr := c.Query("idx")
	if passportsStr == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Passport parameter is required",
		})
		return
	}

	// Разделяем строку по запятой и очищаем пробелы
	passports := utils.ParseCommaSeparatedString(passportsStr)

	if len(passports) == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Passport parameter is required",
		})
		return
	}

	deleteCount, err := h.service.DeleteMultiple(passports)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Error deleting records: " + err.Error(),
		})
		return
	}

	if deleteCount == 0 {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "No records found for the provided passports",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: fmt.Sprintf("Successfully deleted %d records", deleteCount),
	})

}

// Create - POST /api/miinstance
// Создает новое средство измерения
//
// Тело запроса (JSON):
// {
//   "passport": "12345",
//   "name": "СИ-1",
//   "type": "Тип 1",
//   "state_condition": "Хорошее",
//   "tech_condition": "Исправно",
//   "issue_date": "2026-01-01",
//   "commissioning_date": "2026-02-01",
//   "is_fit": true,
//   "mpi": 100
// }

func (h *MIInstanceHandler) Create(c *gin.Context) {
	var req models.MiInstance
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Валидация данных
	if !req.Passport.Valid || req.Passport.String == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Passport is required",
		})
		return
	}

	// Вызов сервиса для создания новой записи
	err := h.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Error creating record: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Record created successfully",
		Data:    req,
	})

}

// DELETE /api/miinstance/passport/:passport
func (h *MIInstanceHandler) Delete(c *gin.Context) {

	passport := c.Param("passport")
	if passport == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Passport parameter is required",
		})
		return
	}

	if err := h.service.Delete(passport); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Error deleting record: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Record deleted successfully",
	})

}

// PUT /api/miinstance/passport/:passport
// обновляет существующую запись по паспорту
func (h *MIInstanceHandler) Update(c *gin.Context) {
	passport := c.Param("passport")
	if passport == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Passport parameter is required",
		})
		return
	}

	var req models.MiInstance
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if err := h.service.Update(passport, &req); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Error updating record: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Record updated successfully",
	})

}

// HEAD /api/miinstance/passport/:passport
// проверяет существование средства измерения по паспорту
func (h *MIInstanceHandler) Head(c *gin.Context) {
	passport := c.Param("passport")

	if passport == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "passport is required",
		})
		return
	}

	exists, err := h.service.Exist(passport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "record not found",
		})
		return
	}

	c.Header("X-Record-Exists", "true")
	c.Header("X-Record-Passport", passport)
	c.Status(http.StatusOK)
}
