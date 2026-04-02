package laslig

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/term"
)

// Format identifies the output representation used for a render operation.
type Format string

const (
	// FormatAuto resolves to a human-oriented format on terminals and plain text otherwise.
	FormatAuto Format = "auto"
	// FormatHuman renders human-oriented structured output.
	FormatHuman Format = "human"
	// FormatPlain renders plain text without terminal styling.
	FormatPlain Format = "plain"
	// FormatJSON renders machine-readable JSON payloads.
	FormatJSON Format = "json"
)

// StylePolicy controls whether ANSI styling is enabled for human output.
type StylePolicy string

const (
	// StyleAuto enables styling only when output is attached to a terminal.
	StyleAuto StylePolicy = "auto"
	// StyleAlways forces styling for human output.
	StyleAlways StylePolicy = "always"
	// StyleNever disables styling.
	StyleNever StylePolicy = "never"
)

// Policy describes the requested output behavior before writer capabilities are
// resolved.
type Policy struct {
	// Format selects the overall render format.
	Format Format
	// Style controls ANSI styling for human output.
	Style StylePolicy
	// Layout overrides the default document spacing and indentation rules.
	Layout *Layout
	// Theme overrides the printer-wide lipgloss style roles.
	Theme *Theme
	// SpinnerStyle selects the built-in transient spinner frame set used by
	// Printer.NewSpinner. Supported values are braille, dot, line, pulse, and
	// meter. The default is braille.
	SpinnerStyle SpinnerStyle
	// GlamourStyle selects the built-in Glamour preset used for Markdown and
	// code-block rendering. Supported values are dark, light, pink, dracula,
	// tokyo-night, ascii, and notty. The default is dracula.
	GlamourStyle GlamourStyle
}

// SpinnerStyle identifies one supported built-in spinner frame set.
//
// Supported built-ins are braille, dot, line, pulse, and meter.
type SpinnerStyle string

const (
	// SpinnerStyleBraille renders the default braille-dot spinner.
	SpinnerStyleBraille SpinnerStyle = "braille"
	// SpinnerStyleDot renders a larger dot spinner.
	SpinnerStyleDot SpinnerStyle = "dot"
	// SpinnerStyleLine renders a compact ASCII-friendly line spinner.
	SpinnerStyleLine SpinnerStyle = "line"
	// SpinnerStylePulse renders a pulsing dot spinner.
	SpinnerStylePulse SpinnerStyle = "pulse"
	// SpinnerStyleMeter renders a meter-like spinner.
	SpinnerStyleMeter SpinnerStyle = "meter"
)

// DefaultSpinnerStyle returns the default built-in spinner style used by
// laslig, which is braille.
func DefaultSpinnerStyle() SpinnerStyle {
	return SpinnerStyleBraille
}

// Valid reports whether the style matches one of laslig's supported built-in
// spinner styles.
func (s SpinnerStyle) Valid() bool {
	switch s {
	case SpinnerStyleBraille, SpinnerStyleDot, SpinnerStyleLine, SpinnerStylePulse, SpinnerStyleMeter:
		return true
	default:
		return false
	}
}

// GlamourStyle identifies one supported built-in Glamour style preset.
//
// Supported built-ins are dark, light, pink, dracula, tokyo-night, ascii, and
// notty.
type GlamourStyle string

const (
	// GlamourStyleDark renders Markdown with Glamour's dark preset.
	GlamourStyleDark GlamourStyle = "dark"
	// GlamourStyleLight renders Markdown with Glamour's light preset.
	GlamourStyleLight GlamourStyle = "light"
	// GlamourStylePink renders Markdown with Glamour's pink preset.
	GlamourStylePink GlamourStyle = "pink"
	// GlamourStyleDracula renders Markdown with Glamour's Dracula preset.
	GlamourStyleDracula GlamourStyle = "dracula"
	// GlamourStyleTokyoNight renders Markdown with Glamour's Tokyo Night preset.
	GlamourStyleTokyoNight GlamourStyle = "tokyo-night"
	// GlamourStyleASCII renders Markdown with Glamour's ASCII preset.
	GlamourStyleASCII GlamourStyle = "ascii"
	// GlamourStyleNoTTY renders Markdown with Glamour's no-TTY preset.
	GlamourStyleNoTTY GlamourStyle = "notty"
)

// DefaultGlamourStyle returns the default built-in Glamour style used by
// laslig, which is dracula.
func DefaultGlamourStyle() GlamourStyle {
	return GlamourStyleDracula
}

// Valid reports whether the style matches one of laslig's supported built-in
// Glamour presets.
func (s GlamourStyle) Valid() bool {
	switch s {
	case GlamourStyleDark, GlamourStyleLight, GlamourStylePink, GlamourStyleDracula, GlamourStyleTokyoNight, GlamourStyleASCII, GlamourStyleNoTTY:
		return true
	default:
		return false
	}
}

// ListMarker identifies the marker shape used for unordered and ordered list output.
type ListMarker string

const (
	// ListMarkerDash renders list items with a dash marker.
	ListMarkerDash ListMarker = "dash"
	// ListMarkerBullet renders list items with a bullet marker.
	ListMarkerBullet ListMarker = "bullet"
	// ListMarkerNumber renders list items with an ordinal marker.
	ListMarkerNumber ListMarker = "number"
)

// Layout describes the high-level document rhythm used by one printer.
//
// Use DefaultLayout as a base, then override individual values with the
// builder-style helpers when a command wants a different shape.
type Layout struct {
	leadingGap    int
	blockGap      int
	sectionGap    int
	sectionIndent int
	listMarker    ListMarker
}

// DefaultLayout returns the opinionated default document layout used by laslig.
func DefaultLayout() Layout {
	return Layout{
		leadingGap:    1,
		blockGap:      1,
		sectionGap:    2,
		sectionIndent: 2,
		listMarker:    ListMarkerDash,
	}
}

// WithLeadingGap returns one layout with an updated leading gap.
func (l Layout) WithLeadingGap(count int) Layout {
	l.leadingGap = clampNonNegative(count)
	return l
}

// WithBlockGap returns one layout with an updated ordinary block gap.
func (l Layout) WithBlockGap(count int) Layout {
	l.blockGap = clampNonNegative(count)
	return l
}

// WithSectionGap returns one layout with an updated section gap.
func (l Layout) WithSectionGap(count int) Layout {
	l.sectionGap = clampNonNegative(count)
	return l
}

// WithSectionIndent returns one layout with an updated section-body indent.
func (l Layout) WithSectionIndent(count int) Layout {
	l.sectionIndent = clampNonNegative(count)
	return l
}

// WithListMarker returns one layout with an updated list marker style.
func (l Layout) WithListMarker(marker ListMarker) Layout {
	if marker == "" {
		marker = ListMarkerDash
	}
	l.listMarker = marker
	return l
}

// Mode describes the resolved output behavior for one writer, including the
// detected terminal width when available.
type Mode struct {
	Format Format
	Styled bool
	Width  int
}

// ResolveMode resolves one writer and policy into a concrete output mode.
func ResolveMode(out io.Writer, policy Policy) Mode {
	isTTY := false
	width := 0
	if file, ok := out.(term.File); ok {
		isTTY = term.IsTerminal(file.Fd())
		if isTTY {
			if terminalWidth, _, err := term.GetSize(file.Fd()); err == nil {
				width = terminalWidth
			}
		}
	}

	if width <= 0 {
		if rawColumns := strings.TrimSpace(os.Getenv("COLUMNS")); rawColumns != "" {
			if parsed, err := strconv.Atoi(rawColumns); err == nil && parsed > 0 {
				width = parsed
			}
		}
	}

	format := policy.Format
	if format == "" {
		format = FormatAuto
	}
	if format == FormatAuto {
		if isTTY {
			format = FormatHuman
		} else {
			format = FormatPlain
		}
	}

	styled := false
	if format == FormatHuman {
		switch policy.Style {
		case StyleAlways:
			styled = true
		case StyleNever:
			styled = false
		default:
			styled = isTTY
		}
	}

	return Mode{
		Format: format,
		Styled: styled,
		Width:  width,
	}
}

// resolveLayout resolves the requested layout policy into one concrete layout.
func resolveLayout(policy Policy) Layout {
	if policy.Layout == nil {
		return DefaultLayout()
	}
	layout := *policy.Layout
	if layout.listMarker == "" {
		layout.listMarker = ListMarkerDash
	}
	layout.leadingGap = clampNonNegative(layout.leadingGap)
	layout.blockGap = clampNonNegative(layout.blockGap)
	layout.sectionGap = clampNonNegative(layout.sectionGap)
	layout.sectionIndent = clampNonNegative(layout.sectionIndent)
	return layout
}

// resolveTheme resolves the requested theme policy into one concrete theme.
func resolveTheme(policy Policy, mode Mode) Theme {
	if policy.Theme == nil {
		return DefaultTheme(mode)
	}
	return *policy.Theme
}

// resolveGlamourStyle resolves the requested Glamour style to one supported
// built-in preset, falling back to the library default when unset or invalid.
func resolveGlamourStyle(policy Policy) GlamourStyle {
	if !policy.GlamourStyle.Valid() {
		return DefaultGlamourStyle()
	}
	return policy.GlamourStyle
}

// resolveSpinnerStyle resolves the requested spinner style to one supported
// built-in frame set.
func resolveSpinnerStyle(policy Policy) SpinnerStyle {
	if !policy.SpinnerStyle.Valid() {
		return DefaultSpinnerStyle()
	}
	return policy.SpinnerStyle
}

// clampNonNegative keeps layout counts at zero or above.
func clampNonNegative(count int) int {
	if count < 0 {
		return 0
	}
	return count
}
