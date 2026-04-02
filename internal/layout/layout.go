package layout

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// WrapText wraps each paragraph in one string to the requested visible width.
func WrapText(value string, width int) string {
	if width <= 0 {
		return value
	}
	paragraphs := strings.Split(value, "\n")
	for index, paragraph := range paragraphs {
		paragraphs[index] = wrapWords(paragraph, width)
	}
	return strings.Join(paragraphs, "\n")
}

// IndentBlock prefixes each line in one multi-line string with the provided prefix.
func IndentBlock(prefix string, value string) string {
	lines := strings.Split(value, "\n")
	for index, line := range lines {
		lines[index] = prefix + line
	}
	return strings.Join(lines, "\n")
}

// wrapWords wraps one single-line paragraph to the requested visible width.
func wrapWords(value string, width int) string {
	if width <= 0 || lipgloss.Width(value) <= width {
		return value
	}

	words := strings.Fields(value)
	if len(words) == 0 {
		return ""
	}

	lines := make([]string, 0, len(words))
	for _, chunk := range splitWideToken(words[0], width) {
		lines = append(lines, chunk)
	}

	for _, word := range words[1:] {
		current := lines[len(lines)-1]
		candidate := current + " " + word
		if lipgloss.Width(candidate) <= width {
			lines[len(lines)-1] = candidate
			continue
		}
		for _, chunk := range splitWideToken(word, width) {
			lines = append(lines, chunk)
		}
	}
	return strings.Join(lines, "\n")
}

func splitWideToken(value string, width int) []string {
	if width <= 0 || lipgloss.Width(value) <= width {
		return []string{value}
	}

	parts := []string{}
	current := strings.Builder{}
	currentWidth := 0
	for _, r := range value {
		segment := string(r)
		segmentWidth := lipgloss.Width(segment)

		if segmentWidth == 0 {
			current.WriteRune(r)
			continue
		}

		if currentWidth > 0 && currentWidth+segmentWidth > width {
			parts = append(parts, current.String())
			current.Reset()
			currentWidth = 0
		}
		current.WriteRune(r)
		currentWidth += segmentWidth
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	if len(parts) == 0 {
		parts = []string{value}
	}
	return parts
}
