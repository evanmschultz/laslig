package laslig

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strings"
	"testing"
)

var spinnerANSIPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// TestSpinnerHumanStyledTransient verifies styled human spinners redraw in place
// and stop with one durable semantic status line.
func TestSpinnerHumanStyledTransient(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: true, Width: 72})
	spin := printer.NewSpinner()
	spin.forceAnimation = true

	if err := spin.Start("Waiting for rollout"); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := spin.Update("Waiting for rollout (2/3)"); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if err := spin.Stop("Rollout ready", NoticeSuccessLevel); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "\r") {
		t.Fatalf("spinner output = %q, want carriage-return redraws", got)
	}
	if !strings.Contains(got, "\x1b[") {
		t.Fatalf("spinner output = %q, want ANSI styling", got)
	}
	plain := spinnerANSIPattern.ReplaceAllString(got, "")
	if !strings.Contains(plain, "Waiting for rollout") {
		t.Fatalf("spinner output missing running text:\n%s", plain)
	}
	if !strings.Contains(plain, "Rollout ready") {
		t.Fatalf("spinner output missing final text:\n%s", plain)
	}
}

// TestSpinnerHumanNoStyleFallback verifies unstyled human output stays
// human-formatted but emits stable plain status rows with no ANSI.
func TestSpinnerHumanNoStyleFallback(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatHuman, Styled: false})
	spin := printer.NewSpinner()

	if err := spin.Start("Waiting for rollout"); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := spin.Update("Waiting for rollout (2/3)"); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if err := spin.Stop("Rollout ready", NoticeSuccessLevel); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	want := "[RUNNING] Waiting for rollout\n[SUCCESS] Rollout ready\n"
	if got := buf.String(); got != want {
		t.Fatalf("spinner output = %q, want %q", got, want)
	}
	if strings.Contains(buf.String(), "\x1b[") {
		t.Fatalf("spinner output = %q, want no ANSI", buf.String())
	}
}

// TestSpinnerPlainFallback verifies plain output stays stable and non-animated.
func TestSpinnerPlainFallback(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatPlain})
	spin := printer.NewSpinner()

	if err := spin.Start("Waiting for rollout"); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := spin.Stop("Rollout ready", NoticeSuccessLevel); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	want := "[RUNNING] Waiting for rollout\n[SUCCESS] Rollout ready\n"
	if got := buf.String(); got != want {
		t.Fatalf("spinner output = %q, want %q", got, want)
	}
}

// TestSpinnerJSON verifies JSON mode emits stable start and finish status-line
// records without transient update frames.
func TestSpinnerJSON(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatJSON})
	spin := printer.NewSpinner()

	if err := spin.Start("Waiting for rollout"); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := spin.Update("Waiting for rollout (2/3)"); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if err := spin.Stop("Rollout ready", NoticeSuccessLevel); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	decoder := json.NewDecoder(&buf)
	var start map[string]any
	if err := decoder.Decode(&start); err != nil {
		t.Fatalf("decode start = %v", err)
	}
	var stop map[string]any
	if err := decoder.Decode(&stop); err != nil {
		t.Fatalf("decode stop = %v", err)
	}
	if err := decoder.Decode(&map[string]any{}); err != io.EOF {
		t.Fatalf("extra JSON records: %v", err)
	}

	if got := start["type"]; got != "status_line" {
		t.Fatalf("start type = %v, want status_line", got)
	}
	if got := stop["type"]; got != "status_line" {
		t.Fatalf("stop type = %v, want status_line", got)
	}
	startPayload := start["payload"].(map[string]any)
	stopPayload := stop["payload"].(map[string]any)
	if got := startPayload["label"]; got != "running" {
		t.Fatalf("start label = %v, want running", got)
	}
	if got := startPayload["text"]; got != "Waiting for rollout" {
		t.Fatalf("start text = %v, want waiting text", got)
	}
	if got := stopPayload["text"]; got != "Rollout ready" {
		t.Fatalf("stop text = %v, want final text", got)
	}
	if _, ok := stopPayload["label"]; ok {
		t.Fatalf("stop payload label = %v, want omitted", stopPayload["label"])
	}
}

// TestSpinnerStopWithoutStart verifies delayed-start callers can stop safely
// when the spinner never actually started.
func TestSpinnerStopWithoutStart(t *testing.T) {
	var buf bytes.Buffer
	printer := newTestPrinter(&buf, Mode{Format: FormatPlain})
	spin := printer.NewSpinner()

	if err := spin.Stop("Done", NoticeSuccessLevel); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	if got := buf.String(); got != "" {
		t.Fatalf("spinner output = %q, want empty output", got)
	}
}

// TestTruncateVisible verifies animated spinner text truncates to one visible line.
func TestTruncateVisible(t *testing.T) {
	got := truncateVisible("alpha beta gamma", 9)
	if got != "alpha be…" {
		t.Fatalf("truncateVisible() = %q, want %q", got, "alpha be…")
	}
}
