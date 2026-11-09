package layout

import (
	"fmt"
	"strings"
)

type Boundary struct {
	Start int
	End   int
}

type LineLayout struct {
	Leading  int
	Words    []string
	Gaps     []int
	Trailing int
}

func CalculateCharacterBoundaries(shapes map[rune][]string, text string) ([][]Boundary, error) {
	boundaries := [][]Boundary{}

	for row := 0; row < 8; row++ {
		currentColumn := 0
		rowBoundaries := []Boundary{}

		for _, char := range text {
			shape, ok := shapes[char]
			if !ok {
				return nil, fmt.Errorf("unsupported character: %q", char)
			}

			width := len(shape[row])

			boundary := Boundary{
				Start: currentColumn,
				End:   currentColumn + width,
			}

			rowBoundaries = append(rowBoundaries, boundary)
			currentColumn = boundary.End
		}

		boundaries = append(boundaries, rowBoundaries)
	}

	return boundaries, nil
}

func CharacterRangeToRenderedRange(rowBoundaries []Boundary, start int, end int) (Boundary, error) {
	if start < 0 || end > len(rowBoundaries) || start >= end {
		return Boundary{}, fmt.Errorf("invalid character range")
	}

	return Boundary{
		Start: rowBoundaries[start].Start,
		End:   rowBoundaries[end-1].End,
	}, nil
}

func CalculateMultilineBoundaries(shapes map[rune][]string, text string) ([][][]Boundary, error) {
	lines := strings.Split(text, "\n")
	boundaries := [][][]Boundary{}

	for _, line := range lines {
		lineBoundaries, err := CalculateCharacterBoundaries(shapes, line)
		if err != nil {
			return nil, err
		}
		boundaries = append(boundaries, lineBoundaries)
	}
	return boundaries, nil
}

func ParseLine(text string) LineLayout {
	var result LineLayout
	i := 0

	// Count leading spaces
	for i < len(text) && text[i] == ' ' {
		result.Leading++
		i++
	}
	for i < len(text) {
		startWord := i
		for i < len(text) && text[i] != ' ' {
			i++
		}
		result.Words = append(result.Words, text[startWord:i])
		spaceCount := 0
		for i < len(text) {
			if text[i] == ' ' {
				spaceCount++
				i++
			} else {
				break
			}
		}
		if i < len(text) {
			result.Gaps = append(result.Gaps, spaceCount)
		} else {
			result.Trailing = spaceCount
		}
	}
	return result
}
