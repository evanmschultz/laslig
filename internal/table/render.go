package table

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// Mode describes the render conditions used for one table render.
type Mode struct {
	Human    bool
	Width    int
	WrapMode WrapMode
}

// Styles contains the styles used for human table rendering.
type Styles struct {
	Header lipgloss.Style
	Rule   lipgloss.Style
	Even   lipgloss.Style
	Odd    lipgloss.Style
}

// WrapMode controls how table cells are compacted when width is constrained.
type WrapMode string

const (
	// WrapAuto wraps long cell values and rebalances column widths.
	WrapAuto WrapMode = "auto"
	// WrapNever truncates long values without wrapping.
	WrapNever WrapMode = "never"
	// WrapTruncate truncates long values with ellipsis.
	WrapTruncate WrapMode = "truncate"
)

func normalizeWrapMode(value WrapMode) WrapMode {
	switch value {
	case WrapNever, WrapTruncate:
		return value
	default:
		return WrapAuto
	}
}

// Render renders a terminal table for either human or plain output.
func Render(header []string, rows [][]string, mode Mode, styles Styles) string {
	allRows := make([][]string, 0, len(rows)+1)
	if len(header) > 0 {
		allRows = append(allRows, header)
	}
	allRows = append(allRows, rows...)

	widths := calcColumnWidths(allRows)
	widths = normalizeMinColumnWidths(widths)

	if mode.Human && mode.Width > 0 {
		widths = rebalanceWidths(widths, tableContentWidth(mode.Width))
	}

	if !mode.Human {
		return renderPlain(header, rows, widths)
	}

	if mode.Width > 0 {
		contentBudget := tableContentWidth(mode.Width)
		if contentBudget > 0 && sum(widths)+separatorWidth(len(widths)) > contentBudget {
			widths = rebalanceWidths(widths, contentBudget)
		}
	}

	lines := []string{}
	if len(header) > 0 {
		lines = append(lines, renderLineSet(header, widths, styles.Header, normalizeWrapMode(mode.WrapMode), styles.Rule)...)
		lines = append(lines, renderRule(widths, styles.Rule))
	}

	for index, row := range rows {
		rowStyle := styles.Even
		if index%2 == 1 {
			rowStyle = styles.Odd
		}
		lines = append(lines, renderLineSet(row, widths, rowStyle, normalizeWrapMode(mode.WrapMode), styles.Rule)...)
	}

	style := tableStyle()
	return style.Render(strings.Join(lines, "\n"))
}

func renderPlain(header []string, rows [][]string, widths []int) string {
	joinLine := func(row []string) string {
		cells := make([]string, 0, len(widths))
		for index, width := range widths {
			value := ""
			if index < len(row) {
				value = row[index]
			}
			if index < len(widths)-1 {
				value = fmt.Sprintf("%-*s", width, value)
			}
			cells = append(cells, value)
		}
		return strings.Join(cells, " | ")
	}

	lines := []string{}
	if len(header) > 0 {
		lines = append(lines, joinLine(header))
		rule := make([]string, len(widths))
		for index, width := range widths {
			rule[index] = strings.Repeat("-", width)
		}
		lines = append(lines, strings.Join(rule, "-+-"))
	}

	for _, row := range rows {
		lines = append(lines, joinLine(row))
	}
	return strings.Join(lines, "\n")
}

func renderRule(widths []int, ruleStyle lipgloss.Style) string {
	if len(widths) == 0 {
		return ""
	}
	ruleParts := make([]string, len(widths))
	for index, width := range widths {
		ruleParts[index] = strings.Repeat("─", width)
		ruleParts[index] = ruleStyle.Render(ruleParts[index])
	}
	return strings.Join(ruleParts, ruleStyle.Render("─┼─"))
}

func renderLineSet(row []string, widths []int, rowStyle lipgloss.Style, wrapMode WrapMode, ruleStyle lipgloss.Style) []string {
	cellParts := make([][]string, len(widths))
	maxHeight := 1
	for index, width := range widths {
		value := ""
		if index < len(row) {
			value = row[index]
		}
		parts := wrapCell(value, width, wrapMode)
		if len(parts) == 0 {
			parts = []string{""}
		}
		cellParts[index] = parts
		if len(parts) > maxHeight {
			maxHeight = len(parts)
		}
	}

	separator := ruleStyle.Render(" │ ")
	lines := make([]string, 0, maxHeight)
	for index := 0; index < maxHeight; index++ {
		cells := make([]string, len(widths))
		for cellIndex, width := range widths {
			value := ""
			if index < len(cellParts[cellIndex]) {
				value = cellParts[cellIndex][index]
			}
			cells[cellIndex] = rowStyle.Render(lipgloss.NewStyle().Width(width).Render(value))
		}
		lines = append(lines, strings.Join(cells, separator))
	}
	return lines
}

func wrapCell(value string, width int, mode WrapMode) []string {
	switch normalizeWrapMode(mode) {
	case WrapNever, WrapTruncate:
		if width <= 0 {
			return []string{value}
		}
		return []string{truncateVisible(value, width)}
	default:
		return wrapParagraph(value, width)
	}
}

func wrapParagraph(value string, width int) []string {
	if width <= 0 || lipgloss.Width(value) <= width {
		return []string{value}
	}

	paragraphs := strings.Split(value, "\n")
	lines := make([]string, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		lines = append(lines, wrapSingleParagraph(paragraph, width)...)
	}
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

func wrapSingleParagraph(value string, width int) []string {
	if width <= 0 {
		return []string{value}
	}
	if strings.TrimSpace(value) == "" {
		return []string{""}
	}
	words := strings.Fields(value)
	if len(words) == 0 {
		return []string{value}
	}

	lines := []string{}
	current := ""
	for _, word := range words {
		chunks := splitWideToken(word, width)
		for _, chunk := range chunks {
			if candidate, ok := candidateLine(current, chunk, width); ok {
				current = candidate
				continue
			}
			if current != "" {
				lines = append(lines, truncateVisible(current, width))
			}
			current = chunk
		}
	}
	if current != "" {
		lines = append(lines, truncateVisible(current, width))
	}
	if len(lines) == 0 {
		lines = []string{""}
	}
	return lines
}

func candidateLine(current, chunk string, width int) (string, bool) {
	if current == "" {
		return chunk, true
	}
	candidate := current + " " + chunk
	if lipgloss.Width(candidate) <= width {
		return candidate, true
	}
	return "", false
}

func splitWideToken(value string, width int) []string {
	if width <= 0 || lipgloss.Width(value) <= width {
		return []string{value}
	}

	parts := []string{}
	current := strings.Builder{}
	currentWidth := 0
	for _, runeValue := range value {
		segment := string(runeValue)
		runeWidth := lipgloss.Width(segment)

		if runeWidth == 0 {
			current.WriteRune(runeValue)
			continue
		}

		if currentWidth > 0 && currentWidth+runeWidth > width {
			parts = append(parts, current.String())
			current.Reset()
			currentWidth = 0
		}
		// If a single rune is wider than width, keep it as-is; final truncation
		// will preserve layout stability without panicking the cell budget.
		current.WriteRune(runeValue)
		currentWidth += runeWidth
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	if len(parts) == 0 {
		return []string{value}
	}
	return parts
}

func truncateVisible(value string, width int) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(value) <= width {
		return value
	}

	const ellipsis = "…"
	if width == 1 {
		return ellipsis
	}

	var builder strings.Builder
	for _, runeValue := range value {
		candidate := builder.String() + string(runeValue)
		if lipgloss.Width(candidate+ellipsis) > width {
			break
		}
		builder.WriteRune(runeValue)
	}
	return strings.TrimRight(builder.String(), " ") + ellipsis
}

func calcColumnWidths(rows [][]string) []int {
	widths := make([]int, 0)
	for _, row := range rows {
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
	return widths
}

func normalizeMinColumnWidths(widths []int) []int {
	normalized := make([]int, len(widths))
	for index, width := range widths {
		if width <= 0 {
			width = 1
		}
		normalized[index] = width
	}
	return normalized
}

func rebalanceWidths(widths []int, budget int) []int {
	if len(widths) == 0 || budget <= 0 {
		return widths
	}

	current := sum(widths) + separatorWidth(len(widths))
	target := budget
	excess := current - target
	if excess <= 0 {
		return widths
	}

	adjusted := append([]int(nil), widths...)
	for excess > 0 {
		index := maxWidthIndex(adjusted)
		if index < 0 || adjusted[index] <= 1 {
			break
		}
		adjusted[index]--
		excess--
	}
	return adjusted
}

func maxWidthIndex(widths []int) int {
	index := -1
	maxWidth := -1
	for i, width := range widths {
		if width <= 1 {
			continue
		}
		if width > maxWidth {
			maxWidth = width
			index = i
		}
	}
	return index
}

func separatorWidth(columns int) int {
	if columns <= 1 {
		return 0
	}
	return 3 * (columns - 1)
}

func sum(values []int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}

func tableContentWidth(totalWidth int) int {
	style := tableStyle()
	frameWidth, _ := style.GetFrameSize()
	contentWidth := totalWidth - frameWidth
	if contentWidth <= 0 {
		return 0
	}
	return contentWidth
}

func tableStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1)
}
