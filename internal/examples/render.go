package examples

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	charmlog "charm.land/log/v2"
	"github.com/charmbracelet/x/term"

	"github.com/evanmschultz/laslig"
	"github.com/evanmschultz/laslig/gotestout"
)

const demoSpinnerStepDelay = 450 * time.Millisecond

// RenderAll writes the aggregate walkthrough used by mage demo.
func RenderAll(out io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("Läslig demo"); err != nil {
		return fmt.Errorf("render demo heading: %w", err)
	}
	if err := printer.Notice(laslig.Notice{
		Level: laslig.NoticeInfoLevel,
		Title: "A guided primitive walkthrough",
		Body:  "This combined demo assembles the same focused examples that live under examples/ so the README GIFs and mage demo stay aligned.",
		Detail: []string{
			"Use Läslig for ordinary CLI output, not command parsing or application logging.",
			"Run any focused example directly with go run ./examples/<name> --format human --style always, or use mage demo to assemble them all.",
		},
	}); err != nil {
		return fmt.Errorf("render demo notice: %w", err)
	}

	renderers := []Renderer{
		RenderSection,
		RenderNotice,
		RenderRecord,
		RenderKV,
		RenderList,
		RenderTable,
		RenderPanel,
		RenderParagraph,
		RenderStatusLine,
		RenderSpinner,
		RenderMarkdown,
		RenderCodeBlock,
		RenderLogBlock,
		RenderGotestout,
		RenderMageCheckPreview,
	}
	for _, render := range renderers {
		if err := render(out, printer); err != nil {
			return err
		}
	}
	return nil
}

// RenderSection demonstrates how Section creates document headings and owned
// blocks.
func RenderSection(_ io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("Section"); err != nil {
		return fmt.Errorf("render section heading: %w", err)
	}
	if err := printer.Paragraph(laslig.Paragraph{
		Body:   "Use Section to start a new document region. The blocks that follow are indented under that heading until the next section begins.",
		Footer: "Section is the document-level primitive that gives the rest of the output clear ownership.",
	}); err != nil {
		return fmt.Errorf("render section explanation: %w", err)
	}
	if err := printer.Record(laslig.Record{
		Title: "Owned blocks",
		Fields: []laslig.Field{
			{Label: "what", Value: "Paragraphs, records, lists, tables, and panels inherit the active section indent."},
			{Label: "why", Value: "That default makes CLI output read more like a document than a stream of unrelated prints."},
		},
	}); err != nil {
		return fmt.Errorf("render section record: %w", err)
	}
	if err := printer.Section("Next section"); err != nil {
		return fmt.Errorf("render next section heading: %w", err)
	}
	if err := printer.Notice(laslig.Notice{
		Level: laslig.NoticeInfoLevel,
		Title: "Section boundaries reset ownership",
		Body:  "Starting another Section makes the boundary visible without requiring callers to manage indentation by hand.",
	}); err != nil {
		return fmt.Errorf("render section boundary notice: %w", err)
	}
	return nil
}

// RenderNotice demonstrates the Notice primitive.
func RenderNotice(_ io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("Notice"); err != nil {
		return fmt.Errorf("render notice section: %w", err)
	}
	return printer.Notice(laslig.Notice{
		Level: laslig.NoticeInfoLevel,
		Title: "Use Notice for semantic user-facing diagnostics",
		Body:  "Notice is the default surface for info, success, warning, and error output that should stand out without turning into logging.",
		Detail: []string{
			"When: validation feedback, partial success, next-step guidance, dependency checks, or release milestones.",
		},
	})
}

// RenderRecord demonstrates the Record primitive.
func RenderRecord(_ io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("Record"); err != nil {
		return fmt.Errorf("render record section: %w", err)
	}
	return printer.Record(laslig.Record{
		Title: "Record",
		Fields: []laslig.Field{
			{Label: "what", Value: "One object or result rendered as labeled facts."},
			{Label: "when", Value: "Build metadata, artifact details, environment summaries."},
			{Label: "example", Value: "module github.com/evanmschultz/laslig", Identifier: true},
		},
	})
}

// RenderKV demonstrates the KV primitive.
func RenderKV(_ io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("KV"); err != nil {
		return fmt.Errorf("render kv section: %w", err)
	}
	return printer.KV(laslig.KV{
		Title: "KV",
		Pairs: []laslig.Field{
			{Label: "what", Value: "Compact aligned configuration or status."},
			{Label: "when", Value: "Resolved policy, flags, execution mode."},
			{Label: "format", Value: string(printer.Mode().Format)},
			{Label: "styled", Value: fmt.Sprintf("%t", printer.Mode().Styled), Muted: true},
		},
	})
}

// RenderList demonstrates the List primitive.
func RenderList(_ io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("List"); err != nil {
		return fmt.Errorf("render list section: %w", err)
	}
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
}

// RenderTable demonstrates the Table primitive.
func RenderTable(_ io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("Table"); err != nil {
		return fmt.Errorf("render table section: %w", err)
	}
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
}

// RenderPanel demonstrates the Panel primitive.
func RenderPanel(_ io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("Panel"); err != nil {
		return fmt.Errorf("render panel section: %w", err)
	}
	return printer.Panel(laslig.Panel{
		Title:  "Panel",
		Body:   "Use Panel for rationale, next steps, and larger callouts that should stand apart from the rest of the document.",
		Footer: "Panels are stronger than Paragraph and lighter than inventing a custom layout.",
	})
}

// RenderParagraph demonstrates the Paragraph primitive.
func RenderParagraph(_ io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("Paragraph"); err != nil {
		return fmt.Errorf("render paragraph section: %w", err)
	}
	return printer.Paragraph(laslig.Paragraph{
		Title:  "Paragraph",
		Body:   "Use Paragraph for readable rationale, release context, and longer help text when a Notice or Panel would be too heavy.",
		Footer: "Paragraph is the simplest way to let a CLI explain itself in normal sentence-style text.",
	})
}

// RenderStatusLine demonstrates the StatusLine primitive.
func RenderStatusLine(_ io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("StatusLine"); err != nil {
		return fmt.Errorf("render statusline section: %w", err)
	}
	return printer.StatusLine(laslig.StatusLine{
		Level:  laslig.NoticeSuccessLevel,
		Text:   "Use StatusLine for one compact semantic result",
		Detail: "When: build ready, cache hit, package passed",
	})
}

// RenderSpinner demonstrates the transient Spinner helper.
func RenderSpinner(out io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("Spinner"); err != nil {
		return fmt.Errorf("render spinner section: %w", err)
	}
	if err := printer.Paragraph(laslig.Paragraph{
		Body:   "Use Spinner when long-running work might otherwise stay quiet for several seconds. Prefer StatusLine or Notice when the operation starts and finishes quickly enough that durable output is enough.",
		Footer: "This focused example falls back to stable start and finish status rows in plain, JSON, and non-interactive output.",
	}); err != nil {
		return fmt.Errorf("render spinner intro: %w", err)
	}

	spin := printer.NewSpinner()
	if err := spin.Start("Waiting for remote rollout"); err != nil {
		return fmt.Errorf("start spinner: %w", err)
	}
	pauseForAnimatedPreview(out, demoSpinnerStepDelay)
	if err := spin.Update("Waiting for remote rollout (2/3)"); err != nil {
		return fmt.Errorf("update spinner: %w", err)
	}
	pauseForAnimatedPreview(out, demoSpinnerStepDelay)
	if err := spin.Update("Waiting for remote rollout (3/3)"); err != nil {
		return fmt.Errorf("update spinner final step: %w", err)
	}
	pauseForAnimatedPreview(out, demoSpinnerStepDelay)
	if err := spin.Stop("Rollout ready", laslig.NoticeSuccessLevel); err != nil {
		return fmt.Errorf("stop spinner: %w", err)
	}
	return nil
}

// RenderMarkdown demonstrates the Markdown primitive.
func RenderMarkdown(_ io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("Markdown"); err != nil {
		return fmt.Errorf("render markdown section: %w", err)
	}
	if err := printer.Paragraph(laslig.Paragraph{
		Body:   "Use Markdown when your CLI already has release notes, changelog text, or generated help content to show.",
		Footer: "The block below is real terminal-rendered Markdown.",
	}); err != nil {
		return fmt.Errorf("render markdown intro: %w", err)
	}
	return printer.Markdown(laslig.Markdown{
		Body: "# Release Notes\n\n## Highlights\n\n- one renderer\n- three output surfaces\n- caller-owned logging",
	})
}

// RenderCodeBlock demonstrates the CodeBlock primitive.
func RenderCodeBlock(_ io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("CodeBlock"); err != nil {
		return fmt.Errorf("render code block section: %w", err)
	}
	if err := printer.Paragraph(laslig.Paragraph{
		Body:   "Use CodeBlock for commands, snippets, generated files, and config examples.",
		Footer: "The block below shows a Go snippet rendered through Glamour.",
	}); err != nil {
		return fmt.Errorf("render code block intro: %w", err)
	}
	return printer.CodeBlock(laslig.CodeBlock{
		Title:    "Go snippet",
		Language: "go",
		Body:     "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello from laslig\")\n}",
		Footer:   "Use CodeBlock when code should stay visibly distinct from prose.",
	})
}

// RenderLogBlock demonstrates the LogBlock primitive.
func RenderLogBlock(_ io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("LogBlock"); err != nil {
		return fmt.Errorf("render log block section: %w", err)
	}
	if err := printer.Paragraph(laslig.Paragraph{
		Body:   "Use LogBlock for selected stderr or log excerpts while the application keeps owning logging.",
		Footer: "The block below captures real charm/log output and renders it through Läslig.",
	}); err != nil {
		return fmt.Errorf("render log block intro: %w", err)
	}
	return printer.LogBlock(laslig.LogBlock{
		Title:  "Captured charm/log transcript",
		Body:   transcript(),
		Footer: "Use LogBlock for selected transcripts while the application keeps owning the logger.",
	})
}

// RenderGotestout demonstrates the focused gotestout example.
func RenderGotestout(out io.Writer, printer *laslig.Printer) error {
	if printer.Mode().Format != laslig.FormatJSON {
		if err := printer.Section("gotestout"); err != nil {
			return fmt.Errorf("render gotestout section: %w", err)
		}
		if err := printer.Paragraph(laslig.Paragraph{
			Body:   "Use gotestout for attractive, structured go test output when your task runner, CLI command, or Go helper behind make/just should keep owning process control.",
			Footer: "This focused example intentionally mixes passing, skipped, and failing test events plus one package build failure.",
		}); err != nil {
			return fmt.Errorf("render gotestout intro: %w", err)
		}
		if err := printer.Notice(laslig.Notice{
			Level: laslig.NoticeInfoLevel,
			Title: "Mixed fixture demo",
			Body:  "The example command itself is expected to exit successfully so you can inspect the output shape.",
			Detail: []string{
				"Use mage test for the real repository task-runner path.",
			},
		}); err != nil {
			return fmt.Errorf("render gotestout notice: %w", err)
		}
	}

	_, err := gotestout.Render(out, strings.NewReader(focusedGotestoutSampleStream), gotestout.Options{
		Policy: laslig.Policy{
			Format: printer.Mode().Format,
			Style:  StylePolicyForMode(printer.Mode()),
		},
		View: gotestout.ViewDetailed,
	})
	if err != nil {
		return fmt.Errorf("render gotestout stream: %w", err)
	}
	return nil
}

// RenderMageCheckPreview demonstrates the repository-style Mage integration.
func RenderMageCheckPreview(out io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("gotestout + Mage"); err != nil {
		return fmt.Errorf("render mage preview section: %w", err)
	}
	if err := printer.Paragraph(laslig.Paragraph{
		Body:   "Use gotestout inside Mage or small Go helpers behind make, just, or task when you want caller-owned process control with a readable test stream.",
		Footer: "The preview below matches this repository's mage check and mage test shape, including the recommended spinner handoff before the live test stream starts.",
	}); err != nil {
		return fmt.Errorf("render mage preview intro: %w", err)
	}
	return renderMageCheckPreview(out, printer)
}

// focusedGotestoutSampleStream is the deterministic mixed test fixture used by
// the focused gotestout example and the aggregate walkthrough.
const focusedGotestoutSampleStream = `{"Action":"run","Package":"example/pkg","Test":"TestPass"}
{"Action":"output","Package":"example/pkg","Test":"TestPass","Output":"=== RUN   TestPass\n"}
{"Action":"output","Package":"example/pkg","Test":"TestPass","Output":"note: useful output\n"}
{"Action":"output","Package":"example/pkg","Test":"TestPass","Output":"--- PASS: TestPass (0.01s)\n"}
{"Action":"pass","Package":"example/pkg","Test":"TestPass","Elapsed":0.01}
{"Action":"run","Package":"example/pkg","Test":"TestSkip"}
{"Action":"output","Package":"example/pkg","Test":"TestSkip","Output":"--- SKIP: TestSkip (0.00s)\n"}
{"Action":"skip","Package":"example/pkg","Test":"TestSkip","Elapsed":0}
{"Action":"run","Package":"example/pkg","Test":"TestFail"}
{"Action":"output","Package":"example/pkg","Test":"TestFail","Output":"main_test.go:42: expected boom\n"}
{"Action":"output","Package":"example/pkg","Test":"TestFail","Output":"--- FAIL: TestFail (0.02s)\n"}
{"Action":"fail","Package":"example/pkg","Test":"TestFail","Elapsed":0.02}
{"Action":"output","Package":"example/pkg","Output":"FAIL\texample/pkg [build failed]\n","FailedBuild":"example/pkg"}
{"Action":"fail","Package":"example/pkg","Elapsed":0.03}
`

// renderMageCheckPreview renders a deterministic passing Mage-shaped flow.
func renderMageCheckPreview(out io.Writer, printer *laslig.Printer) error {
	if err := printer.Section("Build"); err != nil {
		return fmt.Errorf("render build section: %w", err)
	}
	if err := printer.StatusLine(laslig.StatusLine{
		Level:  laslig.NoticeInfoLevel,
		Text:   "Building example packages",
		Detail: "./examples/...",
	}); err != nil {
		return fmt.Errorf("render build start: %w", err)
	}
	if err := printer.StatusLine(laslig.StatusLine{
		Level:  laslig.NoticeSuccessLevel,
		Text:   "Built example packages",
		Detail: "./examples/...",
	}); err != nil {
		return fmt.Errorf("render build success: %w", err)
	}
	spin := printer.NewSpinner()
	if err := spin.Start("Waiting for first test event"); err != nil {
		return fmt.Errorf("start mage spinner: %w", err)
	}
	pauseForAnimatedPreview(out, demoSpinnerStepDelay)
	if err := spin.Update("Waiting for first test event from go test -json"); err != nil {
		return fmt.Errorf("update mage spinner: %w", err)
	}
	pauseForAnimatedPreview(out, demoSpinnerStepDelay)
	if err := spin.Stop("Test stream detected", laslig.NoticeSuccessLevel); err != nil {
		return fmt.Errorf("stop mage spinner: %w", err)
	}

	if err := printer.Section("Tests"); err != nil {
		return fmt.Errorf("render tests section: %w", err)
	}
	if _, err := gotestout.Render(out, strings.NewReader(mageCheckSampleStream()), gotestout.Options{
		Policy: laslig.Policy{
			Format: printer.Mode().Format,
			Style:  StylePolicyForMode(printer.Mode()),
		},
		View: gotestout.ViewCompact,
	}); err != nil {
		return fmt.Errorf("render mage gotestout stream: %w", err)
	}

	if err := printer.Section("Coverage"); err != nil {
		return fmt.Errorf("render coverage section: %w", err)
	}
	if err := printer.Table(laslig.Table{
		Header: []string{"package", "cover"},
		Rows: [][]string{
			{"github.com/evanmschultz/laslig", "73.7%"},
			{"github.com/evanmschultz/laslig/examples/all", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/codeblock", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/gotestout", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/kv", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/list", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/logblock", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/magecheck", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/markdown", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/notice", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/panel", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/paragraph", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/record", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/section", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/spinner", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/statusline", "100.0%"},
			{"github.com/evanmschultz/laslig/examples/table", "100.0%"},
			{"github.com/evanmschultz/laslig/gotestout", "82.0%"},
			{"github.com/evanmschultz/laslig/internal/examples", "70.5%"},
			{"github.com/evanmschultz/laslig/internal/exampletestutil", "75.5%"},
			{"github.com/evanmschultz/laslig/internal/glamrender", "86.7%"},
			{"github.com/evanmschultz/laslig/internal/layout", "87.5%"},
			{"github.com/evanmschultz/laslig/internal/table", "96.6%"},
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

// mageCheckSampleStream mirrors the package-level shape this repository prints
// through mage test.
func mageCheckSampleStream() string {
	packages := []struct {
		name  string
		tests int
	}{
		{"github.com/evanmschultz/laslig", 60},
		{"github.com/evanmschultz/laslig/examples/all", 8},
		{"github.com/evanmschultz/laslig/examples/codeblock", 4},
		{"github.com/evanmschultz/laslig/examples/gotestout", 7},
		{"github.com/evanmschultz/laslig/examples/kv", 4},
		{"github.com/evanmschultz/laslig/examples/list", 4},
		{"github.com/evanmschultz/laslig/examples/logblock", 4},
		{"github.com/evanmschultz/laslig/examples/magecheck", 4},
		{"github.com/evanmschultz/laslig/examples/markdown", 4},
		{"github.com/evanmschultz/laslig/examples/notice", 4},
		{"github.com/evanmschultz/laslig/examples/panel", 4},
		{"github.com/evanmschultz/laslig/examples/paragraph", 4},
		{"github.com/evanmschultz/laslig/examples/record", 4},
		{"github.com/evanmschultz/laslig/examples/section", 4},
		{"github.com/evanmschultz/laslig/examples/spinner", 4},
		{"github.com/evanmschultz/laslig/examples/statusline", 4},
		{"github.com/evanmschultz/laslig/examples/table", 4},
		{"github.com/evanmschultz/laslig/gotestout", 10},
		{"github.com/evanmschultz/laslig/internal/examples", 26},
		{"github.com/evanmschultz/laslig/internal/exampletestutil", 5},
		{"github.com/evanmschultz/laslig/internal/glamrender", 2},
		{"github.com/evanmschultz/laslig/internal/layout", 2},
		{"github.com/evanmschultz/laslig/internal/table", 2},
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

// ttyBuffer reports a file descriptor so charm/log keeps its normal text
// formatting when supported.
type ttyBuffer struct {
	bytes.Buffer
}

// Fd reports the current stdout descriptor.
func (ttyBuffer) Fd() uintptr {
	return os.Stdout.Fd()
}

func writerSupportsAnimation(out io.Writer) bool {
	file, ok := out.(term.File)
	if !ok {
		return false
	}
	return term.IsTerminal(file.Fd())
}

func pauseForAnimatedPreview(out io.Writer, delay time.Duration) {
	if writerSupportsAnimation(out) {
		time.Sleep(delay)
	}
}

// transcript captures one real charm/log transcript for the LogBlock demo.
func transcript() string {
	var writer ttyBuffer

	logger := charmlog.NewWithOptions(&writer, charmlog.Options{
		Formatter:       charmlog.TextFormatter,
		ReportTimestamp: false,
		Prefix:          "demo",
	})

	logger.Info("boot complete", "component", "cache")
	logger.Warn("retry scheduled", "after", "3s")
	logger.Error("dependency missing", "name", "git")

	return strings.TrimRight(writer.String(), "\n")
}
