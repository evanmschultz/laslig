package laslig

import "charm.land/lipgloss/v2"

// Theme contains the styles used by one printer.
type Theme struct {
	Section       lipgloss.Style
	Label         lipgloss.Style
	Value         lipgloss.Style
	Identifier    lipgloss.Style
	Muted         lipgloss.Style
	Badge         lipgloss.Style
	Panel         lipgloss.Style
	TableHeader   lipgloss.Style
	TableRule     lipgloss.Style
	NoticeInfo    lipgloss.Style
	NoticeSuccess lipgloss.Style
	NoticeWarning lipgloss.Style
	NoticeError   lipgloss.Style
}

// DefaultTheme returns the default theme for one resolved output mode.
func DefaultTheme(mode Mode) Theme {
	base := lipgloss.NewStyle()
	if !mode.Styled {
		return Theme{
			Section:       base,
			Label:         base,
			Value:         base,
			Identifier:    base,
			Muted:         base,
			Badge:         base,
			Panel:         base,
			TableHeader:   base,
			TableRule:     base,
			NoticeInfo:    base,
			NoticeSuccess: base,
			NoticeWarning: base,
			NoticeError:   base,
		}
	}

	primary := lipgloss.Color("#7D56F4")
	primarySoft := lipgloss.Color("#B7A2FF")
	primaryDeep := lipgloss.Color("63")
	text := lipgloss.Color("#FAFAFA")
	muted := lipgloss.Color("#A1A1AA")
	success := lipgloss.Color("#04B575")
	warning := lipgloss.Color("#FFCC00")
	errorColor := lipgloss.Color("#FF5F87")

	return Theme{
		Section:       lipgloss.NewStyle().Bold(true).Foreground(primary),
		Label:         lipgloss.NewStyle().Bold(true).Foreground(primarySoft),
		Value:         lipgloss.NewStyle().Foreground(text),
		Identifier:    lipgloss.NewStyle().Bold(true).Foreground(primarySoft),
		Muted:         lipgloss.NewStyle().Foreground(muted),
		Badge:         lipgloss.NewStyle().Bold(true).Foreground(text).Background(primary).Padding(0, 1),
		Panel:         lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(primaryDeep).Padding(1, 2),
		TableHeader:   lipgloss.NewStyle().Bold(true).Foreground(text),
		TableRule:     lipgloss.NewStyle().Foreground(primaryDeep),
		NoticeInfo:    lipgloss.NewStyle().Bold(true).Foreground(text).Background(primary).Padding(0, 1),
		NoticeSuccess: lipgloss.NewStyle().Bold(true).Foreground(text).Background(success).Padding(0, 1),
		NoticeWarning: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#221400")).Background(warning).Padding(0, 1),
		NoticeError:   lipgloss.NewStyle().Bold(true).Foreground(text).Background(errorColor).Padding(0, 1),
	}
}
