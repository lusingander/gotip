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
	allList         list.Model
	matchFilterType matchFilterType
	statusMsgType   statusMsgType
	w, h            int

	beforeSelected int
	tmpTarget      *tip.Target
	retTarget      *tip.Target
}

func newModel(items []list.Item) model {
	allList := list.New(items, testCaseItemDelegate{}, 0, 0)
	allList.SetShowTitle(false)
	allList.SetShowFilter(false)
	allList.SetShowStatusBar(false)
	allList.SetShowHelp(false)
	allList.SetShowPagination(false)
	allList.FilterInput.Prompt = "Filtering: "
	allList.FilterInput.PromptStyle = lipgloss.NewStyle()
	allList.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(cursorColor)
	allList.Filter = fuzzyMatchFilter
	matchFilterType := fuzzyMatchFilterType

	return model{
		allList:         allList,
		matchFilterType: matchFilterType,
		statusMsgType:   noneStatusMsgType,
		beforeSelected:  -1,
		tmpTarget:       nil,
		retTarget:       nil,
	}
}

func (m *model) setSize(w, h int) {
	m.w, m.h = w, h
	m.allList.SetSize(w, h-5)
}

func (m *model) toggleMatchFilter() {
	switch m.matchFilterType {
	case fuzzyMatchFilterType:
		m.allList.Filter = exactMatchFilter
		m.matchFilterType = exactMatchFilterType
		m.statusMsgType = exactMatchFilteredStatusMsgType
	case exactMatchFilterType:
		m.allList.Filter = fuzzyMatchFilter
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

		if m.allList.FilterState() == list.Filtering {
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
			if m.allList.FilterState() == list.Unfiltered {
				m.toggleMatchFilter()
			}
		}
	}

	newList, cmd := m.allList.Update(msg)
	m.allList = newList
	cmds = append(cmds, cmd)

	if m.beforeSelected != m.allList.GlobalIndex() && m.allList.SelectedItem() != nil {
		selected := m.allList.SelectedItem().(*testCaseItem)
		m.tmpTarget = tip.NewTarget(selected.path, selected.name, selected.isUnresolved)
		m.beforeSelected = m.allList.GlobalIndex()
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.w == 0 || m.h == 0 {
		return ""
	}

	currentList := m.allList

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
		switch currentList.FilterState() {
		case list.Filtering:
			footerStatus = strings.TrimRight(currentList.FilterInput.View(), " ")
		case list.FilterApplied:
			footerStatus = footerFilteredStyle.
				Render(fmt.Sprintf("Filtered: %d items [Query: %s]", len(currentList.VisibleItems()), currentList.FilterValue()))
		}
	case fuzzyMatchFilteredStatusMsgType:
		footerStatus = footerMsgStyle.
			Render("Filter mode: Fuzzy match")
	case exactMatchFilteredStatusMsgType:
		footerStatus = footerMsgStyle.
			Render("Filter mode: Exact match")
	}

	var footerSelectedIndex string
	if len(currentList.VisibleItems()) > 0 {
		footerSelectedIndex = footerSelectedIndexStyle.
			Render(fmt.Sprintf("%d / %d", currentList.Index()+1, len(currentList.VisibleItems())))
	}

	footerSpaceWidth := m.w - lipgloss.Width(footerStatus) - lipgloss.Width(footerSelectedIndex) - 2 /* padding */
	footerSpace := strings.Repeat(" ", footerSpaceWidth)

	footer := footerStyle.Width(m.w).Render(footerStatus + footerSpace + footerSelectedIndex)

	return lipgloss.JoinVertical(lipgloss.Left, header, currentList.View(), footer)
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
