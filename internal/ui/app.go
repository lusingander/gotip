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
	footerDividerStyle       = lipgloss.NewStyle().Foreground(borderColor)

	footerStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(borderColor)
)

type view int

const (
	allView view = iota
	historyView
)

type statusMsgType int

const (
	noneStatusMsgType statusMsgType = iota
	fuzzyMatchFilteredStatusMsgType
	exactMatchFilteredStatusMsgType
)

type model struct {
	allList         list.Model
	historyList     list.Model
	currentView     view
	matchFilterType matchFilterType
	statusMsgType   statusMsgType
	w, h            int

	allBeforeSelected     int
	historyBeforeSelected int
	tmpTarget             *tip.Target
	retTarget             *tip.Target
}

func newModel(allTestItems, historyItems []list.Item, defaultView view, defaultFilterType matchFilterType) model {
	allList := newList(allTestItems, testCaseItemDelegate{}, defaultFilterType)
	historyList := newList(historyItems, historyItemDelegate{}, defaultFilterType)
	return model{
		allList:               allList,
		historyList:           historyList,
		currentView:           defaultView,
		matchFilterType:       defaultFilterType,
		statusMsgType:         noneStatusMsgType,
		allBeforeSelected:     -1,
		historyBeforeSelected: -1,
		tmpTarget:             nil,
		retTarget:             nil,
	}
}

func newList(items []list.Item, delegate list.ItemDelegate, defaultFilterType matchFilterType) list.Model {
	l := list.New(items, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.FilterInput.Prompt = "Filtering: "
	l.FilterInput.PromptStyle = lipgloss.NewStyle()
	l.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(cursorColor)
	switch defaultFilterType {
	case fuzzyMatchFilterType:
		l.Filter = fuzzyMatchFilter
	case exactMatchFilterType:
		l.Filter = exactMatchFilter
	}
	return l
}

func (m *model) setSize(w, h int) {
	m.w, m.h = w, h
	m.allList.SetSize(w, h-5)
	m.historyList.SetSize(w, h-5)
}

func (m *model) toggleMatchFilter() {
	switch m.matchFilterType {
	case fuzzyMatchFilterType:
		m.allList.Filter = exactMatchFilter
		m.historyList.Filter = exactMatchFilter
		m.matchFilterType = exactMatchFilterType
		m.statusMsgType = exactMatchFilteredStatusMsgType
	case exactMatchFilterType:
		m.allList.Filter = fuzzyMatchFilter
		m.historyList.Filter = fuzzyMatchFilter
		m.matchFilterType = fuzzyMatchFilterType
		m.statusMsgType = fuzzyMatchFilteredStatusMsgType
	}
}

func (m *model) toggleView() {
	switch m.currentView {
	case allView:
		m.currentView = historyView
		m.updateCurrentSelectedHistoryItem()
	case historyView:
		m.currentView = allView
		m.updateCurrentSelectedAllItem()
	}
}

func (m *model) updateCurrentSelectedAllItem() {
	if m.allList.SelectedItem() != nil {
		selected := m.allList.SelectedItem().(*testCaseItem)
		m.tmpTarget = tip.NewTarget(selected.path, selected.name, selected.isUnresolved)
		m.allBeforeSelected = m.allList.GlobalIndex()
	}
}

func (m *model) updateCurrentSelectedHistoryItem() {
	if m.historyList.SelectedItem() != nil {
		selected := m.historyList.SelectedItem().(*historyItem)
		m.tmpTarget = tip.NewTarget(selected.path, selected.name, selected.isUnresolved)
		m.historyBeforeSelected = m.historyList.GlobalIndex()
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

		if m.allList.FilterState() == list.Filtering || m.historyList.FilterState() == list.Filtering {
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
		case "tab", "shift+tab":
			m.toggleView()
		case "ctrl+x":
			if m.allList.FilterState() == list.Unfiltered || m.historyList.FilterState() == list.Unfiltered {
				m.toggleMatchFilter()
			}
		}
	}

	switch m.currentView {
	case allView:
		newList, cmd := m.allList.Update(msg)
		m.allList = newList
		cmds = append(cmds, cmd)

		if m.allBeforeSelected != m.allList.GlobalIndex() {
			m.updateCurrentSelectedAllItem()
		}
	case historyView:
		newList, cmd := m.historyList.Update(msg)
		m.historyList = newList
		cmds = append(cmds, cmd)

		if m.historyBeforeSelected != m.historyList.GlobalIndex() {
			m.updateCurrentSelectedHistoryItem()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.w == 0 || m.h == 0 {
		return ""
	}

	var currentList list.Model
	switch m.currentView {
	case allView:
		currentList = m.allList
	case historyView:
		currentList = m.historyList
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

	var footerView string
	switch m.currentView {
	case allView:
		footerView = footerDividerStyle.Render(" | ") + footerMsgStyle.Render("All Tests")
	case historyView:
		footerView = footerDividerStyle.Render(" | ") + footerMsgStyle.Render("History  ")
	}

	footerSpaceWidth := m.w - lipgloss.Width(footerStatus) - lipgloss.Width(footerSelectedIndex) - lipgloss.Width(footerView) - 2 /* padding */
	footerSpace := strings.Repeat(" ", footerSpaceWidth)

	footer := footerStyle.Width(m.w).Render(footerStatus + footerSpace + footerSelectedIndex + footerView)

	return lipgloss.JoinVertical(lipgloss.Left, header, currentList.View(), footer)
}

func Start(
	tests map[string][]*tip.TestFunction,
	histories *tip.Histories,
	conf *tip.Config,
	defaultViewStr string,
	defaultFilterTypeStr string,
) (*tip.Target, error) {
	allTestItems := toTestCaseItems(tests)
	historyItems := toHistoryItems(histories, conf.History.DateFormat)
	defaultView := allView
	if defaultViewStr == "history" {
		defaultView = historyView
	}
	defaultFilterType := fuzzyMatchFilterType
	if defaultFilterTypeStr == "exact" {
		defaultFilterType = exactMatchFilterType
	}
	m := newModel(allTestItems, historyItems, defaultView, defaultFilterType)
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
