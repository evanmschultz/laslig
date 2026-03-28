# laslig

`laslig` helps Go CLIs print structured, human-readable output with Charm-native styling and Go-idiomatic ergonomics.

The name comes from the Swedish `läslig`, meaning `legible`.

![laslig demo](docs/vhs/showcase.gif)

## Why

Charm already gives Go developers strong building blocks:

- Lip Gloss for styling and layout
- Fang for help, usage, and CLI error presentation

What is still missing is a narrow, reusable layer for ordinary command output: results, notices, summaries, tables, and diagnostics that should look intentional without forcing an application into a framework.

`laslig` is that layer.

## Status

The first core wave is live. Today the package includes:

- output policy and mode resolution
- a `Printer`
- sections
- notices and diagnostics
- records and lists
- aligned key-value blocks
- paragraph blocks
- compact status lines
- tables
- panels and boxes
- Glamour-backed Markdown blocks
- Glamour-backed code blocks
- boxed log/transcript blocks for caller-provided output
- compact and detailed `testjson` rendering for `go test -json`
- caller-tunable `testjson` summary and output sections

The next wave is focused on theme configuration, deeper `testjson` classification, and tightening the docs/examples further.

## Principles

- small, composable helpers instead of a framework
- writers in, errors out
- no hidden process control
- Charm-native output without depending on Fang or `charm/log`
- easy adoption in Fang, Cobra, Mage, and plain Go commands
- explicit rendering of caller-provided log excerpts without becoming a logger

## Non-Goals

- replacing application logging
- replacing command frameworks
- shipping interactive prompt widgets in v1
- becoming a kitchen-sink terminal toolkit

## Install

```bash
go get github.com/evanmschultz/laslig
```

## Quick Start

```go
package main

import (
	"os"

	"github.com/evanmschultz/laslig"
)

func main() {
	printer := laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.FormatAuto,
		Style:  laslig.StyleAuto,
	})

	_ = printer.Section("release")
	_ = printer.Notice(laslig.Notice{
		Level: laslig.NoticeSuccessLevel,
		Title: "All checks passed",
		Body:  "The CLI can now print structured output with one small helper.",
	})
	_ = printer.Table(laslig.Table{
		Title:  "artifacts",
		Header: []string{"name", "status"},
		Rows: [][]string{
			{"darwin-arm64", "ready"},
			{"linux-amd64", "ready"},
		},
	})
}
```

## Current Surface

```go
printer.Section("Deploy")
printer.Notice(laslig.Notice{Level: laslig.NoticeWarningLevel, Title: "Partial success"})
printer.Record(laslig.Record{Title: "Build"})
printer.KV(laslig.KV{Title: "Config"})
printer.List(laslig.List{Title: "Packages"})
printer.Paragraph(laslig.Paragraph{Title: "Why", Body: "Readable defaults matter."})
printer.StatusLine(laslig.StatusLine{Level: laslig.NoticeSuccessLevel, Text: "Build ready"})
printer.Table(laslig.Table{Title: "Results"})
printer.Panel(laslig.Panel{Title: "Next step", Body: "Run mage check."})
printer.Box(laslig.Panel{Body: "Box is an alias for Panel."})
printer.Markdown(laslig.Markdown{Body: "# Notes\n\n- first\n- second"})
printer.CodeBlock(laslig.CodeBlock{Title: "Example", Language: "go", Body: `fmt.Println("hi")`})
printer.LogBlock(laslig.LogBlock{Title: "stderr excerpt", Body: "INFO boot complete\nWARN retry scheduled"})
```

`FormatAuto` resolves to human output on a terminal and plain text otherwise. `StyleAuto` enables ANSI styling only when the writer is attached to a TTY.

## JSON Mode

The same primitives can render machine-readable payloads:

```go
printer := laslig.New(os.Stdout, laslig.Policy{
	Format: laslig.FormatJSON,
})
```

That makes it practical to keep one semantic output path while exposing human, plain, and JSON surfaces from the same command.

## Rich Text, Code, And Logs

`laslig` can now render wrapped prose, Markdown, code, and caller-provided log excerpts without taking over logging itself:

```go
_ = printer.Paragraph(laslig.Paragraph{
	Title:  "Why",
	Body:   "Readable long-form CLI output should not require a hand-built Lip Gloss layout every time.",
	Footer: "Writers in, errors out.",
})

_ = printer.StatusLine(laslig.StatusLine{
	Level:  laslig.NoticeSuccessLevel,
	Text:   "Build ready",
	Detail: "mage check",
})

_ = printer.Markdown(laslig.Markdown{
	Title: "Release notes",
	Body:  "## Highlights\n\n- one renderer\n- three output surfaces",
})

_ = printer.CodeBlock(laslig.CodeBlock{
	Title:    "example.go",
	Language: "go",
	Body:     `fmt.Println("hello from laslig")`,
	Footer:   "Rendered through Glamour for terminal output.",
})

_ = printer.LogBlock(laslig.LogBlock{
	Title: "stderr excerpt",
	Body:  "INFO boot complete\nWARN retry scheduled\nERROR dependency missing",
})
```

`Markdown` and `CodeBlock` use Glamour for rich terminal rendering when styled human output is active. `LogBlock` is for explicit caller-provided excerpts, transcripts, and stderr captures. `laslig` still does not replace the application's logger.

## Structured Test Output

The `testjson` subpackage parses and renders `go test -json` streams without taking over command execution:

```go
cmd := exec.Command("go", "test", "-json", "./...")
stdout, err := cmd.StdoutPipe()
if err != nil {
	return err
}
cmd.Stderr = os.Stderr

if err := cmd.Start(); err != nil {
	return err
}

summary, err := testjson.Render(os.Stdout, stdout, testjson.Options{
	Policy: laslig.Policy{
		Format: laslig.FormatAuto,
		Style:  laslig.StyleAuto,
	},
	View: testjson.ViewCompact,
	DisabledSections: []testjson.Section{
		testjson.SectionSkippedTests,
	},
})
if err != nil {
	return err
}

if err := cmd.Wait(); err != nil {
	return err
}
if summary.HasFailures() {
	return errors.New("tests failed")
}
```

That shape works well in ordinary CLIs and in Mage targets. `laslig` stays responsible for rendering, while the caller stays responsible for process control. Callers can also disable grouped failed-test, skipped-test, package-error, or captured-output sections when they want a tighter stream.

This repository already dogfoods that pattern in [`magefile.go`](/Users/evanschultz/Documents/Code/hylla/laslig/main/magefile.go): `mage test` runs `go test -json ./...`, renders compact package and failure output through `testjson`, and still returns a normal Mage error on failure.

## Demo

The tracked demo command lives in [cmd/laslig-demo/main.go](/Users/evanschultz/Documents/Code/hylla/laslig/main/cmd/laslig-demo/main.go).
Small verified Go doc examples live in [example_test.go](/Users/evanschultz/Documents/Code/hylla/laslig/main/example_test.go).
The main demo groups the library into structured blocks and rich-text blocks so the visual showcase stays readable.

Run it locally:

```bash
mage demo
go run ./cmd/laslig-demo --format human --style always
go run ./cmd/laslig-demo --format json
```

The README GIF is generated from [docs/vhs/showcase.tape](/Users/evanschultz/Documents/Code/hylla/laslig/main/docs/vhs/showcase.tape).

## Planned Next

- theme configuration and preset flow
- richer `testjson` failure classification and subtest rollups
- compact prefix-style helpers beyond `StatusLine`
- more README visuals and side-by-side comparisons

## Development

This repository uses Mage for local automation.

```bash
mage check
mage test
mage build
mage demo
mage vhs
```

Structural terminal output is also covered by Charm `x/exp/golden` snapshots in the demo and `testjson` packages. Update them intentionally with:

```bash
go test ./cmd/laslig-demo -run TestRunArgsPlainGolden -args -update
go test ./testjson -run TestRenderPlainCompactGolden -args -update
```

README examples and terminal GIFs are generated from the tracked demo app and VHS tapes under [docs/vhs/](/Users/evanschultz/Documents/Code/hylla/laslig/main/docs/vhs).

## Plan

The tracked execution plan lives in [PLAN.md](/Users/evanschultz/Documents/Code/hylla/laslig/main/PLAN.md).
