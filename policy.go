package laslig

import (
	"io"

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

// Policy describes the requested output behavior before writer capabilities are resolved.
type Policy struct {
	Format Format
	Style  StylePolicy
	Layout *Layout
	Theme  *Theme
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

// clampNonNegative keeps layout counts at zero or above.
func clampNonNegative(count int) int {
	if count < 0 {
		return 0
	}
	return count
}
