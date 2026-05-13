// Package repository содержит интерфейсы для работы с хранилищем
// автомобильных номеров (Plate).
package repository

import (
	models "encomers/license/internal/domain/valueObjects"
)

// IRepository — основной интерфейс репозитория для хранения
// и работы с сгенерированными автомобильными номерами.
//
// Используется генераторами номеров для сохранения состояния,
// проверки дублей и восстановления последнего сгенерированного номера.
type IRepository interface {
	// GetLastPlate возвращает последний сохранённый номер.
	//
	// Возвращаемые значения:
	//   - models.Plate — последний номер
	//   - bool — true, если номер найден (репозиторий не пуст)
	//   - error — ошибка при обращении к хранилищу
	GetLastPlate() (models.Plate, bool, error)

	// SavePlate сохраняет номер в репозиторий.
	SavePlate(plate models.Plate) error

	// CheckPlate проверяет, существует ли уже такой номер в репозитории.
	//
	// Возвращает true, если номер уже присутствует (дубликат).
	CheckPlate(plate models.Plate) (bool, error)
}
