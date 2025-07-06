package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/gotip/internal/tip"
)

var (
	selectedColor = lipgloss.Color("#00ADD8")
	cursorColor   = lipgloss.Color("#00ADD8")
	borderColor   = lipgloss.Color("240")
)

var (
	selectedLabelStyle = lipgloss.NewStyle().Foreground(selectedColor)
	selectedNameStyle  = lipgloss.NewStyle().Foreground(selectedColor).Bold(true)
	selectedPathStyle  = lipgloss.NewStyle().Foreground(selectedColor).Bold(true)

	headerStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(borderColor)

	footerMsgStyle           = lipgloss.NewStyle()
	footerFilteredStyle      = lipgloss.NewStyle()
	footerSelectedIndexStyle = lipgloss.NewStyle()

	footerStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(borderColor)
)

type statusMsgType int

const (
	noneStatusMsgType statusMsgType = iota
	fuzzyMatchFilteredStatusMsgType
	exactMatchFilteredStatusMsgType
)

type model struct {
	list            list.Model
	matchFilterType matchFilterType
	statusMsgType   statusMsgType
	w, h            int

	beforeSelected int
	tmpTarget      *tip.Target
	retTarget      *tip.Target
}

func newModel(items []list.Item) model {
	list := list.New(items, itemDelegate{}, 0, 0)
	list.SetShowTitle(false)
	list.SetShowFilter(false)
	list.SetShowStatusBar(false)
	list.SetShowHelp(false)
	list.SetShowPagination(false)
	list.FilterInput.Prompt = "Filtering: "
	list.FilterInput.PromptStyle = lipgloss.NewStyle()
	list.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(cursorColor)
	list.Filter = fuzzyMatchFilter
	matchFilterType := fuzzyMatchFilterType

	return model{
		list:            list,
		matchFilterType: matchFilterType,
		statusMsgType:   noneStatusMsgType,
		beforeSelected:  -1,
		tmpTarget:       nil,
		retTarget:       nil,
	}
}

func (m *model) setSize(w, h int) {
	m.w, m.h = w, h
	m.list.SetSize(w, h-5)
}

func (m *model) toggleMatchFilter() {
	switch m.matchFilterType {
	case fuzzyMatchFilterType:
		m.list.Filter = exactMatchFilter
		m.matchFilterType = exactMatchFilterType
		m.statusMsgType = exactMatchFilteredStatusMsgType
	case exactMatchFilterType:
		m.list.Filter = fuzzyMatchFilter
		m.matchFilterType = fuzzyMatchFilterType
		m.statusMsgType = fuzzyMatchFilteredStatusMsgType
	}
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
		// clear status message
		m.statusMsgType = noneStatusMsgType

		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.retTarget = m.tmpTarget
			return m, tea.Quit
		case "backspace", "ctrl+h":
			if m.tmpTarget != nil {
				m.tmpTarget.DropLastSegment()
			}
		case "ctrl+x":
			if m.list.FilterState() == list.Unfiltered {
				m.toggleMatchFilter()
			}
		}
	}

	newList, cmd := m.list.Update(msg)
	m.list = newList
	cmds = append(cmds, cmd)

	if m.beforeSelected != m.list.GlobalIndex() && m.list.SelectedItem() != nil {
		selected := m.list.SelectedItem().(*testCaseItem)
		m.tmpTarget = tip.NewTarget(selected.path, selected.name, selected.isUnresolved)
		m.beforeSelected = m.list.GlobalIndex()
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.w == 0 || m.h == 0 {
		return ""
	}

	var headerContent string
	if m.tmpTarget != nil {
		name := selectedLabelStyle.Render("Selected: ") + selectedNameStyle.Render(m.tmpTarget.TestNamePattern)
		if m.tmpTarget.IsPrefix {
			name += selectedLabelStyle.Render("*")
		}
		pack := selectedLabelStyle.Render(" Package: ") + selectedPathStyle.Render(m.tmpTarget.PackageName)
		headerContent = name + "\n" + pack
	} else {
		headerContent = "\n"
	}

	header := headerStyle.Width(m.w).Render(headerContent)

	var footerStatus string
	switch m.statusMsgType {
	case noneStatusMsgType:
		switch m.list.FilterState() {
		case list.Filtering:
			footerStatus = strings.TrimRight(m.list.FilterInput.View(), " ")
		case list.FilterApplied:
			footerStatus = footerFilteredStyle.
				Render(fmt.Sprintf("Filtered: %d items [Query: %s]", len(m.list.VisibleItems()), m.list.FilterValue()))
		}
	case fuzzyMatchFilteredStatusMsgType:
		footerStatus = footerMsgStyle.
			Render("Filter mode: Fuzzy match")
	case exactMatchFilteredStatusMsgType:
		footerStatus = footerMsgStyle.
			Render("Filter mode: Exact match")
	}

	var footerSelectedIndex string
	if len(m.list.VisibleItems()) > 0 {
		footerSelectedIndex = footerSelectedIndexStyle.
			Render(fmt.Sprintf("%d / %d", m.list.Index()+1, len(m.list.VisibleItems())))
	}

	footerSpaceWidth := m.w - lipgloss.Width(footerStatus) - lipgloss.Width(footerSelectedIndex) - 2 /* padding */
	footerSpace := strings.Repeat(" ", footerSpaceWidth)

	footer := footerStyle.Width(m.w).Render(footerStatus + footerSpace + footerSelectedIndex)

	return lipgloss.JoinVertical(lipgloss.Left, header, m.list.View(), footer)
}

func Start(tests map[string][]*tip.TestFunction) (*tip.Target, error) {
	items := toTestCaseItems(tests)
	m := newModel(items)
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithOutput(os.Stderr),
	)
	ret, err := p.Run()
	if err != nil {
		return nil, err
	}
	return ret.(model).retTarget, nil
}
