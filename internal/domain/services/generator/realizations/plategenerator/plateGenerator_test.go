package plategenerator

import (
	"fmt"
	"testing"

	"encomers/license/internal/domain/services/generator/realizations/lettergenerator"
	"encomers/license/internal/domain/services/generator/realizations/numbergenerator"
	"encomers/license/internal/domain/services/generator/realizations/regiongenerator"
	repository "encomers/license/internal/domain/services/repository/realizations"
	models "encomers/license/internal/domain/valueObjects"
)

var testVocab = []rune("АЕТ")

func TestPlateGenerator_RealImplementations(t *testing.T) {
	repo := repository.New()

	letterGen, err := lettergenerator.New(testVocab)
	if err != nil {
		t.Fatal(err)
	}

	numberGen := numbergenerator.New(nil)
	region, _ := models.NewRegionPart("116", "RUS")
	regionGen := regiongenerator.New(region)

	plateGen, err := New(numberGen, letterGen, regionGen, repo)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetRandom returns valid plate", func(t *testing.T) {
		plate, err := plateGen.GetRandom()
		if err != nil {
			t.Fatal(err)
		}
		if plate.String() == "" {
			t.Error("plate.String() returned empty")
		}
		fmt.Println("Random plate:", plate.String())
	})

	t.Run("GetNext returns sequential plates", func(t *testing.T) {
		seen := make(map[string]bool)
		for i := 0; i < 5; i++ {
			plate, err := plateGen.GetNext()
			if err != nil {
				t.Fatal(err)
			}
			s := plate.String()
			if seen[s] {
				t.Errorf("duplicate plate: %s", s)
			}
			seen[s] = true
			fmt.Println("Next plate:", s)
		}
	})

	t.Run("IsOverflow works", func(t *testing.T) {
		if plateGen.IsOverflow() {
			t.Error("should not be overflow after few generations")
		}
	})
}
