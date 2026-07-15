package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/gotip/internal/theme"
)

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

func newAppStyles(colorTheme theme.ColorTheme) appStyles {
	return appStyles{
		selectedLabel: lipgloss.NewStyle().Foreground(colorTheme.Accent),
		selectedName:  lipgloss.NewStyle().Foreground(colorTheme.Accent).Bold(true),
		selectedPath:  lipgloss.NewStyle().Foreground(colorTheme.Accent).Bold(true),

		header: lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(colorTheme.Border),

		footerMsg:           lipgloss.NewStyle().Foreground(colorTheme.Text),
		footerFiltered:      lipgloss.NewStyle().Foreground(colorTheme.Text),
		footerSelectedIndex: lipgloss.NewStyle().Foreground(colorTheme.Text),
		footerDivider:       lipgloss.NewStyle().Foreground(colorTheme.Border),

		footer: lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(colorTheme.Border),

		helpHeader:   lipgloss.NewStyle().Foreground(colorTheme.Accent),
		helpContent:  lipgloss.NewStyle().Padding(0, 2),
		helpText:     lipgloss.NewStyle().Foreground(colorTheme.Text),
		helpKey:      lipgloss.NewStyle().Foreground(colorTheme.Highlight).Bold(true),
		filterPrompt: lipgloss.NewStyle().Foreground(colorTheme.Text),
		filterText:   lipgloss.NewStyle().Foreground(colorTheme.Text),
		filterCursor: lipgloss.NewStyle().Foreground(colorTheme.Accent),
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

func newListStyles(colorTheme theme.ColorTheme) listStyles {
	normalTitle := lipgloss.NewStyle().
		Foreground(colorTheme.Text).
		Padding(0, 0, 0, 2)

	normalDesc := normalTitle.
		Foreground(colorTheme.Muted)

	selectedTitle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(colorTheme.Highlight).
		Foreground(colorTheme.Highlight).
		Padding(0, 0, 0, 1)

	selectedDesc := selectedTitle.
		Foreground(colorTheme.Highlight)

	dimmedTitle := lipgloss.NewStyle().
		Foreground(colorTheme.Muted).
		Padding(0, 0, 0, 2)

	dimmedDesc := dimmedTitle.
		Foreground(colorTheme.Dimmed)

	return listStyles{
		normalTitle:   normalTitle,
		normalDesc:    normalDesc,
		selectedTitle: selectedTitle,
		selectedDesc:  selectedDesc,
		dimmedTitle:   dimmedTitle,
		dimmedDesc:    dimmedDesc,
		matchedColor:  colorTheme.Match,
	}
}
