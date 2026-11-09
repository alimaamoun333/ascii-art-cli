package renderer

import (
	"testing"
)

func buildTestShapes() map[rune][]string {
	return map[rune][]string{
		'A': {"A0", "A1", "A2", "A3", "A4", "A5", "A6", "A7"},
		'B': {"B0", "B1", "B2", "B3", "B4", "B5", "B6", "B7"},
		' ': {"S0", "S1", "S2", "S3", "S4", "S5", "S6", "S7"},
	}
}

func TestRenderCharacterExistingCharacter(t *testing.T) {
	shapes := buildTestShapes()

	shape, err := RenderCharacter(shapes, 'A')
	if err != nil {
		t.Fatalf("the character not found %v", err)
	}
	if len(shape) != 8 {
		t.Fatalf("expected 8 rows, got %d", len(shape))
	}
	if shape[0] != "A0" {
		t.Fatalf("expected first row A0, got %q", shape[0])
	}
}

func TestRenderCharacterNotExistingCharacter(t *testing.T) {
	shapes := buildTestShapes()
	shape, err := RenderCharacter(shapes, 'Z')
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if shape != nil {
		t.Fatalf("expected nil shape, got %v", shape)
	}

}

func TestRenderLineMultipleCharacters(t *testing.T) {
	shapes := buildTestShapes()

	rows, err := RenderLine(shapes, "AB")

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(rows) != 8 {
		t.Fatalf("expected 8 rows, got %d", len(rows))
	}

	if rows[0] != "A0B0" {
		t.Fatalf("expected first row A0B0, got %q", rows[0])
	}
}

func TestRenderLineUnsupportedCharacter(t *testing.T) {
	shapes := buildTestShapes()

	rows, err := RenderLine(shapes, "AZ")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if rows != nil {
		t.Fatalf("expected nil rows, got %v", rows)
	}
}

func TestRenderLineEmptyCharacter(t *testing.T) {
	shapes := buildTestShapes()
	rows, err := RenderLine(shapes, "")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if rows != nil {
		t.Fatalf("expected nil rows, got %v", rows)
	}
}

func TestRenderLineSpaceCharacter(t *testing.T) {
	shapes := buildTestShapes()
	rows, err := RenderLine(shapes, " ")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(rows) != 8 {
		t.Fatalf("expected 8 rows, got %d", len(rows))
	}

	if rows[0] != "S0" {
		t.Fatalf("expected first row S0, got %q", rows[0])
	}
}

func TestRenderMultipleLines(t *testing.T) {
	shapes := buildTestShapes()

	rows, err := Render(shapes, "A\nB")

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(rows) != 16 {
		t.Fatalf("expected 16 rows, got %d", len(rows))
	}

	if rows[0] != "A0" {
		t.Fatalf("expected row 0 to be A0, got %q", rows[0])
	}

	if rows[8] != "B0" {
		t.Fatalf("expected row 8 to be B0, got %q", rows[8])
	}
}

func TestRenderConsecutiveNewlines(t *testing.T) {
	shapes := buildTestShapes()

	rows, err := Render(shapes, "A\n\nB")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(rows) != 24 {
		t.Fatalf("expected 24 rows, got %d", len(rows))
	}

	if rows[0] != "A0" {
		t.Fatalf("expected row 0 to be A0, got %q", rows[0])
	}

	if rows[8] != "" {
		t.Fatalf("expected row 8 to be empty, got %q", rows[8])
	}

	if rows[15] != "" {
		t.Fatalf("expected row 15 to be empty, got %q", rows[15])
	}

	if rows[16] != "B0" {
		t.Fatalf("expected row 16 to be B0, got %q", rows[16])
	}
}

func TestRenderUnsupportedCharacter(t *testing.T) {
	shapes := buildTestShapes()
	rows, err := Render(shapes, "A\nZ")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if rows != nil {
		t.Fatalf("expected nil rows, got %v", rows)
	}
}

func TestRenderEmptyText(t *testing.T) {
	shapes := buildTestShapes()

	rows, err := Render(shapes, "")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if rows != nil {
		t.Fatalf("expected nil rows, got %v", rows)
	}
}
