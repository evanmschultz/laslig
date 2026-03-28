package laslig

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// Paragraph writes one wrapped long-form text block.
func (p *Printer) Paragraph(paragraph Paragraph) error {
	if p.mode.Format == FormatJSON {
		return p.writeJSON("paragraph", paragraph)
	}
	if err := p.beginBlock(blockKindContent); err != nil {
		return fmt.Errorf("prepare paragraph: %w", err)
	}

	lines := []string{}
	if trimmed := strings.TrimSpace(paragraph.Title); trimmed != "" {
		lines = append(lines, p.renderHeading(p.wrapText(trimmed, p.maxTextWidth())))
	}

	body := paragraph.Body
	if p.mode.Format == FormatHuman {
		body = p.wrapText(body, p.maxTextWidth())
	}
	if strings.TrimSpace(body) != "" {
		lines = append(lines, body)
	}

	if trimmed := strings.TrimSpace(paragraph.Footer); trimmed != "" {
		footer := trimmed
		if p.mode.Format == FormatHuman {
			footer = p.wrapText(footer, p.maxTextWidth())
			footer = p.theme.Muted.Render(footer)
		}
		lines = append(lines, footer)
	}

	if _, err := fmt.Fprintln(p.out, strings.Join(lines, "\n\n")); err != nil {
		return fmt.Errorf("write paragraph: %w", err)
	}
	return nil
}

// StatusLine writes one compact semantic single-line status row.
func (p *Printer) StatusLine(line StatusLine) error {
	if line.Level == "" {
		line.Level = NoticeInfoLevel
	}
	if p.mode.Format == FormatJSON {
		return p.writeJSON("status_line", line)
	}
	if err := p.beginBlock(blockKindContent); err != nil {
		return fmt.Errorf("prepare status line: %w", err)
	}

	label := strings.TrimSpace(line.Label)
	if label == "" {
		label = strings.ToUpper(string(line.Level))
	}

	if p.mode.Format != FormatHuman {
		rendered := fmt.Sprintf("[%s] %s", strings.ToUpper(label), line.Text)
		if detail := strings.TrimSpace(line.Detail); detail != "" {
			rendered += " (" + detail + ")"
		}
		if _, err := fmt.Fprintln(p.out, rendered); err != nil {
			return fmt.Errorf("write status line: %w", err)
		}
		return nil
	}

	parts := []string{
		p.renderStatusLabel(label, line.Level),
		" ",
		p.theme.Value.Render(line.Text),
	}
	if detail := strings.TrimSpace(line.Detail); detail != "" {
		parts = append(parts, " ", p.theme.Muted.Render("("+detail+")"))
	}
	if _, err := fmt.Fprintln(p.out, lipgloss.JoinHorizontal(lipgloss.Top, parts...)); err != nil {
		return fmt.Errorf("write status line: %w", err)
	}
	return nil
}

// renderStatusLabel renders one compact status badge using a notice-level palette.
func (p *Printer) renderStatusLabel(label string, level NoticeLevel) string {
	plain := "[" + strings.ToUpper(strings.TrimSpace(label)) + "]"
	if p.mode.Format != FormatHuman {
		return plain
	}

	text := strings.ToUpper(strings.TrimSpace(label))
	switch level {
	case NoticeSuccessLevel:
		return p.theme.NoticeSuccess.Render(text)
	case NoticeWarningLevel:
		return p.theme.NoticeWarning.Render(text)
	case NoticeErrorLevel:
		return p.theme.NoticeError.Render(text)
	default:
		return p.theme.NoticeInfo.Render(text)
	}
}
