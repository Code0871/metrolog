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

// DeleteMultiple - удаляет записи по нескольким паспортам
func (r *MiInstanceRepository) DeleteMultiple(passports []string) (int64, error) {
	if len(passports) == 0 {
		return 0, fmt.Errorf("passports list is empty")
	}

	query, args, err := sqlx.In(`
	DELETE FROM miinstance WHERE miinstance_passport IN (?)`,
		passports)

	if err != nil {
		return 0, fmt.Errorf("sqlx.In error: %w", err)
	}

	query = db.DB.Rebind(query)
	result, err := db.DB.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("delete error: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected error: %w", err)
	}
	return rowsAffected, nil

}

// createMiInstance - создает новую запись в таблице miinstance
func (r *MiInstanceRepository) Create(inst *models.MiInstance) error {
	query := `
		INSERT INTO miinstance (
			miinstance_passport,
			miinstance_name,
			miinstance_type,
			miinstance_state_condition,
			miinstance_tech_condition,
			issue_date,
			commissioning_date,
			is_fit,
			mpi
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := db.DB.Exec(query,
		inst.Passport,
		inst.Name,
		inst.Type,
		inst.StateCondition,
		inst.TechCondition,
		inst.IssueDate,
		inst.CommissioningDate,
		inst.IsFit,
		inst.MPI,
	)
	if err != nil {
		return fmt.Errorf("insert error: %w", err)
	}

	return nil
}

// PUT - обновляет существующую запись по паспорту
func (r *MiInstanceRepository) Update(passport string, inst *models.MiInstance) error {
	query := `
		UPDATE miinstance SET
			miinstance_name = $1,
			miinstance_type = $2,
			miinstance_state_condition = $3,
			miinstance_tech_condition = $4,
			issue_date = $5,
			commissioning_date = $6,
			is_fit = $7,
			mpi = $8
		WHERE miinstance_passport = $9
	`

	result, err := db.DB.Exec(query,
		inst.Name,              // $1
		inst.Type,              // $2
		inst.StateCondition,    // $3
		inst.TechCondition,     // $4
		inst.IssueDate,         // $5
		inst.CommissioningDate, // $6
		inst.IsFit,             // $7
		inst.MPI,               // $8
		passport,               // $9
	)
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("record with passport %s not found", passport)
	}

	return nil
}

// DELETE - удаляет запись по паспорту
func (r *MiInstanceRepository) Delete(passport string) error {
	query := `DELETE FROM miinstance WHERE miinstance_passport = $1`

	result, err := db.DB.Exec(query, passport)
	if err != nil {
		return fmt.Errorf("delete error: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("record with passport %s not found", passport)
	}

	return nil
}

// Exists - проверяет существование записи по паспорту
func (r *MiInstanceRepository) Exists(passport string) (bool, error) {
	query := `SELECT COUNT(*) FROM miinstance WHERE miinstance_passport = $1`

	var exists bool
	err := db.DB.Get(&exists, query, passport)
	if err != nil {
		return false, fmt.Errorf("exists error: %w", err)
	}

	return exists, nil

}
