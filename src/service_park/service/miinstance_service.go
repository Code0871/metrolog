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
func (s *MiInstanceService) GetAll(limit, offset int) ([]models.MiInstance, int, error) {
	return s.repo.GetAll(limit, offset)
}

// GetByPassport - возвращает средство измерения по его паспорту
// Параметры:
//   - passport: уникальный идентификатор (номер паспорта)
//
// Возвращает:
//   - *models.MiInstance: найденная запись
//   - error: ошибка если запись не найдена
func (s *MiInstanceService) GetByPassport(passport []string) ([]models.MiInstance, error) {
	return s.repo.GetByPassports(passport)
}
