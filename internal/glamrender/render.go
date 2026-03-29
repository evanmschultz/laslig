package glamrender

import (
	"fmt"
	"strings"

	glamour "charm.land/glamour/v2"
	"charm.land/glamour/v2/styles"
)

// Render renders Markdown using one supported built-in Glamour style with
// optional word wrapping.
func Render(markdown string, width int, style string) (string, error) {
	options := []glamour.TermRendererOption{
		glamour.WithStandardStyle(resolveStandardStyle(style)),
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

func resolveStandardStyle(style string) string {
	switch style {
	case styles.DarkStyle, styles.LightStyle, styles.PinkStyle, styles.DraculaStyle, styles.TokyoNightStyle, styles.AsciiStyle, styles.NoTTYStyle:
		return style
	default:
		return styles.DraculaStyle
	}
}

// FencedCodeBlock returns one Markdown fenced code block string.
func FencedCodeBlock(language string, body string) string {
	trimmed := strings.TrimRight(body, "\n")
	return "```" + language + "\n" + trimmed + "\n```"
}
