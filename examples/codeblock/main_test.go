package main

import (
	"testing"

	"github.com/evanmschultz/laslig/internal/exampletestutil"
)

func TestRunArgsPlain(t *testing.T) {
	exampletestutil.AssertRunArgsPlainContains(t, runArgs, "CodeBlock", "Go snippet", "fmt.Println(\"hello from laslig\")")
}

func TestRunArgsHumanStyled(t *testing.T) {
	exampletestutil.AssertRunArgsHumanStyled(t, runArgs, "Use CodeBlock when code should stay visibly distinct from prose")
}

func TestRunArgsInvalidFlag(t *testing.T) {
	exampletestutil.AssertRunArgsInvalidFlag(t, runArgs)
}

func TestMain(t *testing.T) {
	exampletestutil.AssertMainContains(t, main, "codeblock-example", "Go snippet")
}
