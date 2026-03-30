package main

import (
	"testing"

	"github.com/evanmschultz/laslig/internal/exampletestutil"
)

func TestRunArgsPlain(t *testing.T) {
	exampletestutil.AssertRunArgsPlainContains(t, runArgs, "gotestout + Mage", "[RUNNING] Waiting for first test event", "[SUCCESS] Test stream detected", "Coverage threshold met")
}

func TestRunArgsHumanStyled(t *testing.T) {
	exampletestutil.AssertRunArgsHumanStyled(t, runArgs, "Waiting for first test event", "Test stream detected", "All tests passed")
}

func TestRunArgsInvalidFlag(t *testing.T) {
	exampletestutil.AssertRunArgsInvalidFlag(t, runArgs)
}

func TestMain(t *testing.T) {
	exampletestutil.AssertMainContains(t, main, "magecheck-example", "Test stream detected", "Coverage threshold met")
}
