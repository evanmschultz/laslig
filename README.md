# laslig

`laslig` helps Go CLIs print structured, human-readable output with Charm-native styling and Go-idiomatic ergonomics.

The package and module name stay `laslig`. The product branding is `Läslig`, from the Swedish `läslig`, meaning `legible`.

![Läslig demo](docs/vhs/showcase.gif)

## Why

Charm already gives Go developers strong building blocks:

- Lip Gloss for styling and layout
- Fang for help, usage, and CLI error presentation

What is still missing is a narrow, reusable layer for ordinary command output: results, notices, summaries, tables, warnings, and errors that should look intentional without forcing an application into a framework.

`laslig` is that layer.

## Status

The first core wave is live. Today the package includes:

- output policy and mode resolution
- document-layout defaults with caller-tunable spacing, section indentation, and list markers
- a `Printer`
- sections
- notices for info, success, warning, and error output
- records and lists
- aligned key-value blocks
- paragraph blocks
- compact status lines
- tables
- panels
- Glamour-backed Markdown blocks
- Glamour-backed code blocks
- boxed log/transcript blocks for caller-provided output
- compact and detailed `gotestout` rendering for `go test -json`
- caller-tunable `gotestout` summary and output sections

The next wave is focused on theme configuration, deeper `gotestout` classification, and tightening the docs/examples further.

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
printer.Markdown(laslig.Markdown{Body: "# Notes\n\n- first\n- second"})
printer.CodeBlock(laslig.CodeBlock{Title: "Example", Language: "go", Body: `fmt.Println("hi")`})
printer.LogBlock(laslig.LogBlock{Title: "stderr excerpt", Body: "INFO boot complete\nWARN retry scheduled"})
```

`FormatAuto` resolves to human output on a terminal and plain text otherwise. `StyleAuto` enables ANSI styling only when the writer is attached to a TTY.

## Layout

Läslig now treats output more like a document by default:

- one leading blank line before the first rendered block
- one blank line between ordinary blocks
- stronger separation before new sections
- section-owned indentation for blocks that follow a `Section`

Commands can tune that shape when they need something flatter:

```go
layout := laslig.DefaultLayout().
	WithLeadingGap(0).
	WithSectionIndent(0).
	WithListMarker(laslig.ListMarkerBullet)

printer := laslig.New(os.Stdout, laslig.Policy{
	Format: laslig.FormatAuto,
	Style:  laslig.StyleAuto,
	Layout: &layout,
})
```

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
	Body: "# Release Notes\n\n## Highlights\n\n- one renderer\n- three output surfaces",
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

The `gotestout` subpackage parses and renders `go test -json` streams without taking over command execution:

![gotestout example](docs/vhs/gotestout.gif)

```go
import (
	"errors"
	"os"
	"os/exec"

	"github.com/evanmschultz/laslig"
	"github.com/evanmschultz/laslig/gotestout"
)

cmd := exec.Command("go", "test", "-json", "./...")
stdout, err := cmd.StdoutPipe()
if err != nil {
	return err
}
cmd.Stderr = os.Stderr

if err := cmd.Start(); err != nil {
	return err
}

summary, err := gotestout.Render(os.Stdout, stdout, gotestout.Options{
	Policy: laslig.Policy{
		Format: laslig.FormatAuto,
		Style:  laslig.StyleAuto,
	},
	View: gotestout.ViewCompact,
	DisabledSections: []gotestout.Section{
		gotestout.SectionSkippedTests,
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

That shape works well in ordinary Go CLIs, Mage targets, Cobra/Fang commands, and small Go helpers invoked from tools like `make`, `just`, or `task`. `laslig` stays responsible for rendering, while the caller stays responsible for process control. Callers can also disable grouped failed-test, skipped-test, package-error, or captured-output sections when they want a tighter stream.

This repository already dogfoods that pattern in [`magefile.go`](/Users/evanschultz/Documents/Code/hylla/laslig/main/magefile.go): `mage test` runs `go test -json ./...`, renders compact package and failure output through `gotestout`, and still returns a normal Mage error on failure.
The focused runnable example for that package lives in [examples/gotestout/main.go](/Users/evanschultz/Documents/Code/hylla/laslig/main/examples/gotestout/main.go).

Common ways to try that surface locally:

```bash
go run ./examples/gotestout --format human --style always
mage test
```

## Demo

The tracked all-in-one showcase lives in [examples/all/main.go](/Users/evanschultz/Documents/Code/hylla/laslig/main/examples/all/main.go).
The focused logging example package that uses `charm.land/log/v2` as a demo-only dependency lives in [examples/logging/logging.go](/Users/evanschultz/Documents/Code/hylla/laslig/main/examples/logging/logging.go) and is imported directly by the main showcase.
The focused `gotestout` example lives in [examples/gotestout/main.go](/Users/evanschultz/Documents/Code/hylla/laslig/main/examples/gotestout/main.go).
Small verified Go doc examples live in [example_test.go](/Users/evanschultz/Documents/Code/hylla/laslig/main/example_test.go).
The main showcase is a guided walkthrough: it names each primitive directly and explains what it is for and when to use it, then closes with an explicit `gotestout` section that renders a real Mage-style Build, Tests, and Coverage preview inline.

Run it locally:

```bash
mage demo
go run ./examples/all --format human --style always
go run ./examples/all --format json
go run ./examples/gotestout --format human --style always
mage test
```

`mage demo` is the normal primitive walkthrough entrypoint. `mage test` is the real Mage-facing `gotestout` dogfood path. The `go run` forms above are the same examples with explicit format/style flags.

The README GIF is generated from [docs/vhs/showcase.tape](/Users/evanschultz/Documents/Code/hylla/laslig/main/docs/vhs/showcase.tape).

## Planned Next

- theme configuration and preset flow
- richer `gotestout` failure classification and subtest rollups
- compact prefix-style helpers beyond `StatusLine`
- more README visuals and side-by-side comparisons

## Development

This repository uses Mage for local automation.

Install the same Mage version used in CI:

```bash
go install github.com/magefile/mage@v1.17.0
```

```bash
mage check
mage test
mage build
mage demo
mage vhs
```

See [CONTRIBUTING.md](/Users/evanschultz/Documents/Code/hylla/laslig/main/CONTRIBUTING.md) for contributor workflow details and [SECURITY.md](/Users/evanschultz/Documents/Code/hylla/laslig/main/SECURITY.md) for vulnerability reporting guidance.

Structural terminal output is also covered by Charm `x/exp/golden` snapshots in the demo and `gotestout` packages. Update them intentionally with:

```bash
go test ./examples/all -run TestRunArgsPlainGolden -args -update
go test ./examples/gotestout -run TestRunArgsPlainGolden -args -update
go test ./gotestout -run 'TestRenderPlainCompactGolden|TestRenderHumanStyledCompactGolden' -args -update
```

README examples and terminal GIFs are generated from the tracked demo app and VHS tapes under [docs/vhs/](/Users/evanschultz/Documents/Code/hylla/laslig/main/docs/vhs).

## License

`laslig` is licensed under [Apache-2.0](/Users/evanschultz/Documents/Code/hylla/laslig/main/LICENSE).

## Plan

The tracked execution plan lives in [PLAN.md](/Users/evanschultz/Documents/Code/hylla/laslig/main/PLAN.md).
