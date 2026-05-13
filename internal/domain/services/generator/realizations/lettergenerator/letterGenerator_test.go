package lettergenerator

import (
	"testing"

	models "encomers/license/internal/domain/valueObjects"
)

func TestNew_DeduplicatesAndUppercasesVocabulary(t *testing.T) {
	g, err := New([]rune{'b', 'A', 'a', 'c', 'C'})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	wantVocab := []rune{'A', 'B', 'C'}
	if len(g.vocabulary) != len(wantVocab) {
		t.Fatalf("vocabulary len = %d, want %d", len(g.vocabulary), len(wantVocab))
	}
	for i := range wantVocab {
		if g.vocabulary[i] != wantVocab[i] {
			t.Fatalf("vocabulary[%d] = %q, want %q", i, g.vocabulary[i], wantVocab[i])
		}
	}

	if len(g.last) != models.LETTERS_COUNT {
		t.Fatalf("len(last) = %d, want %d", len(g.last), models.LETTERS_COUNT)
	}

	if g.last[0] != -1 {
		t.Fatalf("last[0] = %d, want -1", g.last[0])
	}
	for i := 1; i < len(g.last); i++ {
		if g.last[i] != 0 {
			t.Fatalf("last[%d] = %d, want 0", i, g.last[i])
		}
	}
}

func TestNew_RejectsEmptyVocabulary(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty vocabulary, got nil")
	}
}

func TestNew_RejectsInvalidLastLettersCount(t *testing.T) {
	badLen := models.LETTERS_COUNT + 1

	last := make([]rune, badLen)
	for i := range last {
		last[i] = 'A'
	}

	_, err := New([]rune{'A', 'B', 'C'}, last...)
	if err == nil {
		t.Fatal("expected error for invalid last letters count, got nil")
	}
}

func TestGetNext_FromFreshGenerator_ReturnsAAA(t *testing.T) {
	g, err := New([]rune{'A', 'B'})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	got := g.GetNext()

	if got.String() != "AAA" {
		t.Fatalf("GetNext() = %q, want %q", got.String(), "AAA")
	}
}

func TestGetNext_CarriesCorrectly(t *testing.T) {
	g, err := New([]rune{'A', 'B'})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	first := g.GetNext()
	if first.String() != "AAA" {
		t.Fatalf("first GetNext() = %q, want %q", first.String(), "AAA")
	}

	second := g.GetNext()
	if second.String() != "AAB" {
		t.Fatalf("second GetNext() = %q, want %q", second.String(), "AAB")
	}

	third := g.GetNext()
	if third.String() != "ABA" {
		t.Fatalf("third GetNext() = %q, want %q", third.String(), "ABA")
	}

	fourth := g.GetNext()
	if fourth.String() != "ABB" {
		t.Fatalf("fourth GetNext() = %q, want %q", fourth.String(), "ABB")
	}
}

func TestGetNext_FromProvidedLastLetters(t *testing.T) {
	g, err := New([]rune{'A', 'B', 'C'}, 'A', 'A', 'A')
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	got := g.GetNext()

	if got.String() != "AAB" {
		t.Fatalf("GetNext() = %q, want %q", got.String(), "AAB")
	}
}

func TestGetRandom_StaysWithinVocabulary(t *testing.T) {
	g, err := New([]rune{'A', 'B', 'C'})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	allowed := map[rune]struct{}{
		'A': {},
		'B': {},
		'C': {},
	}

	for n := 0; n < 10_000; n++ {
		got := g.GetRandom()

		if len(got.Letters()) != models.LETTERS_COUNT {
			t.Fatalf("GetRandom() returned length %d, want %d", len(got.Letters()), models.LETTERS_COUNT)
		}

		for _, r := range got.Letters() {
			if _, ok := allowed[r]; !ok {
				t.Fatalf("GetRandom() returned rune %q not in vocabulary", r)
			}
		}

		if got.String() != g.getCurrent().String() {
			t.Fatalf("GetRandom() = %q, getCurrent() = %q", got.String(), g.getCurrent().String())
		}
	}
}

func TestLetterGenerator_Overflow(t *testing.T) {
	g, err := New([]rune{'A', 'B'})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	for range 10 {
		g.GetNext()
	}

	if !g.IsOverflow() {
		t.Fatal("expected overflow after exhausting all combinations, but IsOverflow() returned false")
	}
}
