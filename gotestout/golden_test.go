package gotestout

import (
	"bytes"
	"strings"
	"testing"

	"github.com/charmbracelet/x/exp/golden"

	"github.com/evanmschultz/laslig"
)

// TestRenderPlainCompactGolden verifies the compact plain renderer structure against a golden snapshot.
func TestRenderPlainCompactGolden(t *testing.T) {
	var buf bytes.Buffer
	_, err := Render(&buf, strings.NewReader(sampleStream), Options{
		Policy: laslig.Policy{
			Format: laslig.FormatPlain,
			Style:  laslig.StyleNever,
		},
		View: ViewCompact,
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	golden.RequireEqual(t, buf.Bytes())
}

// TestRenderHumanStyledCompactGolden verifies fixed-width human/styled output structure.
func TestRenderHumanStyledCompactGolden(t *testing.T) {
	var buf bytes.Buffer
	renderer := NewRenderer(&buf, Options{
		Policy: laslig.Policy{
			Format: laslig.FormatHuman,
			Style:  laslig.StyleAlways,
		},
		View: ViewCompact,
	})
	renderer.mode = laslig.Mode{
		Format: laslig.FormatHuman,
		Styled: true,
		Width:  80,
	}
	renderer.theme = laslig.DefaultTheme(renderer.mode)
	renderer.printer = laslig.NewWithMode(&buf, renderer.mode)

	events, err := Parse(strings.NewReader(sampleStream))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	for _, event := range events {
		if err := renderer.WriteEvent(event); err != nil {
			t.Fatalf("renderer.WriteEvent() error = %v", err)
		}
	}
	if err := renderer.Finish(); err != nil {
		t.Fatalf("renderer.Finish() error = %v", err)
	}

	golden.RequireEqual(t, []byte(stripANSI(buf.String())))
}

// stripANSI removes ANSI escape sequences from one string for stable golden snapshots.
func stripANSI(value string) string {
	return ansiPattern.ReplaceAllString(value, "")
}
