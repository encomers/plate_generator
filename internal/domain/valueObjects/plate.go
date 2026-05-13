// Package valueObjects содержит основные value objects,
// представляющие составные части автомобильного номера.
package valueObjects

import "fmt"

// Plate представляет полный автомобильный номер (лицензионную пластину).
//
// Состоит из трёх основных частей:
//   - Letters — буквенная часть (3 буквы)
//   - Numbers — числовая часть (3 цифры)
//   - Region  — региональная часть (регион + страна)
type Plate struct {
	Numbers NumberPart
	Letters LetterPart
	Region  RegionPart
}

// String возвращает строковое представление номера в читаемом формате.
//
// Пример вывода: `A 042 BC 77 RUS`
func (p Plate) String() string {
	letters := p.Letters.Letters()

	return fmt.Sprintf("%c %03d %c%c %s",
		letters[0],
		p.Numbers.Number(),
		letters[1],
		letters[2],
		p.Region.String(),
	)
}

// Format — удобный метод для форматирования номера в разных вариантах.
// Поддерживает стандартный вывод через %s, %v и %+v.
func (p Plate) Format(f fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		f.Write([]byte(p.String()))
	default:
		fmt.Fprintf(f, "Plate%+v", struct {
			Letters string
			Numbers string
			Region  string
		}{
			Letters: p.Letters.String(),
			Numbers: p.Numbers.String(),
			Region:  p.Region.String(),
		})
	}
}
