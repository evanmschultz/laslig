package laslig

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
)

var defaultSpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

const defaultSpinnerInterval = 100 * time.Millisecond

// Spinner renders one opt-in transient progress line for long-running work.
//
// Use Spinner when a CLI might otherwise be silent for several seconds and a
// caller wants a lightweight "still running" signal. Prefer StatusLine or
// Notice when work starts and finishes quickly enough that durable structured
// output is enough.
//
// Spinner does not own process lifecycle, cancellation, or retries. Callers
// should stop a spinner before writing other laslig blocks to the same writer.
type Spinner struct {
	printer  *Printer
	interval time.Duration
	frames   []string

	mu             sync.Mutex
	text           string
	active         bool
	animated       bool
	forceAnimation bool
	frame          int
	lastWidth      int
	prefix         string
	stopCh         chan struct{}
	doneCh         chan struct{}
	err            error
}

// NewSpinner constructs one opt-in transient progress helper bound to the
// printer's resolved mode, layout, and theme.
//
// Styled human terminals animate one transient line in place. Plain output,
// human output with StyleNever, and JSON output degrade to stable start/finish
// status records without transient frames.
func (p *Printer) NewSpinner() *Spinner {
	return &Spinner{
		printer:  p,
		interval: defaultSpinnerInterval,
		frames:   append([]string(nil), defaultSpinnerFrames...),
	}
}

// Start begins rendering one transient progress line or one stable fallback
// start record when animation is not appropriate for the resolved output mode.
//
// Styled human output animates one in-place line. Plain output, unstyled human
// output, and JSON output emit one durable start record instead.
func (s *Spinner) Start(text string) error {
	if s == nil || s.printer == nil {
		return fmt.Errorf("start spinner: nil spinner")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.active {
		return fmt.Errorf("start spinner: spinner already active")
	}

	s.resetLocked(text)

	if s.printer.mode.Format == FormatJSON {
		if err := s.printer.writeJSON("status_line", StatusLine{
			Level: NoticeInfoLevel,
			Label: "running",
			Text:  text,
		}); err != nil {
			return err
		}
		s.active = true
		return nil
	}

	if err := s.printer.beginBlock(blockKindContent); err != nil {
		return fmt.Errorf("prepare spinner: %w", err)
	}

	s.prefix = strings.Repeat(" ", s.printer.currentContentIndent())
	if !s.shouldAnimateLocked() {
		if err := s.printer.writeContentString(renderPlainStatusLineString(StatusLine{
			Level: NoticeInfoLevel,
			Label: "running",
			Text:  text,
		})); err != nil {
			return err
		}
		s.active = true
		return nil
	}

	s.active = true
	s.animated = true
	s.stopCh = make(chan struct{})
	s.doneCh = make(chan struct{})
	if err := s.writeAnimatedLineLocked(s.renderAnimatedLineLocked(), false); err != nil {
		s.active = false
		s.animated = false
		s.stopCh = nil
		s.doneCh = nil
		return fmt.Errorf("write spinner start: %w", err)
	}

	stopCh := s.stopCh
	doneCh := s.doneCh
	interval := s.interval
	go s.run(stopCh, doneCh, interval)
	return nil
}

// Update replaces the current spinner text while work is still running.
//
// Updates redraw the transient line on styled human terminals. In plain,
// unstyled human, and JSON modes, Update records the latest text but does not
// emit additional output. Call Start before Update.
func (s *Spinner) Update(text string) error {
	if s == nil || s.printer == nil {
		return fmt.Errorf("update spinner: nil spinner")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.err != nil {
		err := s.err
		s.err = nil
		return err
	}
	if !s.active {
		return fmt.Errorf("update spinner: spinner not active")
	}

	s.text = strings.TrimSpace(text)
	if !s.animated {
		return nil
	}
	if err := s.writeAnimatedLineLocked(s.renderAnimatedLineLocked(), false); err != nil {
		s.err = fmt.Errorf("write spinner update: %w", err)
		return s.err
	}
	return nil
}

// Stop finalizes one running spinner with one durable status line or JSON
// status record.
//
// When message is empty, Stop reuses the most recent spinner text. If level is
// empty, Stop defaults to NoticeSuccessLevel. Stop is safe to call even when a
// delayed-start caller never started the spinner.
func (s *Spinner) Stop(message string, level NoticeLevel) error {
	if s == nil || s.printer == nil {
		return nil
	}

	s.mu.Lock()
	if level == "" {
		level = NoticeSuccessLevel
	}
	if strings.TrimSpace(message) == "" {
		message = s.text
	}
	if s.err != nil {
		err := s.err
		s.err = nil
		s.active = false
		s.animated = false
		s.mu.Unlock()
		return err
	}
	if !s.active {
		s.mu.Unlock()
		return nil
	}

	if s.printer.mode.Format == FormatJSON {
		s.active = false
		s.mu.Unlock()
		return s.printer.writeJSON("status_line", StatusLine{
			Level: level,
			Text:  message,
		})
	}

	if !s.animated {
		s.active = false
		s.mu.Unlock()
		return s.printer.writeContentString(renderPlainStatusLineString(StatusLine{
			Level: level,
			Text:  message,
		}))
	}

	close(s.stopCh)
	s.active = false
	s.animated = false
	doneCh := s.doneCh
	err := s.writeAnimatedLineLocked(s.prefix+s.printer.renderStatusLineString(StatusLine{
		Level: level,
		Text:  message,
	}), true)
	s.lastWidth = 0
	s.mu.Unlock()

	if doneCh != nil {
		<-doneCh
	}
	if err != nil {
		return fmt.Errorf("write spinner stop: %w", err)
	}
	return nil
}

func (s *Spinner) run(stopCh <-chan struct{}, doneCh chan<- struct{}, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(doneCh)

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			s.mu.Lock()
			if !s.active || !s.animated {
				s.mu.Unlock()
				return
			}
			s.frame = (s.frame + 1) % len(s.frames)
			if err := s.writeAnimatedLineLocked(s.renderAnimatedLineLocked(), false); err != nil {
				s.err = fmt.Errorf("write spinner frame: %w", err)
				s.active = false
				s.animated = false
				s.mu.Unlock()
				return
			}
			s.mu.Unlock()
		}
	}
}

func (s *Spinner) resetLocked(text string) {
	s.text = text
	s.active = false
	s.animated = false
	s.frame = 0
	s.lastWidth = 0
	s.prefix = ""
	s.stopCh = nil
	s.doneCh = nil
	s.err = nil
}

func (s *Spinner) shouldAnimateLocked() bool {
	if s.printer.mode.Format != FormatHuman || !s.printer.mode.Styled {
		return false
	}
	if s.forceAnimation {
		return true
	}
	return writerIsTerminal(s.printer.out)
}

func (s *Spinner) renderAnimatedLineLocked() string {
	frame := s.frames[s.frame%len(s.frames)]
	text := strings.TrimSpace(s.text)

	if width := s.printer.availableWidth(); width > 0 {
		textWidth := width - lipgloss.Width(frame) - 1
		if textWidth < 0 {
			textWidth = 0
		}
		text = truncateVisible(text, textWidth)
	}

	parts := []string{s.printer.theme.Identifier.Render(frame)}
	if text != "" {
		parts = append(parts, " ", s.printer.theme.Value.Render(text))
	}
	return s.prefix + lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

func (s *Spinner) writeAnimatedLineLocked(line string, newline bool) error {
	padding := ""
	width := lipgloss.Width(line)
	if s.lastWidth > width {
		padding = strings.Repeat(" ", s.lastWidth-width)
	}

	suffix := ""
	if newline {
		s.lastWidth = 0
		suffix = "\n"
	} else {
		s.lastWidth = width
	}

	if _, err := io.WriteString(s.printer.out, "\r"+line+padding+suffix); err != nil {
		return err
	}
	return nil
}

func truncateVisible(value string, width int) string {
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
