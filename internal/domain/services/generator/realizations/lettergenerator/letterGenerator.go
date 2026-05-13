// Package lettergenerator предоставляет генератор комбинаций букв
// для формирования частей номеров.
//
// Генератор работает как счётчик в системе счисления, основанием которой
// является переданный vocabulary (уникальные символы). Поддерживает
// последовательную генерацию (GetNext) и случайную (GetRandom).
package lettergenerator

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"sync"
	"unicode"

	models "encomers/license/internal/domain/valueObjects"
)

// LetterGenerator — генератор фиксированной длины буквенных комбинаций.
//
// Потокобезопасен.
type LetterGenerator struct {
	vocabulary []rune
	last       []int
	overflow   bool
	mu         sync.Mutex
}

// New создаёт новый генератор.
//
// Параметры:
//   - vocabulary — алфавит, из которого будут генерироваться буквы.
//     Дубликаты автоматически удаляются, регистр приводится к верхнему.
//   - lastLetters — (опционально) начальное состояние генератора.
//     Должно содержать ровно models.LETTERS_COUNT символов.
//
// Возвращает ошибку, если vocabulary пустой или lastLetters имеет неверную длину/содержит символы вне алфавита.
func New(vocabulary []rune, lastLetters ...rune) (*LetterGenerator, error) {
	if len(lastLetters) != 0 &&
		len(lastLetters) != models.LETTERS_COUNT {
		return nil, fmt.Errorf(
			"invalid last letters count. Need %d and got %d",
			models.LETTERS_COUNT,
			len(lastLetters),
		)
	}
	if len(vocabulary) == 0 {
		return nil, fmt.Errorf("vocabulary cannot be empty")
	}
	seen := make(map[rune]struct{})
	for _, let := range vocabulary {
		seen[unicode.ToUpper(let)] = struct{}{}
	}
	uniqVocab := make([]rune, 0, len(seen))
	for k := range seen {
		uniqVocab = append(uniqVocab, k)
	}
	slices.Sort(uniqVocab)

	last := make([]int, models.LETTERS_COUNT)

	if len(lastLetters) == 0 {
		last[0] = -1
	} else {
		for i, let := range lastLetters {
			pos := slices.Index(uniqVocab, unicode.ToUpper(let))
			if pos == -1 {
				return nil, fmt.Errorf(
					"last letter %c is not in vocabulary",
					let,
				)
			}

			last[models.LETTERS_COUNT-1-i] = pos
		}
	}

	return &LetterGenerator{
		vocabulary: uniqVocab,
		last:       last,
		overflow:   false,
	}, nil
}

// IsOverflow возвращает true, если генератор прошёл все возможные комбинации
// и начал новый цикл (переполнение).
func (l *LetterGenerator) IsOverflow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.overflow
}

// Size возвращает размер используемого алфавита (количество уникальных символов).
func (l *LetterGenerator) Size() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.vocabulary)
}

// inc увеличивает текущее состояние генератора (как одометр).
// При переполнении устанавливает overflow = true.
func (l *LetterGenerator) inc() {
	for i := range l.last {
		l.last[i]++

		if l.last[i] < len(l.vocabulary) {
			return
		}

		l.last[i] = 0
	}
	for _, v := range l.last {
		if v != 0 {
			return
		}
		l.overflow = true
	}
}

// getCurrent возвращает текущую комбинацию как LetterPart.
func (l *LetterGenerator) getCurrent() models.LetterPart {
	letters := make([]rune, len(l.last))

	for i, idx := range l.last {
		letters[models.LETTERS_COUNT-1-i] = l.vocabulary[idx]
	}

	part, _ := models.NewLetterPart(letters)
	return part
}

// GetNext возвращает следующую комбинацию и сдвигает внутреннее состояние.
//
// Потокобезопасен.
func (l *LetterGenerator) GetNext() models.LetterPart {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.inc()
	current := l.getCurrent()
	return current
}

// GetRandom возвращает случайную комбинацию из алфавита и обновляет
// внутреннее состояние генератора.
//
// Потокобезопасен.
func (l *LetterGenerator) GetRandom() models.LetterPart {
	l.mu.Lock()
	defer l.mu.Unlock()
	for i := range l.last {
		l.last[i] = rand.IntN(len(l.vocabulary))
	}
	return l.getCurrent()
}
