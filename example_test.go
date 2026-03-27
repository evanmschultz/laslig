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
