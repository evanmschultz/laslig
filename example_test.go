package laslig_test

import (
	"os"

	"github.com/evanmschultz/laslig"
)

// newExamplePrinter constructs one plain printer without the default leading gap
// so package examples stay compact in rendered Go docs.
func newExamplePrinter() *laslig.Printer {
	layout := laslig.DefaultLayout().WithLeadingGap(0)
	return laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.FormatPlain,
		Style:  laslig.StyleNever,
		Layout: &layout,
	})
}

// ExamplePrinter_Record shows one simple record render with plain output.
func ExamplePrinter_Record() {
	printer := newExamplePrinter()

	_ = printer.Record(laslig.Record{
		Title: "Build",
		Fields: []laslig.Field{
			{Label: "status", Value: "pass", Badge: true},
			{Label: "runner", Value: "mage", Muted: true},
		},
	})

	// Output:
	// Build
	//   status: [PASS]
	//   runner: mage
}

// ExamplePrinter_Section shows one section heading owning the indentation of
// the blocks that follow it.
func ExamplePrinter_Section() {
	printer := newExamplePrinter()

	_ = printer.Section("Overview")
	_ = printer.Record(laslig.Record{
		Title: "Record",
		Fields: []laslig.Field{
			{Label: "what", Value: "Use Section when later blocks should clearly read as one document group."},
		},
	})

	// Output:
	// Overview
	//
	//   Record
	//     what: Use Section when later blocks should clearly read as one document group.
}

// ExamplePrinter_KV shows aligned key-value rendering.
func ExamplePrinter_KV() {
	printer := newExamplePrinter()

	_ = printer.KV(laslig.KV{
		Title: "Config",
		Pairs: []laslig.Field{
			{Label: "module", Value: "github.com/evanmschultz/laslig", Identifier: true},
			{Label: "style", Value: "auto", Badge: true},
			{Label: "runner", Value: "mage", Muted: true},
		},
	})

	// Output:
	// Config
	//   module  github.com/evanmschultz/laslig
	//   style   [AUTO]
	//   runner  mage
}

// ExamplePrinter_Notice shows one warning notice rendered without ANSI styling.
func ExamplePrinter_Notice() {
	printer := newExamplePrinter()

	_ = printer.Notice(laslig.Notice{
		Level: laslig.NoticeWarningLevel,
		Title: "Coverage dropped",
		Body:  "Package coverage fell below the configured threshold.",
		Detail: []string{
			"Previous: 84.2%",
			"Current:  79.8%",
		},
	})

	// Output:
	// [WARNING] Coverage dropped
	//   Package coverage fell below the configured threshold.
	//   Previous: 84.2%
	//   Current:  79.8%
}

// ExamplePrinter_List shows grouped list items with badges and detail fields.
func ExamplePrinter_List() {
	printer := newExamplePrinter()

	_ = printer.List(laslig.List{
		Title: "Targets",
		Items: []laslig.ListItem{
			{
				Title: "check",
				Badge: "ready",
				Fields: []laslig.Field{
					{Label: "when", Value: "Run verification before handoff."},
				},
			},
			{
				Title: "demo",
				Badge: "live",
				Fields: []laslig.Field{
					{Label: "what", Value: "Show the all-in-one walkthrough."},
				},
			},
		},
	})

	// Output:
	// Targets
	// - check [READY]
	//   when: Run verification before handoff.
	// - demo [LIVE]
	//   what: Show the all-in-one walkthrough.
}

// ExamplePrinter_Table shows one plain table render for stable text output.
func ExamplePrinter_Table() {
	printer := newExamplePrinter()

	_ = printer.Table(laslig.Table{
		Title:  "Targets",
		Header: []string{"name", "status"},
		Rows: [][]string{
			{"check", "ready"},
			{"demo", "ready"},
		},
		Caption: "One policy, three surfaces.",
	})

	// Output:
	// Targets
	// name  | status
	// ------+-------
	// check | ready
	// demo  | ready
	// One policy, three surfaces.
}

// ExamplePrinter_Paragraph shows one wrapped long-form text block in plain mode.
func ExamplePrinter_Paragraph() {
	printer := newExamplePrinter()

	_ = printer.Paragraph(laslig.Paragraph{
		Title:  "Why",
		Body:   "Laslig keeps ordinary command output readable without forcing a framework.",
		Footer: "Writers in, errors out.",
	})

	// Output:
	// Why
	//
	// Laslig keeps ordinary command output readable without forcing a framework.
	//
	// Writers in, errors out.
}

// ExamplePrinter_Panel shows one stronger callout block in plain mode.
func ExamplePrinter_Panel() {
	printer := newExamplePrinter()

	_ = printer.Panel(laslig.Panel{
		Title:  "Next step",
		Body:   "Run mage check before pushing.",
		Footer: "Use Panel when the note should stand apart from the surrounding document.",
	})

	// Output:
	// Next step
	//
	// Run mage check before pushing.
	//
	// Use Panel when the note should stand apart from the surrounding document.
}

// ExamplePrinter_StatusLine shows one compact status row in plain mode.
func ExamplePrinter_StatusLine() {
	printer := newExamplePrinter()

	_ = printer.StatusLine(laslig.StatusLine{
		Level:  laslig.NoticeSuccessLevel,
		Text:   "Build ready",
		Detail: "mage check",
	})

	// Output:
	// [SUCCESS] Build ready (mage check)
}

// ExamplePrinter_Markdown shows one Markdown block rendered as source in plain mode.
func ExamplePrinter_Markdown() {
	printer := newExamplePrinter()

	_ = printer.Markdown(laslig.Markdown{
		Body: "# Notes\n\n- first\n- second",
	})

	// Output:
	// # Notes
	//
	// - first
	// - second
}

// ExamplePrinter_CodeBlock shows one plain code block without ANSI styling.
func ExamplePrinter_CodeBlock() {
	printer := newExamplePrinter()

	_ = printer.CodeBlock(laslig.CodeBlock{
		Title:    "Snippet",
		Language: "go",
		Body:     "fmt.Println(\"hello from laslig\")",
		Footer:   "Use CodeBlock for commands or code samples.",
	})

	// Output:
	// Snippet
	//
	// fmt.Println("hello from laslig")
	//
	// Use CodeBlock for commands or code samples.
}

// ExamplePrinter_LogBlock shows one plain boxed-log surface without owning logging.
func ExamplePrinter_LogBlock() {
	printer := newExamplePrinter()

	_ = printer.LogBlock(laslig.LogBlock{
		Title: "stderr excerpt",
		Body:  "INFO boot complete\nWARN retry scheduled",
	})

	// Output:
	// stderr excerpt
	//
	// INFO boot complete
	// WARN retry scheduled
}
