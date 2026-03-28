package logging

import (
	"strings"
	"testing"
)

// TestTranscript verifies the example transcript contains the expected log events.
func TestTranscript(t *testing.T) {
	got := Transcript()
	if !strings.Contains(got, "boot complete") {
		t.Fatalf("Transcript() = %q, want boot complete", got)
	}
	if !strings.Contains(got, "retry scheduled") {
		t.Fatalf("Transcript() = %q, want retry scheduled", got)
	}
	if !strings.Contains(got, "dependency missing") {
		t.Fatalf("Transcript() = %q, want dependency missing", got)
	}
}
