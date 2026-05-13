// Package numbergenerator предоставляет генератор числовых частей номеров.
//
// Генератор работает как счётчик в заданном диапазоне и поддерживает
// последовательную и случайную генерацию чисел.
package numbergenerator

import (
	"math/rand/v2"
	"sync"

	models "encomers/license/internal/domain/valueObjects"
)

// NumberGenerator — генератор числовой части номера.
//
// Поддерживает инкрементальное увеличение и генерацию случайных значений.
// Потокобезопасен.
type NumberGenerator struct {
	lastNumber models.NumberPart
	overflow   bool
	mu         sync.Mutex
}

// New создаёт новый генератор числовой части.
//
// Если lastNumber == nil, начинается с 0.
func New(lastNumber *models.NumberPart) *NumberGenerator {
	var number models.NumberPart
	if lastNumber != nil {
		number = *lastNumber
	} else {
		number, _ = models.NewNumberPart(0)
	}
	result := &NumberGenerator{
		lastNumber: number,
		overflow:   false,
	}

	return result
}

// IsOverflow возвращает true, если генератор достиг максимального значения
// и произошло переполнение (переход через максимум).
func (n *NumberGenerator) IsOverflow() bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.overflow
}

// GetNext возвращает следующее число и обновляет внутреннее состояние.
//
// При достижении максимального значения (MAX_NUMBER_PART_VALUE):
// - устанавливается флаг overflow = true;
// - следующее значение становится 1.
func (n *NumberGenerator) GetNext() models.NumberPart {
	n.mu.Lock()
	defer n.mu.Unlock()
	next := n.lastNumber.Add(1)

	if next.Number() == 0 {
		n.overflow = true
		next, _ = models.NewNumberPart(1)
	}

	n.lastNumber = next
	return n.lastNumber
}

// GetRandom генерирует случайное число в диапазоне [0, MAX_NUMBER_PART_VALUE]
// и обновляет внутреннее состояние генератора.
func (n *NumberGenerator) GetRandom() models.NumberPart {
	n.mu.Lock()
	defer n.mu.Unlock()
	number := rand.IntN(models.MAX_NUMBER_PART_VALUE + 1)
	n.lastNumber, _ = models.NewNumberPart(number)
	return n.lastNumber
}
