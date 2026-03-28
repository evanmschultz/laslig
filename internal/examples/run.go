package examples

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/evanmschultz/laslig"
)

// Renderer writes one runnable example to the provided writer with one prepared
// printer.
type Renderer func(io.Writer, *laslig.Printer) error

// Run parses the common example flags and renders one shared example.
func Run(out io.Writer, args []string, name string, render Renderer) error {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	format := flags.String("format", string(laslig.FormatAuto), "output format: auto, human, plain, json")
	style := flags.String("style", string(laslig.StyleAuto), "style policy: auto, always, never")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	printer := laslig.New(out, laslig.Policy{
		Format: laslig.Format(strings.ToLower(*format)),
		Style:  laslig.StylePolicy(strings.ToLower(*style)),
	})
	if err := render(out, printer); err != nil {
		return fmt.Errorf("render %s example: %w", name, err)
	}
	return nil
}

// Main runs one focused example, reporting any error and delegating exit handling.
func Main(out io.Writer, errOut io.Writer, args []string, exitFn func(int), name string, render Renderer) {
	if err := Run(out, args, name, render); err != nil {
		fmt.Fprintln(errOut, err)
		exitFn(1)
	}
}

// StylePolicyForMode converts one resolved printer mode back into a style
// policy for helpers that only accept a Policy.
func StylePolicyForMode(mode laslig.Mode) laslig.StylePolicy {
	if mode.Styled {
		return laslig.StyleAlways
	}
	return laslig.StyleNever
}
