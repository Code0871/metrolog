// ============================================================================
// Пакет models определяет структуры данных для работы с API и базой данных
//
// Назначение:
//   - MiInstance - основная бизнес-сущность (средство измерения)
//   - Response - стандартный формат ответа API
//   - ListRequest - параметры пагинации для запросов
//
// Теги:
//   - json: сериализация в JSON
//   - db: маппинг на колонки PostgreSQL
//   - form: привязка параметров HTTP запроса (Gin)
//
// ============================================================================
package models

import (
	"service_park/types"
	"time"
)

// MiInstance - полная модель (соответствует таблице miinstance)
type MiInstance struct {
	Passport          types.NullString `json:"passport" db:"miinstance_passport"`
	Name              types.NullString `json:"name" db:"miinstance_name"`
	Type              types.NullString `json:"type" db:"miinstance_type"`
	StateCondition    types.NullString `json:"state_condition" db:"miinstance_state_condition"`
	TechCondition     types.NullString `json:"tech_condition" db:"miinstance_tech_condition"`
	IssueDate         *time.Time       `json:"issue_date,omitempty" db:"issue_date"`
	CommissioningDate *time.Time       `json:"commissioning_date,omitempty" db:"commissioning_date"`
	IsFit             types.NullBool   `json:"is_fit" db:"is_fit"`
	MPI               types.NullInt32  `json:"mpi" db:"mpi"`
}

// Response - универсальный формат ответа API
//
// Пример успешного ответа:
//
//	{
//	  "success": true,
//	  "data": [...],
//	  "total": 265336
//	}
//
// Пример ответа с ошибкой:
//
//	{
//	  "success": false,
//	  "error": "record not found"
//	}
type Response struct {
	Success bool        `json:"success"`         // Статус выполнения запроса (true/false)
	Data    interface{} `json:"data,omitempty"`  // Полезные данные (любого типа)
	Total   int         `json:"total,omitempty"` // Общее количество записей (для пагинации)
	Error   string      `json:"error,omitempty"` // Текст ошибки, если success = false
}

// ListRequest - параметры пагинации для GET /api/miinstance
//
// Пример использования:
//
//	GET /api/miinstance?limit=20&offset=40
//
// Это вернёт 20 записей, начиная с 41-й (пропустив первые 40)
//
// Параметры:
//   - limit: количество записей на странице (max 200)
//   - offset: количество пропускаемых записей
type ListRequest struct {
	Limit  int `json:"limit" form:"limit"`   // Максимум записей в ответе
	Offset int `json:"offset" form:"offset"` // Смещение от начала
}
