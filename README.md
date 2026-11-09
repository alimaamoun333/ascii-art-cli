# ASCII Art Generator & Justify

A robust, high-performance Command-Line Interface (CLI) application written in Go that transforms standard text into stylized ASCII art banners. Featuring dynamic terminal-width alignment (left, center, right, justify), ANSI color styling with substring targeting, multiple banner templates, multiline input support, and file export functionality.

---

## Table of Contents

- [Features](#features)
- [Project Architecture](#project-architecture)
- [Prerequisites & Installation](#prerequisites--installation)
- [Usage Guide](#usage-guide)
  - [1. Normal Mode](#1-normal-mode)
  - [2. Terminal Alignment & Justification](#2-terminal-alignment--justification)
  - [3. ANSI Color Styling & Target Highlighting](#3-ansi-color-styling--target-highlighting)
  - [4. File Output Export](#4-file-output-export)
  - [5. Multi-line Text](#5-multi-line-text)
- [Reference Tables](#reference-tables)
  - [Supported Banners](#supported-banners)
  - [Supported Colors](#supported-colors)
  - [Alignment Modes](#alignment-modes)
- [CLI Syntax & Options](#cli-syntax--options)
- [Testing & Quality Assurance](#testing--quality-assurance)

---

## Features

- **Terminal Width Alignment & Justification**: Dynamically calculates terminal columns on Windows, Linux, and macOS to center, left-align, right-align, or justify ASCII art blocks across the screen.
- **Multiple Banner Fonts**: Full support for `standard` (default), `shadow`, and `thinkertoy` banner fonts.
- **ANSI Color Styling**: Apply ANSI colors to entire text outputs or highlight specific target substrings with automated range merging.
- **File Export**: Direct ASCII art output to text files with strict filename validation (`[A-Za-z0-9_-]+\.txt`).
- **Multi-line Support**: Accurately splits and renders multiline strings (`\n`), preserving empty lines and line boundaries.
- **Comprehensive Validation**: Validates CLI arguments, flags, font files, and character sets, displaying contextual usage guidance on errors.
- **Modular & Decoupled Architecture**: Clean Go package design with 100% test coverage across all internal modules.

---

## Project Architecture

```
ascii-art-justify/
├── main.go               # Application entry point & orchestration
├── main_test.go          # End-to-end unit and CLI integration tests
├── go.mod                # Go module definition
├── banners/              # Banner font templates
│   ├── standard.txt      # Standard ASCII block font
│   ├── shadow.txt        # Shadow stylized font
│   └── thinkertoy.txt    # Thinkertoy line font
└── internal/             # Internal modular packages
    ├── alignment/        # Terminal alignment logic (left, center, right, justify)
    ├── banner/           # Banner locator and parser
    ├── cli/              # Command line argument parser & usage documentation
    ├── layout/           # Character and multiline rendered boundary calculator
    ├── output/           # Output file writer & persistence engine
    ├── renderer/         # ASCII art glyph and line rendering engine
    ├── style/            # ANSI color styling, substring matcher & range merger
    └── terminal/         # Cross-platform terminal dimension detection
```

---

## Prerequisites & Installation

### Prerequisites

- [Go](https://go.dev/doc/install) (version 1.20 or higher recommended)

### Clone & Build

```bash
git clone <repository-url>
cd ascii-art-cli
```

Verify the installation by running:

```bash
go run . "Hello World"
```

---

## Usage Guide

### 1. Normal Mode

Generate ASCII art directly in the terminal using the default `standard` banner:

```bash
go run . "Hello"
```

Specify a custom banner style (`shadow` or `thinkertoy`):

```bash
go run . "Hello" shadow
go run . "Hello" thinkertoy
```

---

### 2. Terminal Alignment & Justification

Align ASCII art horizontally based on your terminal's current width using `--align=<mode>`:

```bash
# Center alignment
go run . --align=center "Hello"

# Right alignment
go run . --align=right "Hello"

# Left alignment (padded to terminal width)
go run . --align=left "Hello"

# Justify alignment (evenly distributes spacing between words)
go run . --align=justify "Hello World"

# Alignment with custom banner font
go run . --align=center "Hello" shadow
go run . --align=justify "Clean Code" thinkertoy
```

---

### 3. ANSI Color Styling & Target Highlighting

#### Full Text Coloring
Apply color across all characters in the output:

```bash
go run . --color=red "Hello"
go run . --color=green "Hello World" shadow
```

#### Targeted Substring Coloring
Highlight specific characters or words within the text:

```bash
# Color all occurrences of "ell" in red
go run . --color=red "ell" "Hello"

# Color specific substring with explicit flag and custom banner font
go run . --color=blue --target=World "Hello World" shadow
go run . --color=yellow --target=art "ascii-art" thinkertoy
```

---

### 4. File Output Export

Save generated ASCII art directly into a text file using `--output=<fileName.txt>`:

```bash
# Save to banner.txt using standard font
go run . --output=banner.txt "Hello"

# Save with a custom banner style
go run . --output=result.txt "Hello World" shadow
go run . --output=output_test.txt "Multiline\nText" thinkertoy
```

> **Note**: Output filenames must follow the pattern `[A-Za-z0-9_-]+\.txt` (letters, numbers, hyphens, and underscores ending with `.txt`).

---

### 5. Multi-line Text

Render multi-line text by including the newline sequence `\n`:

```bash
# Normal multiline
go run . "Hello\nWorld"

# Multiline with alignment
go run . --align=center "Welcome\nTo\nASCII Art"

# Multiline with color
go run . --color=cyan "Line One\nLine Two"

# Multiline to file
go run . --output=multiline.txt "Hello\nWorld" shadow
```

---

## Reference Tables

### Supported Banners

| Banner Font | Description | Default |
| :--- | :--- | :---: |
| `standard` | Classic block-style ASCII characters | Yes |
| `shadow` | Characters with drop-shadow effect | No |
| `thinkertoy` | Minimalist line art characters | No |

### Supported Colors

| Color Name | ANSI Escape Sequence | Preview |
| :--- | :--- | :--- |
| `red` | `\033[31m` | Red |
| `green` | `\033[32m` | Green |
| `yellow` | `\033[33m` | Yellow |
| `blue` | `\033[34m` | Blue |
| `magenta` | `\033[35m` | Magenta |
| `cyan` | `\033[36m` | Cyan |

### Alignment Modes

| Alignment | Description |
| :--- | :--- |
| `center` | Centers the rendered ASCII art block within the terminal width |
| `left` | Left-aligns the ASCII art block with trailing space padding |
| `right` | Right-aligns the ASCII art block with leading space padding |
| `justify` | Distributes space gaps evenly between individual words across the line |

---

## CLI Syntax & Options

```text
Usage:
  ascii-art "text"
  ascii-art "text" <banner>
  ascii-art --align=<type> "text"
  ascii-art --align=<type> "text" <banner>
  ascii-art --color=<color> "text"
  ascii-art --color=<color> <target> "text"
  ascii-art --color=<color> --target=<target> "text" <banner>
  ascii-art --output=<fileName.txt> "text"
  ascii-art --output=<fileName.txt> "text" <banner>
```

---

## Testing & Quality Assurance

Run the comprehensive unit test suite across all packages:

```bash
# Run all unit tests
go test ./...

# Run all unit tests with verbose output
go test -v ./...
```
