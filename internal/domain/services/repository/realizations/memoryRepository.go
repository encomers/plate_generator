// Package realizations содержит реализации основных компонентов
// для работы с генерацией и хранением автомобильных номеров.
package realizations

import (
	"sync"

	models "encomers/license/internal/domain/valueObjects"
)

// MemoryRepository — in-memory реализация репозитория.
//
// Хранит все номера в слайсе. Предназначена в первую очередь
// для тестирования и разработки.
type MemoryRepository struct {
	plates []models.Plate
	mu     sync.Mutex
}

// New создаёт новый экземпляр in-memory репозитория.
func New() *MemoryRepository {
	return &MemoryRepository{
		plates: make([]models.Plate, 0, 10_000),
	}
}

// GetLastPlate возвращает последний сохранённый номер.
//
// В текущей реализации возвращает первый элемент (поведение может быть изменено).
func (r *MemoryRepository) GetLastPlate() (models.Plate, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.plates) == 0 {
		return models.Plate{}, false, nil
	}

	plate := r.plates[len(r.plates)-1]
	return plate, true, nil
}

// SavePlate сохраняет номер в репозиторий.
func (r *MemoryRepository) SavePlate(plate models.Plate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.plates = append(r.plates, plate)
	return nil
}

// CheckPlate проверяет, существует ли уже такой номер в репозитории.
func (r *MemoryRepository) CheckPlate(plate models.Plate) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range r.plates {
		if p.String() == plate.String() {
			return true, nil
		}
	}
	return false, nil
}
