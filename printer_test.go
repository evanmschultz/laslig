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
