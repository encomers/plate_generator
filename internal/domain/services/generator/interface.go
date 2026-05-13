// Package generator содержит интерфейсы для генераторов автомобильных номеров.
//
// Основная задача пакета — предоставить унифицированные контракты
// для генерации номеров в формате: Буквы + Цифры + Регион.
package generator

import (
	models "encomers/license/internal/domain/valueObjects"
)

// INumberGenerator — интерфейс генератора цифровой части номера.
type INumberGenerator interface {
	// GetNext возвращает следующее число по порядку.
	GetNext() models.NumberPart

	// GetRandom возвращает случайное число в допустимом диапазоне.
	GetRandom() models.NumberPart

	// IsOverflow возвращает true, если достигнут максимум и произошло переполнение.
	IsOverflow() bool
}

// ILetterGenerator — интерфейс генератора буквенной части номера.
type ILetterGenerator interface {
	// GetNext возвращает следующую буквенную комбинацию в лексикографическом порядке.
	GetNext() models.LetterPart

	// GetRandom возвращает случайную буквенную комбинацию.
	GetRandom() models.LetterPart

	// IsOverflow возвращает true, если все возможные буквенные комбинации исчерпаны.
	IsOverflow() bool

	// Size возвращает размер алфавита (количество доступных символов).
	Size() int
}

// IRegionGenerator — интерфейс генератора региональной (код региона) части номера.
type IRegionGenerator interface {
	// SetRegion позволяет изменить регион (поддерживается не всеми реализациями).
	SetRegion(region models.RegionPart)

	// Get возвращает текущий регион.
	Get() models.RegionPart

	// IsOverflow возвращает true в случае проблем с регионом.
	// Для статического генератора всегда false.
	IsOverflow() bool
}

// IPlateGenerator — основной интерфейс генератора полного автомобильного номера.
//
// Генерирует номера в формате: [Буквы][Цифры][Регион]
type IPlateGenerator interface {
	// GetNext возвращает следующий номер по порядку (инкрементальный).
	// Возвращает ошибку, если генерация невозможна (например, общее переполнение).
	GetNext() (models.Plate, error)

	// GetRandom возвращает полностью случайный валидный номер.
	// Возвращает ошибку, если генерация по какой-либо причине невозможна.
	GetRandom() (models.Plate, error)

	// IsOverflow возвращает true, если все возможные номера в текущей конфигурации
	// (буквы + цифры + регион) исчерпаны.
	IsOverflow() bool
}
