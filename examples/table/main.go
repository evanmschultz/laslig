package main

import (
	"io"
	"os"

	demoexamples "github.com/evanmschultz/laslig/internal/examples"
)

// main runs the focused table example and exits non-zero on failure.
func main() {
	execute(os.Stdout, os.Stderr, os.Args[1:], os.Exit)
}

// execute renders the focused table example and reports errors through one writer.
func execute(out io.Writer, errOut io.Writer, args []string, exitFn func(int)) {
	demoexamples.Main(out, errOut, args, exitFn, "table-example", demoexamples.RenderTable)
}

// runArgs renders the focused table example to one writer.
func runArgs(out io.Writer, args []string) error {
	return demoexamples.Run(out, args, "table-example", demoexamples.RenderTable)
}
