package main

import (
	"testing"

	"github.com/evanmschultz/laslig/internal/exampletestutil"
)

func TestRunArgsPlain(t *testing.T) {
	exampletestutil.AssertRunArgsPlainContains(t, runArgs, "gotestout + Mage", "[INFO] Started go test -json (./...)", "Coverage threshold met")
}

func TestRunArgsHumanStyled(t *testing.T) {
	exampletestutil.AssertRunArgsHumanStyled(t, runArgs, "Started go test -json", "All tests passed")
}

func TestRunArgsInvalidFlag(t *testing.T) {
	exampletestutil.AssertRunArgsInvalidFlag(t, runArgs)
}

func TestMain(t *testing.T) {
	exampletestutil.AssertMainContains(t, main, "magecheck-example", "Started go test -json", "Coverage threshold met")
}
