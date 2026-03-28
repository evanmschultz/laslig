package laslig

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"charm.land/lipgloss/v2"

	internallayout "github.com/evanmschultz/laslig/internal/layout"
	internaltable "github.com/evanmschultz/laslig/internal/table"
)

// Printer renders structured output to one writer.
type Printer struct {
	out         io.Writer
	mode        Mode
	theme       Theme
	wroteBlocks bool
	lastBlock   blockKind
}

// blockKind identifies one top-level rendered block family for spacing rules.
type blockKind int

const (
	// blockKindContent identifies an ordinary content block.
	blockKindContent blockKind = iota
	// blockKindSection identifies a section heading block.
	blockKindSection
)

// New constructs one printer by resolving a writer against the provided policy.
func New(out io.Writer, policy Policy) *Printer {
	return NewWithMode(out, ResolveMode(out, policy))
}

// NewWithMode constructs one printer using an already-resolved output mode.
func NewWithMode(out io.Writer, mode Mode) *Printer {
	if out == nil {
		out = io.Discard
	}
	return &Printer{
		out:   out,
		mode:  mode,
		theme: DefaultTheme(mode),
	}
}

// Mode returns the resolved output mode used by the printer.
func (p *Printer) Mode() Mode {
	return p.mode
}

// Section writes one section heading.
func (p *Printer) Section(title string) error {
	if p.mode.Format == FormatJSON {
		return p.writeJSON("section", map[string]any{
			"title": title,
		})
	}
	if err := p.beginBlock(blockKindSection); err != nil {
		return fmt.Errorf("prepare section: %w", err)
	}

	value := title
	if p.mode.Format == FormatHuman {
		value = p.theme.Section.Render(title)
	}
	_, err := fmt.Fprintf(p.out, "%s\n", value)
	if err != nil {
		return fmt.Errorf("write section: %w", err)
	}
	return nil
}

// Notice writes one user-facing notice block.
func (p *Printer) Notice(notice Notice) error {
	if notice.Level == "" {
		notice.Level = NoticeInfoLevel
	}
	if p.mode.Format == FormatJSON {
		return p.writeJSON("notice", notice)
	}
	if err := p.beginBlock(blockKindContent); err != nil {
		return fmt.Errorf("prepare notice: %w", err)
	}

	if p.mode.Format == FormatHuman {
		headline := p.noticeBadge(notice.Level)
		if notice.Title != "" {
			headline += " " + p.theme.Value.Render(notice.Title)
		}
		lines := []string{headline}
		textWidth := p.maxTextWidth()
		if body := strings.TrimSpace(notice.Body); body != "" {
			lines = append(lines, internallayout.IndentBlock("  ", p.wrapText(body, textWidth)))
		}
		for _, detail := range notice.Detail {
			lines = append(lines, internallayout.IndentBlock("  ", p.theme.Muted.Render(p.wrapText(detail, textWidth))))
		}
		_, err := fmt.Fprintln(p.out, strings.Join(lines, "\n"))
		if err != nil {
			return fmt.Errorf("write notice: %w", err)
		}
		return nil
	}

	headline := "[" + strings.ToUpper(string(notice.Level)) + "]"
	if notice.Title != "" {
		headline += " " + notice.Title
	}
	lines := []string{headline}
	if body := strings.TrimSpace(notice.Body); body != "" {
		lines = append(lines, "  "+body)
	}
	for _, detail := range notice.Detail {
		lines = append(lines, "  "+detail)
	}
	_, err := fmt.Fprintln(p.out, strings.Join(lines, "\n"))
	if err != nil {
		return fmt.Errorf("write notice: %w", err)
	}
	return nil
}

// Record writes one titled record block.
func (p *Printer) Record(record Record) error {
	if p.mode.Format == FormatJSON {
		return p.writeJSON("record", record)
	}
	if err := p.beginBlock(blockKindContent); err != nil {
		return fmt.Errorf("prepare record: %w", err)
	}

	if _, err := fmt.Fprintln(p.out, p.renderHeading(record.Title)); err != nil {
		return fmt.Errorf("write record heading: %w", err)
	}
	if len(record.Fields) == 0 {
		return p.writeEmpty("  (none)")
	}
	for _, field := range record.Fields {
		if _, err := fmt.Fprintf(p.out, "%s\n", p.renderField(field)); err != nil {
			return fmt.Errorf("write record field: %w", err)
		}
	}
	return nil
}

// KV writes one aligned key-value block.
func (p *Printer) KV(kv KV) error {
	if p.mode.Format == FormatJSON {
		return p.writeJSON("kv", kv)
	}
	if err := p.beginBlock(blockKindContent); err != nil {
		return fmt.Errorf("prepare kv: %w", err)
	}

	if strings.TrimSpace(kv.Title) != "" {
		if _, err := fmt.Fprintln(p.out, p.renderHeading(kv.Title)); err != nil {
			return fmt.Errorf("write kv heading: %w", err)
		}
	}
	if len(kv.Pairs) == 0 {
		empty := kv.Empty
		if strings.TrimSpace(empty) == "" {
			empty = "(none)"
		}
		return p.writeEmpty("  " + empty)
	}

	width := 0
	for _, pair := range kv.Pairs {
		if cellWidth := lipgloss.Width(pair.Label); cellWidth > width {
			width = cellWidth
		}
	}

	for _, pair := range kv.Pairs {
		label := pair.Label
		if p.mode.Format == FormatHuman {
			label = p.theme.Label.Render(lipgloss.NewStyle().Width(width).Render(label))
		} else {
			label = fmt.Sprintf("%-*s", width, label)
		}
		if _, err := fmt.Fprintf(p.out, "  %s  %s\n", label, p.renderFieldValue(pair)); err != nil {
			return fmt.Errorf("write kv pair: %w", err)
		}
	}
	return nil
}

// List writes one titled list block.
func (p *Printer) List(list List) error {
	if p.mode.Format == FormatJSON {
		return p.writeJSON("list", list)
	}
	if err := p.beginBlock(blockKindContent); err != nil {
		return fmt.Errorf("prepare list: %w", err)
	}

	if _, err := fmt.Fprintln(p.out, p.renderHeading(list.Title)); err != nil {
		return fmt.Errorf("write list heading: %w", err)
	}
	if len(list.Items) == 0 {
		empty := list.Empty
		if strings.TrimSpace(empty) == "" {
			empty = "(none)"
		}
		return p.writeEmpty("- " + empty)
	}
	for _, item := range list.Items {
		title := item.Title
		if strings.TrimSpace(item.Badge) != "" {
			title += " " + p.renderBadge(item.Badge)
		}
		if _, err := fmt.Fprintf(p.out, "- %s\n", p.renderValue(title)); err != nil {
			return fmt.Errorf("write list title: %w", err)
		}
		for _, field := range item.Fields {
			if _, err := fmt.Fprintf(p.out, "%s\n", p.renderField(field)); err != nil {
				return fmt.Errorf("write list field: %w", err)
			}
		}
	}
	return nil
}

// Table writes one titled table block.
func (p *Printer) Table(table Table) error {
	if p.mode.Format == FormatJSON {
		return p.writeJSON("table", table)
	}
	if err := p.beginBlock(blockKindContent); err != nil {
		return fmt.Errorf("prepare table: %w", err)
	}

	if strings.TrimSpace(table.Title) != "" {
		if _, err := fmt.Fprintln(p.out, p.renderHeading(table.Title)); err != nil {
			return fmt.Errorf("write table heading: %w", err)
		}
	}
	if len(table.Rows) == 0 {
		empty := table.Empty
		if strings.TrimSpace(empty) == "" {
			empty = "(none)"
		}
		return p.writeEmpty("  " + empty)
	}

	rendered := internaltable.Render(table.Header, table.Rows, internaltable.Mode{
		Human: p.mode.Format == FormatHuman,
		Width: p.mode.Width,
	}, internaltable.Styles{
		Header: p.theme.TableHeader,
		Rule:   p.theme.TableRule,
		Even:   lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		Odd:    lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
	})
	if _, err := fmt.Fprintln(p.out, rendered); err != nil {
		return fmt.Errorf("write table body: %w", err)
	}
	if strings.TrimSpace(table.Caption) != "" {
		caption := table.Caption
		if p.mode.Format == FormatHuman {
			caption = p.theme.Muted.Render(caption)
		}
		if _, err := fmt.Fprintf(p.out, "%s\n", caption); err != nil {
			return fmt.Errorf("write table caption: %w", err)
		}
	}
	return nil
}

// Panel writes one titled panel block.
func (p *Printer) Panel(panel Panel) error {
	if p.mode.Format == FormatJSON {
		return p.writeJSON("panel", panel)
	}
	if err := p.beginBlock(blockKindContent); err != nil {
		return fmt.Errorf("prepare panel: %w", err)
	}

	lines := []string{}
	if strings.TrimSpace(panel.Title) != "" {
		lines = append(lines, p.renderHeading(p.wrapText(panel.Title, p.maxTextWidth())))
	}
	lines = append(lines, p.wrapText(panel.Body, p.maxTextWidth()))
	if strings.TrimSpace(panel.Footer) != "" {
		footer := panel.Footer
		if p.mode.Format == FormatHuman {
			footer = p.theme.Muted.Render(p.wrapText(footer, p.maxTextWidth()))
		}
		lines = append(lines, footer)
	}

	content := strings.Join(lines, "\n\n")
	if p.mode.Format == FormatHuman {
		if maxWidth := p.maxPanelWidth(); maxWidth > 0 && p.mode.Styled {
			content = p.theme.Panel.MaxWidth(maxWidth).Render(content)
		} else if p.mode.Styled {
			content = p.theme.Panel.Render(content)
		}
	}
	if _, err := fmt.Fprintln(p.out, content); err != nil {
		return fmt.Errorf("write panel: %w", err)
	}
	return nil
}

// Box writes one panel-style block.
func (p *Printer) Box(panel Panel) error {
	return p.Panel(panel)
}

func (p *Printer) writeJSON(kind string, payload any) error {
	encoder := json.NewEncoder(p.out)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(map[string]any{
		"type":    kind,
		"payload": payload,
	}); err != nil {
		return fmt.Errorf("write %s json: %w", kind, err)
	}
	return nil
}

// beginBlock applies the default flow spacing between consecutive rendered blocks.
func (p *Printer) beginBlock(kind blockKind) error {
	if p.mode.Format == FormatJSON {
		return nil
	}
	if !p.wroteBlocks {
		p.wroteBlocks = true
		p.lastBlock = kind
		return nil
	}

	gapLines := 1
	if kind == blockKindSection && p.lastBlock != blockKindSection {
		gapLines = 2
	}
	if err := p.writeBlankLines(gapLines); err != nil {
		return err
	}
	p.lastBlock = kind
	return nil
}

// writeBlankLines writes one or more empty separator lines to the printer output.
func (p *Printer) writeBlankLines(count int) error {
	for range count {
		if _, err := fmt.Fprintln(p.out); err != nil {
			return fmt.Errorf("write separator line: %w", err)
		}
	}
	return nil
}

func (p *Printer) writeEmpty(value string) error {
	switch p.mode.Format {
	case FormatHuman:
		trimmed := strings.TrimSpace(value)
		if strings.HasPrefix(trimmed, "-") {
			_, err := fmt.Fprintf(p.out, "- %s\n", p.theme.Muted.Render(strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))))
			if err != nil {
				return fmt.Errorf("write empty state: %w", err)
			}
			return nil
		}
		_, err := fmt.Fprintf(p.out, "%s\n", p.theme.Muted.Render(strings.TrimSpace(value)))
		if err != nil {
			return fmt.Errorf("write empty state: %w", err)
		}
		return nil
	default:
		_, err := fmt.Fprintln(p.out, value)
		if err != nil {
			return fmt.Errorf("write empty state: %w", err)
		}
		return nil
	}
}

func (p *Printer) renderHeading(value string) string {
	if p.mode.Format != FormatHuman {
		return value
	}
	return p.theme.Section.Render(value)
}

func (p *Printer) renderField(field Field) string {
	label := field.Label + ":"
	if p.mode.Format == FormatHuman {
		label = p.theme.Label.Render(label)
	}
	return "  " + label + " " + p.renderFieldValue(field)
}

func (p *Printer) renderFieldValue(field Field) string {
	switch {
	case field.Badge:
		return p.renderBadge(field.Value)
	case field.Identifier:
		if p.mode.Format == FormatHuman {
			return p.theme.Identifier.Render(field.Value)
		}
	case field.Muted:
		if p.mode.Format == FormatHuman {
			return p.theme.Muted.Render(field.Value)
		}
	}
	return p.renderValue(field.Value)
}

func (p *Printer) renderValue(value string) string {
	if p.mode.Format != FormatHuman {
		return value
	}
	return p.theme.Value.Render(value)
}

func (p *Printer) renderBadge(value string) string {
	trimmed := strings.ToUpper(strings.TrimSpace(value))
	if p.mode.Format != FormatHuman {
		return "[" + trimmed + "]"
	}

	style := p.theme.Badge
	// Mirror blick's compact state-chip palette so only semantic states get strong fills.
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "approved", "active", "success", "pass", "ready", "live":
		style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230")).Background(lipgloss.Color("#04B575")).Padding(0, 1)
	case "pending", "warn", "warning":
		style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("232")).Background(lipgloss.Color("214")).Padding(0, 1)
	case "denied", "revoked", "error", "fail", "failed":
		style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230")).Background(lipgloss.Color("160")).Padding(0, 1)
	case "canceled", "cancelled", "disabled":
		style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255")).Background(lipgloss.Color("240")).Padding(0, 1)
	case "expired":
		style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255")).Background(lipgloss.Color("238")).Padding(0, 1)
	}
	return style.Render(trimmed)
}

func (p *Printer) noticeBadge(level NoticeLevel) string {
	plain := "[" + strings.ToUpper(string(level)) + "]"
	if p.mode.Format != FormatHuman {
		return plain
	}
	switch level {
	case NoticeSuccessLevel:
		return p.theme.NoticeSuccess.Render("SUCCESS")
	case NoticeWarningLevel:
		return p.theme.NoticeWarning.Render("WARNING")
	case NoticeErrorLevel:
		return p.theme.NoticeError.Render("ERROR")
	default:
		return p.theme.NoticeInfo.Render("INFO")
	}
}

func (p *Printer) maxTextWidth() int {
	if p.mode.Format != FormatHuman || p.mode.Width <= 0 {
		return 0
	}
	width := p.mode.Width - 8
	if width < 32 {
		return 32
	}
	if width > 76 {
		return 76
	}
	return width
}

func (p *Printer) maxPanelWidth() int {
	if p.mode.Format != FormatHuman || p.mode.Width <= 0 {
		return 0
	}
	width := p.mode.Width - 4
	if width < 48 {
		return 48
	}
	if width > 88 {
		return 88
	}
	return width
}

func (p *Printer) wrapText(value string, width int) string {
	if width <= 0 || p.mode.Format != FormatHuman {
		return value
	}
	return internallayout.WrapText(value, width)
}
