package main

import (
	"io"
	"os"

	demoexamples "github.com/evanmschultz/laslig/internal/examples"
)

// main runs the focused section example and exits non-zero on failure.
func main() {
	execute(os.Stdout, os.Stderr, os.Args[1:], os.Exit)
}

// execute renders the focused section example and reports errors through one writer.
func execute(out io.Writer, errOut io.Writer, args []string, exitFn func(int)) {
	demoexamples.Main(out, errOut, args, exitFn, "section-example", demoexamples.RenderSection)
}

// runArgs renders the focused section example to one writer.
func runArgs(out io.Writer, args []string) error {
	return demoexamples.Run(out, args, "section-example", demoexamples.RenderSection)
}
