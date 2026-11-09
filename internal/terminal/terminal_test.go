package terminal

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestParseWindowsOutput(t *testing.T) {
	tests := []struct {
		name        string
		output      []byte
		expected    int
		expectError bool
		errContains string
	}{
		{
			name:        "valid standard width",
			output:      []byte("80"),
			expected:    80,
			expectError: false,
		},
		{
			name:        "valid width with whitespace and newlines",
			output:      []byte("  120 \r\n"),
			expected:    120,
			expectError: false,
		},
		{
			name:        "valid single column width",
			output:      []byte("1"),
			expected:    1,
			expectError: false,
		},
		{
			name:        "empty output",
			output:      []byte(""),
			expected:    -1,
			expectError: true,
			errContains: "failed to parse Windows terminal width",
		},
		{
			name:        "non-numeric output",
			output:      []byte("invalid"),
			expected:    -1,
			expectError: true,
			errContains: "failed to parse Windows terminal width",
		},
		{
			name:        "float output",
			output:      []byte("80.5"),
			expected:    -1,
			expectError: true,
			errContains: "failed to parse Windows terminal width",
		},
		{
			name:        "zero width",
			output:      []byte("0"),
			expected:    -1,
			expectError: true,
			errContains: "invalid terminal width: 0",
		},
		{
			name:        "negative width",
			output:      []byte("-10"),
			expected:    -1,
			expectError: true,
			errContains: "invalid terminal width: -10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseWindowsOutput(tt.output)
			if tt.expectError {
				if err == nil {
					t.Fatalf("parseWindowsOutput(%q) expected error, got nil", string(tt.output))
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("parseWindowsOutput(%q) error %q does not contain %q", string(tt.output), err.Error(), tt.errContains)
				}
				if got != -1 {
					t.Errorf("parseWindowsOutput(%q) = %d, want -1 on error", string(tt.output), got)
				}
			} else {
				if err != nil {
					t.Fatalf("parseWindowsOutput(%q) unexpected error: %v", string(tt.output), err)
				}
				if got != tt.expected {
					t.Errorf("parseWindowsOutput(%q) = %d, want %d", string(tt.output), got, tt.expected)
				}
			}
		})
	}
}

func TestParseSttyOutput(t *testing.T) {
	tests := []struct {
		name        string
		output      []byte
		expected    int
		expectError bool
		errContains string
	}{
		{
			name:        "valid standard rows and columns",
			output:      []byte("24 80"),
			expected:    80,
			expectError: false,
		},
		{
			name:        "valid output with extra whitespace and newline",
			output:      []byte("   30    120   \n"),
			expected:    120,
			expectError: false,
		},
		{
			name:        "empty output",
			output:      []byte(""),
			expected:    -1,
			expectError: true,
			errContains: "unexpected output from stty size",
		},
		{
			name:        "single field output",
			output:      []byte("80"),
			expected:    -1,
			expectError: true,
			errContains: "unexpected output from stty size",
		},
		{
			name:        "three fields output",
			output:      []byte("24 80 extra"),
			expected:    -1,
			expectError: true,
			errContains: "unexpected output from stty size",
		},
		{
			name:        "non-numeric width",
			output:      []byte("24 abc"),
			expected:    -1,
			expectError: true,
			errContains: "failed to parse stty terminal width",
		},
		{
			name:        "float width",
			output:      []byte("24 80.5"),
			expected:    -1,
			expectError: true,
			errContains: "failed to parse stty terminal width",
		},
		{
			name:        "zero width",
			output:      []byte("24 0"),
			expected:    -1,
			expectError: true,
			errContains: "invalid terminal width: 0",
		},
		{
			name:        "negative width",
			output:      []byte("24 -5"),
			expected:    -1,
			expectError: true,
			errContains: "invalid terminal width: -5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSttyOutput(tt.output)
			if tt.expectError {
				if err == nil {
					t.Fatalf("parseSttyOutput(%q) expected error, got nil", string(tt.output))
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("parseSttyOutput(%q) error %q does not contain %q", string(tt.output), err.Error(), tt.errContains)
				}
				if got != -1 {
					t.Errorf("parseSttyOutput(%q) = %d, want -1 on error", string(tt.output), got)
				}
			} else {
				if err != nil {
					t.Fatalf("parseSttyOutput(%q) unexpected error: %v", string(tt.output), err)
				}
				if got != tt.expected {
					t.Errorf("parseSttyOutput(%q) = %d, want %d", string(tt.output), got, tt.expected)
				}
			}
		})
	}
}

func TestWidth_Mocked(t *testing.T) {
	origOS := currentOS
	origRunner := cmdRunner
	defer func() {
		currentOS = origOS
		cmdRunner = origRunner
	}()

	tests := []struct {
		name         string
		targetOS     string
		expectedCmd  string
		expectedArgs []string
		mockOutput   []byte
		mockErr      error
		expected     int
		expectError  bool
		errContains  string
	}{
		// Windows tests
		{
			name:         "windows valid output",
			targetOS:     "windows",
			expectedCmd:  "powershell",
			expectedArgs: []string{"-command", "$Host.UI.RawUI.WindowSize.Width"},
			mockOutput:   []byte("100\r\n"),
			mockErr:      nil,
			expected:     100,
			expectError:  false,
		},
		{
			name:         "windows command execution failure",
			targetOS:     "windows",
			expectedCmd:  "powershell",
			expectedArgs: []string{"-command", "$Host.UI.RawUI.WindowSize.Width"},
			mockOutput:   nil,
			mockErr:      errors.New("powershell not found"),
			expected:     -1,
			expectError:  true,
			errContains:  "failed to get Windows terminal width",
		},
		{
			name:         "windows parse failure",
			targetOS:     "windows",
			expectedCmd:  "powershell",
			expectedArgs: []string{"-command", "$Host.UI.RawUI.WindowSize.Width"},
			mockOutput:   []byte("not_a_number"),
			mockErr:      nil,
			expected:     -1,
			expectError:  true,
			errContains:  "failed to parse Windows terminal width",
		},
		{
			name:         "windows invalid non-positive width",
			targetOS:     "windows",
			expectedCmd:  "powershell",
			expectedArgs: []string{"-command", "$Host.UI.RawUI.WindowSize.Width"},
			mockOutput:   []byte("0"),
			mockErr:      nil,
			expected:     -1,
			expectError:  true,
			errContains:  "invalid terminal width: 0",
		},

		// Linux tests
		{
			name:         "linux valid output",
			targetOS:     "linux",
			expectedCmd:  "stty",
			expectedArgs: []string{"size"},
			mockOutput:   []byte("25 80\n"),
			mockErr:      nil,
			expected:     80,
			expectError:  false,
		},
		{
			name:         "linux command execution failure",
			targetOS:     "linux",
			expectedCmd:  "stty",
			expectedArgs: []string{"size"},
			mockOutput:   nil,
			mockErr:      errors.New("stty command failed"),
			expected:     -1,
			expectError:  true,
			errContains:  "failed to get terminal width using stty",
		},
		{
			name:         "linux unexpected output format",
			targetOS:     "linux",
			expectedCmd:  "stty",
			expectedArgs: []string{"size"},
			mockOutput:   []byte("single_number"),
			mockErr:      nil,
			expected:     -1,
			expectError:  true,
			errContains:  "unexpected output from stty size",
		},
		{
			name:         "linux parse failure",
			targetOS:     "linux",
			expectedCmd:  "stty",
			expectedArgs: []string{"size"},
			mockOutput:   []byte("25 abc"),
			mockErr:      nil,
			expected:     -1,
			expectError:  true,
			errContains:  "failed to parse stty terminal width",
		},
		{
			name:         "linux invalid non-positive width",
			targetOS:     "linux",
			expectedCmd:  "stty",
			expectedArgs: []string{"size"},
			mockOutput:   []byte("25 -4"),
			mockErr:      nil,
			expected:     -1,
			expectError:  true,
			errContains:  "invalid terminal width: -4",
		},

		// Darwin (macOS) tests
		{
			name:         "darwin valid output",
			targetOS:     "darwin",
			expectedCmd:  "stty",
			expectedArgs: []string{"size"},
			mockOutput:   []byte("30 140\n"),
			mockErr:      nil,
			expected:     140,
			expectError:  false,
		},

		// Unsupported OS tests
		{
			name:         "unsupported operating system",
			targetOS:     "freebsd",
			expectedCmd:  "",
			expectedArgs: nil,
			mockOutput:   nil,
			mockErr:      nil,
			expected:     -1,
			expectError:  true,
			errContains:  "unsupported OS: freebsd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentOS = tt.targetOS
			var calledCmd string
			var calledArgs []string
			cmdRunner = func(name string, arg ...string) ([]byte, error) {
				calledCmd = name
				calledArgs = arg
				return tt.mockOutput, tt.mockErr
			}

			got, err := Width()

			if tt.targetOS == "windows" || tt.targetOS == "linux" || tt.targetOS == "darwin" {
				if calledCmd != tt.expectedCmd {
					t.Errorf("Width() executed command %q, want %q", calledCmd, tt.expectedCmd)
				}
				if !reflect.DeepEqual(calledArgs, tt.expectedArgs) {
					t.Errorf("Width() passed args %v, want %v", calledArgs, tt.expectedArgs)
				}
			} else {
				if calledCmd != "" {
					t.Errorf("Width() unexpectedly executed command %q for unsupported OS", calledCmd)
				}
			}

			if tt.expectError {
				if err == nil {
					t.Fatalf("Width() expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Width() error %q does not contain %q", err.Error(), tt.errContains)
				}
				if got != -1 {
					t.Errorf("Width() = %d, want -1 on error", got)
				}
			} else {
				if err != nil {
					t.Fatalf("Width() unexpected error: %v", err)
				}
				if got != tt.expected {
					t.Errorf("Width() = %d, want %d", got, tt.expected)
				}
			}
		})
	}
}

// TestWidth_Live is an integration-style test for live environment execution.
// It verifies that calling Width() in the current runtime environment adheres to the API contract:
// returning a positive width (>0) with no error, or -1 with a non-nil error if no terminal is available.
func TestWidth_Live(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live terminal integration test in -short mode")
	}

	// Restore default runner and os to ensure real execution is tested
	origOS := currentOS
	origRunner := cmdRunner
	defer func() {
		currentOS = origOS
		cmdRunner = origRunner
	}()

	currentOS = origOS
	cmdRunner = origRunner

	width, err := Width()
	if err != nil {
		t.Logf("Live Width() returned error (expected in non-interactive/headless environments): %v", err)
		if width != -1 {
			t.Errorf("Width() on error returned %d, want -1", width)
		}
	} else {
		t.Logf("Live Width() returned width: %d", width)
		if width <= 0 {
			t.Errorf("Width() on success returned invalid width %d (must be > 0)", width)
		}
	}
}
