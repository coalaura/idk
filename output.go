package main

import (
	"io"
	"strings"
	"unicode/utf8"
)

const (
	descriptionWidth        = 80
	descriptionIndent       = "  "
	descriptionContentWidth = descriptionWidth - len(descriptionIndent)
)

func (answer CommandAnswer) render(writer io.Writer) error {
	var output strings.Builder

	output.Grow(len(answer.Command) + len(answer.Description) + 6)

	output.WriteString("$ ")
	output.WriteString(answer.Command)
	output.WriteByte('\n')

	writeDescription(&output, answer.Description)

	_, err := io.WriteString(writer, output.String())
	return err
}

func writeDescription(output *strings.Builder, description string) {
	for line := range strings.SplitSeq(strings.TrimSpace(description), "\n") {
		var column int

		for word := range strings.FieldsSeq(line) {
			width := utf8.RuneCountInString(word)
			if column > 0 && column+1+width > descriptionWidth {
				output.WriteByte('\n')
				column = 0
			}

			if column == 0 {
				output.WriteString(descriptionIndent)
				column = len(descriptionIndent)
			} else {
				output.WriteByte(' ')
				column++
			}

			for width > descriptionContentWidth {
				end := 0

				for range descriptionContentWidth {
					_, size := utf8.DecodeRuneInString(word[end:])
					end += size
				}

				output.WriteString(word[:end])
				output.WriteByte('\n')
				output.WriteString(descriptionIndent)

				word = word[end:]

				width -= descriptionContentWidth
			}

			output.WriteString(word)

			column += width
		}

		output.WriteByte('\n')
	}
}
