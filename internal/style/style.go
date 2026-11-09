package style

import (
	"ascii-art-justify/internal/layout"
	"fmt"
	"strings"
)

type Range struct {
	Start int
	End   int
}

func FindRanges(text string, target string) []Range {
	ranges := []Range{}
	if target == "" {
		return ranges
	}
	searchStart := 0
	for {
		relativeIndex := strings.Index(text[searchStart:], target)
		if relativeIndex == -1 {
			return ranges
		}
		matchIndex := searchStart + relativeIndex
		rangeFound := Range{
			Start: matchIndex,
			End:   matchIndex + len(target),
		}
		ranges = append(ranges, rangeFound)
		searchStart = matchIndex + 1
	}

}

func MergeRanges(ranges []Range) []Range {
	if len(ranges) == 0 {
		return []Range{}
	}

	merged := []Range{}
	current := ranges[0]

	for _, next := range ranges[1:] {
		if next.Start <= current.End {
			if next.End > current.End {
				current.End = next.End
			}
			continue
		}

		merged = append(merged, current)
		current = next
	}

	merged = append(merged, current)

	return merged
}

func ANSIColor(color string) (string, error) {
	switch color {
	case "red":
		return "\033[31m", nil
	case "green":
		return "\033[32m", nil
	case "yellow":
		return "\033[33m", nil
	case "blue":
		return "\033[34m", nil
	case "magenta":
		return "\033[35m", nil
	case "cyan":
		return "\033[36m", nil
	default:
		return "", fmt.Errorf("invalid or unsupported request color: %s", color)
	}
}

func ApplyColorToRange(row string, renderedStart int, renderedEnd int, colorCode string) (string, error) {

	const ansiReset = "\033[0m"
	if renderedStart < 0 ||
		renderedEnd > len(row) ||
		renderedStart >= renderedEnd {
		return "", fmt.Errorf("invalid rendered range")
	}

	styled := row[:renderedStart] +
		colorCode +
		row[renderedStart:renderedEnd] +
		ansiReset +
		row[renderedEnd:]

	return styled, nil
}

func ApplyColorToRanges(row string, ranges []layout.Boundary, colorCode string) (string, error) {
	styled := row

	for i := len(ranges) - 1; i >= 0; i-- {
		var err error

		styled, err = ApplyColorToRange(
			styled,
			ranges[i].Start,
			ranges[i].End,
			colorCode,
		)
		if err != nil {
			return "", err
		}
	}

	return styled, nil
}

func ApplyColor(rows []string, text string, target string, color string, boundaries [][]layout.Boundary) ([]string, error) {
	colorCode, err := ANSIColor(color)
	if err != nil {
		return nil, err
	}
	var ranges []Range

	if target == "" {
		if len(text) > 0 {
			ranges = []Range{
				{
					Start: 0,
					End:   len(text),
				},
			}
		}
	} else {
		ranges = FindRanges(text, target)
		ranges = MergeRanges(ranges)
	}

	styledRows := make([]string, len(rows))

	for rowIndex, row := range rows {
		renderedRanges := []layout.Boundary{}

		for _, characterRange := range ranges {
			renderedRange, err := layout.CharacterRangeToRenderedRange(
				boundaries[rowIndex],
				characterRange.Start,
				characterRange.End,
			)
			if err != nil {
				return nil, err
			}

			renderedRanges = append(renderedRanges, renderedRange)
		}

		styledRow, err := ApplyColorToRanges(
			row,
			renderedRanges,
			colorCode,
		)
		if err != nil {
			return nil, err
		}

		styledRows[rowIndex] = styledRow
	}

	return styledRows, nil
}

func ApplyColorMultiline(rows []string, text string, target string, color string, boundaries [][][]layout.Boundary) ([]string, error) {
	lines := strings.Split(text, "\n")
	if len(lines) != len(boundaries) {
		return nil, fmt.Errorf("line count does not match boundaries count")
	}
	if len(rows) != len(lines)*8 {
		return nil, fmt.Errorf("rows count mismatch")
	}

	styledRows := make([]string, 0, len(rows))

	for i, line := range lines {
		lineRows := rows[i*8 : (i+1)*8]
		lineBoundaries := boundaries[i]

		styledLineRows, err := ApplyColor(lineRows, line, target, color, lineBoundaries)
		if err != nil {
			return nil, err
		}

		styledRows = append(styledRows, styledLineRows...)
	}

	return styledRows, nil
}
