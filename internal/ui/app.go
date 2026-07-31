package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/lusingander/gotip/internal/theme"
	"github.com/lusingander/gotip/internal/tip"
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
	styles          appStyles
	currentView     view
	showHelp        bool
	helpOffset      int
	matchFilterType matchFilterType
	fuzzyFilter     list.FilterFunc
	statusMsgType   statusMsgType
	w, h            int

	allBeforeSelected     int
	historyBeforeSelected int
	tmpTarget             *tip.Target
	retTarget             *tip.Target
}

func newModel(
	allTestItems, historyItems []list.Item,
	defaultView view,
	defaultFilterType matchFilterType,
	fuzzyFilter list.FilterFunc,
	colorTheme theme.ColorTheme,
) model {
	styles := newAppStyles(colorTheme)
	listStyles := newListStyles(colorTheme)
	allList := newList(allTestItems, testCaseItemDelegate{styles: listStyles}, defaultFilterType, fuzzyFilter, styles)
	historyList := newList(historyItems, historyItemDelegate{styles: listStyles}, defaultFilterType, fuzzyFilter, styles)
	return model{
		allList:               allList,
		historyList:           historyList,
		styles:                styles,
		currentView:           defaultView,
		showHelp:              false,
		helpOffset:            0,
		matchFilterType:       defaultFilterType,
		fuzzyFilter:           fuzzyFilter,
		statusMsgType:         noneStatusMsgType,
		allBeforeSelected:     -1,
		historyBeforeSelected: -1,
		tmpTarget:             nil,
		retTarget:             nil,
	}
}

func newList(
	items []list.Item,
	delegate list.ItemDelegate,
	defaultFilterType matchFilterType,
	fuzzyFilter list.FilterFunc,
	styles appStyles,
) list.Model {
	l := list.New(items, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.FilterInput.Prompt = "Filtering: "
	l.FilterInput.PromptStyle = styles.filterPrompt
	l.FilterInput.TextStyle = styles.filterText
	l.FilterInput.Cursor.Style = styles.filterCursor
	switch defaultFilterType {
	case fuzzyMatchFilterType:
		l.Filter = fuzzyFilter
	case exactMatchFilterType:
		l.Filter = exactMatchFilter
	}
	defaultKeyMap().applyToList(&l)
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
		m.allList.Filter = m.fuzzyFilter
		m.historyList.Filter = m.fuzzyFilter
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
		selectedLabel := "Selected: "
		selectedNameWidth := m.w - m.styles.header.GetHorizontalFrameSize() - lipgloss.Width(selectedLabel)
		if m.tmpTarget.IsPrefix {
			selectedNameWidth -= lipgloss.Width("*")
		}
		selectedName := ansi.Truncate(m.tmpTarget.TestNamePattern, selectedNameWidth, ellipsis)

		name := m.styles.selectedLabel.Render(selectedLabel) + m.styles.selectedName.Render(selectedName)
		if m.tmpTarget.IsPrefix {
			name += m.styles.selectedLabel.Render("*")
		}
		pack := m.styles.selectedLabel.Render(" Package: ") + m.styles.selectedPath.Render(m.tmpTarget.PackageName)
		headerContent = name + "\n" + pack
	} else {
		headerContent = "\n"
	}

	header := m.styles.header.Width(m.w).Render(headerContent)

	var footerStatus string
	switch m.statusMsgType {
	case noneStatusMsgType:
		switch currentList.FilterState() {
		case list.Filtering:
			footerStatus = trimRightSpace(currentList.FilterInput.View())
		case list.FilterApplied:
			footerStatus = m.styles.footerFiltered.
				Render(fmt.Sprintf("Filtered: %d items [Query: %s]", len(currentList.VisibleItems()), currentList.FilterValue()))
		}
	case fuzzyMatchFilteredStatusMsgType:
		footerStatus = m.styles.footerMsg.
			Render("Filter mode: Fuzzy match")
	case exactMatchFilteredStatusMsgType:
		footerStatus = m.styles.footerMsg.
			Render("Filter mode: Exact match")
	}

	var footerSelectedIndex string
	if len(currentList.VisibleItems()) > 0 {
		footerSelectedIndex = m.styles.footerSelectedIndex.
			Render(fmt.Sprintf("%d / %d", currentList.Index()+1, len(currentList.VisibleItems())))
	}

	var footerView string
	switch m.currentView {
	case allView:
		footerView = m.styles.footerDivider.Render(" | ") + m.styles.footerMsg.Render("All Tests")
	case historyView:
		footerView = m.styles.footerDivider.Render(" | ") + m.styles.footerMsg.Render("History  ")
	}

	footerSpaceWidth := max(m.w-lipgloss.Width(footerStatus)-lipgloss.Width(footerSelectedIndex)-lipgloss.Width(footerView)-2 /* padding */, 0)
	footerSpace := strings.Repeat(" ", footerSpaceWidth)

	footer := m.styles.footer.Width(m.w).Render(footerStatus + footerSpace + footerSelectedIndex + footerView)

	return lipgloss.JoinVertical(lipgloss.Left, header, currentList.View(), footer)
}

func (m model) helpView() string {
	headerProgramName := m.styles.helpHeader.Render(tip.ProgramName)
	headerVersion := m.styles.helpHeader.Render("Version: " + tip.AppVersion)
	header := m.styles.header.Width(m.w).Render(headerProgramName + "\n" + headerVersion)

	contentHeight := m.h - 5
	keyLines := []string{}
	descLines := []string{}
	for _, h := range helpItems() {
		keys := make([]string, 0, len(h.keys))
		for _, k := range h.keys {
			keys = append(keys, m.styles.helpText.Render("<")+m.styles.helpKey.Render(k)+m.styles.helpText.Render(">"))
		}
		keyLine := strings.Join(keys, m.styles.helpText.Render(", ")) + m.styles.helpText.Render(" : ")
		keyLines = append(keyLines, keyLine)
		descLines = append(descLines, m.styles.helpText.Render(h.desc))
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
	content := m.styles.helpContent.Render(strings.Join(lines, "\n") + padLines)

	footerView := m.styles.footerDivider.Render(" | ") + m.styles.footerMsg.Render("Help     ")

	footerSpaceWidth := max(m.w-lipgloss.Width(footerView)-2 /* padding */, 0)
	footerSpace := strings.Repeat(" ", footerSpaceWidth)

	footer := m.styles.footer.Width(m.w).Render(footerSpace + footerView)

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
	colorTheme theme.ColorTheme,
) (*tip.Target, error) {
	allTestItems := toTestCaseItems(tests)
	historyItems := toHistoryItems(histories, conf.History.DateFormat)
	defaultView := viewFromStr(defaultViewStr)
	defaultFilterType := matchFilterTypeFromStr(defaultFilterTypeStr)
	fuzzyFilter := fuzzyMatchFilterFromStr(conf.Filter.FuzzyMatcher)
	m := newModel(allTestItems, historyItems, defaultView, defaultFilterType, fuzzyFilter, colorTheme)
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
