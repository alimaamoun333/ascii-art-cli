package main

import (
	"ascii-art-justify/internal/alignment"
	"ascii-art-justify/internal/banner"
	"ascii-art-justify/internal/cli"
	"ascii-art-justify/internal/layout"
	"ascii-art-justify/internal/output"
	"ascii-art-justify/internal/renderer"
	"ascii-art-justify/internal/style"
	"ascii-art-justify/internal/terminal"
	"fmt"
	"strings"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)

	}
}

func run() error {
	rows := []string{}
	args := cli.ReadArgs()

	options, err := cli.ParseArgs(args)
	if err != nil {
		return err
	}
	path, err := banner.Path(options.Banner)
	if err != nil {
		return err
	}
	shapes, err := banner.Load(path)
	if err != nil {
		return err
	}
	if options.AlignRequested {
		// align mode
		spaceBlock, err := renderer.RenderLine(shapes, " ")
		if err != nil {
			return err
		}

		width, err := terminal.Width()
		if err != nil {
			return err
		}

		lines := strings.Split(options.Text, "\n")

		for _, line := range lines {
			lineLayout := layout.ParseLine(line)

			wordBlocks := [][]string{}

			for _, word := range lineLayout.Words {
				wordRows, err := renderer.RenderLine(shapes, word)
				if err != nil {
					return err
				}

				wordBlocks = append(wordBlocks, wordRows)
			}

			content, err := alignment.NewContent(
				spaceBlock,
				wordBlocks,
				lineLayout.Leading,
				lineLayout.Gaps,
				lineLayout.Trailing,
			)
			if err != nil {
				return err
			}
			mode := alignment.Mode(options.Align)
			alignedRows, err := alignment.Align(content, width, mode)
			if err != nil {
				return err
			}
			rows = append(rows, alignedRows...)
		}
	} else {
		rows, err = renderer.Render(shapes, options.Text)
		if err != nil {
			return err
		}
		if options.ColorRequested {
			boundaries, err := layout.CalculateMultilineBoundaries(
				shapes,
				options.Text,
			)
			if err != nil {
				return err
			}

			styledRows, err := style.ApplyColorMultiline(
				rows,
				options.Text,
				options.Target,
				options.Color,
				boundaries,
			)
			if err != nil {
				return err
			}

			rows = styledRows
		}
	}

	// normal mode
	text := strings.Join(rows, "\n") + "\n"
	if options.OutputRequested {
		// write to file ( output mode )
		if err := output.Write(options.Output, text); err != nil {
			return err
		}
		return nil
	}
	fmt.Print(text)

	return nil
}
