package style

import (
	"ascii-art-justify/internal/layout"
	"reflect"
	"testing"
)

func TestFindRanges(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		target   string
		expected []Range
	}{
		{
			name:   "exact match",
			text:   "kit",
			target: "kit",
			expected: []Range{
				{Start: 0, End: 3},
			},
		},
		{
			name:     "no exact match",
			text:     "hello",
			target:   "kit",
			expected: []Range{},
		},
		{
			name:     "empty target",
			text:     "hello",
			target:   "",
			expected: []Range{},
		},
		{
			name:   "multiple occurrences",
			text:   "banana",
			target: "an",
			expected: []Range{
				{Start: 1, End: 3},
				{Start: 3, End: 5},
			},
		},
		{
			name:   "overlapping occurrences",
			text:   "banana",
			target: "ana",
			expected: []Range{
				{Start: 1, End: 4},
				{Start: 3, End: 6},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindRanges(tt.text, tt.target)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestANSIColor(t *testing.T) {
	tests := []struct {
		name      string
		color     string
		expected  string
		expectErr bool
	}{
		{
			name:      "red",
			color:     "red",
			expected:  "\033[31m",
			expectErr: false,
		},
		{
			name:      "invalid color",
			color:     "orange123",
			expected:  "",
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ANSIColor(tt.color)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if got != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestMergeRanges(t *testing.T) {
	tests := []struct {
		name     string
		input    []Range
		expected []Range
	}{
		{
			name: "overlapping ranges",
			input: []Range{
				{Start: 1, End: 4},
				{Start: 3, End: 6},
			},
			expected: []Range{
				{Start: 1, End: 6},
			},
		},
		{
			name: "touching ranges",
			input: []Range{
				{Start: 1, End: 3},
				{Start: 3, End: 5},
			},
			expected: []Range{
				{Start: 1, End: 5},
			},
		},
		{
			name: "contained range",
			input: []Range{
				{Start: 1, End: 6},
				{Start: 2, End: 4},
			},
			expected: []Range{
				{Start: 1, End: 6},
			},
		},
		{
			name: "separate ranges",
			input: []Range{
				{Start: 1, End: 3},
				{Start: 5, End: 7},
			},
			expected: []Range{
				{Start: 1, End: 3},
				{Start: 5, End: 7},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeRanges(tt.input)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestApplyColorToRange(t *testing.T) {
	row := "kkkkiiittttt"

	got, err := ApplyColorToRange(row, 4, 12, "\033[31m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "kkkk\033[31miiittttt\033[0m"

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestApplyColorToRanges(t *testing.T) {
	row := "AAAABBBBCCCCDDDD"

	ranges := []layout.Boundary{
		{Start: 4, End: 8},
		{Start: 12, End: 16},
	}

	got, err := ApplyColorToRanges(row, ranges, "\033[31m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "AAAA\033[31mBBBB\033[0mCCCC\033[31mDDDD\033[0m"

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestApplyColor(t *testing.T) {
	rows := []string{
		"AAAABBBBCCCC",
	}

	boundaries := [][]layout.Boundary{
		{
			{Start: 0, End: 4},
			{Start: 4, End: 8},
			{Start: 8, End: 12},
		},
	}

	got, err := ApplyColor(
		rows,
		"ABC",
		"B",
		"red",
		boundaries,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"AAAA\033[31mBBBB\033[0mCCCC",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}

}

func TestApplyColorNotFound(t *testing.T) {
	rows := []string{
		"AAAABBBBCCCC",
	}

	boundaries := [][]layout.Boundary{
		{
			{Start: 0, End: 4},
			{Start: 4, End: 8},
			{Start: 8, End: 12},
		},
	}

	got, err := ApplyColor(
		rows,
		"ABC",
		"Z",
		"red",
		boundaries,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"AAAABBBBCCCC",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestApplyColorMultiline(t *testing.T) {
	rows := make([]string, 16)
	for i := 0; i < 8; i++ {
		rows[i] = "AAAABBBBCCCC"
	}
	for i := 8; i < 16; i++ {
		rows[i] = "DDDDEEEEFFFF"
	}

	line1Boundaries := make([][]layout.Boundary, 8)
	for i := 0; i < 8; i++ {
		line1Boundaries[i] = []layout.Boundary{
			{Start: 0, End: 4},
			{Start: 4, End: 8},
			{Start: 8, End: 12},
		}
	}

	line2Boundaries := make([][]layout.Boundary, 8)
	for i := 0; i < 8; i++ {
		line2Boundaries[i] = []layout.Boundary{
			{Start: 0, End: 4},
			{Start: 4, End: 8},
			{Start: 8, End: 12},
		}
	}

	boundaries := [][][]layout.Boundary{line1Boundaries, line2Boundaries}

	got, err := ApplyColorMultiline(
		rows,
		"ABC\nDEF",
		"B",
		"red",
		boundaries,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := make([]string, 16)
	for i := 0; i < 8; i++ {
		expected[i] = "AAAA\033[31mBBBB\033[0mCCCC"
	}
	for i := 8; i < 16; i++ {
		expected[i] = "DDDDEEEEFFFF"
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestApplyColorMultilineBoundariesMismatch(t *testing.T) {
	rows := make([]string, 16)
	for i := 0; i < 16; i++ {
		rows[i] = "AAAA"
	}

	line1Boundaries := make([][]layout.Boundary, 8)
	for i := 0; i < 8; i++ {
		line1Boundaries[i] = []layout.Boundary{
			{Start: 0, End: 4},
		}
	}

	// only 1 line group for a 2-line text ("A\nB")
	boundaries := [][][]layout.Boundary{line1Boundaries}

	got, err := ApplyColorMultiline(
		rows,
		"A\nB",
		"A",
		"red",
		boundaries,
	)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if got != nil {
		t.Fatalf("expected nil rows, got %v", got)
	}
}

func TestApplyColorMultilineRowsMismatch(t *testing.T) {
	rows := make([]string, 15)

	line1Boundaries := make([][]layout.Boundary, 8)
	line2Boundaries := make([][]layout.Boundary, 8)

	boundaries := [][][]layout.Boundary{
		line1Boundaries,
		line2Boundaries,
	}

	got, err := ApplyColorMultiline(
		rows,
		"A\nB",
		"A",
		"red",
		boundaries,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if got != nil {
		t.Fatalf("expected nil rows, got %v", got)
	}
}

func TestApplyColorMultilineEmptyLine(t *testing.T) {
	// 3 lines: "A", "" (empty line), "B", each line has 8 rows (total 24 rows)
	rows := make([]string, 24)
	for i := 0; i < 8; i++ {
		rows[i] = "AAAA"
	}
	for i := 8; i < 16; i++ {
		rows[i] = ""
	}
	for i := 16; i < 24; i++ {
		rows[i] = "BBBB"
	}

	line1Boundaries := make([][]layout.Boundary, 8)
	for i := 0; i < 8; i++ {
		line1Boundaries[i] = []layout.Boundary{
			{Start: 0, End: 4},
		}
	}

	line2Boundaries := make([][]layout.Boundary, 8)
	for i := 0; i < 8; i++ {
		line2Boundaries[i] = []layout.Boundary{}
	}

	line3Boundaries := make([][]layout.Boundary, 8)
	for i := 0; i < 8; i++ {
		line3Boundaries[i] = []layout.Boundary{
			{Start: 0, End: 4},
		}
	}

	boundaries := [][][]layout.Boundary{line1Boundaries, line2Boundaries, line3Boundaries}

	got, err := ApplyColorMultiline(
		rows,
		"A\n\nB",
		"B",
		"red",
		boundaries,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := make([]string, 24)
	for i := 0; i < 8; i++ {
		expected[i] = "AAAA"
	}
	for i := 8; i < 16; i++ {
		expected[i] = ""
	}
	for i := 16; i < 24; i++ {
		expected[i] = "\033[31mBBBB\033[0m"
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestApplyColorMultilineWholeTextWithEmptyLine(t *testing.T) {
	// 3 lines: "A", "" (empty line), "B", each line has 8 rows (total 24 rows)
	rows := make([]string, 24)
	for i := 0; i < 8; i++ {
		rows[i] = "AAAA"
	}
	for i := 8; i < 16; i++ {
		rows[i] = ""
	}
	for i := 16; i < 24; i++ {
		rows[i] = "BBBB"
	}

	line1Boundaries := make([][]layout.Boundary, 8)
	for i := 0; i < 8; i++ {
		line1Boundaries[i] = []layout.Boundary{
			{Start: 0, End: 4},
		}
	}

	line2Boundaries := make([][]layout.Boundary, 8)
	for i := 0; i < 8; i++ {
		line2Boundaries[i] = []layout.Boundary{}
	}

	line3Boundaries := make([][]layout.Boundary, 8)
	for i := 0; i < 8; i++ {
		line3Boundaries[i] = []layout.Boundary{
			{Start: 0, End: 4},
		}
	}

	boundaries := [][][]layout.Boundary{line1Boundaries, line2Boundaries, line3Boundaries}

	// target == "" applies color to the entire line text
	got, err := ApplyColorMultiline(
		rows,
		"A\n\nB",
		"",
		"red",
		boundaries,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := make([]string, 24)
	for i := 0; i < 8; i++ {
		expected[i] = "\033[31mAAAA\033[0m"
	}
	for i := 8; i < 16; i++ {
		expected[i] = ""
	}
	for i := 16; i < 24; i++ {
		expected[i] = "\033[31mBBBB\033[0m"
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}
