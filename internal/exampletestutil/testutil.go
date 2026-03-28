package exampletestutil

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

// ansiPattern matches ANSI escape sequences for stable assertions.
var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// StripANSI removes ANSI escape sequences for stable assertions.
func StripANSI(value string) string {
	return ansiPattern.ReplaceAllString(value, "")
}

// AssertRunArgsPlainContains verifies one example renders expected plain output.
func AssertRunArgsPlainContains(t *testing.T, runArgs func(io.Writer, []string) error, want ...string) {
	t.Helper()

	var buf bytes.Buffer
	if err := runArgs(&buf, []string{"-format", "plain", "-style", "never"}); err != nil {
		t.Fatalf("runArgs() error = %v", err)
	}
	got := buf.String()
	for _, needle := range want {
		if !strings.Contains(got, needle) {
			t.Fatalf("runArgs() output missing %q:\n%s", needle, got)
		}
	}
}

// AssertRunArgsHumanStyled verifies one example emits styled human output.
func AssertRunArgsHumanStyled(t *testing.T, runArgs func(io.Writer, []string) error, want ...string) {
	t.Helper()

	var buf bytes.Buffer
	if err := runArgs(&buf, []string{"-format", "human", "-style", "always"}); err != nil {
		t.Fatalf("runArgs() error = %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "\x1b[") {
		t.Fatalf("runArgs() output missing ANSI styling: %q", got)
	}
	plain := StripANSI(got)
	for _, needle := range want {
		if !strings.Contains(plain, needle) {
			t.Fatalf("runArgs() stripped output missing %q:\n%s", needle, plain)
		}
	}
}

// AssertRunArgsInvalidFlag verifies one example returns a parse error.
func AssertRunArgsInvalidFlag(t *testing.T, runArgs func(io.Writer, []string) error) {
	t.Helper()

	err := runArgs(io.Discard, []string{"-unknown"})
	if err == nil {
		t.Fatal("runArgs() error = nil, want parse error")
	}
	if !strings.Contains(err.Error(), "parse flags") {
		t.Fatalf("runArgs() error = %v, want parse flags prefix", err)
	}
}

// AssertMainContains verifies one example main entrypoint writes expected text.
func AssertMainContains(t *testing.T, mainFunc func(), command string, want ...string) {
	t.Helper()

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
	os.Args = []string{command, "-format", "plain", "-style", "never"}
	os.Stdout = writer

	mainFunc()

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

	got := string(data)
	for _, needle := range want {
		if !strings.Contains(got, needle) {
			t.Fatalf("main() output missing %q:\n%s", needle, got)
		}
	}
}
