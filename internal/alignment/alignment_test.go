package alignment

import (
	"reflect"
	"strings"
	"testing"
)

func TestBlockWidth(t *testing.T) {
	tests := []struct {
		name     string
		rows     []string
		expected int
	}{
		{
			name:     "nil slice",
			rows:     nil,
			expected: 0,
		},
		{
			name:     "empty slice",
			rows:     []string{},
			expected: 0,
		},
		{
			name:     "single empty row",
			rows:     []string{""},
			expected: 0,
		},
		{
			name:     "multiple empty rows",
			rows:     []string{"", "", ""},
			expected: 0,
		},
		{
			name:     "single row",
			rows:     []string{"hello"},
			expected: 5,
		},
		{
			name:     "uniform row lengths",
			rows:     []string{"abc", "def", "ghi"},
			expected: 3,
		},
		{
			name:     "max width at first row",
			rows:     []string{"longest line", "short", "mid len"},
			expected: 12,
		},
		{
			name:     "max width in middle row",
			rows:     []string{"short", "the longest row here", "tiny"},
			expected: 20,
		},
		{
			name:     "max width at last row",
			rows:     []string{"a", "ab", "abc"},
			expected: 3,
		},
		{
			name: "ascii art style block",
			rows: []string{
				" _  _     _ ",
				"| || |___| |",
				"| __ / -_) |",
				"|_||_\\___|_|",
			},
			expected: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := blockWidth(tt.rows)
			if got != tt.expected {
				t.Errorf("blockWidth(%v) = %d, want %d", tt.rows, got, tt.expected)
			}
		})
	}
}

func TestNormalizeBlock(t *testing.T) {
	tests := []struct {
		name     string
		rows     []string
		expected []string
	}{
		{
			name:     "nil slice",
			rows:     nil,
			expected: []string{},
		},
		{
			name:     "empty slice",
			rows:     []string{},
			expected: []string{},
		},
		{
			name:     "single empty row",
			rows:     []string{""},
			expected: []string{""},
		},
		{
			name:     "multiple empty rows",
			rows:     []string{"", "", ""},
			expected: []string{"", "", ""},
		},
		{
			name:     "single row",
			rows:     []string{"hello"},
			expected: []string{"hello"},
		},
		{
			name:     "uniform row lengths",
			rows:     []string{"abc", "def", "123"},
			expected: []string{"abc", "def", "123"},
		},
		{
			name:     "first row longest",
			rows:     []string{"longest row", "short", "mid"},
			expected: []string{"longest row", "short      ", "mid        "},
		},
		{
			name:     "middle row longest",
			rows:     []string{"a", "longer", "b"},
			expected: []string{"a     ", "longer", "b     "},
		},
		{
			name:     "last row longest",
			rows:     []string{"mid", "short", "longest row"},
			expected: []string{"mid        ", "short      ", "longest row"},
		},
		{
			name:     "varying lengths with mixed padding",
			rows:     []string{"a", "abc", "ab"},
			expected: []string{"a  ", "abc", "ab "},
		},
		{
			name: "ascii art block with jagged lines",
			rows: []string{
				"/\\",
				"/__\\",
				"||",
			},
			expected: []string{
				"/\\  ",
				"/__\\",
				"||  ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeBlock(tt.rows)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("normalizeBlock(%v) = %#v, want %#v", tt.rows, got, tt.expected)
			}
		})
	}
}

func TestNewContent(t *testing.T) {
	validSpace := []string{" ", " ", " ", " ", " ", " ", " ", " "}
	validWord1 := []string{"#", "#", "#", "#", "#", "#", "#", "#"}
	validWord2 := []string{"@", "@", "@", "@", "@", "@", "@", "@"}

	tests := []struct {
		name      string
		space     []string
		words     [][]string
		leading   int
		gaps      []int
		trailing  int
		want      Content
		wantErr   bool
		errSubstr string
	}{
		{
			name:     "valid single word without gaps",
			space:    validSpace,
			words:    [][]string{validWord1},
			leading:  0,
			gaps:     []int{},
			trailing: 0,
			want: Content{
				space:    validSpace,
				words:    [][]string{validWord1},
				leading:  0,
				gaps:     []int{},
				trailing: 0,
			},
			wantErr: false,
		},
		{
			name:     "valid multiple words with gaps and padding",
			space:    validSpace,
			words:    [][]string{validWord1, validWord2},
			leading:  2,
			gaps:     []int{3},
			trailing: 1,
			want: Content{
				space:    validSpace,
				words:    [][]string{validWord1, validWord2},
				leading:  2,
				gaps:     []int{3},
				trailing: 1,
			},
			wantErr: false,
		},
		{
			name:     "valid no words with 0 gaps",
			space:    validSpace,
			words:    [][]string{},
			leading:  5,
			gaps:     []int{},
			trailing: 3,
			want: Content{
				space:    validSpace,
				words:    [][]string{},
				leading:  5,
				gaps:     []int{},
				trailing: 3,
			},
			wantErr: false,
		},
		{
			name:      "invalid space block row count (less than 8)",
			space:     []string{" ", " "},
			words:     [][]string{validWord1},
			leading:   0,
			gaps:      []int{},
			trailing:  0,
			wantErr:   true,
			errSubstr: "space block must contain 8 rows",
		},
		{
			name:      "invalid space block row count (more than 8)",
			space:     append(validSpace, " extra "),
			words:     [][]string{validWord1},
			leading:   0,
			gaps:      []int{},
			trailing:  0,
			wantErr:   true,
			errSubstr: "space block must contain 8 rows",
		},
		{
			name:      "invalid word block row count (less than 8)",
			space:     validSpace,
			words:     [][]string{{"only", "two"}},
			leading:   0,
			gaps:      []int{},
			trailing:  0,
			wantErr:   true,
			errSubstr: "word block must contain 8 rows",
		},
		{
			name:      "invalid word block row count (one of multiple words invalid)",
			space:     validSpace,
			words:     [][]string{validWord1, {"short"}},
			leading:   0,
			gaps:      []int{1},
			trailing:  0,
			wantErr:   true,
			errSubstr: "word block must contain 8 rows",
		},
		{
			name:      "mismatched gaps count (fewer gaps than words-1)",
			space:     validSpace,
			words:     [][]string{validWord1, validWord2},
			leading:   0,
			gaps:      []int{},
			trailing:  0,
			wantErr:   true,
			errSubstr: "expected 1 gaps for 2 words, got 0",
		},
		{
			name:      "mismatched gaps count (more gaps than words-1)",
			space:     validSpace,
			words:     [][]string{validWord1, validWord2},
			leading:   0,
			gaps:      []int{1, 2},
			trailing:  0,
			wantErr:   true,
			errSubstr: "expected 1 gaps for 2 words, got 2",
		},
		{
			name:      "mismatched gaps count when words empty",
			space:     validSpace,
			words:     [][]string{},
			leading:   0,
			gaps:      []int{1},
			trailing:  0,
			wantErr:   true,
			errSubstr: "expected 0 gaps for 0 words, got 1",
		},
		{
			name:      "gap value zero",
			space:     validSpace,
			words:     [][]string{validWord1, validWord2},
			leading:   0,
			gaps:      []int{0},
			trailing:  0,
			wantErr:   true,
			errSubstr: "gap must be greater than 0, got 0",
		},
		{
			name:      "gap value negative",
			space:     validSpace,
			words:     [][]string{validWord1, validWord2},
			leading:   0,
			gaps:      []int{-1},
			trailing:  0,
			wantErr:   true,
			errSubstr: "gap must be greater than 0, got -1",
		},
		{
			name:      "negative leading",
			space:     validSpace,
			words:     [][]string{validWord1},
			leading:   -1,
			gaps:      []int{},
			trailing:  0,
			wantErr:   true,
			errSubstr: "leading must be greater than or equal to 0, got -1",
		},
		{
			name:      "negative trailing",
			space:     validSpace,
			words:     [][]string{validWord1},
			leading:   0,
			gaps:      []int{},
			trailing:  -2,
			wantErr:   true,
			errSubstr: "trailing must be greater than or equal to 0, got -2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewContent(tt.space, tt.words, tt.leading, tt.gaps, tt.trailing)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewContent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if tt.errSubstr != "" && !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("NewContent() error = %q, expected substring %q", err.Error(), tt.errSubstr)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContent() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestBuildNormalBlock(t *testing.T) {
	spaceBlock := []string{" ", " ", " ", " ", " ", " ", " ", " "}
	space3Block := []string{"...", "...", "...", "...", "...", "...", "...", "..."}

	wordA := []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8"}
	wordB := []string{"B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8"}
	wordC := []string{"C1", "C2", "C3", "C4", "C5", "C6", "C7", "C8"}

	jaggedSpace := []string{" ", "  ", " ", "   ", " ", " ", " ", " "}
	jaggedWord1 := []string{"X", "XX", "XXX", "X", "X", "X", "X", "X"}

	tests := []struct {
		name     string
		content  Content
		expected []string
	}{
		{
			name: "single word without leading or trailing spaces",
			content: Content{
				space:    spaceBlock,
				words:    [][]string{wordA},
				leading:  0,
				gaps:     []int{},
				trailing: 0,
			},
			expected: []string{
				"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8",
			},
		},
		{
			name: "single word with leading and trailing spaces",
			content: Content{
				space:    spaceBlock,
				words:    [][]string{wordA},
				leading:  2,
				gaps:     []int{},
				trailing: 3,
			},
			expected: []string{
				"  A1   ",
				"  A2   ",
				"  A3   ",
				"  A4   ",
				"  A5   ",
				"  A6   ",
				"  A7   ",
				"  A8   ",
			},
		},
		{
			name: "multiple words with gaps and no padding",
			content: Content{
				space:    spaceBlock,
				words:    [][]string{wordA, wordB},
				leading:  0,
				gaps:     []int{2},
				trailing: 0,
			},
			expected: []string{
				"A1  B1",
				"A2  B2",
				"A3  B3",
				"A4  B4",
				"A5  B5",
				"A6  B6",
				"A7  B7",
				"A8  B8",
			},
		},
		{
			name: "multiple words with leading, multiple gaps, and trailing",
			content: Content{
				space:    space3Block,
				words:    [][]string{wordA, wordB, wordC},
				leading:  1,
				gaps:     []int{2, 1},
				trailing: 2,
			},
			expected: []string{
				"...A1......B1...C1......",
				"...A2......B2...C2......",
				"...A3......B3...C3......",
				"...A4......B4...C4......",
				"...A5......B5...C5......",
				"...A6......B6...C6......",
				"...A7......B7...C7......",
				"...A8......B8...C8......",
			},
		},
		{
			name: "normalizes jagged space and jagged word blocks",
			content: Content{
				space:    jaggedSpace,
				words:    [][]string{jaggedWord1},
				leading:  1,
				gaps:     []int{},
				trailing: 1,
			},
			expected: []string{
				"   X     ",
				"   XX    ",
				"   XXX   ",
				"   X     ",
				"   X     ",
				"   X     ",
				"   X     ",
				"   X     ",
			},
		},
		{
			name: "empty words slice with only leading and trailing spaces",
			content: Content{
				space:    spaceBlock,
				words:    [][]string{},
				leading:  3,
				gaps:     []int{},
				trailing: 2,
			},
			expected: []string{
				"     ",
				"     ",
				"     ",
				"     ",
				"     ",
				"     ",
				"     ",
				"     ",
			},
		},
		{
			name: "empty words slice with zero spaces",
			content: Content{
				space:    spaceBlock,
				words:    [][]string{},
				leading:  0,
				gaps:     []int{},
				trailing: 0,
			},
			expected: []string{
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildNormalBlock(tt.content)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("buildNormalBlock() =\n%#v\nwant:\n%#v", got, tt.expected)
			}
		})
	}
}

func TestAlignRight(t *testing.T) {
	tests := []struct {
		name     string
		rows     []string
		width    int
		expected []string
	}{
		{
			name:     "nil slice",
			rows:     nil,
			width:    10,
			expected: []string{},
		},
		{
			name:     "empty slice",
			rows:     []string{},
			width:    10,
			expected: []string{},
		},
		{
			name:     "single empty row",
			rows:     []string{""},
			width:    5,
			expected: []string{"     "},
		},
		{
			name:     "single row exact width",
			rows:     []string{"hello"},
			width:    5,
			expected: []string{"hello"},
		},
		{
			name:     "single row with left padding",
			rows:     []string{"hello"},
			width:    10,
			expected: []string{"     hello"},
		},
		{
			name:     "uniform rows with left padding",
			rows:     []string{"abc", "def", "ghi"},
			width:    7,
			expected: []string{"    abc", "    def", "    ghi"},
		},
		{
			name:     "jagged rows normalized and right-aligned",
			rows:     []string{"a", "abcde", "ab"},
			width:    9,
			expected: []string{"    a    ", "    abcde", "    ab   "},
		},
		{
			name: "ascii art block",
			rows: []string{
				" _  _ ",
				"| || |",
				"| __ |",
				"|_||_|",
			},
			width: 10,
			expected: []string{
				"     _  _ ",
				"    | || |",
				"    | __ |",
				"    |_||_|",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := alignRight(tt.rows, tt.width)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("alignRight(%v, %d) =\n%#v\nwant:\n%#v", tt.rows, tt.width, got, tt.expected)
			}
		})
	}
}

func TestAlignLeft(t *testing.T) {
	tests := []struct {
		name     string
		rows     []string
		width    int
		expected []string
	}{
		{
			name:     "nil slice",
			rows:     nil,
			width:    10,
			expected: []string{},
		},
		{
			name:     "empty slice",
			rows:     []string{},
			width:    10,
			expected: []string{},
		},
		{
			name:     "single empty row",
			rows:     []string{""},
			width:    5,
			expected: []string{"     "},
		},
		{
			name:     "single row exact width",
			rows:     []string{"hello"},
			width:    5,
			expected: []string{"hello"},
		},
		{
			name:     "single row with right padding",
			rows:     []string{"hello"},
			width:    10,
			expected: []string{"hello     "},
		},
		{
			name:     "uniform rows with right padding",
			rows:     []string{"abc", "def", "ghi"},
			width:    7,
			expected: []string{"abc    ", "def    ", "ghi    "},
		},
		{
			name:     "jagged rows normalized and left-aligned",
			rows:     []string{"a", "abcde", "ab"},
			width:    9,
			expected: []string{"a        ", "abcde    ", "ab       "},
		},
		{
			name: "ascii art block",
			rows: []string{
				" _  _ ",
				"| || |",
				"| __ |",
				"|_||_|",
			},
			width: 10,
			expected: []string{
				" _  _     ",
				"| || |    ",
				"| __ |    ",
				"|_||_|    ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := alignLeft(tt.rows, tt.width)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("alignLeft(%v, %d) =\n%#v\nwant:\n%#v", tt.rows, tt.width, got, tt.expected)
			}
		})
	}
}

func TestAlignCenter(t *testing.T) {
	tests := []struct {
		name     string
		rows     []string
		width    int
		expected []string
	}{
		{
			name:     "nil slice",
			rows:     nil,
			width:    10,
			expected: []string{},
		},
		{
			name:     "empty slice",
			rows:     []string{},
			width:    10,
			expected: []string{},
		},
		{
			name:     "single empty row even width",
			rows:     []string{""},
			width:    6,
			expected: []string{"      "},
		},
		{
			name:     "single empty row odd width",
			rows:     []string{""},
			width:    5,
			expected: []string{"     "},
		},
		{
			name:     "single row exact width",
			rows:     []string{"hello"},
			width:    5,
			expected: []string{"hello"},
		},
		{
			name:     "single row even padding",
			rows:     []string{"word"},
			width:    10,
			// contentWidth=4, padding=6 -> left=3, right=3
			expected: []string{"   word   "},
		},
		{
			name:     "single row odd padding (extra space on right)",
			rows:     []string{"word"},
			width:    9,
			// contentWidth=4, padding=5 -> left=2, right=3
			expected: []string{"  word   "},
		},
		{
			name:     "uniform rows centered",
			rows:     []string{"abc", "def", "ghi"},
			width:    7,
			// contentWidth=3, padding=4 -> left=2, right=2
			expected: []string{"  abc  ", "  def  ", "  ghi  "},
		},
		{
			name:     "jagged rows normalized and centered",
			rows:     []string{"a", "abcde", "ab"},
			width:    9,
			// max row is "abcde" (len 5). normalized: "a    ", "abcde", "ab   ".
			// padding=9-5=4 -> left=2, right=2
			expected: []string{"  a      ", "  abcde  ", "  ab     "},
		},
		{
			name: "ascii art block",
			rows: []string{
				" _  _ ",
				"| || |",
				"| __ |",
				"|_||_|",
			},
			width: 10,
			// contentWidth=6, padding=4 -> left=2, right=2
			expected: []string{
				"   _  _   ",
				"  | || |  ",
				"  | __ |  ",
				"  |_||_|  ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := alignCenter(tt.rows, tt.width)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("alignCenter(%v, %d) =\n%#v\nwant:\n%#v", tt.rows, tt.width, got, tt.expected)
			}
		})
	}
}

func TestJustify(t *testing.T) {
	spaceBlock := []string{" ", " ", " ", " ", " ", " ", " ", " "}
	wordA := []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8"}
	wordB := []string{"B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8"}
	wordC := []string{"C1", "C2", "C3", "C4", "C5", "C6", "C7", "C8"}

	cSingle, _ := NewContent(spaceBlock, [][]string{wordA}, 0, []int{}, 0)
	cTwoExact, _ := NewContent(spaceBlock, [][]string{wordA, wordB}, 0, []int{1}, 0)
	cTwoExtra, _ := NewContent(spaceBlock, [][]string{wordA, wordB}, 0, []int{1}, 0)
	cThreeEven, _ := NewContent(spaceBlock, [][]string{wordA, wordB, wordC}, 0, []int{1, 1}, 0)
	cThreeOdd, _ := NewContent(spaceBlock, [][]string{wordA, wordB, wordC}, 0, []int{1, 1}, 0)
	cLeadingTrailing, _ := NewContent(spaceBlock, [][]string{wordA, wordB}, 2, []int{1}, 3)
	cNoWords, _ := NewContent(spaceBlock, [][]string{}, 2, []int{}, 1)

	tests := []struct {
		name     string
		content  Content
		width    int
		expected []string
	}{
		{
			name:    "single word without gaps behaves like alignLeft",
			content: cSingle,
			width:   6,
			expected: []string{
				"A1    ", "A2    ", "A3    ", "A4    ",
				"A5    ", "A6    ", "A7    ", "A8    ",
			},
		},
		{
			name:    "two words exact width without extra space",
			content: cTwoExact,
			width:   5,
			expected: []string{
				"A1 B1", "A2 B2", "A3 B3", "A4 B4",
				"A5 B5", "A6 B6", "A7 B7", "A8 B8",
			},
		},
		{
			name:    "two words with extra spaces in single gap",
			content: cTwoExtra,
			width:   8,
			expected: []string{
				"A1    B1", "A2    B2", "A3    B3", "A4    B4",
				"A5    B5", "A6    B6", "A7    B7", "A8    B8",
			},
		},
		{
			name:    "three words with evenly divided extra spaces across gaps",
			content: cThreeEven,
			width:   12,
			expected: []string{
				"A1   B1   C1", "A2   B2   C2", "A3   B3   C3", "A4   B4   C4",
				"A5   B5   C5", "A6   B6   C6", "A7   B7   C7", "A8   B8   C8",
			},
		},
		{
			name:    "three words with odd remainder distributed from the right",
			content: cThreeOdd,
			width:   11,
			expected: []string{
				"A1  B1   C1", "A2  B2   C2", "A3  B3   C3", "A4  B4   C4",
				"A5  B5   C5", "A6  B6   C6", "A7  B7   C7", "A8  B8   C8",
			},
		},
		{
			name:    "multiple words preserving leading and trailing spaces",
			content: cLeadingTrailing,
			width:   13,
			expected: []string{
				"  A1    B1   ", "  A2    B2   ", "  A3    B3   ", "  A4    B4   ",
				"  A5    B5   ", "  A6    B6   ", "  A7    B7   ", "  A8    B8   ",
			},
		},
		{
			name:    "no words with leading and trailing spaces behaves like alignLeft",
			content: cNoWords,
			width:   6,
			expected: []string{
				"      ", "      ", "      ", "      ",
				"      ", "      ", "      ", "      ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := justify(tt.content, tt.width)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("justify() =\n%#v\nwant:\n%#v", got, tt.expected)
			}
		})
	}
}

func TestJustifyUsesRenderedSpaceForOriginalGap(t *testing.T) {
	spaceBlock := []string{
		"...", "...", "...", "...",
		"...", "...", "...", "...",
	}
	wordA := []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8"}
	wordB := []string{"B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8"}

	content, err := NewContent(spaceBlock, [][]string{wordA, wordB}, 1, []int{1}, 1)
	if err != nil {
		t.Fatalf("NewContent() unexpected error: %v", err)
	}

	// Leading: "..." (3)
	// WordA: "A1" (2)
	// Original Gap: "..." (3)
	// WordB: "B1" (2)
	// Trailing: "..." (3)
	// Base content width = 13.
	// Target width = 16 -> extra justification space = 16 - 13 = 3 spaces added to gap.
	targetWidth := 16
	expected := []string{
		"...A1...   B1...",
		"...A2...   B2...",
		"...A3...   B3...",
		"...A4...   B4...",
		"...A5...   B5...",
		"...A6...   B6...",
		"...A7...   B7...",
		"...A8...   B8...",
	}

	got := justify(content, targetWidth)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("justify() =\n%#v\nwant:\n%#v", got, expected)
	}
}

func TestAlign(t *testing.T) {
	spaceBlock := []string{" ", " ", " ", " ", " ", " ", " ", " "}
	wordA := []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8"}
	wordB := []string{"B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8"}

	cTwoWords, err := NewContent(spaceBlock, [][]string{wordA, wordB}, 0, []int{1}, 0)
	if err != nil {
		t.Fatalf("NewContent() unexpected error: %v", err)
	}

	cSingleWord, err := NewContent(spaceBlock, [][]string{wordA}, 0, []int{}, 0)
	if err != nil {
		t.Fatalf("NewContent() unexpected error: %v", err)
	}

	tests := []struct {
		name      string
		content   Content
		width     int
		mode      Mode
		expected  []string
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "left alignment with padding",
			content: cTwoWords,
			width:   8,
			mode:    Left,
			expected: []string{
				"A1 B1   ", "A2 B2   ", "A3 B3   ", "A4 B4   ",
				"A5 B5   ", "A6 B6   ", "A7 B7   ", "A8 B8   ",
			},
			wantErr: false,
		},
		{
			name:    "right alignment with padding",
			content: cTwoWords,
			width:   8,
			mode:    Right,
			expected: []string{
				"   A1 B1", "   A2 B2", "   A3 B3", "   A4 B4",
				"   A5 B5", "   A6 B6", "   A7 B7", "   A8 B8",
			},
			wantErr: false,
		},
		{
			name:    "center alignment with even padding",
			content: cTwoWords,
			width:   9,
			mode:    Center,
			expected: []string{
				"  A1 B1  ", "  A2 B2  ", "  A3 B3  ", "  A4 B4  ",
				"  A5 B5  ", "  A6 B6  ", "  A7 B7  ", "  A8 B8  ",
			},
			wantErr: false,
		},
		{
			name:    "center alignment with odd padding",
			content: cTwoWords,
			width:   8,
			mode:    Center,
			expected: []string{
				" A1 B1  ", " A2 B2  ", " A3 B3  ", " A4 B4  ",
				" A5 B5  ", " A6 B6  ", " A7 B7  ", " A8 B8  ",
			},
			wantErr: false,
		},
		{
			name:    "justify alignment with extra spaces in gap",
			content: cTwoWords,
			width:   8,
			mode:    Justify,
			expected: []string{
				"A1    B1", "A2    B2", "A3    B3", "A4    B4",
				"A5    B5", "A6    B6", "A7    B7", "A8    B8",
			},
			wantErr: false,
		},
		{
			name:    "justify single word behaves like left alignment",
			content: cSingleWord,
			width:   6,
			mode:    Justify,
			expected: []string{
				"A1    ", "A2    ", "A3    ", "A4    ",
				"A5    ", "A6    ", "A7    ", "A8    ",
			},
			wantErr: false,
		},
		{
			name:    "exact width left alignment",
			content: cTwoWords,
			width:   5,
			mode:    Left,
			expected: []string{
				"A1 B1", "A2 B2", "A3 B3", "A4 B4",
				"A5 B5", "A6 B6", "A7 B7", "A8 B8",
			},
			wantErr: false,
		},
		{
			name:    "exact width right alignment",
			content: cTwoWords,
			width:   5,
			mode:    Right,
			expected: []string{
				"A1 B1", "A2 B2", "A3 B3", "A4 B4",
				"A5 B5", "A6 B6", "A7 B7", "A8 B8",
			},
			wantErr: false,
		},
		{
			name:    "exact width center alignment",
			content: cTwoWords,
			width:   5,
			mode:    Center,
			expected: []string{
				"A1 B1", "A2 B2", "A3 B3", "A4 B4",
				"A5 B5", "A6 B6", "A7 B7", "A8 B8",
			},
			wantErr: false,
		},
		{
			name:    "exact width justify alignment",
			content: cTwoWords,
			width:   5,
			mode:    Justify,
			expected: []string{
				"A1 B1", "A2 B2", "A3 B3", "A4 B4",
				"A5 B5", "A6 B6", "A7 B7", "A8 B8",
			},
			wantErr: false,
		},
		{
			name:      "error width is zero",
			content:   cTwoWords,
			width:     0,
			mode:      Left,
			wantErr:   true,
			errSubstr: "width must be greater than 0",
		},
		{
			name:      "error width is negative",
			content:   cTwoWords,
			width:     -5,
			mode:      Left,
			wantErr:   true,
			errSubstr: "width must be greater than 0",
		},
		{
			name:      "error width is smaller than content width",
			content:   cTwoWords,
			width:     4,
			mode:      Left,
			wantErr:   true,
			errSubstr: "width 4 is too small for content width 5",
		},
		{
			name:      "error invalid alignment mode",
			content:   cTwoWords,
			width:     10,
			mode:      Mode("invalid"),
			wantErr:   true,
			errSubstr: "invalid alignment mode: invalid",
		},
		{
			name:      "error empty alignment mode",
			content:   cTwoWords,
			width:     10,
			mode:      Mode(""),
			wantErr:   true,
			errSubstr: "invalid alignment mode: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Align(tt.content, tt.width, tt.mode)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Align() =\n%#v\nwant:\n%#v", got, tt.expected)
			}
		})
	}
}
