package banner

import (
	"strings"
	"testing"
)

func buildValidBanner() string {
	var builder strings.Builder

	builder.WriteByte('\n')

	for i := 32; i <= 126; i++ {
		for j := 0; j < 8; j++ {
			builder.WriteByte(byte(i))

			if !(i == 126 && j == 7) {
				builder.WriteByte('\n')
			}
		}

		if i != 126 {
			builder.WriteByte('\n')
		}
	}

	return builder.String()
}

func TestParseBannerValidBanner(t *testing.T) {
	banner := buildValidBanner()
	shapes, err := parseBanner(banner)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(shapes) != 95 {
		t.Errorf("expected 95 shapes, got %d", len(shapes))
	}

	for character, shape := range shapes {
		if len(shape) != 8 {
			t.Errorf(
				"character %q: expected 8 rows, got %d",
				character,
				len(shape),
			)
		}
	}

	if shapes['A'][0] != "A" {
		t.Errorf("expected A shape to contain A, got %q", shapes['A'][0])
	}
}

func TestParseBannerInvalidBanner(t *testing.T) {
	banner := buildValidBanner()
	banner = banner[1:] // Remove the initial newline.
	_, err := parseBanner(banner)
	if err == nil {
		t.Errorf("expected error for invalid banner, got nil")
	}
}

func TestParseBannerMissingRows(t *testing.T) {
	banner := buildValidBanner()

	banner = banner[:len(banner)-2]

	shapes, err := parseBanner(banner)

	if err == nil {
		t.Fatalf("expected error for missing rows, got nil")
	}

	if shapes != nil {
		t.Errorf("expected nil shapes for missing rows, got %v", shapes)
	}
}

func TestParseBannerMissingSeparator(t *testing.T) {
	banner := buildValidBanner()
	lines := strings.Split(banner, "\n")
	lines = append(lines[:9], lines[10:]...)
	banner = strings.Join(lines, "\n")
	shapes, err := parseBanner(banner)
	if err == nil {
		t.Fatalf("expected error for missing separator, got nil")
	}

	if shapes != nil {
		t.Errorf("expected nil shapes for missing separator, got %v", shapes)
	}
}

func TestPath(t *testing.T) {
	tests := []struct {
		name         string
		banner       string
		expectedPath string
		expectError  bool
	}{
		{
			name:         "standard",
			banner:       "standard",
			expectedPath: "banners/standard.txt",
			expectError:  false,
		}, {
			name:         "shadow",
			banner:       "shadow",
			expectedPath: "banners/shadow.txt",
			expectError:  false,
		}, {
			name:         "thinkertoy",
			banner:       "thinkertoy",
			expectedPath: "banners/thinkertoy.txt",
			expectError:  false,
		}, {
			name:         "unsupported banner",
			banner:       "banana",
			expectedPath: "",
			expectError:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path, err := Path(test.banner)
			if !test.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if test.expectError && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if path != test.expectedPath {
				t.Errorf("expected path %q, got %q", test.expectedPath, path)
			}
		})
	}
}
