package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/evanmschultz/laslig"
	"github.com/evanmschultz/laslig/gotestout"
)

// sampleStream exercises passing, skipped, failing, and build-error events.
const sampleStream = `{"Action":"run","Package":"example/pkg","Test":"TestPass"}
{"Action":"output","Package":"example/pkg","Test":"TestPass","Output":"=== RUN   TestPass\n"}
{"Action":"output","Package":"example/pkg","Test":"TestPass","Output":"note: useful output\n"}
{"Action":"output","Package":"example/pkg","Test":"TestPass","Output":"--- PASS: TestPass (0.01s)\n"}
{"Action":"pass","Package":"example/pkg","Test":"TestPass","Elapsed":0.01}
{"Action":"run","Package":"example/pkg","Test":"TestSkip"}
{"Action":"output","Package":"example/pkg","Test":"TestSkip","Output":"--- SKIP: TestSkip (0.00s)\n"}
{"Action":"skip","Package":"example/pkg","Test":"TestSkip","Elapsed":0}
{"Action":"run","Package":"example/pkg","Test":"TestFail"}
{"Action":"output","Package":"example/pkg","Test":"TestFail","Output":"main_test.go:42: expected boom\n"}
{"Action":"output","Package":"example/pkg","Test":"TestFail","Output":"--- FAIL: TestFail (0.02s)\n"}
{"Action":"fail","Package":"example/pkg","Test":"TestFail","Elapsed":0.02}
{"Action":"output","Package":"example/pkg","Output":"FAIL\texample/pkg [build failed]\n","FailedBuild":"example/pkg"}
{"Action":"fail","Package":"example/pkg","Elapsed":0.03}
`

// main runs the focused gotestout example and exits non-zero on failure.
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run renders the example with process arguments.
func run() error {
	return runArgs(os.Stdout, os.Args[1:])
}

// runArgs renders the focused gotestout example to one writer.
func runArgs(out io.Writer, args []string) error {
	flags := flag.NewFlagSet("gotestout-example", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	format := flags.String("format", string(laslig.FormatAuto), "output format: auto, human, plain, json")
	style := flags.String("style", string(laslig.StyleAuto), "style policy: auto, always, never")
	view := flags.String("view", string(gotestout.ViewCompact), "view: compact, detailed")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	_, err := gotestout.Render(out, strings.NewReader(sampleStream), gotestout.Options{
		Policy: laslig.Policy{
			Format: laslig.Format(*format),
			Style:  laslig.StylePolicy(*style),
		},
		View: gotestout.View(*view),
	})
	if err != nil {
		return fmt.Errorf("render gotestout example: %w", err)
	}
	return nil
}
