package ui

import (
	"fmt"
	"io"

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

type testCaseItemDelegate struct{}

func (d testCaseItemDelegate) Height() int {
	return 2
}

func (d testCaseItemDelegate) Spacing() int {
	return 1
}

func (d testCaseItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d testCaseItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i := item.(*testCaseItem)
	title := i.name
	desc := i.path

	if m.Width() <= 0 {
		return
	}

	textwidth := m.Width() - listNormalTitleStyle.GetPaddingLeft() - listNormalTitleStyle.GetPaddingRight()
	title = ansi.Truncate(title, textwidth, ellipsis)
	desc = ansi.Truncate(desc, textwidth, ellipsis)

	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	var matchedRunes []int
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

type historyItemDelegate struct{}

func (d historyItemDelegate) Height() int {
	return 3
}

func (d historyItemDelegate) Spacing() int {
	return 1
}

func (d historyItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d historyItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i := item.(*historyItem)
	title := i.nameForView
	desc := i.path
	runAt := i.runAt

	if m.Width() <= 0 {
		return
	}

	textwidth := m.Width() - listNormalTitleStyle.GetPaddingLeft() - listNormalTitleStyle.GetPaddingRight()
	title = ansi.Truncate(title, textwidth, ellipsis)
	desc = ansi.Truncate(desc, textwidth, ellipsis)
	runAt = ansi.Truncate(runAt, textwidth, ellipsis)

	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	var matchedRunes []int
	if isFiltered && index < len(m.VisibleItems()) {
		matchedRunes = m.MatchesForItem(index)
	}

	if emptyFilter {
		title = listDimmedTitleStyle.Render(title)
		desc = listDimmedDescStyle.Render(desc)
		runAt = listDimmedDescStyle.Render(runAt)
	} else {
		if isSelected && m.FilterState() != list.Filtering {
			if isFiltered {
				unmatched := listSelectedTitleStyle.Inline(true)
				matched := unmatched.Foreground(listMatchedColor)
				title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
			}
			title = listSelectedTitleStyle.Render(title)
			desc = listSelectedDescStyle.Render(desc)
			runAt = listSelectedDescStyle.Render(runAt)
		} else {
			if m.FilterState() == list.Filtering {
				if isFiltered {
					unmatched := listDimmedTitleStyle.Inline(true)
					matched := unmatched.Foreground(listMatchedColor)
					title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
				}
				title = listDimmedTitleStyle.Render(title)
				desc = listDimmedDescStyle.Render(desc)
				runAt = listDimmedDescStyle.Render(runAt)
			} else {
				if isFiltered {
					unmatched := listNormalTitleStyle.Inline(true)
					matched := unmatched.Foreground(listMatchedColor)
					title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
				}
				title = listNormalTitleStyle.Render(title)
				desc = listNormalDescStyle.Render(desc)
				runAt = listNormalDescStyle.Render(runAt)
			}
		}
	}

	fmt.Fprintf(w, "%s\n%s\n%s", title, desc, runAt)
}
