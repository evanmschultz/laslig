package main

import (
	"io"
	"os"

	demoexamples "github.com/evanmschultz/laslig/internal/examples"
)

// main runs the focused spinner example and exits non-zero on failure.
func main() {
	execute(os.Stdout, os.Stderr, os.Args[1:], os.Exit)
}

// execute renders the focused spinner example and reports errors through one writer.
func execute(out io.Writer, errOut io.Writer, args []string, exitFn func(int)) {
	demoexamples.Main(out, errOut, args, exitFn, "spinner-example", demoexamples.RenderSpinner)
}

// runArgs renders the focused spinner example to one writer.
func runArgs(out io.Writer, args []string) error {
	return demoexamples.Run(out, args, "spinner-example", demoexamples.RenderSpinner)
}
