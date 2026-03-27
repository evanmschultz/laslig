package main

import (
	"flag"
	"fmt"
	"io"
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
	return runArgs(os.Stdout, os.Args[1:])
}

// runArgs parses flags and renders the demo output to one writer.
func runArgs(out io.Writer, args []string) error {
	flags := flag.NewFlagSet("laslig-demo", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	format := flags.String("format", string(laslig.FormatAuto), "output format: auto, human, plain, json")
	style := flags.String("style", string(laslig.StyleAuto), "style policy: auto, always, never")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	printer := laslig.New(out, laslig.Policy{
		Format: laslig.Format(*format),
		Style:  laslig.StylePolicy(*style),
	})

	steps := []struct {
		name   string
		render func() error
	}{
		{
			name: "section",
			render: func() error {
				return printer.Section("laslig demo")
			},
		},
		{
			name: "notice",
			render: func() error {
				return printer.Notice(laslig.Notice{
					Level: laslig.NoticeInfoLevel,
					Title: "Readable by default",
					Body:  "Structured output should look intentional without forcing a framework.",
					Detail: []string{
						"Use laslig for results and diagnostics.",
						"Keep logging and CLI orchestration separate.",
					},
				})
			},
		},
		{
			name: "record",
			render: func() error {
				return printer.Record(laslig.Record{
					Title: "Project",
					Fields: []laslig.Field{
						{Label: "module", Value: "github.com/evanmschultz/laslig", Identifier: true},
						{Label: "runtime deps", Value: "Charm + stdlib"},
						{Label: "task runner", Value: "Mage", Muted: true},
					},
				})
			},
		},
		{
			name: "kv",
			render: func() error {
				return printer.KV(laslig.KV{
					Title: "Policy",
					Pairs: []laslig.Field{
						{Label: "format", Value: string(printer.Mode().Format)},
						{Label: "styled", Value: fmt.Sprintf("%t", printer.Mode().Styled), Muted: true},
					},
				})
			},
		},
		{
			name: "list",
			render: func() error {
				return printer.List(laslig.List{
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
							Badge: "live",
							Fields: []laslig.Field{
								{Label: "use", Value: "Charm-native go test output", Muted: true},
							},
						},
					},
				})
			},
		},
		{
			name: "table",
			render: func() error {
				return printer.Table(laslig.Table{
					Title:  "Formats",
					Header: []string{"format", "goal"},
					Rows: [][]string{
						{"human", "pleasant terminal output"},
						{"plain", "stable no-ANSI text"},
						{"json", "machine-readable payloads"},
					},
					Caption: "One policy, three surfaces.",
				})
			},
		},
		{
			name: "panel",
			render: func() error {
				return printer.Panel(laslig.Panel{
					Title:  "Why this shape",
					Body:   "Fang should own help and command errors. Laslig should own ordinary output blocks and summaries.",
					Footer: "Next up: richer test summaries and more small helpers.",
				})
			},
		},
	}
	for _, step := range steps {
		if err := step.render(); err != nil {
			return fmt.Errorf("render %s: %w", step.name, err)
		}
	}
	return nil
}
