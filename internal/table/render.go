package table

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// Mode describes the render conditions used for one table render.
type Mode struct {
	Human bool
	Width int
}

// Styles contains the styles used for human table rendering.
type Styles struct {
	Header lipgloss.Style
	Rule   lipgloss.Style
	Even   lipgloss.Style
	Odd    lipgloss.Style
}

// Render renders a terminal table for either human or plain output.
func Render(header []string, rows [][]string, mode Mode, styles Styles) string {
	allRows := make([][]string, 0, len(rows)+1)
	if len(header) > 0 {
		allRows = append(allRows, header)
	}
	allRows = append(allRows, rows...)

	widths := make([]int, 0)
	for _, row := range allRows {
		for index, cell := range row {
			width := lipgloss.Width(cell)
			if len(widths) <= index {
				widths = append(widths, width)
				continue
			}
			if width > widths[index] {
				widths[index] = width
			}
		}
	}

	joinRow := func(row []string, style lipgloss.Style) string {
		cells := make([]string, 0, len(widths))
		for index, width := range widths {
			value := ""
			if index < len(row) {
				value = row[index]
			}
			cell := value
			if mode.Human {
				cell = lipgloss.NewStyle().Width(width).Render(value)
				cell = style.Render(cell)
			} else if index < len(widths)-1 {
				cell = fmt.Sprintf("%-*s", width, value)
			}
			cells = append(cells, cell)
		}
		separator := " | "
		if mode.Human {
			separator = styles.Rule.Render(" │ ")
		}
		return strings.Join(cells, separator)
	}

	lines := []string{}
	if len(header) > 0 {
		lines = append(lines, joinRow(header, styles.Header))
		ruleParts := make([]string, 0, len(widths))
		for _, width := range widths {
			ruleParts = append(ruleParts, strings.Repeat("─", width))
		}
		ruleSeparator := "─┼─"
		if !mode.Human {
			ruleSeparator = "-+-"
			for index, width := range widths {
				ruleParts[index] = strings.Repeat("-", width)
			}
		} else {
			ruleSeparator = styles.Rule.Render("─┼─")
			for index, width := range widths {
				ruleParts[index] = styles.Rule.Render(strings.Repeat("─", width))
			}
		}
		lines = append(lines, strings.Join(ruleParts, ruleSeparator))
	}

	for index, row := range rows {
		style := styles.Even
		if index%2 == 1 {
			style = styles.Odd
		}
		lines = append(lines, joinRow(row, style))
	}

	rendered := strings.Join(lines, "\n")
	if !mode.Human {
		return rendered
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1)
	if mode.Width > 0 {
		maxWidth := mode.Width - 4
		if maxWidth > 0 {
			style = style.MaxWidth(maxWidth)
		}
	}
	return style.Render(rendered)
}
