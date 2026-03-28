package main

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/x/exp/golden"
)

// TestRunArgsPlainGolden verifies the plain aggregate walkthrough snapshot.
func TestRunArgsPlainGolden(t *testing.T) {
	var buf bytes.Buffer
	if err := runArgs(&buf, []string{"-format", "plain", "-style", "never"}); err != nil {
		t.Fatalf("runArgs() error = %v", err)
	}

	golden.RequireEqual(t, buf.Bytes())
}

// TestRunArgsHumanStyledGolden verifies the fixed-width styled aggregate snapshot.
func TestRunArgsHumanStyledGolden(t *testing.T) {
	var buf bytes.Buffer
	if err := runArgs(&buf, []string{"-format", "human", "-style", "always"}); err != nil {
		t.Fatalf("runArgs() error = %v", err)
	}

	golden.RequireEqual(t, []byte(stripANSI(buf.String())))
}
