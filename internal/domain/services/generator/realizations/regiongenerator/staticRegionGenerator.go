// Package regiongenerator предоставляет генераторы для региональной части номера
package regiongenerator

import (
	models "encomers/license/internal/domain/valueObjects"
)

// StaticRegionGenerator — статический генератор региональной части номера.
//
// Всегда возвращает одно и то же фиксированное значение региона.
// Реализует общий интерфейс генераторов, но не поддерживает изменение состояния.
type StaticRegionGenerator struct {
	staticValue models.RegionPart
}

// New создаёт новый статический генератор региона.
//
// Переданное значение региона будет возвращаться при каждом вызове Get().
func New(value models.RegionPart) *StaticRegionGenerator {
	return &StaticRegionGenerator{
		staticValue: value,
	}
}

// IsOverflow всегда возвращает false, так как генератор статический
// и не может переполниться.
func (s *StaticRegionGenerator) IsOverflow() bool {
	return false
}

// SetRegion — пустая реализация метода для совместимости с общим интерфейсом.
// В статическом генераторе значение региона нельзя изменить после создания.
func (s *StaticRegionGenerator) SetRegion(region models.RegionPart) {

}

// Get возвращает фиксированное значение региона.
func (s *StaticRegionGenerator) Get() models.RegionPart {
	return s.staticValue
}
