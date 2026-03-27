package laslig

import (
	"bytes"
	"strings"
	"testing"
)

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

// TestNoticePlain verifies plain notice rendering.
func TestNoticePlain(t *testing.T) {
	var buf bytes.Buffer
	printer := NewWithMode(&buf, Mode{Format: FormatPlain})

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
	printer := NewWithMode(&buf, Mode{Format: FormatHuman, Styled: true})

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
	printer := NewWithMode(&buf, Mode{Format: FormatHuman, Styled: false, Width: 52})

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

// TestSectionJSON verifies JSON section rendering.
func TestSectionJSON(t *testing.T) {
	var buf bytes.Buffer
	printer := NewWithMode(&buf, Mode{Format: FormatJSON})

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

// TestRecordJSON verifies machine-readable record rendering.
func TestRecordJSON(t *testing.T) {
	var buf bytes.Buffer
	printer := NewWithMode(&buf, Mode{Format: FormatJSON})

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
	printer := NewWithMode(&buf, Mode{Format: FormatPlain})

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
	printer := NewWithMode(&buf, Mode{Format: FormatJSON})

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
	printer := NewWithMode(&buf, Mode{Format: FormatHuman, Styled: false})

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

// TestListEmptyPlain verifies plain empty list rendering.
func TestListEmptyPlain(t *testing.T) {
	var buf bytes.Buffer
	printer := NewWithMode(&buf, Mode{Format: FormatPlain})

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
	printer := NewWithMode(&buf, Mode{Format: FormatPlain})

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
	printer := NewWithMode(&buf, Mode{Format: FormatHuman, Styled: false})

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
	printer := NewWithMode(&buf, Mode{Format: FormatHuman, Styled: false})

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
	printer := NewWithMode(&buf, Mode{Format: FormatHuman, Styled: false, Width: 56})

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

// TestBoxHumanNoStyle verifies that Box is an alias for Panel.
func TestBoxHumanNoStyle(t *testing.T) {
	var buf bytes.Buffer
	printer := NewWithMode(&buf, Mode{Format: FormatHuman, Styled: false})

	err := printer.Box(Panel{
		Title:  "Alias",
		Body:   "Box should delegate to Panel.",
		Footer: "Still plain here.",
	})
	if err != nil {
		t.Fatalf("Box() error = %v", err)
	}

	want := "Alias\n\nBox should delegate to Panel.\n\nStill plain here.\n"
	if got := buf.String(); got != want {
		t.Fatalf("Box() output = %q, want %q", got, want)
	}
}
