package layout

import "testing"

// TestWrapText verifies paragraph wrapping preserves paragraph boundaries.
func TestWrapText(t *testing.T) {
	got := WrapText("alpha beta gamma delta\nsecond paragraph", 10)
	want := "alpha beta\ngamma\ndelta\nsecond\nparagraph"
	if got != want {
		t.Fatalf("WrapText() = %q, want %q", got, want)
	}
}

// TestIndentBlock verifies each rendered line receives the same prefix.
func TestIndentBlock(t *testing.T) {
	got := IndentBlock("  ", "one\ntwo")
	want := "  one\n  two"
	if got != want {
		t.Fatalf("IndentBlock() = %q, want %q", got, want)
	}
}
