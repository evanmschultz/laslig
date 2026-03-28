package main

import (
	"testing"

	"github.com/evanmschultz/laslig/internal/exampletestutil"
)

func TestRunArgsPlain(t *testing.T) {
	exampletestutil.AssertRunArgsPlainContains(t, runArgs, "StatusLine", "[SUCCESS] Use StatusLine for one compact semantic result")
}

func TestRunArgsHumanStyled(t *testing.T) {
	exampletestutil.AssertRunArgsHumanStyled(t, runArgs, "cache hit, package passed")
}

func TestRunArgsInvalidFlag(t *testing.T) {
	exampletestutil.AssertRunArgsInvalidFlag(t, runArgs)
}

func TestMain(t *testing.T) {
	exampletestutil.AssertMainContains(t, main, "statusline-example", "Use StatusLine for one compact semantic result")
}
