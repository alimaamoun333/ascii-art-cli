package output

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWrite(t *testing.T) {
	t.Run("success writing new file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "output.txt")
		content := "Hello, World!\n"

		err := Write(filePath, content)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read written file: %v", err)
		}

		if string(data) != content {
			t.Errorf("expected file content %q, got %q", content, string(data))
		}
	})

	t.Run("success overwriting existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "overwrite.txt")

		err := os.WriteFile(filePath, []byte("initial data"), 0644)
		if err != nil {
			t.Fatalf("failed to create initial file: %v", err)
		}

		newContent := "overwritten content"
		err = Write(filePath, newContent)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		if string(data) != newContent {
			t.Errorf("expected file content %q, got %q", newContent, string(data))
		}
	})

	t.Run("success writing empty content", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "empty.txt")
		content := ""

		err := Write(filePath, content)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read written file: %v", err)
		}

		if string(data) != content {
			t.Errorf("expected empty file content, got %q", string(data))
		}
	})

	t.Run("error when directory does not exist", func(t *testing.T) {
		tempDir := t.TempDir()
		invalidPath := filepath.Join(tempDir, "non_existent_folder", "output.txt")

		err := Write(invalidPath, "content")
		if err == nil {
			t.Fatalf("expected error when writing to non-existent directory, got nil")
		}
	})
}
