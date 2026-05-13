package regiongenerator

import (
	"testing"

	models "encomers/license/internal/domain/valueObjects"
)

func TestRegionGenerator_Get(t *testing.T) {
	staticValue, _ := models.NewRegionPart("116", "rus")
	generator := New(staticValue)
	for range 10 {
		value := generator.Get()
		if !staticValue.Equal(value) {
			t.Fatalf("values are different. Original: %s\t Got: %s", staticValue.String(), value.String())
		}
	}
}
