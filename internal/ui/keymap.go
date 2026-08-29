package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/lusingander/gotip/internal/tip"
)

type keyMap struct {
	cursorUp         key.Binding
	cursorDown       key.Binding
	previousPage     key.Binding
	nextPage         key.Binding
	goToStart        key.Binding
	goToEnd          key.Binding
	run              key.Binding
	parent           key.Binding
	startFilter      key.Binding
	clearFilter      key.Binding
	cancelFilter     key.Binding
	confirmFilter    key.Binding
	toggleFilterType key.Binding
	switchView       key.Binding
	showHelp         key.Binding
	closeHelp        key.Binding
	quit             key.Binding
	forceQuit        key.Binding
}

func newKeyMap(config tip.KeybindingsConfig) keyMap {
	return keyMap{
		cursorUp:         key.NewBinding(key.WithKeys(config.SelectPrevious...)),
		cursorDown:       key.NewBinding(key.WithKeys(config.SelectNext...)),
		previousPage:     key.NewBinding(key.WithKeys(config.PreviousPage...)),
		nextPage:         key.NewBinding(key.WithKeys(config.NextPage...)),
		goToStart:        key.NewBinding(key.WithKeys(config.GoToStart...)),
		goToEnd:          key.NewBinding(key.WithKeys(config.GoToEnd...)),
		run:              key.NewBinding(key.WithKeys(config.Run...)),
		parent:           key.NewBinding(key.WithKeys(config.Parent...)),
		startFilter:      key.NewBinding(key.WithKeys(config.StartFilter...)),
		clearFilter:      key.NewBinding(key.WithKeys(config.ClearFilter...)),
		cancelFilter:     key.NewBinding(key.WithKeys(config.CancelFilter...)),
		confirmFilter:    key.NewBinding(key.WithKeys(config.ConfirmFilter...)),
		toggleFilterType: key.NewBinding(key.WithKeys(config.ToggleFilterType...)),
		switchView:       key.NewBinding(key.WithKeys(config.SwitchView...)),
		showHelp:         key.NewBinding(key.WithKeys(config.ShowHelp...)),
		closeHelp:        key.NewBinding(key.WithKeys(config.CloseHelp...)),
		quit:             key.NewBinding(key.WithKeys(config.Quit...)),
		forceQuit:        key.NewBinding(key.WithKeys(config.ForceQuit...)),
	}
}

func (k keyMap) applyToList(l *list.Model) {
	l.KeyMap.CursorUp = k.cursorUp
	l.KeyMap.CursorDown = k.cursorDown
	l.KeyMap.PrevPage = k.previousPage
	l.KeyMap.NextPage = k.nextPage
	l.KeyMap.GoToStart = k.goToStart
	l.KeyMap.GoToEnd = k.goToEnd
	l.KeyMap.Filter = k.startFilter
	l.KeyMap.ClearFilter = k.clearFilter
	l.KeyMap.CancelWhileFiltering = k.cancelFilter
	l.KeyMap.AcceptWhileFiltering = k.confirmFilter
	l.KeyMap.ShowFullHelp = k.showHelp
	l.KeyMap.CloseFullHelp = k.closeHelp
	l.KeyMap.Quit = k.quit
	l.KeyMap.ForceQuit = k.forceQuit
}

func keyLabel(k string) string {
	if k == " " {
		return "Space"
	}

	var prefix strings.Builder
	for {
		matched := false
		for _, modifier := range []string{"ctrl", "alt", "shift"} {
			if remaining, ok := strings.CutPrefix(k, modifier+"+"); ok {
				prefix.WriteString(strings.ToUpper(modifier[:1]) + modifier[1:] + "-")
				k = remaining
				matched = true
				break
			}
		}
		if !matched {
			break
		}
	}
	switch k {
	case "up", "down", "left", "right", "enter", "backspace", "tab", "esc", "home", "end", "delete", "insert":
		k = strings.ToUpper(k[:1]) + k[1:]
	case "pgup":
		k = "PgUp"
	case "pgdown":
		k = "PgDown"
	}
	return prefix.String() + k
}
