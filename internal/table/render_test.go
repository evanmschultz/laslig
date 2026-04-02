package table

import (
	"regexp"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(value string) string {
	return ansiPattern.ReplaceAllString(value, "")
}

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

// TestRenderHumanAutoWrap verifies long values are wrapped and rows can rebalance
// to satisfy the supplied width budget.
func TestRenderHumanAutoWrap(t *testing.T) {
	got := Render([]string{"artifact", "run_id", "created"}, [][]string{{
		"github.com/evanmschultz/hylla-fixture-go-2/pkg/very-long-artifact-reference/module",
		"run_2026-04-01T00:00:00.123456789Z_very_long",
		"2026-04-01T00:00:00Z",
	}}, Mode{
		Human:    true,
		Width:    58,
		WrapMode: WrapAuto,
	}, Styles{
		Header: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69")),
		Rule:   lipgloss.NewStyle().Foreground(lipgloss.Color("63")),
		Even:   lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		Odd:    lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
	})
	value := stripANSI(got)
	if !strings.Contains(value, "╰") {
		t.Fatalf("Render() = %q, want bottom border", value)
	}

	lines := strings.Split(value, "\n")
	for _, line := range lines {
		if lipgloss.Width(line) > 58 {
			t.Fatalf("styled table line exceeds width budget: %q (%d > 58)", line, lipgloss.Width(line))
		}
	}

	if strings.Count(value, "\n") <= 2 {
		t.Fatalf("Render() = %q, expected wrapped rows for constrained width", value)
	}
}

// TestRenderHumanTruncate keeps one-line cells when truncate mode is set.
func TestRenderHumanTruncate(t *testing.T) {
	got := Render([]string{"artifact", "status"}, [][]string{{
		"very-very-long-artifact-reference-should-truncate-when-narrow",
		"ready",
	}}, Mode{
		Human:    true,
		Width:    42,
		WrapMode: WrapTruncate,
	}, Styles{
		Header: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69")),
		Rule:   lipgloss.NewStyle().Foreground(lipgloss.Color("63")),
		Even:   lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		Odd:    lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
	})
	value := stripANSI(got)
	if !strings.Contains(value, "╰") {
		t.Fatalf("Render() = %q, want bottom border", value)
	}
	if !strings.Contains(value, "…") {
		t.Fatalf("Render() = %q, want truncation marker (…)", value)
	}

	lines := strings.Split(value, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if lipgloss.Width(line) > 42 {
			t.Fatalf("truncated table line exceeds width budget: %q (%d > 42)", line, lipgloss.Width(line))
		}
	}
}

// TestRenderHumanNever keeps row cells to a single visual row and truncates if needed.
func TestRenderHumanNever(t *testing.T) {
	got := Render([]string{"artifact", "status"}, [][]string{{
		"one--very-very-very-long-cell-that-does-not-wrap-when-never-mode-is-used",
		"ready",
	}}, Mode{
		Human:    true,
		Width:    54,
		WrapMode: WrapNever,
	}, Styles{
		Header: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69")),
		Rule:   lipgloss.NewStyle().Foreground(lipgloss.Color("63")),
		Even:   lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		Odd:    lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
	})
	value := stripANSI(got)
	if strings.Count(value, "\n") > 6 {
		t.Fatalf("Render() = %q, never mode should avoid multi-line wrapping per cell", value)
	}
	if strings.Contains(value, "…") == false {
		t.Fatalf("Render() = %q, expected truncation marker in never mode", value)
	}
}
