package main

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/x/exp/golden"
	"github.com/evanmschultz/laslig"
)

// TestRunArgsPlainGolden verifies the plain demo structure against a golden snapshot.
func TestRunArgsPlainGolden(t *testing.T) {
	var buf bytes.Buffer
	err := runArgs(&buf, []string{"-format", "plain", "-style", "never"})
	if err != nil {
		t.Fatalf("runArgs() error = %v", err)
	}

	golden.RequireEqual(t, buf.Bytes())
}

// TestRenderShowcaseHumanStyledGolden verifies fixed-width human styled output structure.
func TestRenderShowcaseHumanStyledGolden(t *testing.T) {
	var buf bytes.Buffer
	printer := laslig.NewWithMode(&buf, laslig.Mode{
		Format: laslig.FormatHuman,
		Styled: true,
		Width:  80,
	})

	if err := renderShowcase(printer); err != nil {
		t.Fatalf("renderShowcase() error = %v", err)
	}

	golden.RequireEqual(t, []byte(stripANSI(buf.String())))
}
