package testjson

import (
	"bytes"
	"strings"
	"testing"

	"github.com/evanmschultz/laslig"
)

// sampleStream exercises passing, skipped, failing, and build-error events.
const sampleStream = `{"Action":"run","Package":"example/pkg","Test":"TestPass"}
{"Action":"output","Package":"example/pkg","Test":"TestPass","Output":"=== RUN   TestPass\n"}
{"Action":"output","Package":"example/pkg","Test":"TestPass","Output":"note: useful output\n"}
{"Action":"output","Package":"example/pkg","Test":"TestPass","Output":"--- PASS: TestPass (0.01s)\n"}
{"Action":"pass","Package":"example/pkg","Test":"TestPass","Elapsed":0.01}
{"Action":"run","Package":"example/pkg","Test":"TestSkip"}
{"Action":"output","Package":"example/pkg","Test":"TestSkip","Output":"--- SKIP: TestSkip (0.00s)\n"}
{"Action":"skip","Package":"example/pkg","Test":"TestSkip","Elapsed":0}
{"Action":"run","Package":"example/pkg","Test":"TestFail"}
{"Action":"output","Package":"example/pkg","Test":"TestFail","Output":"=== RUN   TestFail\n"}
{"Action":"output","Package":"example/pkg","Test":"TestFail","Output":"renderer_test.go:42: expected boom\n"}
{"Action":"output","Package":"example/pkg","Test":"TestFail","Output":"--- FAIL: TestFail (0.02s)\n"}
{"Action":"fail","Package":"example/pkg","Test":"TestFail","Elapsed":0.02}
{"Action":"output","Package":"example/pkg","Output":"FAIL\texample/pkg [build failed]\n","FailedBuild":"example/pkg"}
{"Action":"fail","Package":"example/pkg","Elapsed":0.03}
`

// TestParse verifies that Parse decodes all events from a JSON stream.
func TestParse(t *testing.T) {
	events, err := Parse(strings.NewReader(sampleStream))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if got, want := len(events), 15; got != want {
		t.Fatalf("len(Parse()) = %d, want %d", got, want)
	}
	if got, want := events[0].Action, ActionRun; got != want {
		t.Fatalf("events[0].Action = %q, want %q", got, want)
	}
	if got, want := events[len(events)-1].Action, ActionFail; got != want {
		t.Fatalf("events[last].Action = %q, want %q", got, want)
	}
}

// TestRenderPlainCompact verifies the compact plain renderer and summary counts.
func TestRenderPlainCompact(t *testing.T) {
	var buf bytes.Buffer
	summary, err := Render(&buf, strings.NewReader(sampleStream), Options{
		Policy: laslig.Policy{
			Format: laslig.FormatPlain,
			Style:  laslig.StyleNever,
		},
		View: ViewCompact,
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if summary.TestsPassed != 1 || summary.TestsFailed != 1 || summary.TestsSkipped != 1 {
		t.Fatalf("summary tests = %+v, want 1 pass, 1 fail, 1 skip", summary)
	}
	if summary.PackagesFailed != 1 || summary.BuildErrors != 1 {
		t.Fatalf("summary packages = %+v, want 1 failed package with 1 build error", summary)
	}

	got := buf.String()
	if strings.Contains(got, "[PASS] example/pkg :: TestPass (0.01s)") {
		t.Fatalf("Render() compact output unexpectedly included pass test line:\n%s", got)
	}
	if !strings.Contains(got, "[FAIL] example/pkg :: TestFail (0.02s)") {
		t.Fatalf("Render() output missing fail line:\n%s", got)
	}
	if !strings.Contains(got, "[PKG FAIL] example/pkg (0.03s)") {
		t.Fatalf("Render() output missing package fail line:\n%s", got)
	}
	if !strings.Contains(got, "renderer_test.go:42: expected boom") {
		t.Fatalf("Render() output missing failure output:\n%s", got)
	}
	if !strings.Contains(got, "Failed tests") {
		t.Fatalf("Render() output missing failed tests section:\n%s", got)
	}
	if !strings.Contains(got, "Package errors") {
		t.Fatalf("Render() output missing package errors section:\n%s", got)
	}
	if !strings.Contains(got, "Skipped tests") {
		t.Fatalf("Render() output missing skipped tests section:\n%s", got)
	}
	if !strings.Contains(got, "Test summary") {
		t.Fatalf("Render() output missing summary heading:\n%s", got)
	}
	if !strings.Contains(got, "[ERROR] Test failures detected") {
		t.Fatalf("Render() output missing error notice:\n%s", got)
	}
}

// TestRenderPlainDetailed verifies that the detailed renderer keeps useful passing output.
func TestRenderPlainDetailed(t *testing.T) {
	var buf bytes.Buffer
	_, err := Render(&buf, strings.NewReader(sampleStream), Options{
		Policy: laslig.Policy{
			Format: laslig.FormatPlain,
			Style:  laslig.StyleNever,
		},
		View: ViewDetailed,
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "note: useful output") {
		t.Fatalf("Render() detailed output missing passing test output:\n%s", got)
	}
	if !strings.Contains(got, "[PASS] example/pkg :: TestPass (0.01s)") {
		t.Fatalf("Render() detailed output missing pass line:\n%s", got)
	}
	if !strings.Contains(got, "[SKIP] example/pkg :: TestSkip (0.00s)") {
		t.Fatalf("Render() detailed output missing skip line:\n%s", got)
	}
}

// TestRenderPlainDetailedDisabledSections verifies callers can trim live output and grouped sections.
func TestRenderPlainDetailedDisabledSections(t *testing.T) {
	var buf bytes.Buffer
	_, err := Render(&buf, strings.NewReader(sampleStream), Options{
		Policy: laslig.Policy{
			Format: laslig.FormatPlain,
			Style:  laslig.StyleNever,
		},
		View: ViewDetailed,
		DisabledSections: []Section{
			SectionFailedTests,
			SectionSkippedTests,
			SectionPackageErrors,
			SectionOutput,
		},
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "[PASS] example/pkg :: TestPass (0.01s)") {
		t.Fatalf("Render() detailed output missing pass line:\n%s", got)
	}
	if strings.Contains(got, "note: useful output") {
		t.Fatalf("Render() output unexpectedly included captured output:\n%s", got)
	}
	if strings.Contains(got, "renderer_test.go:42: expected boom") {
		t.Fatalf("Render() output unexpectedly included failure detail:\n%s", got)
	}
	if strings.Contains(got, "Failed tests") {
		t.Fatalf("Render() output unexpectedly included failed tests section:\n%s", got)
	}
	if strings.Contains(got, "Package errors") {
		t.Fatalf("Render() output unexpectedly included package errors section:\n%s", got)
	}
	if strings.Contains(got, "Skipped tests") {
		t.Fatalf("Render() output unexpectedly included skipped tests section:\n%s", got)
	}
}

// TestRenderPlainCompactDisabledOutput verifies grouped sections remain while captured output is suppressed.
func TestRenderPlainCompactDisabledOutput(t *testing.T) {
	var buf bytes.Buffer
	_, err := Render(&buf, strings.NewReader(sampleStream), Options{
		Policy: laslig.Policy{
			Format: laslig.FormatPlain,
			Style:  laslig.StyleNever,
		},
		View:             ViewCompact,
		DisabledSections: []Section{SectionOutput},
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "Package errors") {
		t.Fatalf("Render() output missing package errors section:\n%s", got)
	}
	if strings.Contains(got, "detail:") {
		t.Fatalf("Render() output unexpectedly included grouped detail field:\n%s", got)
	}
	if strings.Contains(got, "renderer_test.go:42: expected boom") {
		t.Fatalf("Render() output unexpectedly included failure output:\n%s", got)
	}
}

// TestRenderJSON verifies that JSON mode re-emits events while still returning summary counts.
func TestRenderJSON(t *testing.T) {
	var buf bytes.Buffer
	summary, err := Render(&buf, strings.NewReader(sampleStream), Options{
		Policy: laslig.Policy{Format: laslig.FormatJSON},
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if summary.TotalTests() != 3 {
		t.Fatalf("summary.TotalTests() = %d, want 3", summary.TotalTests())
	}
	if !strings.Contains(buf.String(), `"Action":"pass"`) {
		t.Fatalf("Render() JSON output missing encoded events:\n%s", buf.String())
	}
	if strings.Contains(buf.String(), "Test summary") {
		t.Fatalf("Render() JSON output should not include human summary:\n%s", buf.String())
	}
}
