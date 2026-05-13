package numbergenerator

import (
	"fmt"
	"testing"

	models "encomers/license/internal/domain/valueObjects"
)

const (
	iterations = 100_000
	maxMinSeen = 500
)

func TestNilNumberGenerator_GetNext(t *testing.T) {
	generator := New(nil)
	for i := range iterations {
		next := generator.GetNext()
		generatedNum := next.Number()
		generatedStr := next.String()

		target := (i % models.MAX_NUMBER_PART_VALUE) + 1
		if target == 0 {
			target = 1
		}
		targetStr := fmt.Sprintf("%03d", target)

		if target != generatedNum {
			t.Fatalf("iteration: %d\tGenerated number (%d) not equals target (%d)", i, generatedNum, target)
		} else if targetStr != generatedStr {
			t.Fatalf("Generated string number (%s) not equals target (%s)", generatedStr, targetStr)
		}
	}
}

func TestLastNumberGenerator_GetNext(t *testing.T) {
	val := 101
	lastValue, _ := models.NewNumberPart(val)
	generator := New(&lastValue)
	for i := range iterations {
		next := generator.GetNext()
		generatedNum := next.Number()
		generatedStr := next.String()

		target := (val+i)%models.MAX_NUMBER_PART_VALUE + 1
		targetStr := fmt.Sprintf("%03d", target)

		if target != generatedNum {
			t.Fatalf("Generated number (%d) not equals target (%d)", generatedNum, target)
		} else if targetStr != generatedStr {
			t.Fatalf("Generated string number (%s) not equals target (%s)", generatedStr, targetStr)
		}
	}
}

func TestNumberGenerator_GetRandom(t *testing.T) {
	gen := New(nil)

	seen := make(map[int]struct{})

	for i := range iterations {
		result := gen.GetRandom()

		val := result.Number()

		if val < 0 || val > models.MAX_NUMBER_PART_VALUE {
			t.Errorf("iteration %d: GetRandom() returned %d", i, val)
		}

		seen[val] = struct{}{}
	}

	minSeen := iterations / 10
	if minSeen > maxMinSeen {
		minSeen = maxMinSeen
	}
	if len(seen) < minSeen {
		t.Errorf("GetRandom() returned less than %d unique numbers: %d", minSeen, len(seen))
	}
}

func TestNumberGenerator_Overflow(t *testing.T) {
	gen := New(nil)

	for range 1000 {
		gen.GetNext()
	}

	if !gen.IsOverflow() {
		t.Error("Expected overflow after 1000 increments, but it did not occur")
	}
}
