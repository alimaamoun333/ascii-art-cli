package alignment

import (
	"fmt"
	"strings"
)

type Mode string

const (
	Left    Mode = "left"
	Center  Mode = "center"
	Right   Mode = "right"
	Justify Mode = "justify"
)

type Content struct {
	space    []string
	words    [][]string
	leading  int
	gaps     []int
	trailing int
}

func NewContent(space []string, words [][]string, leading int, gaps []int, trailing int) (Content, error) {
	if len(space) != 8 {
		return Content{}, fmt.Errorf("space block must contain 8 rows, got %d", len(space))
	}
	for _, word := range words {
		if len(word) != 8 {
			return Content{}, fmt.Errorf("word block must contain 8 rows, got %d", len(word))
		}
	}
	expectedGaps := 0
	if len(words) > 0 {
		expectedGaps = len(words) - 1
	}

	if len(gaps) != expectedGaps {
		return Content{}, fmt.Errorf(
			"expected %d gaps for %d words, got %d",
			expectedGaps,
			len(words),
			len(gaps),
		)
	}
	for _, gap := range gaps {
		if gap <= 0 {
			return Content{}, fmt.Errorf(
				"gap must be greater than 0, got %d",
				gap,
			)
		}
	}
	if leading < 0 {
		return Content{}, fmt.Errorf(
			"leading must be greater than or equal to 0, got %d",
			leading,
		)
	}
	if trailing < 0 {
		return Content{}, fmt.Errorf(
			"trailing must be greater than or equal to 0, got %d",
			trailing,
		)
	}
	return Content{
		space:    space,
		words:    words,
		leading:  leading,
		gaps:     gaps,
		trailing: trailing}, nil
}

func blockWidth(rows []string) int {
	maxWidth := 0
	for _, row := range rows {
		if len(row) > maxWidth {
			maxWidth = len(row)
		}
	}
	return maxWidth
}

func normalizeBlock(rows []string) []string {
	width := blockWidth(rows)
	normalized := make([]string, len(rows))
	for i, row := range rows {
		if len(row) < width {
			normalized[i] = row + strings.Repeat(" ", width-len(row))
		} else {
			normalized[i] = row
		}
	}
	return normalized
}

func buildNormalBlock(content Content) []string {
	normalizedSpace := normalizeBlock(content.space)
	normalizedWords := make([][]string, len(content.words))
	for i, word := range content.words {
		normalizedWords[i] = normalizeBlock(word)
	}
	rows := make([]string, 8)
	for row := 0; row < 8; row++ {
		// leading
		for i := 0; i < content.leading; i++ {
			rows[row] += normalizedSpace[row]
		}
		// words and gaps
		for wordIndex, word := range normalizedWords {
			rows[row] += word[row]
			// gaps
			if wordIndex < len(content.gaps) {
				for i := 0; i < content.gaps[wordIndex]; i++ {
					rows[row] += normalizedSpace[row]
				}
			}
		}
		// trailing
		for i := 0; i < content.trailing; i++ {
			rows[row] += normalizedSpace[row]
		}
	}
	return rows
}

func alignRight(rows []string, width int) []string {
	normalized := normalizeBlock(rows)
	contentWidth := blockWidth(normalized)
	leftPadding := width - contentWidth
	aligned := make([]string, len(normalized))
	for i, row := range normalized {
		aligned[i] = strings.Repeat(" ", leftPadding) + row
	}
	return aligned
}

func alignLeft(rows []string, width int) []string {
	normalized := normalizeBlock(rows)
	contentWidth := blockWidth(normalized)
	rightPadding := width - contentWidth

	aligned := make([]string, len(normalized))
	for i, row := range normalized {
		aligned[i] = row + strings.Repeat(" ", rightPadding)
	}

	return aligned
}

// we keep the old implementation of alignCenter commented out for reference. The new implementation normalizes the block first and then calculates the padding needed to center the content, ensuring consistent alignment even with varying row lengths.
/*
func alignCenter(rows []string, width int) []string {
	aligned := make([]string, len(rows))

	for i, row := range rows {
		if len(row) < width {
			padding := width - len(row)
			leftPadding := padding / 2
			rightPadding := padding - leftPadding
			aligned[i] = strings.Repeat(" ", leftPadding) + row + strings.Repeat(" ", rightPadding)
		} else {
			aligned[i] = row
		}
	}

	return aligned
}
*/

// with the new implementation of alignCenter, we first normalize the block to ensure all rows have the same width. Then we calculate the total padding needed to center the content and distribute it evenly on both sides. This approach ensures that the alignment is consistent even if the input rows have varying lengths.
func alignCenter(rows []string, width int) []string {
	normalized := normalizeBlock(rows)

	contentWidth := blockWidth(normalized)
	padding := width - contentWidth

	leftPadding := padding / 2
	rightPadding := padding - leftPadding

	aligned := make([]string, len(normalized))
	for i, row := range normalized {
		aligned[i] = strings.Repeat(" ", leftPadding) + row + strings.Repeat(" ", rightPadding)
	}
	return aligned
}

func justify(content Content, width int) []string {
	// No inter-word gaps means there is nothing to justify.
	// Single-word content behaves like left alignment.
	if len(content.gaps) == 0 {
		normal := buildNormalBlock(content)
		return alignLeft(normal, width)
	}

	// Find how many extra terminal columns we need to fill.
	normal := buildNormalBlock(content)
	contentWidth := blockWidth(normal)
	extraSpace := width - contentWidth

	gapCount := len(content.gaps)
	baseExtra := extraSpace / gapCount
	remainder := extraSpace % gapCount

	// Make every rendered block rectangular before composition.
	normalizedSpace := normalizeBlock(content.space)

	normalizedWords := make([][]string, len(content.words))
	for i, word := range content.words {
		normalizedWords[i] = normalizeBlock(word)
	}

	// One logical ASCII-art line contains 8 rendered rows.
	rows := make([]string, 8)

	for row := 0; row < 8; row++ {
		// Preserve original leading spaces.
		for i := 0; i < content.leading; i++ {
			rows[row] += normalizedSpace[row]
		}

		// Compose word → original gap → extra justified gap.
		for wordIndex, word := range normalizedWords {
			rows[row] += word[row]

			if wordIndex < len(content.gaps) {
				// Preserve the original source gap.
				for i := 0; i < content.gaps[wordIndex]; i++ {
					rows[row] += normalizedSpace[row]
				}

				// Add justification space.
				extra := baseExtra

				// Uneven remainder goes to gaps from the right.
				if wordIndex >= gapCount-remainder {
					extra++
				}

				rows[row] += strings.Repeat(" ", extra)
			}
		}

		// Preserve original trailing spaces.
		for i := 0; i < content.trailing; i++ {
			rows[row] += normalizedSpace[row]
		}
	}

	return rows
}

func Align(content Content, width int, mode Mode) ([]string, error) {
	if width <= 0 {
		return nil, fmt.Errorf("width must be greater than 0, got %d", width)
	}

	normal := buildNormalBlock(content)
	contentWidth := blockWidth(normal)

	if width < contentWidth {
		return nil, fmt.Errorf(
			"width %d is too small for content width %d",
			width,
			contentWidth,
		)
	}

	switch mode {
	case Left:
		return alignLeft(normal, width), nil
	case Right:
		return alignRight(normal, width), nil
	case Center:
		return alignCenter(normal, width), nil
	case Justify:
		return justify(content, width), nil
	default:
		return nil, fmt.Errorf("invalid alignment mode: %v", mode)
	}
}
