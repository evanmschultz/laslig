package laslig_test

import (
	"os"

	"github.com/evanmschultz/laslig"
)

// ExamplePrinter_Record shows one simple record render with plain output.
func ExamplePrinter_Record() {
	printer := laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.FormatPlain,
		Style:  laslig.StyleNever,
	})

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

// ExamplePrinter_KV shows aligned key-value rendering.
func ExamplePrinter_KV() {
	printer := laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.FormatPlain,
		Style:  laslig.StyleNever,
	})

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
	printer := laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.FormatPlain,
		Style:  laslig.StyleNever,
	})

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

// ExamplePrinter_Table shows one plain table render for stable text output.
func ExamplePrinter_Table() {
	printer := laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.FormatPlain,
		Style:  laslig.StyleNever,
	})

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
	printer := laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.FormatPlain,
		Style:  laslig.StyleNever,
	})

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

// ExamplePrinter_StatusLine shows one compact status row in plain mode.
func ExamplePrinter_StatusLine() {
	printer := laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.FormatPlain,
		Style:  laslig.StyleNever,
	})

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
	printer := laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.FormatPlain,
		Style:  laslig.StyleNever,
	})

	_ = printer.Markdown(laslig.Markdown{
		Body: "# Notes\n\n- first\n- second",
	})

	// Output:
	// # Notes
	//
	// - first
	// - second
}

// ExamplePrinter_LogBlock shows one plain boxed-log surface without owning logging.
func ExamplePrinter_LogBlock() {
	printer := laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.FormatPlain,
		Style:  laslig.StyleNever,
	})

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
