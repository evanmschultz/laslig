package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/evanmschultz/laslig"
)

// main runs the demo command and exits non-zero on failure.
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run parses flags and renders the demo output.
func run() error {
	format := flag.String("format", string(laslig.FormatAuto), "output format: auto, human, plain, json")
	style := flag.String("style", string(laslig.StyleAuto), "style policy: auto, always, never")
	flag.Parse()

	printer := laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.Format(*format),
		Style:  laslig.StylePolicy(*style),
	})

	if err := printer.Section("laslig demo"); err != nil {
		return err
	}
	if err := printer.Notice(laslig.Notice{
		Level: laslig.NoticeInfoLevel,
		Title: "Readable by default",
		Body:  "Structured output should look intentional without forcing you into a framework.",
		Detail: []string{
			"Use laslig for results and diagnostics.",
			"Keep logging and CLI orchestration separate.",
		},
	}); err != nil {
		return err
	}
	if err := printer.Record(laslig.Record{
		Title: "Project",
		Fields: []laslig.Field{
			{Label: "module", Value: "github.com/evanmschultz/laslig", Identifier: true},
			{Label: "runtime deps", Value: "Charm + stdlib", Badge: true},
			{Label: "task runner", Value: "Mage", Muted: true},
		},
	}); err != nil {
		return err
	}
	if err := printer.List(laslig.List{
		Title: "Primitives",
		Items: []laslig.ListItem{
			{
				Title: "Notice",
				Badge: "ready",
				Fields: []laslig.Field{
					{Label: "use", Value: "warnings, successes, errors"},
				},
			},
			{
				Title: "Table",
				Badge: "ready",
				Fields: []laslig.Field{
					{Label: "use", Value: "human summaries and reports"},
				},
			},
			{
				Title: "testjson",
				Badge: "next",
				Fields: []laslig.Field{
					{Label: "use", Value: "Charm-native go test output", Muted: true},
				},
			},
		},
	}); err != nil {
		return err
	}
	if err := printer.Table(laslig.Table{
		Title:  "Formats",
		Header: []string{"format", "goal"},
		Rows: [][]string{
			{"human", "pleasant terminal output"},
			{"plain", "stable no-ANSI text"},
			{"json", "machine-readable payloads"},
		},
		Caption: "One policy, three surfaces.",
	}); err != nil {
		return err
	}
	return printer.Panel(laslig.Panel{
		Title:  "Why this shape",
		Body:   "Fang should keep owning help and command errors. Laslig should own ordinary output blocks and structured summaries.",
		Footer: "Next up: go test -json rendering.",
	})
}
