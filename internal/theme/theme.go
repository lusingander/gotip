package theme

import "github.com/charmbracelet/lipgloss"

type ColorTheme struct {
	Text      lipgloss.Color
	Accent    lipgloss.Color
	Highlight lipgloss.Color
	Muted     lipgloss.Color
	Dimmed    lipgloss.Color
	Border    lipgloss.Color
	Match     lipgloss.Color
	Command   lipgloss.Color
}

func DefaultColorTheme() ColorTheme {
	return ColorTheme{
		Text:      lipgloss.Color("#DDDDDD"),
		Accent:    lipgloss.Color("#00ADD8"),
		Highlight: lipgloss.Color("#5DC9E2"),
		Muted:     lipgloss.Color("#777777"),
		Dimmed:    lipgloss.Color("#4D4D4D"),
		Border:    lipgloss.Color("240"),
		Match:     lipgloss.Color("#CE3262"),
		Command:   lipgloss.Color("#00A29C"),
	}
}
