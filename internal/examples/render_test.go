package examples

import (
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/evanmschultz/laslig"
)

// ansiPattern matches ANSI CSI escape sequences for stable styled-output assertions.
var ansiPattern = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)

type failingWriter struct{}

func (failingWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("boom")
}

type fakeFileWriter struct {
	bytes.Buffer
	fd uintptr
}

func (w fakeFileWriter) Fd() uintptr {
	return w.fd
}

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
		{name: "spinner", render: RenderSpinner},
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

// TestRunFramedHumanStyledGolden verifies framed examples across width and wrap
// combinations using the public example flags.
func TestRunFramedHumanStyledGolden(t *testing.T) {
	t.Setenv("COLUMNS", "72")

	tests := []struct {
		name   string
		render Renderer
		args   []string
	}{
		{name: "table_default", render: RenderTable, args: []string{"-format", "human", "-style", "always"}},
		{name: "table_long_auto", render: RenderTable, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "58"}},
		{name: "table_long_truncate", render: RenderTable, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "48", "-wrap-mode", "truncate"}},
		{name: "table_long_never", render: RenderTable, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "48", "-wrap-mode", "never"}},
		{name: "panel_default", render: RenderPanel, args: []string{"-format", "human", "-style", "always"}},
		{name: "panel_long_auto", render: RenderPanel, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "58"}},
		{name: "panel_long_truncate", render: RenderPanel, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "48", "-wrap-mode", "truncate"}},
		{name: "panel_long_never", render: RenderPanel, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "48", "-wrap-mode", "never"}},
		{name: "codeblock_default", render: RenderCodeBlock, args: []string{"-format", "human", "-style", "always"}},
		{name: "codeblock_long_auto", render: RenderCodeBlock, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "58"}},
		{name: "codeblock_long_truncate", render: RenderCodeBlock, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "48", "-wrap-mode", "truncate"}},
		{name: "codeblock_long_never", render: RenderCodeBlock, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "48", "-wrap-mode", "never"}},
		{name: "logblock_default", render: RenderLogBlock, args: []string{"-format", "human", "-style", "always"}},
		{name: "logblock_long_auto", render: RenderLogBlock, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "58"}},
		{name: "logblock_long_truncate", render: RenderLogBlock, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "48", "-wrap-mode", "truncate"}},
		{name: "logblock_long_never", render: RenderLogBlock, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "48", "-wrap-mode", "never"}},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := Run(&buf, tc.args, tc.name, tc.render); err != nil {
				t.Fatalf("Run(%q) error = %v", tc.name, err)
			}

			value := stripANSI(buf.String())
			for _, line := range strings.Split(strings.TrimRight(value, "\n"), "\n") {
				if lipgloss.Width(line) > 72 {
					t.Fatalf("styled example line exceeded width budget: %q (%d > 72)", line, lipgloss.Width(line))
				}
			}

			golden.RequireEqual(t, []byte(value))
		})
	}
}

// TestRunFramedExamplesAdaptAcrossColumns verifies the public example runner
// stays inside varying terminal widths so layout regressions are obvious.
func TestRunFramedExamplesAdaptAcrossColumns(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		render Renderer
		args   []string
	}{
		{name: "table_auto_72", width: 72, render: RenderTable, args: []string{"-format", "human", "-style", "always", "-content", "long", "-wrap-mode", "auto"}},
		{name: "table_auto_56", width: 56, render: RenderTable, args: []string{"-format", "human", "-style", "always", "-content", "long", "-wrap-mode", "auto"}},
		{name: "table_truncate_48", width: 48, render: RenderTable, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "48", "-wrap-mode", "truncate"}},
		{name: "panel_auto_56", width: 56, render: RenderPanel, args: []string{"-format", "human", "-style", "always", "-content", "long", "-wrap-mode", "auto"}},
		{name: "codeblock_truncate_48", width: 48, render: RenderCodeBlock, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "48", "-wrap-mode", "truncate"}},
		{name: "logblock_never_48", width: 48, render: RenderLogBlock, args: []string{"-format", "human", "-style", "always", "-content", "long", "-max-width", "48", "-wrap-mode", "never"}},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("COLUMNS", strconv.Itoa(tc.width))

			var buf bytes.Buffer
			if err := Run(&buf, tc.args, tc.name, tc.render); err != nil {
				t.Fatalf("Run(%q) error = %v", tc.name, err)
			}

			value := stripANSI(buf.String())
			for _, line := range strings.Split(strings.TrimRight(value, "\n"), "\n") {
				if lipgloss.Width(line) > tc.width {
					t.Fatalf("example line exceeded width %d: %q (%d)", tc.width, line, lipgloss.Width(line))
				}
			}
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

func TestRunInvalidWrapMode(t *testing.T) {
	err := Run(&bytes.Buffer{}, []string{"-wrap-mode", "bogus"}, "table", RenderTable)
	if err == nil {
		t.Fatal("Run() error = nil, want wrap mode error")
	}
	if !strings.Contains(err.Error(), `invalid wrap mode "bogus"`) {
		t.Fatalf("Run() error = %v, want invalid wrap mode message", err)
	}
}

func TestRunInvalidContentMode(t *testing.T) {
	err := Run(&bytes.Buffer{}, []string{"-content", "bogus"}, "table", RenderTable)
	if err == nil {
		t.Fatal("Run() error = nil, want content mode error")
	}
	if !strings.Contains(err.Error(), `invalid content mode "bogus"`) {
		t.Fatalf("Run() error = %v, want invalid content mode message", err)
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
		"github.com/evanmschultz/laslig/examples/spinner",
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

// TestWriterSupportsAnimation verifies both the plain-writer and term.File
// branches of the animation probe remain stable.
func TestWriterSupportsAnimation(t *testing.T) {
	if writerSupportsAnimation(&bytes.Buffer{}) {
		t.Fatal("writerSupportsAnimation(bytes.Buffer) = true, want false")
	}

	writer := fakeFileWriter{fd: ^uintptr(0)}
	if writerSupportsAnimation(&writer) {
		t.Fatal("writerSupportsAnimation(fakeFileWriter) = true, want false for non-tty fd")
	}
}

// TestMaybeSleepAnimatedPreview verifies the preview pause only sleeps when
// animation is enabled.
func TestMaybeSleepAnimatedPreview(t *testing.T) {
	calls := 0
	var slept time.Duration
	sleep := func(delay time.Duration) {
		calls++
		slept = delay
	}

	maybeSleepAnimatedPreview(false, 10*time.Millisecond, sleep)
	if calls != 0 {
		t.Fatalf("maybeSleepAnimatedPreview(false) sleep calls = %d, want 0", calls)
	}

	maybeSleepAnimatedPreview(true, 10*time.Millisecond, sleep)
	if calls != 1 {
		t.Fatalf("maybeSleepAnimatedPreview(true) sleep calls = %d, want 1", calls)
	}
	if slept != 10*time.Millisecond {
		t.Fatalf("maybeSleepAnimatedPreview(true) slept = %v, want %v", slept, 10*time.Millisecond)
	}
}

// TestDelayedPreviewStreamReader verifies the animation-only stream helper
// still replays the full input in order.
func TestDelayedPreviewStreamReader(t *testing.T) {
	const raw = "{\"Action\":\"pass\"}\n{\"Action\":\"fail\"}\n"

	data, err := io.ReadAll(delayedPreviewStreamReader(raw, time.Millisecond))
	if err != nil {
		t.Fatalf("ReadAll(delayedPreviewStreamReader()) error = %v", err)
	}
	if got := string(data); got != raw {
		t.Fatalf("delayedPreviewStreamReader() = %q, want %q", got, raw)
	}
}

// TestPreviewStreamReaderNoAnimation verifies the public preview helper falls
// back to one immediate reader when animation is unavailable.
func TestPreviewStreamReaderNoAnimation(t *testing.T) {
	const raw = "{\"Action\":\"pass\"}\n"

	data, err := io.ReadAll(previewStreamReader(&bytes.Buffer{}, raw, time.Second))
	if err != nil {
		t.Fatalf("ReadAll(previewStreamReader()) error = %v", err)
	}
	if got := string(data); got != raw {
		t.Fatalf("previewStreamReader() = %q, want %q", got, raw)
	}
}

// TestRenderSpinnerWriteError verifies the spinner demo wraps underlying write
// failures instead of swallowing them.
func TestRenderSpinnerWriteError(t *testing.T) {
	printer := laslig.NewWithMode(failingWriter{}, laslig.Mode{Format: laslig.FormatPlain})

	err := RenderSpinner(failingWriter{}, printer)
	if err == nil {
		t.Fatal("RenderSpinner() error = nil, want write failure")
	}
	if !strings.Contains(err.Error(), "render spinner section") {
		t.Fatalf("RenderSpinner() error = %v, want wrapped spinner-section prefix", err)
	}
}

// TestRenderMageCheckPreviewWriteError verifies the Mage-focused preview wraps
// early writer failures through the public renderer entrypoint.
func TestRenderMageCheckPreviewWriteError(t *testing.T) {
	printer := laslig.NewWithMode(failingWriter{}, laslig.Mode{Format: laslig.FormatPlain})

	err := RenderMageCheckPreview(io.Discard, printer)
	if err == nil {
		t.Fatal("RenderMageCheckPreview() error = nil, want write failure")
	}
	if !strings.Contains(err.Error(), "render mage preview section") {
		t.Fatalf("RenderMageCheckPreview() error = %v, want wrapped mage-preview prefix", err)
	}
}

// TestRenderMageCheckPreviewStreamError verifies the focused Mage preview
// reports stream-writer failures during the live gotestout render.
func TestRenderMageCheckPreviewStreamError(t *testing.T) {
	var buf bytes.Buffer
	printer := laslig.NewWithMode(&buf, laslig.Mode{Format: laslig.FormatPlain})

	err := renderMageCheckPreview(failingWriter{}, printer)
	if err == nil {
		t.Fatal("renderMageCheckPreview() error = nil, want stream failure")
	}
	if !strings.Contains(err.Error(), "render mage gotestout stream") {
		t.Fatalf("renderMageCheckPreview() error = %v, want wrapped gotestout prefix", err)
	}
}
