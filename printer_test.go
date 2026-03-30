package laslig

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

// newTestPrinter constructs one printer with the default leading gap disabled so
// primitive formatting tests can focus on block content.
func newTestPrinter(out io.Writer, mode Mode) *Printer {
	layout := DefaultLayout().WithLeadingGap(0)
	return newPrinter(out, mode, layout, DefaultTheme(mode), DefaultGlamourStyle())
}

// TestResolveModePlainForBuffer verifies that non-terminal writers resolve to plain output by default.
func TestResolveModePlainForBuffer(t *testing.T) {
	mode := ResolveMode(&bytes.Buffer{}, Policy{Format: FormatAuto, Style: StyleAuto})
	if mode.Format != FormatPlain {
		t.Fatalf("ResolveMode().Format = %q, want %q", mode.Format, FormatPlain)
	}
	if mode.Styled {
		t.Fatal("ResolveMode().Styled = true, want false")
	}
}

// TestResolveModeHumanStyleAlways verifies explicit human styled output resolution.
func TestResolveModeHumanStyleAlways(t *testing.T) {
	mode := ResolveMode(&bytes.Buffer{}, Policy{Format: FormatHuman, Style: StyleAlways})
	if mode.Format != FormatHuman {
		t.Fatalf("ResolveMode().Format = %q, want %q", mode.Format, FormatHuman)
	}
	if !mode.Styled {
		t.Fatal("ResolveMode().Styled = false, want true")
	}
}

// TestNewUsesCustomTheme verifies callers can swap the default theme through Policy.
func TestNewUsesCustomTheme(t *testing.T) {
	var buf bytes.Buffer
	layout := DefaultLayout().WithLeadingGap(0)
	theme := DefaultTheme(Mode{Format: FormatHuman, Styled: true})
	theme.Section = lipgloss.NewStyle()

	printer := New(&buf, Policy{
		Format: FormatHuman,
		Style:  StyleAlways,
		Layout: &layout,
		Theme:  &theme,
	})
	if err := printer.Section("Deploy"); err != nil {
		t.Fatalf("Section() error = %v", err)
	}

	if got := buf.String(); got != "Deploy\n" {
		t.Fatalf("Section() with custom theme = %q, want %q", got, "Deploy\n")
	}
}

// TestNoticePlain verifies plain notice rendering.
func TestNoticePlain(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatPlain})

	err := printer.Notice(Notice{
		Level: NoticeWarningLevel,
		Title: "Careful",
		Body:  "Something needs attention.",
	})
	if err != nil {
		t.Fatalf("Notice() error = %v", err)
	}

	want := "[WARNING] Careful\n  Something needs attention.\n"
	if got := buf.String(); got != want {
		t.Fatalf("Notice() output = %q, want %q", got, want)
	}
}

// TestNoticeHumanStyled verifies styled human notice rendering and default level handling.
func TestNoticeHumanStyled(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: true})

	err := printer.Notice(Notice{
		Title:  "Heads up",
		Body:   "Styled output should render with a default info badge.",
		Detail: []string{"detail line"},
	})
	if err != nil {
		t.Fatalf("Notice() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "Heads up") {
		t.Fatalf("Notice() output missing title: %q", got)
	}
	if !strings.Contains(got, "detail line") {
		t.Fatalf("Notice() output missing detail: %q", got)
	}
}

// TestNoticeHumanWrap verifies human notice wrapping when a width is available.
func TestNoticeHumanWrap(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: false, Width: 52})

	err := printer.Notice(Notice{
		Title: "Heads up",
		Body:  "This notice body should wrap to fit a narrower terminal width by default.",
	})
	if err != nil {
		t.Fatalf("Notice() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "fit a\n  narrower terminal width by default.") {
		t.Fatalf("Notice() output did not wrap as expected:\n%s", got)
	}
}

// TestNoticeHumanNoStyle verifies unstyled human notices preserve human layout without ANSI.
func TestNoticeHumanNoStyle(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: false})

	err := printer.Notice(Notice{
		Level: NoticeWarningLevel,
		Title: "Careful",
		Body:  "Something needs attention.",
	})
	if err != nil {
		t.Fatalf("Notice() error = %v", err)
	}

	want := "[WARNING] Careful\n  Something needs attention.\n"
	if got := buf.String(); got != want {
		t.Fatalf("Notice() output = %q, want %q", got, want)
	}
	if strings.Contains(buf.String(), "\x1b[") {
		t.Fatalf("Notice() output = %q, want no ANSI", buf.String())
	}
}

// TestRenderBadgeHumanStyled verifies semantic badge values use distinct styled chips.
func TestRenderBadgeHumanStyled(t *testing.T) {
	printer := newTestPrinter(&bytes.Buffer{}, Mode{Format: FormatHuman, Styled: true})

	pass := printer.renderBadge("pass")
	custom := printer.renderBadge("custom")
	warn := printer.renderBadge("warning")

	if !strings.Contains(pass, "PASS") {
		t.Fatalf("renderBadge(pass) = %q, want PASS text", pass)
	}
	if !strings.Contains(custom, "CUSTOM") {
		t.Fatalf("renderBadge(custom) = %q, want CUSTOM text", custom)
	}
	if !strings.Contains(warn, "WARNING") {
		t.Fatalf("renderBadge(warning) = %q, want WARNING text", warn)
	}
	if !strings.Contains(pass, "\x1b[") || !strings.Contains(custom, "\x1b[") || !strings.Contains(warn, "\x1b[") {
		t.Fatal("renderBadge() output missing ANSI styling")
	}
	if pass == custom {
		t.Fatalf("renderBadge(pass) = %q, want semantic style distinct from custom badge", pass)
	}
	if warn == custom {
		t.Fatalf("renderBadge(warning) = %q, want semantic style distinct from custom badge", warn)
	}
}

// TestRenderBadgeHumanNoStyle verifies unstyled human badges fall back to plain bracketed text.
func TestRenderBadgeHumanNoStyle(t *testing.T) {
	printer := newTestPrinter(&bytes.Buffer{}, Mode{Format: FormatHuman, Styled: false})

	pass := printer.renderBadge("pass")
	custom := printer.renderBadge("custom")
	warn := printer.renderBadge("warning")

	if got, want := pass, "[PASS]"; got != want {
		t.Fatalf("renderBadge(pass) = %q, want %q", got, want)
	}
	if got, want := custom, "[CUSTOM]"; got != want {
		t.Fatalf("renderBadge(custom) = %q, want %q", got, want)
	}
	if got, want := warn, "[WARNING]"; got != want {
		t.Fatalf("renderBadge(warning) = %q, want %q", got, want)
	}
	if strings.Contains(pass, "\x1b[") || strings.Contains(custom, "\x1b[") || strings.Contains(warn, "\x1b[") {
		t.Fatal("renderBadge() output contained ANSI styling in unstyled human mode")
	}
}

// TestSectionJSON verifies JSON section rendering.
func TestSectionJSON(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatJSON})

	err := printer.Section("release")
	if err != nil {
		t.Fatalf("Section() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"type": "section"`) {
		t.Fatalf("Section() output = %q, want section type", got)
	}
	if !strings.Contains(got, `"title": "release"`) {
		t.Fatalf("Section() output = %q, want section title", got)
	}
}

// TestSectionPlainSpacing verifies sections participate in the default flow-spacing rhythm.
func TestSectionPlainSpacing(t *testing.T) {
	var buf bytes.Buffer
	printer := NewWithMode(&buf, Mode{Format: FormatPlain})

	if err := printer.Section("Intro"); err != nil {
		t.Fatalf("Section() error = %v", err)
	}
	if err := printer.Notice(Notice{Title: "Readable"}); err != nil {
		t.Fatalf("Notice() error = %v", err)
	}
	if err := printer.Section("Next"); err != nil {
		t.Fatalf("Section() error = %v", err)
	}

	want := "\nIntro\n\n  [INFO] Readable\n\n\nNext\n"
	if got := buf.String(); got != want {
		t.Fatalf("flow spacing output = %q, want %q", got, want)
	}
}

// TestDefaultThemeHumanStyled verifies the default styled theme renders visible ANSI styling.
func TestDefaultThemeHumanStyled(t *testing.T) {
	theme := DefaultTheme(Mode{Format: FormatHuman, Styled: true})
	if got := theme.Section.Render("Heading"); !strings.Contains(got, "\x1b[") {
		t.Fatalf("theme.Section.Render() = %q, want ANSI styling", got)
	}
	if got := theme.NoticeInfo.Render("INFO"); !strings.Contains(got, "\x1b[") {
		t.Fatalf("theme.NoticeInfo.Render() = %q, want ANSI styling", got)
	}
}

// TestRecordJSON verifies machine-readable record rendering.
func TestRecordJSON(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatJSON})

	err := printer.Record(Record{
		Title: "Project",
		Fields: []Field{
			{Label: "name", Value: "laslig", Identifier: true},
		},
	})
	if err != nil {
		t.Fatalf("Record() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"type": "record"`) {
		t.Fatalf("Record() output = %q, want record type", got)
	}
	if !strings.Contains(got, `"title": "Project"`) {
		t.Fatalf("Record() output = %q, want record title", got)
	}
}

// TestKVPlain verifies plain aligned key-value rendering.
func TestKVPlain(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatPlain})

	err := printer.KV(KV{
		Title: "Project",
		Pairs: []Field{
			{Label: "module", Value: "github.com/evanmschultz/laslig", Identifier: true},
			{Label: "task runner", Value: "mage", Muted: true},
		},
	})
	if err != nil {
		t.Fatalf("KV() error = %v", err)
	}

	want := "Project\n  module       github.com/evanmschultz/laslig\n  task runner  mage\n"
	if got := buf.String(); got != want {
		t.Fatalf("KV() output = %q, want %q", got, want)
	}
}

// TestKVJSON verifies machine-readable kv rendering.
func TestKVJSON(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatJSON})

	err := printer.KV(KV{
		Title: "Project",
		Pairs: []Field{
			{Label: "module", Value: "github.com/evanmschultz/laslig", Identifier: true},
		},
	})
	if err != nil {
		t.Fatalf("KV() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"type": "kv"`) {
		t.Fatalf("KV() output = %q, want kv type", got)
	}
}

// TestListHumanNoStyle verifies unstyled human list rendering.
func TestListHumanNoStyle(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: false})

	err := printer.List(List{
		Title: "Profiles",
		Items: []ListItem{{
			Title: "default",
			Fields: []Field{
				{Label: "provider", Value: "codex", Muted: true},
			},
		}},
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	want := "Profiles\n- default\n  provider: codex\n"
	if got := buf.String(); got != want {
		t.Fatalf("List() output = %q, want %q", got, want)
	}
}

// TestListHumanNoStyleBadge verifies unstyled human list badges stay plain without ANSI.
func TestListHumanNoStyleBadge(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: false})

	err := printer.List(List{
		Title: "Profiles",
		Items: []ListItem{{
			Title: "dev",
			Badge: "active",
		}},
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	want := "Profiles\n- dev [ACTIVE]\n"
	if got := buf.String(); got != want {
		t.Fatalf("List() output = %q, want %q", got, want)
	}
	if strings.Contains(buf.String(), "\x1b[") {
		t.Fatalf("List() output = %q, want no ANSI", buf.String())
	}
}

// TestRecordHumanNoStyleBadge verifies unstyled human record field badges stay plain without ANSI.
func TestRecordHumanNoStyleBadge(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: false})

	err := printer.Record(Record{
		Title: "Build",
		Fields: []Field{
			{Label: "status", Value: "pass", Badge: true},
		},
	})
	if err != nil {
		t.Fatalf("Record() error = %v", err)
	}

	want := "Build\n  status: [PASS]\n"
	if got := buf.String(); got != want {
		t.Fatalf("Record() output = %q, want %q", got, want)
	}
	if strings.Contains(buf.String(), "\x1b[") {
		t.Fatalf("Record() output = %q, want no ANSI", buf.String())
	}
}

// TestListEmptyPlain verifies plain empty list rendering.
func TestListEmptyPlain(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatPlain})

	err := printer.List(List{Title: "Profiles"})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	want := "Profiles\n- (none)\n"
	if got := buf.String(); got != want {
		t.Fatalf("List() output = %q, want %q", got, want)
	}
}

// TestTablePlain verifies plain table rendering.
func TestTablePlain(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatPlain})

	err := printer.Table(Table{
		Title:  "Packages",
		Header: []string{"name", "status"},
		Rows: [][]string{
			{"github.com/evanmschultz/laslig", "pass"},
		},
	})
	if err != nil {
		t.Fatalf("Table() error = %v", err)
	}

	want := "Packages\nname                           | status\n-------------------------------+-------\ngithub.com/evanmschultz/laslig | pass\n"
	if got := buf.String(); got != want {
		t.Fatalf("Table() output = %q, want %q", got, want)
	}
}

// TestTableEmptyHumanNoStyle verifies human empty table rendering.
func TestTableEmptyHumanNoStyle(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: false})

	err := printer.Table(Table{Title: "Packages"})
	if err != nil {
		t.Fatalf("Table() error = %v", err)
	}

	want := "Packages\n(none)\n"
	if got := buf.String(); got != want {
		t.Fatalf("Table() output = %q, want %q", got, want)
	}
}

// TestPanelHumanNoStyle verifies unstyled human panel rendering.
func TestPanelHumanNoStyle(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: false})

	err := printer.Panel(Panel{
		Title:  "Next step",
		Body:   "Run mage test.",
		Footer: "The repo is ready.",
	})
	if err != nil {
		t.Fatalf("Panel() error = %v", err)
	}

	want := "Next step\n\nRun mage test.\n\nThe repo is ready.\n"
	if got := buf.String(); got != want {
		t.Fatalf("Panel() output = %q, want %q", got, want)
	}
}

// TestPanelHumanWrap verifies panel content wraps when a width is available.
func TestPanelHumanWrap(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: false, Width: 56})

	err := printer.Panel(Panel{
		Title:  "Why this shape",
		Body:   "Panels should avoid stretching across the full terminal when a smaller readable width is more appropriate.",
		Footer: "Readable defaults matter.",
	})
	if err != nil {
		t.Fatalf("Panel() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "across the full\nterminal when a smaller") {
		t.Fatalf("Panel() output did not wrap as expected:\n%s", got)
	}
}

// TestParagraphPlain verifies plain paragraph rendering preserves structure.
func TestParagraphPlain(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatPlain})

	err := printer.Paragraph(Paragraph{
		Title:  "Why",
		Body:   "Laslig keeps ordinary command output readable.",
		Footer: "Writers in, errors out.",
	})
	if err != nil {
		t.Fatalf("Paragraph() error = %v", err)
	}

	want := "Why\n\nLaslig keeps ordinary command output readable.\n\nWriters in, errors out.\n"
	if got := buf.String(); got != want {
		t.Fatalf("Paragraph() output = %q, want %q", got, want)
	}
}

// TestParagraphHumanWrap verifies paragraph bodies wrap for narrower human widths.
func TestParagraphHumanWrap(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: false, Width: 48})

	err := printer.Paragraph(Paragraph{
		Title: "Why",
		Body:  "Paragraph helpers should wrap body text to a readable width by default.",
	})
	if err != nil {
		t.Fatalf("Paragraph() error = %v", err)
	}

	if got := buf.String(); !strings.Contains(got, "wrap body text\nto a readable width by default.") {
		t.Fatalf("Paragraph() output did not wrap as expected:\n%s", got)
	}
}

// TestStatusLinePlain verifies plain status-line rendering is compact and stable.
func TestStatusLinePlain(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatPlain})

	err := printer.StatusLine(StatusLine{
		Level:  NoticeSuccessLevel,
		Text:   "Build ready",
		Detail: "mage check",
	})
	if err != nil {
		t.Fatalf("StatusLine() error = %v", err)
	}

	want := "[SUCCESS] Build ready (mage check)\n"
	if got := buf.String(); got != want {
		t.Fatalf("StatusLine() output = %q, want %q", got, want)
	}
}

// TestStatusLineHumanNoStyle verifies unstyled human output keeps the human
// layout but falls back to plain bracketed labels with no ANSI.
func TestStatusLineHumanNoStyle(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: false})

	err := printer.StatusLine(StatusLine{
		Level:  NoticeSuccessLevel,
		Text:   "Build ready",
		Detail: "mage check",
	})
	if err != nil {
		t.Fatalf("StatusLine() error = %v", err)
	}

	want := "[SUCCESS] Build ready (mage check)\n"
	if got := buf.String(); got != want {
		t.Fatalf("StatusLine() output = %q, want %q", got, want)
	}
	if strings.Contains(buf.String(), "\x1b[") {
		t.Fatalf("StatusLine() output = %q, want no ANSI", buf.String())
	}
}

// TestStatusLineJSON verifies machine-readable status-line rendering.
func TestStatusLineJSON(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatJSON})

	err := printer.StatusLine(StatusLine{
		Level: NoticeWarningLevel,
		Text:  "Coverage dipped",
	})
	if err != nil {
		t.Fatalf("StatusLine() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"type": "status_line"`) {
		t.Fatalf("StatusLine() output = %q, want status_line type", got)
	}
	if !strings.Contains(got, `"text": "Coverage dipped"`) {
		t.Fatalf("StatusLine() output = %q, want text field", got)
	}
}

// TestMarkdownPlain verifies plain Markdown rendering preserves the source text.
func TestMarkdownPlain(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatPlain})

	err := printer.Markdown(Markdown{
		Title:  "Notes",
		Body:   "# Heading\n\n- first\n- second",
		Footer: "Rendered as source in plain mode.",
	})
	if err != nil {
		t.Fatalf("Markdown() error = %v", err)
	}

	want := "Notes\n\n# Heading\n\n- first\n- second\n\nRendered as source in plain mode.\n"
	if got := buf.String(); got != want {
		t.Fatalf("Markdown() output = %q, want %q", got, want)
	}
}

// TestMarkdownHumanStyled verifies styled Markdown rendering flows through Glamour.
func TestMarkdownHumanStyled(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: true, Width: 80})

	err := printer.Markdown(Markdown{
		Body: "# Heading\n\n- first item\n- second item",
	})
	if err != nil {
		t.Fatalf("Markdown() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "Heading") {
		t.Fatalf("Markdown() output missing heading:\n%s", got)
	}
	if !strings.Contains(got, "first") || !strings.Contains(got, "second") {
		t.Fatalf("Markdown() output missing list items:\n%s", got)
	}
	if !strings.Contains(got, "\x1b[") {
		t.Fatalf("Markdown() output missing ANSI styling: %q", got)
	}
}

// TestCodeBlockPlain verifies plain code-block rendering preserves content.
func TestCodeBlockPlain(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatPlain})

	err := printer.CodeBlock(CodeBlock{
		Title:    "Example",
		Language: "go",
		Body:     "fmt.Println(\"hello\")",
		Footer:   "Go snippet.",
	})
	if err != nil {
		t.Fatalf("CodeBlock() error = %v", err)
	}

	want := "Example\n\nfmt.Println(\"hello\")\n\nGo snippet.\n"
	if got := buf.String(); got != want {
		t.Fatalf("CodeBlock() output = %q, want %q", got, want)
	}
}

// TestCodeBlockHumanStyled verifies styled code-block rendering uses ANSI output.
func TestCodeBlockHumanStyled(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: true, Width: 80})

	err := printer.CodeBlock(CodeBlock{
		Title:    "Example",
		Language: "go",
		Body:     "fmt.Println(\"hello\")",
	})
	if err != nil {
		t.Fatalf("CodeBlock() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "Println") {
		t.Fatalf("CodeBlock() output missing code body:\n%s", got)
	}
	if !strings.Contains(got, "\"hello\"") {
		t.Fatalf("CodeBlock() output missing string literal:\n%s", got)
	}
	if !strings.Contains(got, "\x1b[") {
		t.Fatalf("CodeBlock() output missing ANSI styling: %q", got)
	}
}

// TestLogBlockPlain verifies plain log-block rendering preserves newlines exactly.
func TestLogBlockPlain(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatPlain})

	err := printer.LogBlock(LogBlock{
		Title: "Recent logs",
		Body:  "INFO boot complete\nWARN retry scheduled",
	})
	if err != nil {
		t.Fatalf("LogBlock() error = %v", err)
	}

	want := "Recent logs\n\nINFO boot complete\nWARN retry scheduled\n"
	if got := buf.String(); got != want {
		t.Fatalf("LogBlock() output = %q, want %q", got, want)
	}
}

// TestLogBlockJSON verifies machine-readable log-block rendering.
func TestLogBlockJSON(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatJSON})

	err := printer.LogBlock(LogBlock{
		Title: "stderr",
		Body:  "panic: boom",
	})
	if err != nil {
		t.Fatalf("LogBlock() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"type": "log_block"`) {
		t.Fatalf("LogBlock() output = %q, want log_block type", got)
	}
	if !strings.Contains(got, `"title": "stderr"`) {
		t.Fatalf("LogBlock() output = %q, want title", got)
	}
}
