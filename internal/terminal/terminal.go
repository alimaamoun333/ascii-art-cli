package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var (
	currentOS = runtime.GOOS
	cmdRunner = func(name string, arg ...string) ([]byte, error) {
		cmd := exec.Command(name, arg...)
		if currentOS == "linux" || currentOS == "darwin" {
			cmd.Stdin = os.Stdin
		}
		return cmd.Output()
	}
)

func parseWindowsOutput(output []byte) (int, error) {
	number := strings.TrimSpace(string(output))
	width, err := strconv.Atoi(number)
	if err != nil {
		return -1, fmt.Errorf("failed to parse Windows terminal width: %w", err)
	}
	if width <= 0 {
		return -1, fmt.Errorf("invalid terminal width: %d", width)
	}
	return width, nil
}

func parseSttyOutput(output []byte) (int, error) {
	parts := strings.Fields(string(output))
	if len(parts) != 2 {
		return -1, fmt.Errorf("unexpected output from stty size: %s", output)
	}
	width, err := strconv.Atoi(parts[1])
	if err != nil {
		return -1, fmt.Errorf("failed to parse stty terminal width: %w", err)
	}
	if width <= 0 {
		return -1, fmt.Errorf("invalid terminal width: %d", width)
	}
	return width, nil
}

func Width() (int, error) {
	switch currentOS {
	case "windows":
		// Windows
		output, err := cmdRunner("powershell", "-command", "$Host.UI.RawUI.WindowSize.Width")
		if err != nil {
			return -1, fmt.Errorf("failed to get Windows terminal width: %w", err)
		}
		return parseWindowsOutput(output)

	case "linux", "darwin":
		// Linux or macOS
		output, err := cmdRunner("stty", "size")
		if err != nil {
			return -1, fmt.Errorf("failed to get terminal width using stty: %w", err)
		}
		return parseSttyOutput(output)

	default:
		return -1, fmt.Errorf("unsupported OS: %s", currentOS)
	}
}
