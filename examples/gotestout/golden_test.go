package main

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/x/exp/golden"
)

// TestRunArgsPlainGolden verifies the focused gotestout plain snapshot.
func TestRunArgsPlainGolden(t *testing.T) {
	var buf bytes.Buffer
	if err := runArgs(&buf, []string{"-format", "plain", "-style", "never"}); err != nil {
		t.Fatalf("runArgs() error = %v", err)
	}

	golden.RequireEqual(t, buf.Bytes())
}
