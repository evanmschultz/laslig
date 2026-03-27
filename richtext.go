package laslig

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/glamour"
)

// CodeBlock writes one titled code-style block.
func (p *Printer) CodeBlock(block CodeBlock) error {
	if p.mode.Format == FormatJSON {
		return p.writeJSON("code_block", block)
	}

	body := strings.TrimRight(block.Body, "\n")
	if p.mode.Format == FormatHuman && p.mode.Styled {
		rendered, err := p.renderStyledCodeBlock(block)
		if err != nil {
			return fmt.Errorf("render code block: %w", err)
		}
		body = rendered
	}

	if err := p.writeFramedBlock("code block", block.Title, body, block.Footer); err != nil {
		return fmt.Errorf("write code block: %w", err)
	}
	return nil
}

// LogBlock writes one titled boxed transcript or log excerpt.
func (p *Printer) LogBlock(block LogBlock) error {
	if p.mode.Format == FormatJSON {
		return p.writeJSON("log_block", block)
	}

	if err := p.writeFramedBlock("log block", block.Title, strings.TrimRight(block.Body, "\n"), block.Footer); err != nil {
		return fmt.Errorf("write log block: %w", err)
	}
	return nil
}

// writeFramedBlock writes one titled block using plain layout or a styled frame depending on mode.
func (p *Printer) writeFramedBlock(kind string, title string, body string, footer string) error {
	lines := []string{}
	if trimmed := strings.TrimSpace(title); trimmed != "" {
		lines = append(lines, p.renderHeading(p.wrapText(trimmed, p.maxTextWidth())))
	}
	if strings.TrimSpace(body) != "" {
		lines = append(lines, body)
	}
	if trimmed := strings.TrimSpace(footer); trimmed != "" {
		rendered := p.wrapText(trimmed, p.maxTextWidth())
		if p.mode.Format == FormatHuman {
			rendered = p.theme.Muted.Render(rendered)
		}
		lines = append(lines, rendered)
	}

	content := strings.Join(lines, "\n\n")
	if p.mode.Format == FormatHuman && p.mode.Styled {
		content = p.renderFramedContent(content)
	}
	if _, err := fmt.Fprintln(p.out, content); err != nil {
		return fmt.Errorf("write %s content: %w", kind, err)
	}
	return nil
}

// renderFramedContent applies a neutral border and padding around one block body.
func (p *Printer) renderFramedContent(content string) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2)
	if maxWidth := p.maxPanelWidth(); maxWidth > 0 {
		style = style.MaxWidth(maxWidth)
	}
	return style.Render(content)
}

// renderStyledCodeBlock renders one fenced code block through Glamour for ANSI output.
func (p *Printer) renderStyledCodeBlock(block CodeBlock) (string, error) {
	options := []glamour.TermRendererOption{
		glamour.WithStandardStyle("dark"),
	}
	if width := p.maxCodeWidth(); width > 0 {
		options = append(options, glamour.WithWordWrap(width))
	}

	renderer, err := glamour.NewTermRenderer(options...)
	if err != nil {
		return "", fmt.Errorf("create glamour renderer: %w", err)
	}

	markdown := fencedCodeBlock(block.Language, block.Body)
	rendered, err := renderer.Render(markdown)
	if err != nil {
		return "", fmt.Errorf("render glamour markdown: %w", err)
	}
	return strings.TrimRight(rendered, "\n"), nil
}

// maxCodeWidth returns one readable code-render width when a terminal width is known.
func (p *Printer) maxCodeWidth() int {
	if p.mode.Format != FormatHuman || p.mode.Width <= 0 {
		return 0
	}
	width := p.maxPanelWidth() - 6
	if width < 40 {
		return 40
	}
	if width > 100 {
		return 100
	}
	return width
}

// fencedCodeBlock returns one Markdown fenced code block string for Glamour rendering.
func fencedCodeBlock(language string, body string) string {
	trimmed := strings.TrimRight(body, "\n")
	if trimmed == "" {
		return "```\n```"
	}
	return fmt.Sprintf("```%s\n%s\n```", strings.TrimSpace(language), trimmed)
}
