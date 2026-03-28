package main

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/x/exp/golden"
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
