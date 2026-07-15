package ui

import "github.com/charmbracelet/lipgloss"

type ColorTheme struct {
	Text      lipgloss.Color
	Accent    lipgloss.Color
	Highlight lipgloss.Color
	Muted     lipgloss.Color
	Dimmed    lipgloss.Color
	Border    lipgloss.Color
	Match     lipgloss.Color
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
	}
}

type appStyles struct {
	selectedLabel       lipgloss.Style
	selectedName        lipgloss.Style
	selectedPath        lipgloss.Style
	header              lipgloss.Style
	footerMsg           lipgloss.Style
	footerFiltered      lipgloss.Style
	footerSelectedIndex lipgloss.Style
	footerDivider       lipgloss.Style
	footer              lipgloss.Style
	helpHeader          lipgloss.Style
	helpContent         lipgloss.Style
	helpText            lipgloss.Style
	helpKey             lipgloss.Style
	filterPrompt        lipgloss.Style
	filterText          lipgloss.Style
	filterCursor        lipgloss.Style
}

func newAppStyles(theme ColorTheme) appStyles {
	return appStyles{
		selectedLabel: lipgloss.NewStyle().Foreground(theme.Accent),
		selectedName:  lipgloss.NewStyle().Foreground(theme.Accent).Bold(true),
		selectedPath:  lipgloss.NewStyle().Foreground(theme.Accent).Bold(true),

		header: lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(theme.Border),

		footerMsg:           lipgloss.NewStyle().Foreground(theme.Text),
		footerFiltered:      lipgloss.NewStyle().Foreground(theme.Text),
		footerSelectedIndex: lipgloss.NewStyle().Foreground(theme.Text),
		footerDivider:       lipgloss.NewStyle().Foreground(theme.Border),

		footer: lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(theme.Border),

		helpHeader:   lipgloss.NewStyle().Foreground(theme.Accent),
		helpContent:  lipgloss.NewStyle().Padding(0, 2),
		helpText:     lipgloss.NewStyle().Foreground(theme.Text),
		helpKey:      lipgloss.NewStyle().Foreground(theme.Highlight).Bold(true),
		filterPrompt: lipgloss.NewStyle().Foreground(theme.Text),
		filterText:   lipgloss.NewStyle().Foreground(theme.Text),
		filterCursor: lipgloss.NewStyle().Foreground(theme.Accent),
	}
}

type listStyles struct {
	normalTitle   lipgloss.Style
	normalDesc    lipgloss.Style
	selectedTitle lipgloss.Style
	selectedDesc  lipgloss.Style
	dimmedTitle   lipgloss.Style
	dimmedDesc    lipgloss.Style
	matchedColor  lipgloss.Color
}

func newListStyles(theme ColorTheme) listStyles {
	normalTitle := lipgloss.NewStyle().
		Foreground(theme.Text).
		Padding(0, 0, 0, 2)

	normalDesc := normalTitle.
		Foreground(theme.Muted)

	selectedTitle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(theme.Highlight).
		Foreground(theme.Highlight).
		Padding(0, 0, 0, 1)

	selectedDesc := selectedTitle.
		Foreground(theme.Highlight)

	dimmedTitle := lipgloss.NewStyle().
		Foreground(theme.Muted).
		Padding(0, 0, 0, 2)

	dimmedDesc := dimmedTitle.
		Foreground(theme.Dimmed)

	return listStyles{
		normalTitle:   normalTitle,
		normalDesc:    normalDesc,
		selectedTitle: selectedTitle,
		selectedDesc:  selectedDesc,
		dimmedTitle:   dimmedTitle,
		dimmedDesc:    dimmedDesc,
		matchedColor:  theme.Match,
	}
}
