package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/gotip/internal/tip"
)

var (
	selectedColor     = lipgloss.Color("#00ADD8")
	listSelectedColor = lipgloss.Color("#5DC9E2")
	cursorColor       = lipgloss.Color("#00ADD8")
	borderColor       = lipgloss.Color("240")
)

type model struct {
	list list.Model
	w, h int

	target *tip.Target
}

func newModel(items []list.Item) model {
	itemDelegate := list.NewDefaultDelegate()
	itemDelegate.Styles.SelectedTitle = itemDelegate.Styles.SelectedTitle.Foreground(listSelectedColor).BorderForeground(listSelectedColor)
	itemDelegate.Styles.SelectedDesc = itemDelegate.Styles.SelectedDesc.Foreground(listSelectedColor).BorderForeground(listSelectedColor)
	list := list.New(items, itemDelegate, 0, 0)
	list.SetShowTitle(false)
	list.SetShowFilter(false)
	list.SetShowStatusBar(false)
	list.SetShowHelp(false)
	list.SetShowPagination(false)
	list.FilterInput.Prompt = "Filtering: "
	list.FilterInput.PromptStyle = lipgloss.NewStyle()
	list.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(cursorColor)

	return model{
		list:   list,
		target: nil,
	}
}

func (m *model) setSize(w, h int) {
	m.w, m.h = w, h
	m.list.SetSize(w, h-5)
}

var _ tea.Model = (*model)(nil)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			selected := m.list.SelectedItem()
			if selected != nil {
				m.target = selected.(*testCaseItem).ToTarget()
				return m, tea.Quit
			}
		}
	}

	newList, cmd := m.list.Update(msg)
	m.list = newList
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.w == 0 || m.h == 0 {
		return ""
	}

	var headerContent string
	if m.list.SelectedItem() != nil {
		selected := m.list.SelectedItem().(*testCaseItem)
		name := lipgloss.NewStyle().Foreground(selectedColor).Bold(true).Render(selected.name)
		path := lipgloss.NewStyle().Foreground(selectedColor).Render(selected.path)
		headerContent = name + "\n" + path
	} else {
		headerContent = "\n"
	}

	header := lipgloss.NewStyle().
		Padding(0, 2).
		Width(m.w).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(borderColor).
		Render(headerContent)

	var footerStatus string
	switch m.list.FilterState() {
	case list.Filtering:
		footerStatus = strings.TrimRight(m.list.FilterInput.View(), " ")
	case list.FilterApplied:
		footerStatus = lipgloss.NewStyle().
			Render(fmt.Sprintf("Filtered: %d items [Query: %s]", len(m.list.VisibleItems()), m.list.FilterValue()))
	}

	var footerSelected string
	if len(m.list.VisibleItems()) > 0 {
		footerSelected = lipgloss.NewStyle().
			Render(fmt.Sprintf("%d / %d", m.list.Index()+1, len(m.list.VisibleItems())))
	}

	footerSpaceWidth := m.w - lipgloss.Width(footerStatus) - lipgloss.Width(footerSelected) - 2 /* padding */
	footerSpace := strings.Repeat(" ", footerSpaceWidth)

	footer := lipgloss.NewStyle().
		Padding(0, 1).
		Width(m.w).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(borderColor).
		Render(footerStatus + footerSpace + footerSelected)

	return lipgloss.JoinVertical(lipgloss.Left, header, m.list.View(), footer)
}

func Start(tests map[string][]*tip.TestFunction) (*tip.Target, error) {
	items := toTestCaseItems(tests)
	m := newModel(items)
	p := tea.NewProgram(m, tea.WithAltScreen())
	ret, err := p.Run()
	if err != nil {
		return nil, err
	}
	return ret.(model).target, nil
}
