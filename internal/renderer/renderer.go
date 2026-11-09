package renderer

import (
	"errors"
	"strings"
)

// RenderCharacter returns the glyph rows for one character.
func RenderCharacter(shapes map[rune][]string, char rune) ([]string, error) {
	shape, ok := shapes[char]

	if ok == false {
		return nil, errors.New("character does not exist")
	}

	return shape, nil
}

// RenderLine renders one non-empty logical line into 8 ASCII-art rows.
func RenderLine(shapes map[rune][]string, text string) ([]string, error) {
	if text == "" {
		return nil, errors.New("text cannot be empty")
	}
	rows := make([]string, 8)

	for _, char := range text {
		shape, err := RenderCharacter(shapes, char)
		if err != nil {
			return nil, err
		}

		for i := 0; i < 8; i++ {
			rows[i] += shape[i]
		}
	}

	return rows, nil
}

// Render renders complete text, including logical newlines, into ASCII-art rows.
func Render(shapes map[rune][]string, text string) ([]string, error) {
	if text == "" {
		return nil, errors.New("text cannot be empty")
	}

	lines := strings.Split(text, "\n")
	rows := []string{}
	for _, line := range lines {
		if line == "" {
			for i := 0; i < 8; i++ {
				rows = append(rows, "")
			}
		} else {
			renderedRows, err := RenderLine(shapes, line)
			if err != nil {
				return nil, err
			}
			rows = append(rows, renderedRows...)

		}

	}
	return rows, nil
}

// if we need a string edit the function to be Render(...) (string, error)  and add test = strings.Join(rows, "\n") and return text
