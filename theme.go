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

	heading := lipgloss.Color("69")
	label := lipgloss.Color("109")
	identifier := lipgloss.Color("81")
	accent := lipgloss.Color("62")
	border := lipgloss.Color("63")
	text := lipgloss.Color("252")
	muted := lipgloss.Color("245")
	success := lipgloss.Color("28")
	warning := lipgloss.Color("214")
	errorColor := lipgloss.Color("160")
	badgeText := lipgloss.Color("230")
	warningText := lipgloss.Color("232")

	return Theme{
		Section:       lipgloss.NewStyle().Bold(true).Foreground(heading),
		Label:         lipgloss.NewStyle().Bold(true).Foreground(label),
		Value:         lipgloss.NewStyle().Foreground(text),
		Identifier:    lipgloss.NewStyle().Bold(true).Foreground(identifier),
		Muted:         lipgloss.NewStyle().Foreground(muted),
		Badge:         lipgloss.NewStyle().Bold(true).Foreground(badgeText).Background(accent).Padding(0, 1),
		Panel:         lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(border).Foreground(text).Padding(1, 2),
		TableHeader:   lipgloss.NewStyle().Bold(true).Foreground(heading),
		TableRule:     lipgloss.NewStyle().Foreground(border),
		NoticeInfo:    lipgloss.NewStyle().Bold(true).Foreground(badgeText).Background(accent).Padding(0, 1),
		NoticeSuccess: lipgloss.NewStyle().Bold(true).Foreground(badgeText).Background(success).Padding(0, 1),
		NoticeWarning: lipgloss.NewStyle().Bold(true).Foreground(warningText).Background(warning).Padding(0, 1),
		NoticeError:   lipgloss.NewStyle().Bold(true).Foreground(badgeText).Background(errorColor).Padding(0, 1),
	}
}
