package main

import (
	"io"
	"os"

	demoexamples "github.com/evanmschultz/laslig/internal/examples"
)

// main runs the focused notice example and exits non-zero on failure.
func main() {
	execute(os.Stdout, os.Stderr, os.Args[1:], os.Exit)
}

// execute renders the focused notice example and reports errors through one writer.
func execute(out io.Writer, errOut io.Writer, args []string, exitFn func(int)) {
	demoexamples.Main(out, errOut, args, exitFn, "notice-example", demoexamples.RenderNotice)
}

// runArgs renders the focused notice example to one writer.
func runArgs(out io.Writer, args []string) error {
	return demoexamples.Run(out, args, "notice-example", demoexamples.RenderNotice)
}
