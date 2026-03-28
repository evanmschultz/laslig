package glamrender

import (
	"fmt"
	"strings"

	glamour "charm.land/glamour/v2"
	"charm.land/glamour/v2/styles"
)

// Render renders Markdown using Glamour's standard dark style with optional word wrapping.
func Render(markdown string, width int) (string, error) {
	options := []glamour.TermRendererOption{
		glamour.WithStandardStyle(styles.DarkStyle),
	}
	if width > 0 {
		options = append(options, glamour.WithWordWrap(width))
	}

	renderer, err := glamour.NewTermRenderer(options...)
	if err != nil {
		return "", fmt.Errorf("create glamour renderer: %w", err)
	}

	rendered, err := renderer.Render(markdown)
	if err != nil {
		return "", fmt.Errorf("render glamour markdown: %w", err)
	}
	return strings.Trim(rendered, "\n"), nil
}

// FencedCodeBlock returns one Markdown fenced code block string.
func FencedCodeBlock(language string, body string) string {
	trimmed := strings.TrimRight(body, "\n")
	return "```" + language + "\n" + trimmed + "\n```"
}
