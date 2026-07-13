package ui

import "github.com/charmbracelet/lipgloss"

type ColorTheme struct {
	Normal   lipgloss.Color
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
		Normal:   lipgloss.Color("#dddddd"),
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
		selectedLabel: lipgloss.NewStyle().Foreground(theme.Selected),
		selectedName:  lipgloss.NewStyle().Foreground(theme.Selected).Bold(true),
		selectedPath:  lipgloss.NewStyle().Foreground(theme.Selected).Bold(true),

		header: lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(theme.Border),

		footerMsg:           lipgloss.NewStyle().Foreground(theme.Normal),
		footerFiltered:      lipgloss.NewStyle().Foreground(theme.Normal),
		footerSelectedIndex: lipgloss.NewStyle().Foreground(theme.Normal),
		footerDivider:       lipgloss.NewStyle().Foreground(theme.Border),

		footer: lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(theme.Border),

		helpHeader:   lipgloss.NewStyle().Foreground(theme.HelpHeader),
		helpContent:  lipgloss.NewStyle().Padding(0, 2),
		helpText:     lipgloss.NewStyle().Foreground(theme.Normal),
		helpKey:      lipgloss.NewStyle().Foreground(theme.HelpKey).Bold(true),
		filterPrompt: lipgloss.NewStyle().Foreground(theme.Normal),
		filterText:   lipgloss.NewStyle().Foreground(theme.Normal),
		filterCursor: lipgloss.NewStyle().Foreground(theme.Cursor),
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
		Foreground(theme.ListNormalTitle).
		Padding(0, 0, 0, 2)

	normalDesc := normalTitle.
		Foreground(theme.ListNormalDesc)

	selectedTitle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(theme.ListSelected).
		Foreground(theme.ListSelected).
		Padding(0, 0, 0, 1)

	selectedDesc := selectedTitle.
		Foreground(theme.ListSelected)

	dimmedTitle := lipgloss.NewStyle().
		Foreground(theme.ListDimmedTitle).
		Padding(0, 0, 0, 2)

	dimmedDesc := dimmedTitle.
		Foreground(theme.ListDimmedDesc)

	return listStyles{
		normalTitle:   normalTitle,
		normalDesc:    normalDesc,
		selectedTitle: selectedTitle,
		selectedDesc:  selectedDesc,
		dimmedTitle:   dimmedTitle,
		dimmedDesc:    dimmedDesc,
		matchedColor:  theme.ListMatched,
	}
}
