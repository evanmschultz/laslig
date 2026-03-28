package glamrender

import (
	"strings"
	"testing"
)

// TestRender verifies Glamour-backed rendering preserves markdown semantics and emits ANSI output.
func TestRender(t *testing.T) {
	rendered, err := Render("# Heading\n\n- first\n- second", 80)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !strings.Contains(rendered, "Heading") {
		t.Fatalf("Render() = %q, want heading text", rendered)
	}
	if !strings.Contains(rendered, "first") || !strings.Contains(rendered, "second") {
		t.Fatalf("Render() = %q, want list items", rendered)
	}
	if !strings.Contains(rendered, "\x1b[") {
		t.Fatalf("Render() = %q, want ANSI styling", rendered)
	}
}

// TestFencedCodeBlock verifies fenced code creation preserves the language tag and body.
func TestFencedCodeBlock(t *testing.T) {
	got := FencedCodeBlock("go", "fmt.Println(\"hi\")\n")
	want := "```go\nfmt.Println(\"hi\")\n```"
	if got != want {
		t.Fatalf("FencedCodeBlock() = %q, want %q", got, want)
	}
}
