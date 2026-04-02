package examples

import (
	"flag"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/evanmschultz/laslig"
)

// Renderer writes one runnable example to the provided writer with one prepared
// printer.
type Renderer func(io.Writer, *laslig.Printer) error

type exampleRenderOptions struct {
	maxWidth    int
	wrapMode    laslig.TableWrapMode
	contentMode string
}

var activeExampleRenderOptions struct {
	sync.RWMutex
	value exampleRenderOptions
}

var defaultExampleRenderOptions = exampleRenderOptions{
	contentMode: "default",
}

func setExampleRenderOptions(opts exampleRenderOptions) {
	activeExampleRenderOptions.Lock()
	defer activeExampleRenderOptions.Unlock()
	activeExampleRenderOptions.value = opts
}

func getExampleRenderOptions() exampleRenderOptions {
	activeExampleRenderOptions.RLock()
	defer activeExampleRenderOptions.RUnlock()
	return activeExampleRenderOptions.value
}

// Run parses the common example flags and renders one shared example.
func Run(out io.Writer, args []string, name string, render Renderer) error {
	setExampleRenderOptions(defaultExampleRenderOptions)
	defer setExampleRenderOptions(defaultExampleRenderOptions)

	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	format := flags.String("format", string(laslig.FormatAuto), "output format: auto, human, plain, json")
	style := flags.String("style", string(laslig.StyleAuto), "style policy: auto, always, never")
	spinnerStyle := flags.String("spinner-style", string(laslig.DefaultSpinnerStyle()), "spinner style: braille, dot, line, pulse, meter")
	glamourStyle := flags.String("glamour-style", string(laslig.DefaultGlamourStyle()), "glamour markdown style: dark, light, pink, dracula, tokyo-night, ascii, notty")
	maxWidth := flags.Int("max-width", 0, "override framed max width for table/panel/codeblock/logblock examples")
	wrapMode := flags.String("wrap-mode", "", "override table-style wrapping for framed examples: auto, truncate, never")
	contentMode := flags.String("content", "default", "example content mode: default, long")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	resolvedSpinnerStyle := laslig.SpinnerStyle(strings.ToLower(*spinnerStyle))
	if !resolvedSpinnerStyle.Valid() {
		return fmt.Errorf("parse flags: invalid spinner style %q", *spinnerStyle)
	}
	resolvedGlamourStyle := laslig.GlamourStyle(strings.ToLower(*glamourStyle))
	if !resolvedGlamourStyle.Valid() {
		return fmt.Errorf("parse flags: invalid glamour style %q", *glamourStyle)
	}
	resolvedWrapMode := laslig.TableWrapMode(strings.ToLower(strings.TrimSpace(*wrapMode)))
	if *wrapMode != "" && resolvedWrapMode != laslig.TableWrapAuto && resolvedWrapMode != laslig.TableWrapNever && resolvedWrapMode != laslig.TableWrapTruncate {
		return fmt.Errorf("parse flags: invalid wrap mode %q", *wrapMode)
	}
	resolvedContentMode := strings.TrimSpace(strings.ToLower(*contentMode))
	switch resolvedContentMode {
	case "", "default", "long":
	default:
		return fmt.Errorf("parse flags: invalid content mode %q", *contentMode)
	}
	if resolvedContentMode == "" {
		resolvedContentMode = "default"
	}
	setExampleRenderOptions(exampleRenderOptions{
		maxWidth:    *maxWidth,
		wrapMode:    resolvedWrapMode,
		contentMode: resolvedContentMode,
	})

	printer := laslig.New(out, laslig.Policy{
		Format:       laslig.Format(strings.ToLower(*format)),
		Style:        laslig.StylePolicy(strings.ToLower(*style)),
		SpinnerStyle: resolvedSpinnerStyle,
		GlamourStyle: resolvedGlamourStyle,
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
