// Package valueObjects содержит основные value objects,
// представляющие составные части автомобильного номера.
package valueObjects

import (
	"errors"
	"strings"
)

// RegionPart представляет региональную часть автомобильного номера
// (код региона + страна).
type RegionPart struct {
	region  string
	country string
}

// NewRegionPart создаёт новую региональную часть номера.
//
// Параметры:
//   - region  — код региона (например "77", "МО", "АВА")
//   - country — код страны (например "RUS", "UA", "BY")
//
// Оба параметра приводятся к верхнему регистру и очищаются от пробелов.
// Возвращает ошибку, если любой из параметров пустой.
func NewRegionPart(region, country string) (RegionPart, error) {
	if len(strings.TrimSpace(region)) == 0 || len(strings.TrimSpace(country)) == 0 {
		return RegionPart{}, errors.New("invalid region part")
	}

	return RegionPart{
		region:  strings.TrimSpace(strings.ToUpper(region)),
		country: strings.TrimSpace(strings.ToUpper(country)),
	}, nil
}

// String возвращает строковое представление региона в формате "REGION COUNTRY"
// (например: "77 RUS").
func (r RegionPart) String() string {
	return strings.ToUpper(r.region + " " + r.country)
}

// Equal сравнивает два региона на равенство.
// Сравнение производится по country и region.
func (r RegionPart) Equal(other RegionPart) bool {
	return r.country == other.country && r.region == other.region
}

// Region возвращает код региона в верхнем регистре.
func (r RegionPart) Region() string {
	return r.region
}

// Country возвращает код страны в верхнем регистре.
func (r RegionPart) Country() string {
	return r.country
}
