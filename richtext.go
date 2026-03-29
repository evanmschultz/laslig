package laslig

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	internalglamrender "github.com/evanmschultz/laslig/internal/glamrender"
)

// Markdown writes one Markdown block.
func (p *Printer) Markdown(block Markdown) error {
	if p.mode.Format == FormatJSON {
		return p.writeJSON("markdown", block)
	}
	if err := p.beginBlock(blockKindContent); err != nil {
		return fmt.Errorf("prepare markdown: %w", err)
	}

	body := strings.TrimRight(block.Body, "\n")
	if p.mode.Format == FormatHuman && p.mode.Styled {
		rendered, err := p.renderStyledMarkdown(block.Body)
		if err != nil {
			return fmt.Errorf("render markdown: %w", err)
		}
		body = rendered
	}

	lines := []string{}
	if trimmed := strings.TrimSpace(block.Title); trimmed != "" {
		lines = append(lines, p.renderHeading(p.wrapText(trimmed, p.maxTextWidth())))
	}
	if strings.TrimSpace(body) != "" {
		lines = append(lines, body)
	}
	if trimmed := strings.TrimSpace(block.Footer); trimmed != "" {
		rendered := trimmed
		if p.mode.Format == FormatHuman {
			rendered = p.wrapText(trimmed, p.maxTextWidth())
			rendered = p.theme.Muted.Render(rendered)
		}
		lines = append(lines, rendered)
	}

	if err := p.writeContentString(strings.Join(lines, "\n\n")); err != nil {
		return fmt.Errorf("write markdown: %w", err)
	}
	return nil
}

// CodeBlock writes one titled code-style block.
func (p *Printer) CodeBlock(block CodeBlock) error {
	if p.mode.Format == FormatJSON {
		return p.writeJSON("code_block", block)
	}
	if err := p.beginBlock(blockKindContent); err != nil {
		return fmt.Errorf("prepare code block: %w", err)
	}

	body := strings.TrimRight(block.Body, "\n")
	if p.mode.Format == FormatHuman && p.mode.Styled {
		rendered, err := p.renderStyledMarkdown(internalglamrender.FencedCodeBlock(block.Language, block.Body))
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
	if err := p.beginBlock(blockKindContent); err != nil {
		return fmt.Errorf("prepare log block: %w", err)
	}

	body := strings.TrimRight(block.Body, "\n")
	if p.mode.Format == FormatHuman && p.mode.Styled {
		body = p.renderLogBody(body)
	}
	if err := p.writeFramedBlock("log block", block.Title, body, block.Footer); err != nil {
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
	if err := p.writeContentString(content); err != nil {
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
		style = p.constrainStyledBlockWidth(style, maxWidth)
	}
	return style.Render(content)
}

// renderStyledMarkdown renders one Markdown string through Glamour for ANSI output.
func (p *Printer) renderStyledMarkdown(markdown string) (string, error) {
	return internalglamrender.Render(markdown, p.maxCodeWidth(), string(p.glamourStyle))
}

// renderLogBody applies semantic highlighting to explicit caller-provided log excerpts.
func (p *Printer) renderLogBody(body string) string {
	lines := strings.Split(body, "\n")
	for index, line := range lines {
		lines[index] = p.renderLogLine(line)
	}
	return strings.Join(lines, "\n")
}

// renderLogLine applies lightweight level styling to one log-like line when it starts with a known level.
func (p *Printer) renderLogLine(line string) string {
	trimmedLeft := strings.TrimLeft(line, " \t")
	indent := line[:len(line)-len(trimmedLeft)]
	if trimmedLeft == "" {
		return line
	}

	for _, level := range []string{"TRACE", "DEBUG", "INFO", "WARN", "WARNING", "ERROR", "FATAL", "SUCCESS"} {
		bracketed := "[" + level + "]"
		switch {
		case strings.HasPrefix(trimmedLeft, level+" "):
			return indent + p.renderLogLevel(level) + trimmedLeft[len(level):]
		case strings.HasPrefix(trimmedLeft, bracketed+" "):
			return indent + "[" + p.renderLogLevel(level) + "]" + trimmedLeft[len(bracketed):]
		}
	}
	return line
}

// renderLogLevel renders one recognized log level token with a calmer semantic foreground color.
func (p *Printer) renderLogLevel(level string) string {
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("245"))
	switch level {
	case "TRACE", "DEBUG":
		style = style.Foreground(lipgloss.Color("245"))
	case "INFO":
		style = style.Foreground(lipgloss.Color("69"))
	case "WARN", "WARNING":
		style = style.Foreground(lipgloss.Color("214"))
	case "ERROR", "FATAL":
		style = style.Foreground(lipgloss.Color("160"))
	case "SUCCESS":
		style = style.Foreground(lipgloss.Color("#04B575"))
	}
	return style.Render(level)
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
