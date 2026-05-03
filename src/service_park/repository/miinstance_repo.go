// ============================================================================
// Пакет repository реализует слой доступа к данным (Data Access Layer)
//
// Назначение:
//   - Изоляция бизнес-логики от конкретной реализации базы данных
//   - Выполнение SQL запросов к таблице miinstance
//   - Преобразование результатов запросов в модели данных
//
//
// Зависимости:
//   - service_park/db - глобальное соединение с БД
//   - service_park/models - структуры данных
// ============================================================================

package repository

import (
	"fmt"
	"service_park/db"
	"service_park/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

const miInstanceDueDateExpr = "CASE WHEN commissioning_date IS NOT NULL AND mpi IS NOT NULL THEN commissioning_date + make_interval(months => mpi) END"

type MiInstanceRepository struct{}

func NewMiInstanceRepository() *MiInstanceRepository {
	return &MiInstanceRepository{}
}

func buildMiInstanceWhereClause(query, expiringRange string) (string, []interface{}) {
	clauses := make([]string, 0, 2)
	args := make([]interface{}, 0, 3)

	if query != "" {
		likePattern := "%" + query + "%"
		clauses = append(clauses, "(miinstance_passport ILIKE ? OR miinstance_name ILIKE ? OR miinstance_type ILIKE ?)")
		args = append(args, likePattern, likePattern, likePattern)
	}

	if expiringRange != "" {
		switch expiringRange {
		case "week":
			clauses = append(clauses, "("+miInstanceDueDateExpr+" BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '7 days')")
		case "month":
			clauses = append(clauses, "("+miInstanceDueDateExpr+" BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '1 month')")
		case "year", "all":
			clauses = append(clauses, "("+miInstanceDueDateExpr+" BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '1 year')")
		}
	}

	if len(clauses) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(clauses, " AND "), args
}

// GetMiInstances - получает список всех единиц оборудования (MI) из базы данных.
func (r *MiInstanceRepository) GetAll(limit, offset int, query, expiringRange string) ([]models.MiInstance, int, error) {
	whereClause, args := buildMiInstanceWhereClause(query, expiringRange)
	orderByClause := "miinstance_passport"
	if expiringRange != "" {
		orderByClause = "COALESCE(" + miInstanceDueDateExpr + ", DATE '9999-12-31'), miinstance_passport"
	}

	selectQuery := `
		SELECT 
			miinstance_passport,
			miinstance_name,
			miinstance_type,
			miinstance_state_condition,
			miinstance_tech_condition,
			issue_date,
			commissioning_date,
			is_fit,
			mpi
		FROM miinstance
		` + whereClause + `
		ORDER BY ` + orderByClause + `
		LIMIT ? OFFSET ?
	`

	var instances []models.MiInstance
	selectArgs := append(append([]interface{}{}, args...), limit, offset)
	err := db.DB.Select(&instances, db.DB.Rebind(selectQuery), selectArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query error: %w", err)
	}
	// Получаем общее количество записей для пагинации
	var total int
	countQuery := `SELECT COUNT(*) FROM miinstance ` + whereClause
	err = db.DB.Get(&total, db.DB.Rebind(countQuery), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("count error: %w", err)
	}

	return instances, total, nil
}

// GetByPassports - получает список по нескольким паспортам
func (r *MiInstanceRepository) GetByPassports(passports []string) ([]models.MiInstance, error) {
	if len(passports) == 0 {
		return []models.MiInstance{}, nil
	}

	query, args, err := sqlx.In(`
		SELECT 
			miinstance_passport,
			miinstance_name,
			miinstance_type,
			miinstance_state_condition,
			miinstance_tech_condition,
			issue_date,
			commissioning_date,
			is_fit,
			mpi
		FROM miinstance
		WHERE miinstance_passport IN (?)
	`, passports)

	if err != nil {
		return nil, fmt.Errorf("sqlx.In error: %w", err)
	}

	query = db.DB.Rebind(query)

	var instances []models.MiInstance
	err = db.DB.Select(&instances, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return instances, nil
}
