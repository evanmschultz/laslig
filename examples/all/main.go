package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/evanmschultz/laslig"
	loggingexample "github.com/evanmschultz/laslig/examples/logging"
	"github.com/evanmschultz/laslig/gotestout"
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
	return renderShowcase(out, printer)
}

// renderShowcase renders the all-in-one Läslig walkthrough with one prepared printer.
func renderShowcase(out io.Writer, printer *laslig.Printer) error {
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
		{
			name: "specialized packages section",
			render: func() error {
				return printer.Section("Specialized Packages")
			},
		},
		{
			name: "gotestout intro",
			render: func() error {
				return printer.Paragraph(laslig.Paragraph{
					Title:  "gotestout",
					Body:   "Use gotestout for Charm-native go test output when your task runner, CLI command, or Go helper behind make/just should keep owning process control.",
					Footer: "Try go run ./examples/gotestout --format human --style always first, or mage test for the real task-runner path. The output below shows the same Build, Tests, and Coverage shape this repository prints through mage check.",
				})
			},
		},
		{
			name: "gotestout output",
			render: func() error {
				return renderMageCheckShowcase(out, printer)
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

// renderMageCheckShowcase renders one Mage-style flow that uses gotestout for
// the test stream and the core laslig primitives for build and coverage.
func renderMageCheckShowcase(out io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("Build"); err != nil {
		return fmt.Errorf("render build section: %w", err)
	}
	if err := printer.StatusLine(laslig.StatusLine{
		Level:  laslig.NoticeInfoLevel,
		Text:   "Building showcase example",
		Detail: "./examples/all",
	}); err != nil {
		return fmt.Errorf("render build start: %w", err)
	}
	if err := printer.StatusLine(laslig.StatusLine{
		Level:  laslig.NoticeSuccessLevel,
		Text:   "Built showcase example",
		Detail: "bin/laslig-demo",
	}); err != nil {
		return fmt.Errorf("render build success: %w", err)
	}

	if err := printer.Section("Tests"); err != nil {
		return fmt.Errorf("render tests section: %w", err)
	}
	if _, err := gotestout.Render(out, strings.NewReader(mageCheckSampleStream()), gotestout.Options{
		Policy: laslig.Policy{
			Format: printer.Mode().Format,
			Style:  stylePolicyForMode(printer.Mode()),
		},
		View: gotestout.ViewCompact,
	}); err != nil {
		return fmt.Errorf("render gotestout stream: %w", err)
	}

	if err := printer.Section("Coverage"); err != nil {
		return fmt.Errorf("render coverage section: %w", err)
	}
	if err := printer.Table(laslig.Table{
		Header: []string{"package", "cover"},
		Rows: [][]string{
			{"github.com/evanmschultz/laslig", "71.2%"},
			{"github.com/evanmschultz/laslig/examples/all", "84.4%"},
			{"github.com/evanmschultz/laslig/examples/gotestout", "86.7%"},
			{"github.com/evanmschultz/laslig/internal/layout", "87.5%"},
			{"github.com/evanmschultz/laslig/internal/glamrender", "83.3%"},
			{"github.com/evanmschultz/laslig/examples/logging", "75.0%"},
			{"github.com/evanmschultz/laslig/internal/table", "96.6%"},
			{"github.com/evanmschultz/laslig/gotestout", "82.3%"},
		},
		Caption: "Minimum package coverage: 70.0%.",
	}); err != nil {
		return fmt.Errorf("render coverage table: %w", err)
	}
	if err := printer.Notice(laslig.Notice{
		Level: laslig.NoticeSuccessLevel,
		Title: "Coverage threshold met",
		Body:  "All packages are at or above 70.0% coverage.",
	}); err != nil {
		return fmt.Errorf("render coverage success notice: %w", err)
	}
	return nil
}

// mageCheckSampleStream returns one deterministic passing test stream that
// mirrors the package-level shape this repository prints through mage test.
func mageCheckSampleStream() string {
	packages := []struct {
		name  string
		tests int
	}{
		{"github.com/evanmschultz/laslig", 12},
		{"github.com/evanmschultz/laslig/examples/all", 8},
		{"github.com/evanmschultz/laslig/examples/gotestout", 8},
		{"github.com/evanmschultz/laslig/internal/layout", 6},
		{"github.com/evanmschultz/laslig/internal/glamrender", 6},
		{"github.com/evanmschultz/laslig/examples/logging", 5},
		{"github.com/evanmschultz/laslig/internal/table", 12},
		{"github.com/evanmschultz/laslig/gotestout", 14},
	}

	var stream strings.Builder
	for _, pkg := range packages {
		for index := 1; index <= pkg.tests; index++ {
			fmt.Fprintf(&stream, "{\"Action\":\"pass\",\"Package\":\"%s\",\"Test\":\"Test%02d\",\"Elapsed\":0}\n", pkg.name, index)
		}
		fmt.Fprintf(&stream, "{\"Action\":\"pass\",\"Package\":\"%s\",\"Elapsed\":0}\n", pkg.name)
	}
	return stream.String()
}

// stylePolicyForMode converts one resolved mode back into a policy choice for
// renderers that only accept Policy input.
func stylePolicyForMode(mode laslig.Mode) laslig.StylePolicy {
	if mode.Styled {
		return laslig.StyleAlways
	}
	return laslig.StyleNever
}
