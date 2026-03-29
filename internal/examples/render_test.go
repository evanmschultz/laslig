package examples

import (
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/x/exp/golden"
	"github.com/evanmschultz/laslig"
)

// ansiPattern matches ANSI escape sequences for stable styled-output assertions.
var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// stripANSI removes ANSI escape sequences for stable golden snapshots.
func stripANSI(value string) string {
	return ansiPattern.ReplaceAllString(value, "")
}

// TestRunFocusedPlainGolden verifies every focused example against a plain golden snapshot.
func TestRunFocusedPlainGolden(t *testing.T) {
	tests := []struct {
		name   string
		render Renderer
	}{
		{name: "section", render: RenderSection},
		{name: "notice", render: RenderNotice},
		{name: "record", render: RenderRecord},
		{name: "kv", render: RenderKV},
		{name: "list", render: RenderList},
		{name: "table", render: RenderTable},
		{name: "panel", render: RenderPanel},
		{name: "paragraph", render: RenderParagraph},
		{name: "statusline", render: RenderStatusLine},
		{name: "markdown", render: RenderMarkdown},
		{name: "codeblock", render: RenderCodeBlock},
		{name: "logblock", render: RenderLogBlock},
		{name: "gotestout", render: RenderGotestout},
		{name: "magecheck", render: RenderMageCheckPreview},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := Run(&buf, []string{"-format", "plain", "-style", "never"}, tc.name, tc.render); err != nil {
				t.Fatalf("Run(%q) error = %v", tc.name, err)
			}

			golden.RequireEqual(t, buf.Bytes())
		})
	}
}

// TestRunInvalidFlag verifies the shared runner wraps parse failures.
func TestRunInvalidFlag(t *testing.T) {
	err := Run(&bytes.Buffer{}, []string{"-unknown"}, "notice", RenderNotice)
	if err == nil {
		t.Fatal("Run() error = nil, want parse error")
	}
	if !strings.Contains(err.Error(), "parse flags") {
		t.Fatalf("Run() error = %v, want parse flags prefix", err)
	}
}

// TestRunInvalidGlamourStyle verifies the shared runner rejects unsupported
// built-in Glamour style names.
func TestRunInvalidGlamourStyle(t *testing.T) {
	err := Run(&bytes.Buffer{}, []string{"-glamour-style", "bogus"}, "markdown", RenderMarkdown)
	if err == nil {
		t.Fatal("Run() error = nil, want glamour style error")
	}
	if !strings.Contains(err.Error(), `invalid glamour style "bogus"`) {
		t.Fatalf("Run() error = %v, want invalid glamour style message", err)
	}
}

// TestRunRenderError verifies the shared runner wraps renderer failures.
func TestRunRenderError(t *testing.T) {
	err := Run(&bytes.Buffer{}, []string{"-format", "plain", "-style", "never"}, "boom", func(io.Writer, *laslig.Printer) error {
		return errors.New("boom")
	})
	if err == nil {
		t.Fatal("Run() error = nil, want render error")
	}
	if !strings.Contains(err.Error(), "render boom example") {
		t.Fatalf("Run() error = %v, want wrapped render prefix", err)
	}
}

// TestMainExit verifies the shared main helper reports failures through stderr and exit.
func TestMainExit(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exitCode := 0

	Main(&out, &errOut, []string{"-unknown"}, func(code int) {
		exitCode = code
	}, "notice", RenderNotice)

	if exitCode != 1 {
		t.Fatalf("Main() exitCode = %d, want 1", exitCode)
	}
	if !strings.Contains(errOut.String(), "parse flags") {
		t.Fatalf("Main() stderr missing parse failure:\n%s", errOut.String())
	}
}

// TestRunAllPlainGolden verifies the aggregate demo against a plain golden snapshot.
func TestRunAllPlainGolden(t *testing.T) {
	var buf bytes.Buffer
	if err := Run(&buf, []string{"-format", "plain", "-style", "never"}, "all", RenderAll); err != nil {
		t.Fatalf("Run(all) error = %v", err)
	}

	golden.RequireEqual(t, buf.Bytes())
}

// TestRenderAllHumanStyledGolden verifies fixed-width human output for the aggregate demo.
func TestRenderAllHumanStyledGolden(t *testing.T) {
	var buf bytes.Buffer
	printer := laslig.NewWithMode(&buf, laslig.Mode{
		Format: laslig.FormatHuman,
		Styled: true,
		Width:  80,
	})

	if err := RenderAll(&buf, printer); err != nil {
		t.Fatalf("RenderAll() error = %v", err)
	}

	golden.RequireEqual(t, []byte(stripANSI(buf.String())))
}

// TestStylePolicyForMode verifies mode-derived style policy choices.
func TestStylePolicyForMode(t *testing.T) {
	if got := StylePolicyForMode(laslig.Mode{Styled: true}); got != laslig.StyleAlways {
		t.Fatalf("StylePolicyForMode(styled=true) = %q, want %q", got, laslig.StyleAlways)
	}
	if got := StylePolicyForMode(laslig.Mode{Styled: false}); got != laslig.StyleNever {
		t.Fatalf("StylePolicyForMode(styled=false) = %q, want %q", got, laslig.StyleNever)
	}
}

// TestMageCheckSampleStream verifies the Mage preview fixture includes the split example packages.
func TestMageCheckSampleStream(t *testing.T) {
	got := mageCheckSampleStream()
	for _, want := range []string{
		"github.com/evanmschultz/laslig/examples/section",
		"github.com/evanmschultz/laslig/examples/notice",
		"github.com/evanmschultz/laslig/examples/magecheck",
		"github.com/evanmschultz/laslig/internal/examples",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("mageCheckSampleStream() missing %q", want)
		}
	}
}

// TestTranscript verifies the shared LogBlock transcript remains deterministic.
func TestTranscript(t *testing.T) {
	got := transcript()
	for _, want := range []string{
		"INFO demo: boot complete",
		"WARN demo: retry scheduled",
		"ERRO demo: dependency missing",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("transcript() missing %q in %q", want, got)
		}
	}
}

// TestTTYBufferFd verifies the ttyBuffer reports the current stdout descriptor.
func TestTTYBufferFd(t *testing.T) {
	var writer ttyBuffer
	if got, want := writer.Fd(), os.Stdout.Fd(); got != want {
		t.Fatalf("ttyBuffer.Fd() = %d, want %d", got, want)
	}
}
