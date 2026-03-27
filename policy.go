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
}

// Mode describes the resolved output behavior for one writer.
type Mode struct {
	Format Format
	Styled bool
}

// ResolveMode resolves one writer and policy into a concrete output mode.
func ResolveMode(out io.Writer, policy Policy) Mode {
	isTTY := false
	if file, ok := out.(term.File); ok {
		isTTY = term.IsTerminal(file.Fd())
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
	}
}
