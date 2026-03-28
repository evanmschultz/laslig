package main

import (
	"testing"

	"github.com/evanmschultz/laslig/internal/exampletestutil"
)

// TestRunArgsPlain verifies the focused Section example renders its plain output.
func TestRunArgsPlain(t *testing.T) {
	exampletestutil.AssertRunArgsPlainContains(t, runArgs, "Section", "Owned blocks", "Next section")
}

// TestRunArgsHumanStyled verifies the focused Section example emits styled output.
func TestRunArgsHumanStyled(t *testing.T) {
	exampletestutil.AssertRunArgsHumanStyled(t, runArgs, "Section boundaries reset ownership")
}

// TestRunArgsInvalidFlag verifies invalid arguments return a parse error.
func TestRunArgsInvalidFlag(t *testing.T) {
	exampletestutil.AssertRunArgsInvalidFlag(t, runArgs)
}

// TestMain verifies the command entrypoint succeeds for a valid invocation.
func TestMain(t *testing.T) {
	exampletestutil.AssertMainContains(t, main, "section-example", "Owned blocks")
}
