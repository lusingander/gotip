package ui

import "github.com/charmbracelet/lipgloss"

type ColorTheme struct {
	Selected lipgloss.Color
	Cursor   lipgloss.Color
	Border   lipgloss.Color

	HelpHeader lipgloss.Color
	HelpKey    lipgloss.Color

	ListNormalTitle lipgloss.Color
	ListNormalDesc  lipgloss.Color
	ListSelected    lipgloss.Color
	ListMatched     lipgloss.Color
	ListDimmedTitle lipgloss.Color
	ListDimmedDesc  lipgloss.Color
}

func DefaultColorTheme() ColorTheme {
	return ColorTheme{
		Selected: lipgloss.Color("#00ADD8"),
		Cursor:   lipgloss.Color("#00ADD8"),
		Border:   lipgloss.Color("240"),

		HelpHeader: lipgloss.Color("#00ADD8"),
		HelpKey:    lipgloss.Color("#5DC9E2"),

		ListNormalTitle: lipgloss.Color("#dddddd"),
		ListNormalDesc:  lipgloss.Color("#777777"),
		ListSelected:    lipgloss.Color("#5DC9E2"),
		ListMatched:     lipgloss.Color("#CE3262"),
		ListDimmedTitle: lipgloss.Color("#777777"),
		ListDimmedDesc:  lipgloss.Color("#4D4D4D"),
	}
}
