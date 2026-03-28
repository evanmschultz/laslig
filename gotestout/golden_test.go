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
