package exampletestutil

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

// TestStripANSI verifies ANSI escape removal for stable assertions.
func TestStripANSI(t *testing.T) {
	got := StripANSI("\x1b[31mhello\x1b[0m")
	if got != "hello" {
		t.Fatalf("StripANSI() = %q, want hello", got)
	}
}

// TestAssertRunArgsInvalidFlag verifies the shared invalid-flag helper behavior.
func TestAssertRunArgsInvalidFlag(t *testing.T) {
	runArgs := func(_ io.Writer, _ []string) error {
		return errors.New("parse flags: flag provided but not defined")
	}
	AssertRunArgsInvalidFlag(t, runArgs)
}

// TestAssertRunArgsPlainContains verifies the shared plain-output assertion helper.
func TestAssertRunArgsPlainContains(t *testing.T) {
	runArgs := func(out io.Writer, _ []string) error {
		_, err := fmt.Fprint(out, "plain output")
		return err
	}
	AssertRunArgsPlainContains(t, runArgs, "plain output")
}

// TestAssertRunArgsHumanStyled verifies the shared human-output assertion helper.
func TestAssertRunArgsHumanStyled(t *testing.T) {
	runArgs := func(out io.Writer, _ []string) error {
		_, err := fmt.Fprint(out, "\x1b[32mstyled output\x1b[0m")
		return err
	}
	AssertRunArgsHumanStyled(t, runArgs, "styled output")
}

// TestAssertMainContains verifies the shared main-entrypoint assertion helper.
func TestAssertMainContains(t *testing.T) {
	mainFunc := func() {
		_, _ = fmt.Fprint(os.Stdout, "main output")
	}
	AssertMainContains(t, mainFunc, "example", "main output")
}
