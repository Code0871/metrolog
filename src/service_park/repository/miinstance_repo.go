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

	"github.com/jmoiron/sqlx"
)

type MiInstanceRepository struct{}

func NewMiInstanceRepository() *MiInstanceRepository {
	return &MiInstanceRepository{}
}

// GetMiInstances - получает список всех единиц оборудования (MI) из базы данных.
func (r *MiInstanceRepository) GetAll(limit, offset int) ([]models.MiInstance, int, error) {

	query := `
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
		ORDER BY miinstance_passport
		LIMIT $1 OFFSET $2
	`

	var instances []models.MiInstance
	err := db.DB.Select(&instances, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query error: %w", err)
	}
	// Получаем общее количество записей для пагинации
	var total int
	countQuery := `SELECT COUNT(*) FROM miinstance`
	err = db.DB.Get(&total, countQuery)
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
