package main

import (
	"io"
	"os"

	demoexamples "github.com/evanmschultz/laslig/internal/examples"
)

// main runs the aggregate demo command and exits non-zero on failure.
func main() {
	execute(os.Stdout, os.Stderr, os.Args[1:], os.Exit)
}

// execute renders the aggregate walkthrough and reports errors through one writer.
func execute(out io.Writer, errOut io.Writer, args []string, exitFn func(int)) {
	demoexamples.Main(out, errOut, args, exitFn, "laslig-demo", demoexamples.RenderAll)
}

// runArgs renders the aggregate walkthrough to one writer.
func runArgs(out io.Writer, args []string) error {
	return demoexamples.Run(out, args, "laslig-demo", demoexamples.RenderAll)
}
