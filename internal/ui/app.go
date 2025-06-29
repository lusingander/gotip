package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/gotip/internal/tip"
)

type model struct {
	list list.Model

	target *tip.Target
}

func newModel(items []list.Item) model {
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "GOTIP"

	return model{
		list:   list,
		target: nil,
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
		m.list.SetSize(msg.Width, msg.Height)
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
				item := selected.(*testCaseItem)
				m.target = &tip.Target{
					Path:         item.path,
					Name:         item.name,
					IsUnresolved: item.isUnresolved,
				}
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
	return m.list.View()
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
