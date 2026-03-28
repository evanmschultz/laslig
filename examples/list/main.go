package main

import (
	"io"
	"os"

	demoexamples "github.com/evanmschultz/laslig/internal/examples"
)

// main runs the focused list example and exits non-zero on failure.
func main() {
	execute(os.Stdout, os.Stderr, os.Args[1:], os.Exit)
}

// execute renders the focused list example and reports errors through one writer.
func execute(out io.Writer, errOut io.Writer, args []string, exitFn func(int)) {
	demoexamples.Main(out, errOut, args, exitFn, "list-example", demoexamples.RenderList)
}

// runArgs renders the focused list example to one writer.
func runArgs(out io.Writer, args []string) error {
	return demoexamples.Run(out, args, "list-example", demoexamples.RenderList)
}
