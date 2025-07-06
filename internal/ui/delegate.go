package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var (
	listNormalTitleColor = lipgloss.Color("#dddddd")
	listNormalDescColor  = lipgloss.Color("#777777")
	listSelectedColor    = lipgloss.Color("#5DC9E2")
	listMatchedColor     = lipgloss.Color("#CE3262")
	listDimmedTitleColor = lipgloss.Color("#777777")
	listDimmedDescColor  = lipgloss.Color("#4D4D4D")
)

var (
	listNormalTitleStyle = lipgloss.NewStyle().
				Foreground(listNormalTitleColor).
				Padding(0, 0, 0, 2)

	listNormalDescStyle = listNormalTitleStyle.
				Foreground(listNormalDescColor)

	listSelectedTitleStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(listSelectedColor).
				Foreground(listSelectedColor).
				Padding(0, 0, 0, 1)

	listSelectedDescStyle = listSelectedTitleStyle.
				Foreground(listSelectedColor)

	listDimmedTitleStyle = lipgloss.NewStyle().
				Foreground(listDimmedTitleColor).
				Padding(0, 0, 0, 2)

	listDimmedDescStyle = listDimmedTitleStyle.
				Foreground(listDimmedDescColor)
)

const (
	ellipsis = "…"
)

type itemDelegate struct{}

func (d itemDelegate) Height() int {
	return 2
}

func (d itemDelegate) Spacing() int {
	return 1
}

func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		title, desc  string
		matchedRunes []int
	)

	i := item.(list.DefaultItem)
	title = i.Title()
	desc = i.Description()

	if m.Width() <= 0 {
		return
	}

	textwidth := m.Width() - listNormalTitleStyle.GetPaddingLeft() - listNormalTitleStyle.GetPaddingRight()
	title = ansi.Truncate(title, textwidth, ellipsis)

	var lines []string
	for i, line := range strings.Split(desc, "\n") {
		if i >= d.Height()-1 {
			break
		}
		lines = append(lines, ansi.Truncate(line, textwidth, ellipsis))
	}
	desc = strings.Join(lines, "\n")

	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	if isFiltered && index < len(m.VisibleItems()) {
		matchedRunes = m.MatchesForItem(index)
	}

	if emptyFilter {
		title = listDimmedTitleStyle.Render(title)
		desc = listDimmedDescStyle.Render(desc)
	} else {
		if isSelected && m.FilterState() != list.Filtering {
			if isFiltered {
				unmatched := listSelectedTitleStyle.Inline(true)
				matched := unmatched.Foreground(listMatchedColor)
				title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
			}
			title = listSelectedTitleStyle.Render(title)
			desc = listSelectedDescStyle.Render(desc)
		} else {
			if m.FilterState() == list.Filtering {
				if isFiltered {
					unmatched := listDimmedTitleStyle.Inline(true)
					matched := unmatched.Foreground(listMatchedColor)
					title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
				}
				title = listDimmedTitleStyle.Render(title)
				desc = listDimmedDescStyle.Render(desc)
			} else {
				if isFiltered {
					unmatched := listNormalTitleStyle.Inline(true)
					matched := unmatched.Foreground(listMatchedColor)
					title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
				}
				title = listNormalTitleStyle.Render(title)
				desc = listNormalDescStyle.Render(desc)
			}
		}
	}

	fmt.Fprintf(w, "%s\n%s", title, desc)
}
