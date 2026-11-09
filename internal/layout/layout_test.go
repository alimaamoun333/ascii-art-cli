package layout

import (
	"reflect"
	"testing"
)

func TestCalculateCharacterBoundariesSingleCharacter(t *testing.T) {
	shapes := map[rune][]string{
		'A': {
			"AA",
			"AAA",
			"A",
			"AAAA",
			"AA",
			"AAA",
			"A",
			"AAAA",
		},
	}

	got, err := CalculateCharacterBoundaries(shapes, "A")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := [][]Boundary{
		{{Start: 0, End: 2}},
		{{Start: 0, End: 3}},
		{{Start: 0, End: 1}},
		{{Start: 0, End: 4}},
		{{Start: 0, End: 2}},
		{{Start: 0, End: 3}},
		{{Start: 0, End: 1}},
		{{Start: 0, End: 4}},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}
func TestCalculateCharacterBoundariesMultipleCharacters(t *testing.T) {
	shapes := map[rune][]string{
		'A': {
			"AA",
			"AAA",
			"A",
			"AAAA",
			"AA",
			"AAA",
			"A",
			"AAAA",
		},
		'B': {
			"BBB",
			"B",
			"BB",
			"BBB",
			"B",
			"BB",
			"BBB",
			"B",
		},
	}

	got, err := CalculateCharacterBoundaries(shapes, "AB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := [][]Boundary{
		{{Start: 0, End: 2}, {Start: 2, End: 5}},
		{{Start: 0, End: 3}, {Start: 3, End: 4}},
		{{Start: 0, End: 1}, {Start: 1, End: 3}},
		{{Start: 0, End: 4}, {Start: 4, End: 7}},
		{{Start: 0, End: 2}, {Start: 2, End: 3}},
		{{Start: 0, End: 3}, {Start: 3, End: 5}},
		{{Start: 0, End: 1}, {Start: 1, End: 4}},
		{{Start: 0, End: 4}, {Start: 4, End: 5}},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestCalculateCharacterBoundariesUnsupportedCharacter(t *testing.T) {
	shapes := map[rune][]string{
		'A': {
			"AA",
			"AAA",
			"A",
			"AAAA",
			"AA",
			"AAA",
			"A",
			"AAAA",
		},
	}

	got, err := CalculateCharacterBoundaries(shapes, "AZ")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if got != nil {
		t.Fatalf("expected nil boundaries, got %v", got)
	}
}

func TestCharacterRangeToRenderedRange(t *testing.T) {
	rowBoundaries := []Boundary{
		{Start: 0, End: 4},
		{Start: 4, End: 7},
		{Start: 7, End: 12},
	}

	got, err := CharacterRangeToRenderedRange(rowBoundaries, 1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := Boundary{
		Start: 4,
		End:   12,
	}

	if got != expected {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestCalculateMultilineBoundaries(t *testing.T) {
	shapes := map[rune][]string{
		'A': {
			"AA",
			"AAA",
			"A",
			"AAAA",
			"AA",
			"AAA",
			"A",
			"AAAA",
		},
		'B': {
			"BBB",
			"B",
			"BB",
			"BBB",
			"B",
			"BB",
			"BBB",
			"B",
		},
	}

	got, err := CalculateMultilineBoundaries(shapes, "A\nB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := [][][]Boundary{
		{
			{{Start: 0, End: 2}},
			{{Start: 0, End: 3}},
			{{Start: 0, End: 1}},
			{{Start: 0, End: 4}},
			{{Start: 0, End: 2}},
			{{Start: 0, End: 3}},
			{{Start: 0, End: 1}},
			{{Start: 0, End: 4}},
		},
		{
			{{Start: 0, End: 3}},
			{{Start: 0, End: 1}},
			{{Start: 0, End: 2}},
			{{Start: 0, End: 3}},
			{{Start: 0, End: 1}},
			{{Start: 0, End: 2}},
			{{Start: 0, End: 3}},
			{{Start: 0, End: 1}},
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestCalculateMultilineBoundariesUnsupportedCharacter(t *testing.T) {
	shapes := map[rune][]string{
		'A': {
			"AA",
			"AAA",
			"A",
			"AAAA",
			"AA",
			"AAA",
			"A",
			"AAAA",
		},
	}

	got, err := CalculateMultilineBoundaries(shapes, "A\nZ")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if got != nil {
		t.Fatalf("expected nil boundaries, got %v", got)
	}
}

func TestParseLinePreservesSpacing(t *testing.T) {
	input := "  how   are you "
	got := ParseLine(input)

	expected := LineLayout{
		Leading:  2,
		Words:    []string{"how", "are", "you"},
		Gaps:     []int{3, 1},
		Trailing: 1,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %+v, got %+v", expected, got)
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected LineLayout
	}{
		{
			name:     "empty string",
			input:    "",
			expected: LineLayout{Leading: 0, Words: nil, Gaps: nil, Trailing: 0},
		},
		{
			name:     "single space",
			input:    " ",
			expected: LineLayout{Leading: 1, Words: nil, Gaps: nil, Trailing: 0},
		},
		{
			name:     "multiple spaces only",
			input:    "     ",
			expected: LineLayout{Leading: 5, Words: nil, Gaps: nil, Trailing: 0},
		},
		{
			name:     "single word without spaces",
			input:    "hello",
			expected: LineLayout{Leading: 0, Words: []string{"hello"}, Gaps: nil, Trailing: 0},
		},
		{
			name:     "single word with leading spaces",
			input:    "   hello",
			expected: LineLayout{Leading: 3, Words: []string{"hello"}, Gaps: nil, Trailing: 0},
		},
		{
			name:     "single word with trailing spaces",
			input:    "hello   ",
			expected: LineLayout{Leading: 0, Words: []string{"hello"}, Gaps: nil, Trailing: 3},
		},
		{
			name:     "single word with leading and trailing spaces",
			input:    "  hello    ",
			expected: LineLayout{Leading: 2, Words: []string{"hello"}, Gaps: nil, Trailing: 4},
		},
		{
			name:     "two words with single space gap",
			input:    "hello world",
			expected: LineLayout{Leading: 0, Words: []string{"hello", "world"}, Gaps: []int{1}, Trailing: 0},
		},
		{
			name:     "two words with multiple space gap",
			input:    "hello    world",
			expected: LineLayout{Leading: 0, Words: []string{"hello", "world"}, Gaps: []int{4}, Trailing: 0},
		},
		{
			name:     "multiple words with variable gaps",
			input:    "a  b   c    d",
			expected: LineLayout{Leading: 0, Words: []string{"a", "b", "c", "d"}, Gaps: []int{2, 3, 4}, Trailing: 0},
		},
		{
			name:     "complex spacing with leading, gaps, and trailing",
			input:    "   one   two  three    four ",
			expected: LineLayout{Leading: 3, Words: []string{"one", "two", "three", "four"}, Gaps: []int{3, 2, 4}, Trailing: 1},
		},
		{
			name:     "symbols, numbers and punctuation",
			input:    "123 + 456 = 579!",
			expected: LineLayout{Leading: 0, Words: []string{"123", "+", "456", "=", "579!"}, Gaps: []int{1, 1, 1, 1}, Trailing: 0},
		},
		{
			name:     "single character words",
			input:    " x y z ",
			expected: LineLayout{Leading: 1, Words: []string{"x", "y", "z"}, Gaps: []int{1, 1}, Trailing: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLine(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseLine(%q) = %+v, want %+v", tt.input, got, tt.expected)
			}
		})
	}
}



