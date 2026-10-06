package domain

import "testing"

func TestParseArrowScore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		code  string
		value int
	}{
		{code: "X", value: 10},
		{code: "10", value: 10},
		{code: "9", value: 9},
		{code: "8", value: 8},
		{code: "7", value: 7},
		{code: "6", value: 6},
		{code: "5", value: 5},
		{code: "4", value: 4},
		{code: "3", value: 3},
		{code: "2", value: 2},
		{code: "1", value: 1},
		{code: "M", value: 0},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			t.Parallel()
			got, err := ParseArrowScore(tt.code)
			if err != nil {
				t.Fatalf("ParseArrowScore(%q): %v", tt.code, err)
			}
			if got.Code != tt.code || got.NumericValue != tt.value {
				t.Fatalf("got %+v, want code %q value %d", got, tt.code, tt.value)
			}
		})
	}
}

func TestParseArrowScoreRejectsInvalid(t *testing.T) {
	t.Parallel()

	for _, code := range []string{"", "11", "0", "x", "miss", "10 "} {
		_, err := ParseArrowScore(code)
		if err != ErrInvalidScoreCode {
			t.Errorf("ParseArrowScore(%q) error = %v, want ErrInvalidScoreCode", code, err)
		}
	}
}

func TestArrowScoreIsGold(t *testing.T) {
	t.Parallel()

	golds := map[string]bool{"X": true, "10": true, "9": true, "8": false, "M": false, "1": false}
	for code, want := range golds {
		score, err := ParseArrowScore(code)
		if err != nil {
			t.Fatalf("ParseArrowScore(%q): %v", code, err)
		}
		if got := score.IsGold(); got != want {
			t.Errorf("%s IsGold() = %v, want %v", code, got, want)
		}
	}
}

func TestArrowScoreIsHit(t *testing.T) {
	t.Parallel()

	miss, err := ParseArrowScore("M")
	if err != nil {
		t.Fatal(err)
	}
	if miss.IsHit() {
		t.Fatal("M should not count as a hit")
	}
	ten, err := ParseArrowScore("10")
	if err != nil {
		t.Fatal(err)
	}
	if !ten.IsHit() {
		t.Fatal("10 should count as a hit")
	}
}
