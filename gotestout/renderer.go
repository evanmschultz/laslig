package gotestout

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/evanmschultz/laslig"
)

// outputKey identifies one buffered stream of package or test output.
type outputKey struct {
	pkg  string
	test string
}

// Renderer consumes go test -json events and renders them to one writer.
type Renderer struct {
	out          io.Writer
	mode         laslig.Mode
	theme        laslig.Theme
	printer      *laslig.Printer
	options      Options
	summary      Summary
	outputs      map[outputKey][]string
	buildFailed  map[string]bool
	failedTests  []outcome
	skippedTests []outcome
	packageError []outcome
	wroteResults bool
	jsonEncoder  *json.Encoder
}

// NewRenderer constructs one renderer for the provided writer and options.
func NewRenderer(out io.Writer, options Options) *Renderer {
	options = withDefaults(options)
	mode := laslig.ResolveMode(out, options.Policy)
	layout := laslig.DefaultLayout().WithLeadingGap(0)

	renderer := &Renderer{
		out:   out,
		mode:  mode,
		theme: laslig.DefaultTheme(mode),
		printer: laslig.New(out, laslig.Policy{
			Format: options.Policy.Format,
			Style:  options.Policy.Style,
			Layout: &layout,
		}),
		options:     options,
		outputs:     make(map[outputKey][]string),
		buildFailed: make(map[string]bool),
	}
	if mode.Format == laslig.FormatJSON {
		renderer.jsonEncoder = json.NewEncoder(out)
		renderer.jsonEncoder.SetEscapeHTML(false)
	}
	return renderer
}

// Summary returns the counts accumulated so far.
func (r *Renderer) Summary() Summary {
	return r.summary
}

// WriteEvent consumes one parsed event and writes any corresponding output.
func (r *Renderer) WriteEvent(event Event) error {
	if r.mode.Format == laslig.FormatJSON {
		if err := r.jsonEncoder.Encode(event); err != nil {
			return fmt.Errorf("write json event: %w", err)
		}
	}

	switch event.Action {
	case ActionOutput, ActionBuildOutput:
		r.recordOutput(event)
	case ActionPass, ActionFail, ActionSkip:
		r.recordTerminal(event)
		if r.mode.Format != laslig.FormatJSON {
			if err := r.renderTerminal(event); err != nil {
				return err
			}
		}
		r.clearOutput(event)
	}
	return nil
}

// Finish writes the final summary for human and plain output.
func (r *Renderer) Finish() error {
	if r.mode.Format == laslig.FormatJSON {
		return nil
	}
	if r.wroteResults {
		if _, err := fmt.Fprintln(r.out); err != nil {
			return fmt.Errorf("write summary spacing: %w", err)
		}
	}

	record := laslig.Record{
		Title: "Test summary",
		Fields: []laslig.Field{
			{Label: "tests", Value: fmt.Sprintf("%d", r.summary.TotalTests())},
			{Label: "passed", Value: fmt.Sprintf("%d", r.summary.TestsPassed)},
			{Label: "failed", Value: fmt.Sprintf("%d", r.summary.TestsFailed)},
			{Label: "skipped", Value: fmt.Sprintf("%d", r.summary.TestsSkipped)},
			{Label: "packages", Value: fmt.Sprintf("%d", r.summary.TotalPackages())},
			{Label: "pkg passed", Value: fmt.Sprintf("%d", r.summary.PackagesPassed)},
			{Label: "pkg failed", Value: fmt.Sprintf("%d", r.summary.PackagesFailed)},
			{Label: "pkg skipped", Value: fmt.Sprintf("%d", r.summary.PackagesSkipped)},
		},
	}
	if r.summary.BuildErrors > 0 {
		record.Fields = append(record.Fields, laslig.Field{
			Label: "build errors",
			Value: fmt.Sprintf("%d", r.summary.BuildErrors),
		})
	}
	if err := r.printer.Record(record); err != nil {
		return fmt.Errorf("write summary record: %w", err)
	}
	if r.sectionEnabled(SectionFailedTests) && len(r.failedTests) > 0 {
		if err := r.printer.List(r.outcomeList("Failed tests", "fail", r.failedTests, false)); err != nil {
			return fmt.Errorf("write failed tests list: %w", err)
		}
	}
	if r.sectionEnabled(SectionPackageErrors) && len(r.packageError) > 0 {
		if err := r.printer.List(r.outcomeList("Package errors", "error", r.packageError, true)); err != nil {
			return fmt.Errorf("write package errors list: %w", err)
		}
	}
	if r.sectionEnabled(SectionSkippedTests) && len(r.skippedTests) > 0 {
		if err := r.printer.List(r.outcomeList("Skipped tests", "skip", r.skippedTests, false)); err != nil {
			return fmt.Errorf("write skipped tests list: %w", err)
		}
	}

	notice := laslig.Notice{
		Level: laslig.NoticeSuccessLevel,
		Title: "All tests passed",
		Body:  fmt.Sprintf("%d test%s passed across %d package%s.", r.summary.TestsPassed, plural(r.summary.TestsPassed), r.summary.TotalPackages(), plural(r.summary.TotalPackages())),
	}
	switch {
	case r.summary.HasFailures():
		notice.Level = laslig.NoticeErrorLevel
		notice.Title = "Test failures detected"
		notice.Body = fmt.Sprintf("%d test failure%s and %d build error%s across %d package%s.", r.summary.TestsFailed, plural(r.summary.TestsFailed), r.summary.BuildErrors, plural(r.summary.BuildErrors), r.summary.TotalPackages(), plural(r.summary.TotalPackages()))
	case r.summary.TestsSkipped > 0:
		notice.Level = laslig.NoticeWarningLevel
		notice.Title = "Tests passed with skips"
		notice.Body = fmt.Sprintf("%d skipped test%s across %d package%s.", r.summary.TestsSkipped, plural(r.summary.TestsSkipped), r.summary.TotalPackages(), plural(r.summary.TotalPackages()))
	case r.summary.PackagesSkipped > 0:
		notice.Detail = append(notice.Detail, fmt.Sprintf("%d package%s had no tests.", r.summary.PackagesSkipped, plural(r.summary.PackagesSkipped)))
	}
	if err := r.printer.Notice(notice); err != nil {
		return fmt.Errorf("write summary notice: %w", err)
	}
	return nil
}

// Parse decodes a full stream into memory.
func Parse(in io.Reader) ([]Event, error) {
	decoder := json.NewDecoder(bufio.NewReader(in))
	var events []Event
	for {
		var event Event
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				return events, nil
			}
			return nil, fmt.Errorf("decode event: %w", err)
		}
		events = append(events, event)
	}
}

// Render parses a stream, renders it, and returns a summary of terminal events.
func Render(out io.Writer, in io.Reader, options Options) (Summary, error) {
	renderer := NewRenderer(out, options)
	decoder := json.NewDecoder(bufio.NewReader(in))

	for {
		var event Event
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			return Summary{}, fmt.Errorf("decode event: %w", err)
		}
		if err := renderer.WriteEvent(event); err != nil {
			return Summary{}, err
		}
	}

	if err := renderer.Finish(); err != nil {
		return Summary{}, err
	}
	return renderer.Summary(), nil
}

func withDefaults(options Options) Options {
	if options.Policy.Format == "" {
		options.Policy.Format = laslig.FormatAuto
	}
	if options.Policy.Style == "" {
		options.Policy.Style = laslig.StyleAuto
	}
	if options.View == "" {
		options.View = ViewCompact
	}
	return options
}

// sectionEnabled reports whether one optional rendered section is enabled.
func (r *Renderer) sectionEnabled(section Section) bool {
	for _, disabled := range r.options.DisabledSections {
		if disabled == section {
			return false
		}
	}
	return true
}

func (r *Renderer) recordOutput(event Event) {
	key := outputKey{pkg: event.Package, test: event.Test}
	r.outputs[key] = append(r.outputs[key], event.Output)
	if event.FailedBuild != "" || strings.Contains(event.Output, "[build failed]") {
		r.buildFailed[event.Package] = true
	}
}

func (r *Renderer) clearOutput(event Event) {
	delete(r.outputs, outputKey{pkg: event.Package, test: event.Test})
	if event.PackageEvent() {
		delete(r.outputs, outputKey{pkg: event.Package})
	}
}

func (r *Renderer) recordTerminal(event Event) {
	if event.PackageEvent() {
		switch event.Action {
		case ActionPass:
			r.summary.PackagesPassed++
		case ActionFail:
			r.summary.PackagesFailed++
			if lines := r.cleanedOutput(outputKey{pkg: event.Package}); len(lines) > 0 || r.buildFailed[event.Package] {
				r.packageError = append(r.packageError, outcome{
					Package: event.Package,
					Elapsed: event.Elapsed,
					Output:  lines,
				})
			}
			if r.buildFailed[event.Package] {
				r.summary.BuildErrors++
			}
		case ActionSkip:
			r.summary.PackagesSkipped++
		}
		return
	}

	switch event.Action {
	case ActionPass:
		r.summary.TestsPassed++
	case ActionFail:
		r.summary.TestsFailed++
		r.failedTests = append(r.failedTests, outcome{
			Package: event.Package,
			Test:    event.Test,
			Elapsed: event.Elapsed,
			Output:  r.cleanedOutput(outputKey{pkg: event.Package, test: event.Test}),
		})
	case ActionSkip:
		r.summary.TestsSkipped++
		r.skippedTests = append(r.skippedTests, outcome{
			Package: event.Package,
			Test:    event.Test,
			Elapsed: event.Elapsed,
		})
	}
}

func (r *Renderer) renderTerminal(event Event) error {
	if event.PackageEvent() {
		if err := r.writeLine(r.renderPackageLine(event)); err != nil {
			return err
		}
		if r.sectionEnabled(SectionOutput) && (event.Action == ActionFail || r.options.View == ViewDetailed) {
			if err := r.writeOutputLines(outputKey{pkg: event.Package}); err != nil {
				return err
			}
		}
		return nil
	}

	if r.options.View == ViewCompact && event.Action != ActionFail {
		return nil
	}
	if err := r.writeLine(r.renderTestLine(event)); err != nil {
		return err
	}
	if r.sectionEnabled(SectionOutput) && (event.Action == ActionFail || r.options.View == ViewDetailed) {
		if err := r.writeOutputLines(outputKey{pkg: event.Package, test: event.Test}); err != nil {
			return err
		}
	}
	return nil
}

func (r *Renderer) outcomeList(title string, badge string, outcomes []outcome, includeOutput bool) laslig.List {
	items := make([]laslig.ListItem, 0, len(outcomes))
	for _, item := range outcomes {
		titleValue := item.Package
		if item.Test != "" {
			titleValue = item.Test
		}

		fields := []laslig.Field{
			{Label: "package", Value: item.Package, Identifier: true},
			{Label: "elapsed", Value: fmt.Sprintf("%.2fs", item.Elapsed), Muted: true},
		}
		if includeOutput && r.sectionEnabled(SectionOutput) && len(item.Output) > 0 {
			fields = append(fields, laslig.Field{
				Label: "detail",
				Value: item.Output[0],
				Muted: true,
			})
		}

		items = append(items, laslig.ListItem{
			Title:  titleValue,
			Badge:  badge,
			Fields: fields,
		})
	}
	return laslig.List{
		Title: title,
		Items: items,
	}
}

func (r *Renderer) writeOutputLines(key outputKey) error {
	lines := r.cleanedOutput(key)
	for _, line := range lines {
		if _, err := fmt.Fprintf(r.out, "  %s\n", line); err != nil {
			return fmt.Errorf("write event output: %w", err)
		}
	}
	return nil
}

func (r *Renderer) cleanedOutput(key outputKey) []string {
	raw := r.outputs[key]
	lines := make([]string, 0, len(raw))
	for _, chunk := range raw {
		for _, line := range strings.Split(chunk, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || isFramingLine(line, key.pkg, key.test) {
				continue
			}
			lines = append(lines, line)
		}
	}
	return lines
}

func isFramingLine(line string, pkg string, test string) bool {
	return line == "PASS" ||
		line == "FAIL" ||
		strings.HasPrefix(line, "=== RUN") ||
		strings.HasPrefix(line, "=== PAUSE") ||
		strings.HasPrefix(line, "=== CONT") ||
		strings.HasPrefix(line, "--- PASS: "+test) ||
		strings.HasPrefix(line, "--- FAIL: "+test) ||
		strings.HasPrefix(line, "--- SKIP: "+test) ||
		strings.HasPrefix(line, "ok  \t"+pkg) ||
		strings.HasPrefix(line, "FAIL\t"+pkg) ||
		strings.HasPrefix(line, "?   \t"+pkg)
}

func (r *Renderer) renderPackageLine(event Event) string {
	subject := event.Package
	if subject == "" {
		subject = "(package)"
	}
	return r.renderLine("pkg", event.Action, subject, event.Elapsed)
}

func (r *Renderer) renderTestLine(event Event) string {
	subject := event.Package
	if subject == "" {
		subject = "(package)"
	}
	if event.Test != "" {
		subject += " :: " + event.Test
	}
	return r.renderLine("test", event.Action, subject, event.Elapsed)
}

func (r *Renderer) renderLine(kind string, action Action, subject string, elapsed float64) string {
	status := strings.ToUpper(string(action))
	if kind == "pkg" {
		status = "PKG " + status
	}
	duration := fmt.Sprintf("%.2fs", elapsed)

	if r.mode.Format != laslig.FormatHuman {
		return fmt.Sprintf("[%s] %s (%s)", status, subject, duration)
	}

	statusStyle := r.theme.NoticeInfo
	switch action {
	case ActionPass:
		statusStyle = r.theme.NoticeSuccess
	case ActionSkip:
		statusStyle = r.theme.NoticeWarning
	case ActionFail:
		statusStyle = r.theme.NoticeError
	}

	subjectStyle := r.theme.Value
	if kind == "pkg" {
		subjectStyle = r.theme.Identifier
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		statusStyle.Render(status),
		" ",
		subjectStyle.Render(subject),
		" ",
		r.theme.Muted.Render("("+duration+")"),
	)
}

func (r *Renderer) writeLine(value string) error {
	r.wroteResults = true
	if _, err := fmt.Fprintln(r.out, value); err != nil {
		return fmt.Errorf("write event line: %w", err)
	}
	return nil
}

func plural(value int) string {
	if value == 1 {
		return ""
	}
	return "s"
}
