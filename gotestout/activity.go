package gotestout

import (
	"fmt"
	"io"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"

	"github.com/evanmschultz/laslig"
)

const defaultActivityDelay = 750 * time.Millisecond
const defaultActivityInterval = 100 * time.Millisecond
const activityClearLine = "\r\x1b[2K"

type activityState struct {
	text         string
	frames       []string
	frame        int
	delay        time.Duration
	startedAt    time.Time
	currentPkg   string
	currentTest  string
	testsPassed  int
	testsFailed  int
	testsSkipped int
	pkgsPassed   int
	pkgsFailed   int
	pkgsSkipped  int
	height       int
	shown        bool
	stopCh       chan struct{}
	doneCh       chan struct{}
	err          error
}

func (r *Renderer) startActivity() error {
	if !r.shouldShowActivity() {
		return nil
	}

	delay := r.options.Activity.Delay
	if delay == 0 {
		delay = defaultActivityDelay
	}

	r.activity = &activityState{
		text:      strings.TrimSpace(r.options.Activity.Text),
		frames:    activityFrames(r.options.Activity.SpinnerStyle),
		delay:     delay,
		startedAt: time.Now(),
		stopCh:    make(chan struct{}),
		doneCh:    make(chan struct{}),
	}

	go r.runActivity(r.activity.stopCh, r.activity.doneCh)
	return nil
}

func (r *Renderer) stopActivity() error {
	r.writeMu.Lock()
	activity := r.activity
	if activity == nil {
		r.writeMu.Unlock()
		return nil
	}
	r.activity = nil
	stopCh := activity.stopCh
	doneCh := activity.doneCh
	shown := activity.shown
	height := activity.height
	err := activity.err
	r.writeMu.Unlock()

	if stopCh != nil {
		close(stopCh)
	}
	if doneCh != nil {
		<-doneCh
	}

	if shown {
		r.writeMu.Lock()
		if _, clearErr := io.WriteString(r.out, clearActivityBlock(height)); err == nil && clearErr != nil {
			err = fmt.Errorf("clear activity footer: %w", clearErr)
		}
		r.writeMu.Unlock()
	}
	return err
}

func (r *Renderer) activityError() error {
	r.writeMu.Lock()
	defer r.writeMu.Unlock()
	if r.activity == nil {
		return nil
	}
	return r.activity.err
}

func (r *Renderer) withActivityHidden(fn func() error) error {
	r.writeMu.Lock()
	defer r.writeMu.Unlock()

	if r.activity != nil && r.activity.err != nil {
		return r.activity.err
	}
	if err := r.clearActivityLocked(); err != nil {
		return err
	}
	if err := fn(); err != nil {
		return err
	}
	return r.redrawActivityLocked()
}

func (r *Renderer) updateActivity(event Event) {
	if r.activity == nil {
		return
	}

	r.writeMu.Lock()
	defer r.writeMu.Unlock()
	if r.activity == nil {
		return
	}

	switch event.Action {
	case ActionStart:
		if event.Package != "" {
			r.activity.currentPkg = event.Package
			r.activity.currentTest = ""
		}
	case ActionRun:
		r.activity.currentPkg = event.Package
		r.activity.currentTest = event.Test
	case ActionOutput, ActionBuildOutput:
		if event.Package != "" {
			r.activity.currentPkg = event.Package
			if event.Test != "" {
				r.activity.currentTest = event.Test
			}
		}
	}

	r.activity.testsPassed = r.summary.TestsPassed
	r.activity.testsFailed = r.summary.TestsFailed
	r.activity.testsSkipped = r.summary.TestsSkipped
	r.activity.pkgsPassed = r.summary.PackagesPassed
	r.activity.pkgsFailed = r.summary.PackagesFailed
	r.activity.pkgsSkipped = r.summary.PackagesSkipped

	if event.PackageEvent() && event.Action.IsTerminal() {
		r.activity.currentTest = ""
	}
}

func (r *Renderer) shouldShowActivity() bool {
	if r.mode.Format != laslig.FormatHuman || !r.mode.Styled {
		return false
	}

	switch r.options.Activity.Mode {
	case ActivityOff:
		return false
	case ActivityOn:
		return true
	case ActivityAuto:
		return writerIsTerminal(r.out)
	default:
		return false
	}
}

func (r *Renderer) runActivity(stopCh <-chan struct{}, doneCh chan<- struct{}) {
	defer close(doneCh)
	activity := r.activity
	if activity == nil {
		return
	}

	if activity.delay > 0 {
		timer := time.NewTimer(activity.delay)
		defer timer.Stop()
		select {
		case <-stopCh:
			return
		case <-timer.C:
		}
	}

	if err := r.tickActivity(false); err != nil {
		return
	}

	ticker := time.NewTicker(defaultActivityInterval)
	defer ticker.Stop()
	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			if err := r.tickActivity(true); err != nil {
				return
			}
		}
	}
}

func (r *Renderer) tickActivity(advance bool) error {
	r.writeMu.Lock()
	defer r.writeMu.Unlock()

	if r.activity == nil {
		return nil
	}
	if r.activity.err != nil {
		return r.activity.err
	}
	if advance {
		r.activity.frame = (r.activity.frame + 1) % len(r.activity.frames)
	}
	block, height := r.renderActivityBlockLocked()
	output := block
	if r.activity.shown {
		output = clearActivityBlock(r.activity.height) + output
	}
	if _, err := io.WriteString(r.out, output); err != nil {
		r.activity.err = fmt.Errorf("write activity footer: %w", err)
		return r.activity.err
	}
	r.activity.height = height
	r.activity.shown = true
	return nil
}

func (r *Renderer) clearActivityLocked() error {
	if r.activity == nil || !r.activity.shown {
		return nil
	}
	if _, err := io.WriteString(r.out, clearActivityBlock(r.activity.height)); err != nil {
		r.activity.err = fmt.Errorf("clear activity footer: %w", err)
		return r.activity.err
	}
	r.activity.height = 0
	r.activity.shown = false
	return nil
}

func (r *Renderer) redrawActivityLocked() error {
	if r.activity == nil || r.activity.err != nil {
		if r.activity != nil {
			return r.activity.err
		}
		return nil
	}
	block, height := r.renderActivityBlockLocked()
	if _, err := io.WriteString(r.out, block); err != nil {
		r.activity.err = fmt.Errorf("redraw activity footer: %w", err)
		return r.activity.err
	}
	r.activity.height = height
	r.activity.shown = true
	return nil
}

func (r *Renderer) renderActivityBlockLocked() (string, int) {
	frame := r.activity.frames[r.activity.frame%len(r.activity.frames)]
	text := strings.TrimSpace(r.activity.text)
	if text == "" {
		text = "Running go test -json"
	}

	lines := []string{}
	lines = append(lines, r.renderActivityHeaderLines(frame, text)...)

	subject := strings.TrimSpace(r.activity.currentPkg)
	if subject != "" && strings.TrimSpace(r.activity.currentTest) != "" {
		subject += " :: " + strings.TrimSpace(r.activity.currentTest)
	}
	if subject != "" {
		lines = append(lines, r.renderActivityValueLines("- ", subject, r.theme.Identifier)...)
	}

	if pkg := strings.TrimSpace(r.activity.currentPkg); pkg != "" {
		lines = append(lines, r.renderActivityFieldLines("package", pkg, r.theme.Identifier)...)
	}
	if test := strings.TrimSpace(r.activity.currentTest); test != "" {
		lines = append(lines, r.renderActivityFieldLines("test", test, r.theme.Identifier)...)
	}
	lines = append(lines, r.renderActivityCountsLines(
		"tests",
		r.activity.testsPassed,
		r.activity.testsFailed,
		r.activity.testsSkipped,
	)...)
	lines = append(lines, r.renderActivityCountsLines(
		"packages",
		r.activity.pkgsPassed,
		r.activity.pkgsFailed,
		r.activity.pkgsSkipped,
	)...)
	lines = append(lines, r.renderActivityFieldLines("elapsed", formatActivityElapsed(time.Since(r.activity.startedAt)), r.theme.Muted)...)
	return strings.Join(lines, "\n"), len(lines)
}

func clearActivityBlock(height int) string {
	if height <= 0 {
		return activityClearLine
	}

	var builder strings.Builder
	builder.WriteString(activityClearLine)
	for i := 1; i < height; i++ {
		builder.WriteString("\x1b[1A")
		builder.WriteString(activityClearLine)
	}
	return builder.String()
}

func (r *Renderer) renderActivityHeaderLines(frame string, text string) []string {
	lines := wrapActivityText(text, r.activityContentWidth(lipgloss.Width(frame)+1))
	if len(lines) == 0 {
		lines = []string{""}
	}

	rendered := make([]string, 0, len(lines))
	for index, line := range lines {
		if index == 0 {
			rendered = append(rendered, lipgloss.JoinHorizontal(
				lipgloss.Top,
				r.theme.Identifier.Render(frame),
				" ",
				r.theme.Value.Render(line),
			))
			continue
		}
		rendered = append(rendered, strings.Repeat(" ", lipgloss.Width(frame)+1)+r.theme.Value.Render(line))
	}
	return rendered
}

func (r *Renderer) renderActivityValueLines(prefix string, value string, style lipgloss.Style) []string {
	return renderWrappedStyledValueLines(prefix, "", value, style, style, r.activityContentWidth(lipgloss.Width(prefix)))
}

func (r *Renderer) renderActivityFieldLines(label string, value string, style lipgloss.Style) []string {
	prefixPlain := "  " + label + ": "
	prefixStyled := "  " + r.theme.Label.Render(label+":") + " "
	return renderWrappedStyledValueLines(prefixPlain, prefixStyled, value, style, style, r.activityContentWidth(lipgloss.Width(prefixPlain)))
}

func (r *Renderer) renderActivityCountsLines(label string, passed int, failed int, skipped int) []string {
	prefixPlain := "  " + label + ": "
	prefixStyled := "  " + r.theme.Label.Render(label+":") + " "
	plainValue := fmt.Sprintf("%d pass, %d fail, %d skip", passed, failed, skipped)
	valueWidth := r.activityContentWidth(lipgloss.Width(prefixPlain))
	lines := wrapActivityText(plainValue, valueWidth)
	if len(lines) == 0 {
		lines = []string{""}
	}

	rendered := make([]string, 0, len(lines))
	for index, line := range lines {
		prefix := strings.Repeat(" ", lipgloss.Width(prefixPlain))
		if index == 0 {
			prefix = prefixStyled
		}
		rendered = append(rendered, prefix+r.renderActivityCountsValue(line))
	}
	return rendered
}

func (r *Renderer) renderActivityCountsValue(value string) string {
	if r.mode.Format != laslig.FormatHuman {
		return value
	}

	passStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	failStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
	skipStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))

	replaced := value
	replaced = strings.ReplaceAll(replaced, " pass", " "+passStyle.Render("pass"))
	replaced = strings.ReplaceAll(replaced, " fail", " "+failStyle.Render("fail"))
	replaced = strings.ReplaceAll(replaced, " skip", " "+skipStyle.Render("skip"))
	return r.theme.Value.Render(replaced)
}

func (r *Renderer) activityContentWidth(prefixWidth int) int {
	if r.mode.Width <= 0 {
		return 0
	}
	width := r.mode.Width - prefixWidth
	if width < 1 {
		return 1
	}
	return width
}

func activityFrames(style laslig.SpinnerStyle) []string {
	switch style {
	case laslig.SpinnerStyleDot:
		return []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	case laslig.SpinnerStyleLine:
		return []string{"-", "\\", "|", "/"}
	case laslig.SpinnerStylePulse:
		return []string{"∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙●∙"}
	case laslig.SpinnerStyleMeter:
		return []string{"[    ]", "[=   ]", "[==  ]", "[=== ]", "[ ===]", "[  ==]", "[   =]"}
	default:
		return []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	}
}

func formatActivityElapsed(elapsed time.Duration) string {
	if elapsed < time.Second {
		return elapsed.Round(100 * time.Millisecond).String()
	}
	return elapsed.Round(time.Second).String()
}

func wrapActivityText(value string, width int) []string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	if width <= 0 || lipgloss.Width(trimmed) <= width {
		return []string{trimmed}
	}

	lines := []string{}
	remaining := trimmed
	for strings.TrimSpace(remaining) != "" {
		lines = append(lines, sliceVisibleWidth(remaining, width))
		if len(lines[len(lines)-1]) >= len(remaining) {
			break
		}
		remaining = strings.TrimLeft(remaining[len(lines[len(lines)-1]):], " ")
	}
	return lines
}

func sliceVisibleWidth(value string, width int) string {
	if width <= 0 {
		return value
	}

	var builder strings.Builder
	for _, r := range value {
		candidate := builder.String() + string(r)
		if lipgloss.Width(candidate) > width {
			break
		}
		builder.WriteRune(r)
	}
	if builder.Len() == 0 {
		return value
	}
	return builder.String()
}

func renderWrappedStyledValueLines(prefixPlain string, prefixStyled string, value string, firstStyle lipgloss.Style, continuationStyle lipgloss.Style, width int) []string {
	if prefixStyled == "" {
		prefixStyled = prefixPlain
	}
	lines := wrapActivityText(value, width)
	if len(lines) == 0 {
		lines = []string{""}
	}

	rendered := make([]string, 0, len(lines))
	for index, line := range lines {
		prefix := strings.Repeat(" ", lipgloss.Width(prefixPlain))
		style := continuationStyle
		if index == 0 {
			prefix = prefixStyled
			style = firstStyle
		}
		rendered = append(rendered, prefix+style.Render(line))
	}
	return rendered
}

func writerIsTerminal(out io.Writer) bool {
	file, ok := out.(term.File)
	if !ok {
		return false
	}
	return term.IsTerminal(file.Fd())
}
