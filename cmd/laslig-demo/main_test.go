package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

// failWriter is an io.Writer that always fails.
type failWriter struct{}

// Write implements io.Writer by always returning an error.
func (failWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("boom")
}

// TestRunArgsPlain verifies plain demo rendering through the testable entrypoint.
func TestRunArgsPlain(t *testing.T) {
	var buf bytes.Buffer
	err := runArgs(&buf, []string{"-format", "plain", "-style", "never"})
	if err != nil {
		t.Fatalf("runArgs() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "laslig demo") {
		t.Fatalf("runArgs() output missing section:\n%s", got)
	}
	if !strings.Contains(got, "Policy") {
		t.Fatalf("runArgs() output missing kv block:\n%s", got)
	}
	if !strings.Contains(got, "testjson [LIVE]") {
		t.Fatalf("runArgs() output missing live badge:\n%s", got)
	}
}

// TestRunArgsJSON verifies JSON demo rendering through the testable entrypoint.
func TestRunArgsJSON(t *testing.T) {
	var buf bytes.Buffer
	err := runArgs(&buf, []string{"-format", "json"})
	if err != nil {
		t.Fatalf("runArgs() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"type": "section"`) {
		t.Fatalf("runArgs() JSON output missing section event:\n%s", got)
	}
	if !strings.Contains(got, `"type": "kv"`) {
		t.Fatalf("runArgs() JSON output missing kv event:\n%s", got)
	}
	if !strings.Contains(got, `"type": "panel"`) {
		t.Fatalf("runArgs() JSON output missing panel event:\n%s", got)
	}
}

// TestRunArgsHumanStyled verifies forced human styling output.
func TestRunArgsHumanStyled(t *testing.T) {
	var buf bytes.Buffer
	err := runArgs(&buf, []string{"-format", "human", "-style", "always"})
	if err != nil {
		t.Fatalf("runArgs() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "\x1b[") {
		t.Fatalf("runArgs() output missing ANSI styling: %q", got)
	}
	if !strings.Contains(got, "LIVE") {
		t.Fatalf("runArgs() output missing live badge text: %q", got)
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

// TestRunArgsRenderError verifies render failures are wrapped with the step name.
func TestRunArgsRenderError(t *testing.T) {
	err := runArgs(failWriter{}, []string{"-format", "plain", "-style", "never"})
	if err == nil {
		t.Fatal("runArgs() error = nil, want render error")
	}
	if !strings.Contains(err.Error(), "render section") {
		t.Fatalf("runArgs() error = %v, want render section prefix", err)
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
	os.Args = []string{"laslig-demo", "-format", "plain", "-style", "never"}
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

	if !strings.Contains(string(data), "laslig demo") {
		t.Fatalf("main() output missing section:\n%s", string(data))
	}
}
