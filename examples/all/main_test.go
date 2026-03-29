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

// TestRunArgsPlain verifies plain aggregate showcase rendering.
func TestRunArgsPlain(t *testing.T) {
	var buf bytes.Buffer
	if err := runArgs(&buf, []string{"-format", "plain", "-style", "never"}); err != nil {
		t.Fatalf("runArgs() error = %v", err)
	}

	got := buf.String()
	for _, want := range []string{
		"Läslig demo",
		"Section",
		"Notice",
		"Record",
		"KV",
		"List",
		"Table",
		"Panel",
		"Paragraph",
		"StatusLine",
		"Markdown",
		"CodeBlock",
		"LogBlock",
		"gotestout",
		"gotestout + Mage",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("runArgs() output missing %q:\n%s", want, got)
		}
	}
}

// TestRunArgsJSON verifies aggregate JSON rendering through the testable entrypoint.
func TestRunArgsJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := runArgs(&buf, []string{"-format", "json"}); err != nil {
		t.Fatalf("runArgs() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"type": "section"`) {
		t.Fatalf("runArgs() JSON output missing section event:\n%s", got)
	}
	if !strings.Contains(got, `"type": "record"`) {
		t.Fatalf("runArgs() JSON output missing record event:\n%s", got)
	}
	if !strings.Contains(got, `"type": "code_block"`) {
		t.Fatalf("runArgs() JSON output missing code_block event:\n%s", got)
	}
	if !strings.Contains(got, `"Action":"fail"`) {
		t.Fatalf("runArgs() JSON output missing gotestout events:\n%s", got)
	}
}

// TestRunArgsHumanStyled verifies forced human styling output.
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
	if !strings.Contains(plain, "Läslig demo") {
		t.Fatalf("runArgs() output missing intro: %q", plain)
	}
	if !strings.Contains(plain, "Use gotestout for attractive, structured go test output when your task runner") {
		t.Fatalf("runArgs() output missing gotestout section: %q", plain)
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

// TestRunArgsRenderError verifies render failures are wrapped by the shared runner.
func TestRunArgsRenderError(t *testing.T) {
	err := runArgs(failWriter{}, []string{"-format", "plain", "-style", "never"})
	if err == nil {
		t.Fatal("runArgs() error = nil, want render error")
	}
	if !strings.Contains(err.Error(), "render laslig-demo example") {
		t.Fatalf("runArgs() error = %v, want shared runner render prefix", err)
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

	if !strings.Contains(string(data), "Läslig demo") {
		t.Fatalf("main() output missing intro:\n%s", string(data))
	}
}
