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

	helpHeaderColor = lipgloss.Color("#00ADD8")
	helpKeyColor    = lipgloss.Color("#5DC9E2")
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

	helpHeaderStyle = lipgloss.NewStyle().Foreground(helpHeaderColor)

	helpContentStyle = lipgloss.NewStyle().Padding(0, 2)
	helpKeyStyle     = lipgloss.NewStyle().Foreground(helpKeyColor).Bold(true)
)

type view int

const (
	allView view = iota
	historyView
)

func viewFromStr(s string) view {
	switch s {
	case "all":
		return allView
	case "history":
		return historyView
	default:
		panic("unknown view type: " + s)
	}
}

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
	showHelp        bool
	helpOffset      int
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
		showHelp:              false,
		helpOffset:            0,
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

func (m *model) openHelp() {
	m.showHelp = true
	m.helpOffset = 0
}

func (m *model) closeHelp() {
	m.showHelp = false
	m.helpOffset = 0
}

func (m *model) scrollHelpUp() {
	if m.helpOffset > 0 {
		m.helpOffset--
	}
}

func (m *model) scrollHelpDown() {
	if m.helpOffset < len(helpItems())-1 {
		m.helpOffset++
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
		if msg.String() == "ctrl+c" {
			// exit
			return m, tea.Quit
		}

		// clear status message
		m.statusMsgType = noneStatusMsgType

		if m.allList.FilterState() == list.Filtering || m.historyList.FilterState() == list.Filtering {
			break
		}

		if m.showHelp {
			switch msg.String() {
			case "up", "k":
				m.scrollHelpUp()
			case "down", "j":
				m.scrollHelpDown()
			case "?", "backspace", "ctrl+h":
				m.closeHelp()
			}
			return m, nil
		}

		switch msg.String() {
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
		case "?":
			m.openHelp()
			return m, nil
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
	if m.showHelp {
		return m.helpView()
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

func (m model) helpView() string {
	headerProgramName := helpHeaderStyle.Render(tip.ProgramName)
	headerVersion := helpHeaderStyle.Render("Version: " + tip.AppVersion)
	header := headerStyle.Width(m.w).Render(headerProgramName + "\n" + headerVersion)

	contentHeight := m.h - 5
	keyLines := []string{}
	descLines := []string{}
	for _, h := range helpItems() {
		keys := make([]string, 0, len(h.keys))
		for _, k := range h.keys {
			keys = append(keys, "<"+helpKeyStyle.Render(k)+">")
		}
		keyLine := strings.Join(keys, ", ") + " : "
		keyLines = append(keyLines, keyLine)
		descLines = append(descLines, h.desc)
	}
	linesJoined := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Right, keyLines...),
		lipgloss.JoinVertical(lipgloss.Left, descLines...),
	)
	lines := []string{}
	for i, line := range strings.Split(linesJoined, "\n") {
		if i < m.helpOffset {
			continue
		}
		if len(lines) >= contentHeight {
			break
		}
		lines = append(lines, line)
	}

	padLines := strings.Repeat("\n", contentHeight-len(lines))
	content := helpContentStyle.Render(strings.Join(lines, "\n") + padLines)

	footerView := footerDividerStyle.Render(" | ") + footerMsgStyle.Render("Help     ")

	footerSpaceWidth := m.w - lipgloss.Width(footerView) - 2 /* padding */
	footerSpace := strings.Repeat(" ", footerSpaceWidth)

	footer := footerStyle.Width(m.w).Render(footerSpace + footerView)

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

type helpItem struct {
	keys []string
	desc string
}

func helpItems() []helpItem {
	return []helpItem{
		{keys: []string{"Ctrl-c"}, desc: "Quit"},
		{keys: []string{"Down", "j"}, desc: "Select next item"},
		{keys: []string{"Up", "k"}, desc: "Select previous item"},
		{keys: []string{"Right", "l"}, desc: "Select next page"},
		{keys: []string{"Left", "h"}, desc: "Select previous page"},
		{keys: []string{"Enter"}, desc: "Run the selected test / Confirm filter (in filtering mode)"},
		{keys: []string{"Backspace"}, desc: "Select parent test group"},
		{keys: []string{"/"}, desc: "Enter filtering mode"},
		{keys: []string{"Esc"}, desc: "Clear filtering mode"},
		{keys: []string{"Ctrl-x"}, desc: "Toggle filtering type"},
		{keys: []string{"Tab"}, desc: "Switch view"},
		{keys: []string{"?"}, desc: "Show help"},
	}
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
	defaultView := viewFromStr(defaultViewStr)
	defaultFilterType := matchFilterTypeFromStr(defaultFilterTypeStr)
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
