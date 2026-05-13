// Package valueObjects содержит основные value objects (объекты-значения),
// представляющие составные части автомобильного номера.
package valueObjects

import (
	"errors"
	"unicode"
)

const (
	// LETTERS_COUNT — количество букв в буквенной части номера.
	LETTERS_COUNT = 3
)

// LetterPart представляет буквенную часть автомобильного номера (обычно 3 буквы).
//
// Все буквы внутри хранятся в верхнем регистре.
type LetterPart struct {
	letters []rune
}

// NewLetterPart создаёт новую буквенную часть номера.
//
// Параметры:
//   - letters — слайс рун длиной ровно LETTERS_COUNT.
//
// Возвращает ошибку, если количество букв неверное.
func NewLetterPart(letters []rune) (LetterPart, error) {
	if len(letters) != LETTERS_COUNT {
		return LetterPart{}, errors.New("invalid letter part")
	}

	lettersCmp := make([]rune, LETTERS_COUNT)
	for i, l := range letters {
		lettersCmp[i] = unicode.ToUpper(l)
	}

	return LetterPart{
		letters: lettersCmp,
	}, nil
}

// NewLetterPartFromString создаёт LetterPart из строки.
//
// Удобный конструктор для создания из строковых литералов.
func NewLetterPartFromString(letters string) (LetterPart, error) {
	return NewLetterPart([]rune(letters))
}

// String возвращает строковое представление буквенной части номера.
func (l LetterPart) String() string {
	return string(l.letters)
}

// Letters возвращает копию букв в виде слайса рун (все в верхнем регистре).
func (l LetterPart) Letters() []rune {
	// Возвращаем копию, чтобы сохранить иммутабельность value object
	result := make([]rune, len(l.letters))
	copy(result, l.letters)
	return result
}
