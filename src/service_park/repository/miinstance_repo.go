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

// GetByPassport - получает полную информацию о средстве измерения по паспорту
func (r *MiInstanceRepository) GetByPassport(passport string) (*models.MiInstance, error) {
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
		WHERE miinstance_passport = $1
	`
	var instance models.MiInstance
	err := db.DB.Get(&instance, query, passport)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return &instance, nil
}
