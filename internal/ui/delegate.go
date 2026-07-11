package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

const (
	ellipsis = "..."
)

type testCaseItemDelegate struct {
	styles listStyles
}

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

	textwidth := m.Width() - d.styles.normalTitle.GetPaddingLeft() - d.styles.normalTitle.GetPaddingRight()
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
		title = d.styles.dimmedTitle.Render(title)
		desc = d.styles.dimmedDesc.Render(desc)
	} else {
		if isSelected && m.FilterState() != list.Filtering {
			if isFiltered {
				unmatched := d.styles.selectedTitle.Inline(true)
				matched := unmatched.Foreground(d.styles.matchedColor)
				title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
			}
			title = d.styles.selectedTitle.Render(title)
			desc = d.styles.selectedDesc.Render(desc)
		} else {
			if m.FilterState() == list.Filtering {
				if isFiltered {
					unmatched := d.styles.dimmedTitle.Inline(true)
					matched := unmatched.Foreground(d.styles.matchedColor)
					title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
				}
				title = d.styles.dimmedTitle.Render(title)
				desc = d.styles.dimmedDesc.Render(desc)
			} else {
				if isFiltered {
					unmatched := d.styles.normalTitle.Inline(true)
					matched := unmatched.Foreground(d.styles.matchedColor)
					title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
				}
				title = d.styles.normalTitle.Render(title)
				desc = d.styles.normalDesc.Render(desc)
			}
		}
	}

	fmt.Fprintf(w, "%s\n%s", title, desc)
}

type historyItemDelegate struct {
	styles listStyles
}

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

	textwidth := m.Width() - d.styles.normalTitle.GetPaddingLeft() - d.styles.normalTitle.GetPaddingRight()
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
		title = d.styles.dimmedTitle.Render(title)
		desc = d.styles.dimmedDesc.Render(desc)
		runAt = d.styles.dimmedDesc.Render(runAt)
	} else {
		if isSelected && m.FilterState() != list.Filtering {
			if isFiltered {
				unmatched := d.styles.selectedTitle.Inline(true)
				matched := unmatched.Foreground(d.styles.matchedColor)
				title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
			}
			title = d.styles.selectedTitle.Render(title)
			desc = d.styles.selectedDesc.Render(desc)
			runAt = d.styles.selectedDesc.Render(runAt)
		} else {
			if m.FilterState() == list.Filtering {
				if isFiltered {
					unmatched := d.styles.dimmedTitle.Inline(true)
					matched := unmatched.Foreground(d.styles.matchedColor)
					title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
				}
				title = d.styles.dimmedTitle.Render(title)
				desc = d.styles.dimmedDesc.Render(desc)
				runAt = d.styles.dimmedDesc.Render(runAt)
			} else {
				if isFiltered {
					unmatched := d.styles.normalTitle.Inline(true)
					matched := unmatched.Foreground(d.styles.matchedColor)
					title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
				}
				title = d.styles.normalTitle.Render(title)
				desc = d.styles.normalDesc.Render(desc)
				runAt = d.styles.normalDesc.Render(runAt)
			}
		}
	}

	fmt.Fprintf(w, "%s\n%s\n%s", title, desc, runAt)
}
