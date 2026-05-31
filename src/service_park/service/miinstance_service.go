// ============================================================================
// Пакет service реализует слой бизнес-логики приложения
//
// Назначение:
//   - Прослойка между хендлерами и репозиторием
//   - Подготовка и валидация данных перед отправкой в репозиторий
//   - Вычисление производных полей (остаточный срок, статусы) **в планах написать**
//   - Логирование и обработка бизнес-ошибок
//
// Зависимости:
//   - service_park/repository - слой доступа к данным
//   - service_park/models - структуры данных
//
// Примечание:
//   - В текущей версии реализована базовая прослойка и всё
//   - При необходимости расширяется бизнес-логикой и расчётами надеюсь что не нужно
//
// ============================================================================
package service

import (
	"fmt"
	"service_park/models"
	"service_park/repository"
)

// MiInstanceService - сервис для работы со средствами измерения
// Содержит бизнес-логику приложения
type MiInstanceService struct {
	repo *repository.MiInstanceRepository
}

// NewMiInstanceService - конструктор сервиса
// Принимает репозиторий и возвращает новый экземпляр сервиса
func NewMiInstanceService(repo *repository.MiInstanceRepository) *MiInstanceService {
	return &MiInstanceService{repo: repo}
}

// GetAll - возвращает список всех средств измерения с пагинацией
// Параметры:
//   - limit: количество записей на странице
//   - offset: смещение от начала
//
// Возвращает:
//   - []models.MiInstance: список записей
//   - int: общее количество записей в БД
//   - error: ошибка при выполнении запроса
func (s *MiInstanceService) GetAll(limit, offset int, query, expiringRange string) ([]models.MiInstance, int, error) {
	return s.repo.GetAll(limit, offset, query, expiringRange)
}

// GetByPassport - возвращает средство измерения по его паспорту
// Параметры:
//   - passport: уникальный идентификатор (номер паспорта)
//
// Возвращает:
//   - *models.MiInstance: найденная запись
//   - error: ошибка если запись не найдена
func (s *MiInstanceService) GetByPassport(passport []string) ([]models.MiInstance, error) {
	if len(passport) == 0 {
		return nil, fmt.Errorf("passport is required")
	}
	return s.repo.GetByPassports(passport)

}

func (s *MiInstanceService) DeleteMultiple(passports []string) (int64, error) {
	if len(passports) == 0 {
		return 0, fmt.Errorf("passport is required")
	}
	return s.repo.DeleteMultiple(passports)
}

// Create - создает новое средство измерения
// Параметры:
//   - inst: структура данных нового средства измерения
//
// Возвращает:
//   - error: ошибка при создании записи
func (s *MiInstanceService) Create(inst *models.MiInstance) error {
	if !inst.Passport.Valid || inst.Passport.String == "" {
		return fmt.Errorf("passport is required")
	}

	if !inst.Name.Valid || inst.Name.String == "" {
		return fmt.Errorf("name is required")
	}

	return s.repo.Create(inst)

}

// Delete - удаляет средство измерения по паспорту
// Параметры:
//   - passport: уникальный идентификатор (номер паспорта)
//
// Возвращает:
//   - error: ошибка если паспорт не указан
func (s *MiInstanceService) Delete(passport string) error {
	if passport == "" {
		return fmt.Errorf("passport is required")
	}
	return s.repo.Delete(passport)
}

// Exist - проверяет существование средства измерения по паспорту
// Параметры:
//   - passport: уникальный идентификатор (номер паспорта)
//
// Возвращает:
//   - bool: true если средство измерения существует, false если нет
//   - error: ошибка при выполнении запроса
func (s *MiInstanceService) Exist(passport string) (bool, error) {
	if passport == "" {
		return false, fmt.Errorf("passport is required")
	}
	return s.repo.Exists(passport)
}

// Update - обновляет запись по паспорту
func (s *MiInstanceService) Update(passport string, inst *models.MiInstance) error {
	if passport == "" {
		return fmt.Errorf("passport is required")
	}
	return s.repo.Update(passport, inst)
}
