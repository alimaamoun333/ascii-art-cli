package output

import (
	"fmt"
	"os"
)

func Write(filename string, content string) error {
	contentBytes := []byte(content)

	err := os.WriteFile(filename, contentBytes, 0644)
	if err != nil {
		return fmt.Errorf("write output file %q: %w", filename, err)
	}

	return nil
}
