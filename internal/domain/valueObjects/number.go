// Package valueObjects содержит основные value objects,
// представляющие составные части автомобильного номера.
package valueObjects

import (
	"errors"
	"fmt"
)

const (
	// MAX_NUMBER_PART_VALUE — максимальное значение числовой части номера.
	// Обычно используется формат 001-999.
	MAX_NUMBER_PART_VALUE = 999
)

// NumberPart представляет числовую часть автомобильного номера (обычно 3 цифры).
type NumberPart struct {
	number int
}

// NewNumberPart создаёт новую числовую часть номера.
//
// Параметры:
//   - number — целое число от 0 до MAX_NUMBER_PART_VALUE включительно.
//
// Возвращает ошибку, если значение выходит за допустимый диапазон.
func NewNumberPart(number int) (NumberPart, error) {
	if number < 0 || number > MAX_NUMBER_PART_VALUE {
		return NumberPart{}, errors.New("invalid number part")
	}

	return NumberPart{
		number: number,
	}, nil
}

// String возвращает строковое представление числа с ведущими нулями
// до 3 цифр (например: "001", "042", "999").
func (n NumberPart) String() string {
	return fmt.Sprintf("%03d", n.number)
}

// Add увеличивает число на указанное значение и возвращает новый NumberPart.
//
// При переполнении происходит циклический переход (как одометр).
// Отрицательные значения игнорируются (считаются как 0).
func (n NumberPart) Add(value int) NumberPart {
	if value < 0 {
		value = 0
	}

	newNumber := (n.number + value) % (MAX_NUMBER_PART_VALUE + 1)
	return NumberPart{
		number: newNumber,
	}
}

// Number возвращает числовое значение.
func (n NumberPart) Number() int {
	return n.number
}
