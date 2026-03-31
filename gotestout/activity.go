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
		if _, clearErr := io.WriteString(r.out, activityClearLine); err == nil && clearErr != nil {
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
	if _, err := io.WriteString(r.out, activityClearLine+r.renderActivityLineLocked()); err != nil {
		r.activity.err = fmt.Errorf("write activity footer: %w", err)
		return r.activity.err
	}
	r.activity.shown = true
	return nil
}

func (r *Renderer) clearActivityLocked() error {
	if r.activity == nil || !r.activity.shown {
		return nil
	}
	if _, err := io.WriteString(r.out, activityClearLine); err != nil {
		r.activity.err = fmt.Errorf("clear activity footer: %w", err)
		return r.activity.err
	}
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
	if _, err := io.WriteString(r.out, activityClearLine+r.renderActivityLineLocked()); err != nil {
		r.activity.err = fmt.Errorf("redraw activity footer: %w", err)
		return r.activity.err
	}
	r.activity.shown = true
	return nil
}

func (r *Renderer) renderActivityLineLocked() string {
	frame := r.activity.frames[r.activity.frame%len(r.activity.frames)]
	subject := r.activity.currentPkg
	if r.activity.currentTest != "" {
		subject += " :: " + r.activity.currentTest
	}

	text := strings.TrimSpace(r.activity.text)
	if text == "" {
		text = "Running go test -json"
	}

	details := []string{
		fmt.Sprintf("tests: %d/%d/%d", r.activity.testsPassed, r.activity.testsFailed, r.activity.testsSkipped),
		fmt.Sprintf("pkgs: %d/%d/%d", r.activity.pkgsPassed, r.activity.pkgsFailed, r.activity.pkgsSkipped),
		formatActivityElapsed(time.Since(r.activity.startedAt)),
	}
	if subject != "" {
		details = append([]string{"current: " + subject}, details...)
	}
	lineText := text + "  " + strings.Join(details, "  ")

	if r.mode.Width > 0 {
		textWidth := r.mode.Width - lipgloss.Width(frame) - 1
		if textWidth < 0 {
			textWidth = 0
		}
		lineText = truncateActivityText(lineText, textWidth)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		r.theme.Identifier.Render(frame),
		" ",
		r.theme.Value.Render(lineText),
	)
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

func truncateActivityText(value string, width int) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(value) <= width {
		return value
	}

	const ellipsis = "…"
	if width == 1 {
		return ellipsis
	}

	var builder strings.Builder
	for _, r := range value {
		candidate := builder.String() + string(r)
		if lipgloss.Width(candidate+ellipsis) > width {
			break
		}
		builder.WriteRune(r)
	}
	return strings.TrimRight(builder.String(), " ") + ellipsis
}

func writerIsTerminal(out io.Writer) bool {
	file, ok := out.(term.File)
	if !ok {
		return false
	}
	return term.IsTerminal(file.Fd())
}
