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

	return Theme{
		Section:       lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#9DDCFF")),
		Label:         lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#A9C6D8")),
		Value:         lipgloss.NewStyle().Foreground(lipgloss.Color("#EAF6FF")),
		Identifier:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#A2F2D9")),
		Muted:         lipgloss.NewStyle().Foreground(lipgloss.Color("#7A8D99")),
		Badge:         lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0A1E2A")).Background(lipgloss.Color("#9DDCFF")).Padding(0, 1),
		Panel:         lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#5B7688")).Padding(1, 2),
		TableHeader:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#D4EEFF")),
		TableRule:     lipgloss.NewStyle().Foreground(lipgloss.Color("#5B7688")),
		NoticeInfo:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0A1E2A")).Background(lipgloss.Color("#8EDCFF")).Padding(0, 1),
		NoticeSuccess: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F4FFF8")).Background(lipgloss.Color("#1F8F56")).Padding(0, 1),
		NoticeWarning: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#211300")).Background(lipgloss.Color("#F4C95D")).Padding(0, 1),
		NoticeError:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFF4F4")).Background(lipgloss.Color("#B42318")).Padding(0, 1),
	}
}
