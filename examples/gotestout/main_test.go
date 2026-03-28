package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

// ansiPattern matches ANSI escape sequences for stable test assertions.
var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// failWriter is an io.Writer that always fails.
type failWriter struct{}

// Write implements io.Writer by always returning an error.
func (failWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("boom")
}

// stripANSI removes ANSI escape sequences from one string for stable assertions.
func stripANSI(value string) string {
	return ansiPattern.ReplaceAllString(value, "")
}

// TestRunArgsPlain verifies focused plain gotestout rendering.
func TestRunArgsPlain(t *testing.T) {
	var buf bytes.Buffer
	if err := runArgs(&buf, []string{"-format", "plain", "-style", "never"}); err != nil {
		t.Fatalf("runArgs() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "[FAIL] example/pkg :: TestFail (0.02s)") {
		t.Fatalf("runArgs() output missing failed test line:\n%s", got)
	}
	if !strings.Contains(got, "Package errors") {
		t.Fatalf("runArgs() output missing package errors section:\n%s", got)
	}
	if !strings.Contains(got, "Test summary") {
		t.Fatalf("runArgs() output missing summary:\n%s", got)
	}
}

// TestRunArgsJSON verifies focused JSON gotestout rendering.
func TestRunArgsJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := runArgs(&buf, []string{"-format", "json"}); err != nil {
		t.Fatalf("runArgs() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"Action":"fail"`) {
		t.Fatalf("runArgs() JSON output missing encoded fail event:\n%s", got)
	}
	if strings.Contains(got, "Test summary") {
		t.Fatalf("runArgs() JSON output unexpectedly included human summary:\n%s", got)
	}
}

// TestRunArgsHumanStyled verifies focused human/styled gotestout rendering.
func TestRunArgsHumanStyled(t *testing.T) {
	var buf bytes.Buffer
	if err := runArgs(&buf, []string{"-format", "human", "-style", "always"}); err != nil {
		t.Fatalf("runArgs() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "\x1b[") {
		t.Fatalf("runArgs() output missing ANSI styling: %q", got)
	}
	plain := stripANSI(got)
	if !strings.Contains(plain, "Test summary") {
		t.Fatalf("runArgs() stripped output missing summary:\n%s", plain)
	}
	if !strings.Contains(plain, "Failed tests") {
		t.Fatalf("runArgs() stripped output missing failed tests section:\n%s", plain)
	}
}

// TestRunArgsInvalidFlag verifies invalid arguments return a parse error.
func TestRunArgsInvalidFlag(t *testing.T) {
	err := runArgs(io.Discard, []string{"-unknown"})
	if err == nil {
		t.Fatal("runArgs() error = nil, want parse error")
	}
	if !strings.Contains(err.Error(), "parse flags") {
		t.Fatalf("runArgs() error = %v, want parse flags prefix", err)
	}
}

// TestRunArgsRenderError verifies render failures are wrapped.
func TestRunArgsRenderError(t *testing.T) {
	err := runArgs(failWriter{}, []string{"-format", "plain", "-style", "never"})
	if err == nil {
		t.Fatal("runArgs() error = nil, want render error")
	}
	if !strings.Contains(err.Error(), "render gotestout example") {
		t.Fatalf("runArgs() error = %v, want render prefix", err)
	}
}

// TestMain verifies the command entrypoint succeeds for a valid invocation.
func TestMain(t *testing.T) {
	oldArgs := os.Args
	oldStdout := os.Stdout
	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
	}()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Args = []string{"gotestout-example", "-format", "plain", "-style", "never"}
	os.Stdout = writer

	main()

	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close() error = %v", err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("reader.Close() error = %v", err)
	}

	if !strings.Contains(string(data), "Test summary") {
		t.Fatalf("main() output missing summary:\n%s", string(data))
	}
}
