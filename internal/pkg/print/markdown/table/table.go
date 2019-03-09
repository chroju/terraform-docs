package table

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/segmentio/terraform-docs/internal/pkg/doc"
	"github.com/segmentio/terraform-docs/internal/pkg/print"
	"github.com/segmentio/terraform-docs/internal/pkg/print/markdown"
	"github.com/segmentio/terraform-docs/internal/pkg/settings"
)

// Print prints a document as Markdown tables.
func Print(document *doc.Doc, settings settings.Settings) (string, error) {
	var buffer bytes.Buffer

	if document.HasComment() {
		printComment(&buffer, document.Comment, settings)
	}

	if document.HasInputs() {
		if settings.Has(print.WithSortByName) {
			if settings.Has(print.WithSortInputsByRequired) {
				doc.SortInputsByRequired(document.Inputs)
			} else {
				doc.SortInputsByName(document.Inputs)
			}
		}

		printInputs(&buffer, document.Inputs, settings)
	}

	if document.HasOutputs() {
		if settings.Has(print.WithSortByName) {
			doc.SortOutputsByName(document.Outputs)
		}

		if document.HasInputs() {
			buffer.WriteString("\n")
		}

		printOutputs(&buffer, document.Outputs, settings)
	}

	return markdown.Sanitize(buffer.String()), nil
}

func getInputDefaultValue(input *doc.Input, settings settings.Settings) string {
	var result = "n/a"

	if input.HasDefault() {
		result = fmt.Sprintf("`%s`", print.GetPrintableValue(input.Default, settings, false))
	}

	return result
}

func printComment(buffer *bytes.Buffer, comment string, settings settings.Settings) {
	buffer.WriteString(fmt.Sprintf("%s\n", comment))
}

func printInputs(buffer *bytes.Buffer, inputs []doc.Input, settings settings.Settings) {
	buffer.WriteString("## Inputs\n\n")

	table := tablewriter.NewWriter(buffer)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAutoFormatHeaders(false)
	table.SetCenterSeparator("|")

	header := []string{"Name", "Description", "Type", "Default"}
	if settings.Has(print.WithRequired) {
		header = append(header, "Required")
	}
	table.SetHeader(header)

	for _, input := range inputs {
		raw := []string{
			strings.Replace(input.Name, "_", "\\_", -1),
			markdown.ConvertMultiLineText(input.Description),
			input.Type,
			getInputDefaultValue(&input, settings),
		}

		if settings.Has(print.WithRequired) {
			raw = append(raw, printIsInputRequired(&input))
		}

		table.Append(raw)
	}

	table.Render()
}

func printIsInputRequired(input *doc.Input) string {
	if input.IsRequired() {
		return "yes"
	}

	return "no"
}

func printOutputs(buffer *bytes.Buffer, outputs []doc.Output, settings settings.Settings) {
	buffer.WriteString("## Outputs\n\n")

	table := tablewriter.NewWriter(buffer)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAutoFormatHeaders(false)
	table.SetCenterSeparator("|")
	table.SetHeader([]string{"Name", "Description"})

	for _, output := range outputs {
		table.Append([]string{
			strings.Replace(output.Name, "_", "\\_", -1),
			markdown.ConvertMultiLineText(output.Description),
		})
	}

	table.Render()
}
