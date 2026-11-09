package cli

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedOptions Options
		expectError     bool
		expectedError   string
	}{
		// Normal Mode
		{
			name:          "empty arguments",
			args:          []string{},
			expectError:   true,
			expectedError: Help(),
		},
		{
			name:          "missing input (only binary name)",
			args:          []string{"ascii-art"},
			expectError:   true,
			expectedError: Help(),
		},
		{
			name: "normal mode valid single text",
			args: []string{"ascii-art", "Hello"},
			expectedOptions: Options{
				Text:           "Hello",
				Banner:         "standard",
				Color:          "",
				Target:         "",
				ColorRequested: false,
			},
		},
		{
			name: "normal mode text with newlines",
			args: []string{"ascii-art", "Hello\\nWorld"},
			expectedOptions: Options{
				Text:           "Hello\nWorld",
				Banner:         "standard",
				Color:          "",
				Target:         "",
				ColorRequested: false,
			},
		},
		{
			name: "normal mode valid with custom banner",
			args: []string{"ascii-art", "Hello", "shadow"},
			expectedOptions: Options{
				Text:           "Hello",
				Banner:         "shadow",
				Color:          "",
				Target:         "",
				ColorRequested: false,
			},
		},
		{
			name:          "normal mode too many arguments",
			args:          []string{"ascii-art", "Hello", "shadow", "extra"},
			expectError:   true,
			expectedError: NormalUsage(),
		},

		// Color Mode
		{
			name:          "color mode missing text",
			args:          []string{"ascii-art", "--color=red"},
			expectError:   true,
			expectedError: ColorUsage(),
		},
		{
			name: "color mode whole text",
			args: []string{"ascii-art", "--color=red", "Hello"},
			expectedOptions: Options{
				Text:           "Hello",
				Banner:         "standard",
				Color:          "red",
				Target:         "",
				ColorRequested: true,
			},
		},
		{
			name: "color mode whole text with newlines",
			args: []string{"ascii-art", "--color=red", "Hello\\nWorld"},
			expectedOptions: Options{
				Text:           "Hello\nWorld",
				Banner:         "standard",
				Color:          "red",
				Target:         "",
				ColorRequested: true,
			},
		},
		{
			name: "color mode with target substring",
			args: []string{"ascii-art", "--color=red", "ll", "Hello"},
			expectedOptions: Options{
				Text:           "Hello",
				Banner:         "standard",
				Color:          "red",
				Target:         "ll",
				ColorRequested: true,
			},
		},
		{
			name: "color mode with target flag and custom banner",
			args: []string{"ascii-art", "--color=red", "--target=ll", "Hello", "thinkertoy"},
			expectedOptions: Options{
				Text:           "Hello",
				Banner:         "thinkertoy",
				Color:          "red",
				Target:         "ll",
				ColorRequested: true,
			},
		},
		{
			name:          "color mode with 3 args but missing --target= prefix",
			args:          []string{"ascii-art", "--color=red", "ll", "Hello", "thinkertoy"},
			expectError:   true,
			expectedError: ColorUsage(),
		},
		{
			name:          "color mode too many arguments",
			args:          []string{"ascii-art", "--color=red", "--target=ll", "Hello", "thinkertoy", "extra"},
			expectError:   true,
			expectedError: ColorUsage(),
		},
		{
			name:          "color mode empty color",
			args:          []string{"ascii-art", "--color=", "Hello"},
			expectError:   true,
			expectedError: ColorUsage(),
		},
		{
			name:          "color mode empty target flag",
			args:          []string{"ascii-art", "--color=red", "--target=", "Hello", "shadow"},
			expectError:   true,
			expectedError: ColorUsage(),
		},

		// Output Mode
		{
			name:          "output mode missing text (len < 3)",
			args:          []string{"ascii-art", "--output=banner.txt"},
			expectError:   true,
			expectedError: OutputUsage(),
		},
		{
			name: "output mode valid single text",
			args: []string{"ascii-art", "--output=banner.txt", "Hello"},
			expectedOptions: Options{
				Text:            "Hello",
				Banner:          "standard",
				Color:           "",
				Target:          "",
				ColorRequested:  false,
				Output:          "banner.txt",
				OutputRequested: true,
			},
		},
		{
			name: "output mode valid text with newlines",
			args: []string{"ascii-art", "--output=banner.txt", "Hello\\nWorld"},
			expectedOptions: Options{
				Text:            "Hello\nWorld",
				Banner:          "standard",
				Color:           "",
				Target:          "",
				ColorRequested:  false,
				Output:          "banner.txt",
				OutputRequested: true,
			},
		},
		{
			name: "output mode valid with custom banner",
			args: []string{"ascii-art", "--output=banner.txt", "Hello", "shadow"},
			expectedOptions: Options{
				Text:            "Hello",
				Banner:          "shadow",
				Color:           "",
				Target:          "",
				ColorRequested:  false,
				Output:          "banner.txt",
				OutputRequested: true,
			},
		},
		{
			name: "output mode valid with standard banner",
			args: []string{"ascii-art", "--output=banner.txt", "Hello", "standard"},
			expectedOptions: Options{
				Text:            "Hello",
				Banner:          "standard",
				Color:           "",
				Target:          "",
				ColorRequested:  false,
				Output:          "banner.txt",
				OutputRequested: true,
			},
		},
		{
			name: "output mode valid text with newlines and banner",
			args: []string{"ascii-art", "--output=banner.txt", "Hello\\nWorld", "thinkertoy"},
			expectedOptions: Options{
				Text:            "Hello\nWorld",
				Banner:          "thinkertoy",
				Color:           "",
				Target:          "",
				ColorRequested:  false,
				Output:          "banner.txt",
				OutputRequested: true,
			},
		},
		{
			name:          "output mode invalid filename without extension",
			args:          []string{"ascii-art", "--output=banner", "Hello"},
			expectError:   true,
			expectedError: OutputUsage(),
		},
		{
			name:          "output mode invalid filename with wrong extension",
			args:          []string{"ascii-art", "--output=banner.pdf", "Hello"},
			expectError:   true,
			expectedError: OutputUsage(),
		},
		{
			name:          "output mode empty output flag",
			args:          []string{"ascii-art", "--output=", "Hello"},
			expectError:   true,
			expectedError: OutputUsage(),
		},
		{
			name:          "output mode too many arguments (len > 4)",
			args:          []string{"ascii-art", "--output=banner.txt", "Hello", "shadow", "extra"},
			expectError:   true,
			expectedError: OutputUsage(),
		},

		// Align Mode
		{
			name:          "align mode missing text (len < 3)",
			args:          []string{"ascii-art", "--align=center"},
			expectError:   true,
			expectedError: AlignUsage(),
		},
		{
			name: "align mode valid left align single text",
			args: []string{"ascii-art", "--align=left", "Hello"},
			expectedOptions: Options{
				Text:            "Hello",
				Banner:          "standard",
				Color:           "",
				Target:          "",
				ColorRequested:  false,
				Output:          "",
				OutputRequested: false,
				Align:           "left",
				AlignRequested:  true,
			},
		},
		{
			name: "align mode valid center align single text",
			args: []string{"ascii-art", "--align=center", "Hello"},
			expectedOptions: Options{
				Text:            "Hello",
				Banner:          "standard",
				Color:           "",
				Target:          "",
				ColorRequested:  false,
				Output:          "",
				OutputRequested: false,
				Align:           "center",
				AlignRequested:  true,
			},
		},
		{
			name: "align mode valid right align single text",
			args: []string{"ascii-art", "--align=right", "Hello"},
			expectedOptions: Options{
				Text:            "Hello",
				Banner:          "standard",
				Color:           "",
				Target:          "",
				ColorRequested:  false,
				Output:          "",
				OutputRequested: false,
				Align:           "right",
				AlignRequested:  true,
			},
		},
		{
			name: "align mode valid justify align single text",
			args: []string{"ascii-art", "--align=justify", "Hello"},
			expectedOptions: Options{
				Text:            "Hello",
				Banner:          "standard",
				Color:           "",
				Target:          "",
				ColorRequested:  false,
				Output:          "",
				OutputRequested: false,
				Align:           "justify",
				AlignRequested:  true,
			},
		},
		{
			name: "align mode valid text with newlines",
			args: []string{"ascii-art", "--align=center", "Hello\\nWorld"},
			expectedOptions: Options{
				Text:            "Hello\nWorld",
				Banner:          "standard",
				Color:           "",
				Target:          "",
				ColorRequested:  false,
				Output:          "",
				OutputRequested: false,
				Align:           "center",
				AlignRequested:  true,
			},
		},
		{
			name: "align mode valid with custom banner",
			args: []string{"ascii-art", "--align=right", "Hello", "shadow"},
			expectedOptions: Options{
				Text:            "Hello",
				Banner:          "shadow",
				Color:           "",
				Target:          "",
				ColorRequested:  false,
				Output:          "",
				OutputRequested: false,
				Align:           "right",
				AlignRequested:  true,
			},
		},
		{
			name: "align mode valid with thinkertoy banner",
			args: []string{"ascii-art", "--align=justify", "Hello\\nWorld", "thinkertoy"},
			expectedOptions: Options{
				Text:            "Hello\nWorld",
				Banner:          "thinkertoy",
				Color:           "",
				Target:          "",
				ColorRequested:  false,
				Output:          "",
				OutputRequested: false,
				Align:           "justify",
				AlignRequested:  true,
			},
		},
		{
			name:          "align mode invalid alignment type",
			args:          []string{"ascii-art", "--align=diagonal", "Hello"},
			expectError:   true,
			expectedError: AlignUsage(),
		},
		{
			name:          "align mode empty align flag",
			args:          []string{"ascii-art", "--align=", "Hello"},
			expectError:   true,
			expectedError: AlignUsage(),
		},
		{
			name:          "align mode too many arguments (len > 4)",
			args:          []string{"ascii-art", "--align=left", "Hello", "standard", "extra"},
			expectError:   true,
			expectedError: AlignUsage(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			options, err := ParseArgs(test.args)
			if test.expectError {
				if err == nil {
					t.Fatalf("for args %v: expected error, got nil", test.args)
				}
				if test.expectedError != "" && err.Error() != test.expectedError {
					t.Errorf("for args %v:\nexpected error:\n%q\ngot error:\n%q", test.args, test.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("for args %v: unexpected error: %v", test.args, err)
				}
				if options != test.expectedOptions {
					t.Errorf("for args %v:\nexpected options: %+v\ngot options:      %+v", test.args, test.expectedOptions, options)
				}
			}
		})
	}
}

func TestValidateOutputFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		// Valid filenames
		{name: "standard filename", filename: "banner.txt", expected: true},
		{name: "uppercase letters", filename: "BANNER.txt", expected: true},
		{name: "with numbers", filename: "banner123.txt", expected: true},
		{name: "with underscore", filename: "banner_output.txt", expected: true},
		{name: "with hyphen", filename: "banner-output.txt", expected: true},
		{name: "mixed characters", filename: "Test_File-01.txt", expected: true},
		{name: "single char name", filename: "a.txt", expected: true},

		// Invalid filenames
		{name: "empty filename", filename: "", expected: false},
		{name: "extension only", filename: ".txt", expected: false},
		{name: "missing extension", filename: "banner", expected: false},
		{name: "wrong extension", filename: "banner.pdf", expected: false},
		{name: "uppercase extension", filename: "banner.TXT", expected: false},
		{name: "contains spaces", filename: "banner output.txt", expected: false},
		{name: "contains forward slash", filename: "dir/banner.txt", expected: false},
		{name: "contains backward slash", filename: "dir\\banner.txt", expected: false},
		{name: "contains special characters", filename: "banner@test.txt", expected: false},
		{name: "double extension", filename: "banner.txt.txt", expected: false},
		{name: "hidden file with extension", filename: ".banner.txt", expected: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := validateOutputFilename(test.filename)
			if result != test.expected {
				t.Errorf("validateOutputFilename(%q) = %v; expected %v", test.filename, result, test.expected)
			}
		})
	}
}

func TestValidAlignment(t *testing.T) {
	tests := []struct {
		name     string
		align    string
		expected bool
	}{
		{name: "valid left", align: "left", expected: true},
		{name: "valid center", align: "center", expected: true},
		{name: "valid right", align: "right", expected: true},
		{name: "valid justify", align: "justify", expected: true},
		{name: "invalid empty", align: "", expected: false},
		{name: "invalid diagonal", align: "diagonal", expected: false},
		{name: "invalid uppercase", align: "LEFT", expected: false},
		{name: "invalid mixed case", align: "Center", expected: false},
		{name: "invalid space padded", align: " left ", expected: false},
		{name: "invalid random string", align: "top", expected: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := validAlignment(test.align)
			if result != test.expected {
				t.Errorf("validAlignment(%q) = %v; expected %v", test.align, result, test.expected)
			}
		})
	}
}
