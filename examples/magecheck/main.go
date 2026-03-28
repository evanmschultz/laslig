package main

import (
	"io"
	"os"

	demoexamples "github.com/evanmschultz/laslig/internal/examples"
)

// main runs the focused Mage-style integration example and exits non-zero on failure.
func main() {
	execute(os.Stdout, os.Stderr, os.Args[1:], os.Exit)
}

// execute renders the focused Mage-style integration example and reports errors through one writer.
func execute(out io.Writer, errOut io.Writer, args []string, exitFn func(int)) {
	demoexamples.Main(out, errOut, args, exitFn, "magecheck-example", demoexamples.RenderMageCheckPreview)
}

// runArgs renders the focused Mage-style integration example to one writer.
func runArgs(out io.Writer, args []string) error {
	return demoexamples.Run(out, args, "magecheck-example", demoexamples.RenderMageCheckPreview)
}
