package main

import (
	"io"
	"os"

	demoexamples "github.com/evanmschultz/laslig/internal/examples"
)

// main runs the focused codeblock example and exits non-zero on failure.
func main() {
	execute(os.Stdout, os.Stderr, os.Args[1:], os.Exit)
}

// execute renders the focused codeblock example and reports errors through one writer.
func execute(out io.Writer, errOut io.Writer, args []string, exitFn func(int)) {
	demoexamples.Main(out, errOut, args, exitFn, "codeblock-example", demoexamples.RenderCodeBlock)
}

// runArgs renders the focused codeblock example to one writer.
func runArgs(out io.Writer, args []string) error {
	return demoexamples.Run(out, args, "codeblock-example", demoexamples.RenderCodeBlock)
}
