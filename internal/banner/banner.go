package banner

import (
	"errors"
	"os"
	"strings"
)

func Path(name string) (string, error) {
	switch name {
	case "standard":
		return "banners/standard.txt", nil
	case "shadow":
		return "banners/shadow.txt", nil
	case "thinkertoy":
		return "banners/thinkertoy.txt", nil
	default:
		return "", errors.New("unsupported banner")
	}

}

func Load(path string) (map[rune][]string, error) {
	content, err := readFile(path)
	if err != nil {
		return nil, err
	}
	return parseBanner(content)
}

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func parseBanner(content string) (map[rune][]string, error) {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")
	if lines[0] != "" {
		return nil, errors.New("invalid banner format")
	}
	shapes := make(map[rune][]string)

	for i := 32; i <= 126; i++ {
		start := 1 + (i-32)*9
		end := start + 8

		if end > len(lines) {
			return nil, errors.New("invalid banner format")
		}

		if i != 126 {
			if end >= len(lines) || lines[end] != "" {
				return nil, errors.New("invalid banner format")
			}
		}

		shape := lines[start:end]
		shapes[rune(i)] = shape
	}

	return shapes, nil
}
