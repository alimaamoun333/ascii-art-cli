package cli

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

type Options struct {
	Text            string
	Banner          string
	Color           string
	Target          string
	ColorRequested  bool
	Output          string
	OutputRequested bool
	Align           string
	AlignRequested  bool
}

func ReadArgs() []string {
	return os.Args
}

func ParseArgs(args []string) (Options, error) {
	options := Options{
		Text:            "",
		Banner:          "standard",
		Color:           "",
		Target:          "",
		ColorRequested:  false,
		OutputRequested: false,
		Align:           "",
		AlignRequested:  false,
	}

	if len(args) <= 1 {
		return Options{}, errors.New(Help())
	}

	// Color mode
	if strings.HasPrefix(args[1], "--color=") {
		options.Color = strings.TrimPrefix(args[1], "--color=")
		if options.Color == "" {
			return Options{}, errors.New(ColorUsage())
		}
		options.ColorRequested = true

		remaining := len(args) - 2

		switch remaining {
		case 0:
			return Options{}, errors.New(ColorUsage())

		case 1:
			options.Text = strings.ReplaceAll(args[2], "\\n", "\n")
			return options, nil

		case 2:
			options.Target = args[2]
			options.Text = strings.ReplaceAll(args[3], "\\n", "\n")
			return options, nil

		case 3:
			if !strings.HasPrefix(args[2], "--target=") {
				return Options{}, errors.New(ColorUsage())
			}
			options.Target = strings.TrimPrefix(args[2], "--target=")
			if options.Target == "" {
				return Options{}, errors.New(ColorUsage())
			}
			options.Text = strings.ReplaceAll(args[3], "\\n", "\n")
			options.Banner = args[4]

			return options, nil

		default:
			return Options{}, errors.New(ColorUsage())
		}
	}

	if strings.HasPrefix(args[1], "--output=") {
		// output mode
		options.Output = strings.TrimPrefix(args[1], "--output=")
		options.OutputRequested = true

		if !validateOutputFilename(options.Output) {
			return Options{}, errors.New(OutputUsage())
		}
		if len(args) < 3 {
			return Options{}, errors.New(OutputUsage())
		}
		if len(args) > 4 {
			return Options{}, errors.New(OutputUsage())
		}

		options.Text = strings.ReplaceAll(args[2], "\\n", "\n")
		if len(args) == 4 {
			options.Banner = args[3]
		}
		return options, nil

	}
	if strings.HasPrefix(args[1], "--align=") {
		// align mode
		options.Align = strings.TrimPrefix(args[1], "--align=")
		options.AlignRequested = true
		if !validAlignment(options.Align) {
			options.AlignRequested = false
			return Options{}, errors.New(AlignUsage())
		}
		if len(args) < 3 || len(args) > 4 {
			return Options{}, errors.New(AlignUsage())
		}
		if len(args) < 3 {
			return Options{}, errors.New(AlignUsage())
		}
		options.Text = strings.ReplaceAll(args[2], "\\n", "\n")
		if len(args) == 4 {
			options.Banner = args[3]
		}
		return options, nil
	}
	// Normal mode
	switch len(args) {
	case 2:
		options.Text = strings.ReplaceAll(args[1], "\\n", "\n")
		return options, nil

	case 3:
		options.Text = strings.ReplaceAll(args[1], "\\n", "\n")
		options.Banner = args[2]
		return options, nil

	default:
		return Options{}, errors.New(NormalUsage())
	}
}

var outputFilenamePattern = regexp.MustCompile(
	`^[A-Za-z0-9_-]+\.txt$`,
)

func validateOutputFilename(filename string) bool {
	return outputFilenamePattern.MatchString(filename)
}

func validAlignment(value string) bool {
	validAlignments := []string{"left", "center", "right", "justify"}
	for _, alignment := range validAlignments {
		if alignment == value {
			return true
		}
	}
	return false
}

func NormalUsage() string {
	return `Usage:
  ascii-art "text"
  ascii-art "text" <banner>`
}

func ColorUsage() string {
	return `Usage:
  ascii-art --color=<color> "text"
  ascii-art --color=<color> <target> "text"
  ascii-art --color=<color> --target=<target> "text" <banner>`
}

func OutputUsage() string {
	return `Usage:	
  ascii-art --output=<fileName.txt> "text"
  ascii-art --output=<fileName.txt> "text" <banner>`
}

func AlignUsage() string {
	return `Usage:
  ascii-art --align=<type> "text"
  ascii-art --align=<type> "text" <banner>`
}

func Help() string {
	return `Usage:
  ascii-art "text"
  ascii-art "text" <banner>
  ascii-art --color=<color> "text"
  ascii-art --color=<color> <target> "text"
  ascii-art --color=<color> --target=<target> "text" <banner>
  ascii-art --output=<fileName.txt> "text"
  ascii-art --output=<fileName.txt> "text" <banner>`
}
