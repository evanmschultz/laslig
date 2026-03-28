package main

import (
	"io"
	"os"

	demoexamples "github.com/evanmschultz/laslig/internal/examples"
)

// main runs the focused record example and exits non-zero on failure.
func main() {
	execute(os.Stdout, os.Stderr, os.Args[1:], os.Exit)
}

// execute renders the focused record example and reports errors through one writer.
func execute(out io.Writer, errOut io.Writer, args []string, exitFn func(int)) {
	demoexamples.Main(out, errOut, args, exitFn, "record-example", demoexamples.RenderRecord)
}

// runArgs renders the focused record example to one writer.
func runArgs(out io.Writer, args []string) error {
	return demoexamples.Run(out, args, "record-example", demoexamples.RenderRecord)
}
