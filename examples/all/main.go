package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/evanmschultz/laslig"
	loggingexample "github.com/evanmschultz/laslig/examples/logging"
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
	return renderShowcase(printer)
}

// renderShowcase renders the all-in-one Läslig walkthrough with one prepared printer.
func renderShowcase(printer *laslig.Printer) error {
	steps := []struct {
		name   string
		render func() error
	}{
		{
			name: "section",
			render: func() error {
				return printer.Section("Läslig demo")
			},
		},
		{
			name: "notice",
			render: func() error {
				return printer.Notice(laslig.Notice{
					Level: laslig.NoticeInfoLevel,
					Title: "A guided primitive walkthrough",
					Body:  "This showcase names Läslig's public primitives directly so you can see what they render and when to reach for them.",
					Detail: []string{
						"Use Läslig for ordinary CLI output, not command parsing or application logging.",
						"Think of it as the layer between raw Charm styles and your command's result text.",
					},
				})
			},
		},
		{
			name: "structured section",
			render: func() error {
				return printer.Section("Structured Primitives")
			},
		},
		{
			name: "record",
			render: func() error {
				return printer.Record(laslig.Record{
					Title: "Record",
					Fields: []laslig.Field{
						{Label: "what", Value: "One object or result rendered as labeled facts."},
						{Label: "when", Value: "Build metadata, artifact details, environment summaries."},
						{Label: "example", Value: "module github.com/evanmschultz/laslig", Identifier: true},
					},
				})
			},
		},
		{
			name: "kv",
			render: func() error {
				return printer.KV(laslig.KV{
					Title: "KV",
					Pairs: []laslig.Field{
						{Label: "what", Value: "Compact aligned configuration or status."},
						{Label: "when", Value: "Resolved policy, flags, execution mode."},
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
					Title: "List",
					Items: []laslig.ListItem{
						{
							Title: "Grouped items",
							Badge: "default",
							Fields: []laslig.Field{
								{Label: "when", Value: "Packages, tasks, phases, or capabilities that scan better than a table."},
							},
						},
						{
							Title: "Badges stay lightweight",
							Badge: "ready",
							Fields: []laslig.Field{
								{Label: "why", Value: "Use a badge when a quick state is enough and a panel would be too heavy."},
							},
						},
						{
							Title: "Detail fields add context",
							Badge: "live",
							Fields: []laslig.Field{
								{Label: "what", Value: "Each item can carry labeled facts and still stay list-shaped.", Muted: true},
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
					Title:  "Table",
					Header: []string{"compare", "prefer when"},
					Rows: [][]string{
						{"Table", "column alignment matters across many rows"},
						{"List", "items are unordered and short"},
						{"Record", "you are describing one object"},
					},
					Caption: "Use Table when comparison matters more than prose.",
				})
			},
		},
		{
			name: "panel",
			render: func() error {
				return printer.Panel(laslig.Panel{
					Title:  "Panel",
					Body:   "Use Panel for rationale, next steps, and larger callouts that should stand apart from the rest of the document.",
					Footer: "Panels are stronger than Paragraph and lighter than inventing a custom layout.",
				})
			},
		},
		{
			name: "rich text section",
			render: func() error {
				return printer.Section("Rich Text Primitives")
			},
		},
		{
			name: "paragraph",
			render: func() error {
				return printer.Paragraph(laslig.Paragraph{
					Title:  "Paragraph",
					Body:   "Use Paragraph for readable rationale, release context, and longer help text when a Notice or Panel would be too heavy.",
					Footer: "Start the rich-text surface here when your CLI needs to explain, teach, or provide context.",
				})
			},
		},
		{
			name: "status line",
			render: func() error {
				return printer.StatusLine(laslig.StatusLine{
					Level:  laslig.NoticeSuccessLevel,
					Text:   "StatusLine keeps one result compact and semantic",
					Detail: "When: build ready, cache hit, package passed",
				})
			},
		},
		{
			name: "markdown intro",
			render: func() error {
				return printer.Paragraph(laslig.Paragraph{
					Title:  "Markdown",
					Body:   "Use Markdown when your CLI already has release notes, changelog text, or generated help content to show.",
					Footer: "The block below is real terminal-rendered Markdown.",
				})
			},
		},
		{
			name: "markdown",
			render: func() error {
				return printer.Markdown(laslig.Markdown{
					Body: "# Release Notes\n\n## Highlights\n\n- one renderer\n- three output surfaces\n- caller-owned logging",
				})
			},
		},
		{
			name: "code block intro",
			render: func() error {
				return printer.Paragraph(laslig.Paragraph{
					Title:  "CodeBlock",
					Body:   "Use CodeBlock for commands, snippets, generated files, and config examples.",
					Footer: "The block below shows a Go snippet rendered through Glamour.",
				})
			},
		},
		{
			name: "code block",
			render: func() error {
				return printer.CodeBlock(laslig.CodeBlock{
					Title:    "Go snippet",
					Language: "go",
					Body:     "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello from laslig\")\n}",
					Footer:   "Use CodeBlock when code should stay visibly distinct from prose.",
				})
			},
		},
		{
			name: "log block intro",
			render: func() error {
				return printer.Paragraph(laslig.Paragraph{
					Title:  "LogBlock",
					Body:   "Use LogBlock for selected stderr or log excerpts while the application keeps owning logging.",
					Footer: "The block below captures real charm/log output and renders it through Läslig.",
				})
			},
		},
		{
			name: "log block",
			render: func() error {
				return printer.LogBlock(loggingexample.Block())
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
