package table

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

// TestRenderPlain verifies plain tables render stable column separators.
func TestRenderPlain(t *testing.T) {
	got := Render([]string{"name", "status"}, [][]string{{"demo", "ready"}}, Mode{}, Styles{})
	want := "name | status\n-----+-------\ndemo | ready"
	if got != want {
		t.Fatalf("Render() = %q, want %q", got, want)
	}
}

// TestRenderHuman verifies human tables render through Lip Gloss with borders and ANSI styling.
func TestRenderHuman(t *testing.T) {
	got := Render([]string{"name", "status"}, [][]string{{"demo", "ready"}}, Mode{Human: true, Width: 80}, Styles{
		Header: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69")),
		Rule:   lipgloss.NewStyle().Foreground(lipgloss.Color("63")),
		Even:   lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		Odd:    lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
	})
	if !strings.Contains(got, "demo") || !strings.Contains(got, "ready") {
		t.Fatalf("Render() = %q, want row content", got)
	}
	if !strings.Contains(got, "╭") || !strings.Contains(got, "\x1b[") {
		t.Fatalf("Render() = %q, want framed ANSI table", got)
	}
}
