package main

import (
	"testing"

	"github.com/evanmschultz/laslig/internal/exampletestutil"
)

func TestRunArgsPlain(t *testing.T) {
	exampletestutil.AssertRunArgsPlainContains(t, runArgs, "Spinner", "[RUNNING] Waiting for remote rollout", "[SUCCESS] Rollout ready")
}

func TestRunArgsHumanStyled(t *testing.T) {
	exampletestutil.AssertRunArgsHumanStyled(t, runArgs, "Use Spinner when long-running work might otherwise stay quiet", "[RUNNING] Waiting for remote rollout")
}

func TestRunArgsInvalidFlag(t *testing.T) {
	exampletestutil.AssertRunArgsInvalidFlag(t, runArgs)
}

func TestMain(t *testing.T) {
	exampletestutil.AssertMainContains(t, main, "spinner-example", "Spinner", "Rollout ready")
}
