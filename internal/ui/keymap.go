package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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

func defaultKeyMap() keyMap {
	return keyMap{
		cursorUp:         key.NewBinding(key.WithKeys("up", "k")),
		cursorDown:       key.NewBinding(key.WithKeys("down", "j")),
		previousPage:     key.NewBinding(key.WithKeys("left", "h", "pgup", "b", "u")),
		nextPage:         key.NewBinding(key.WithKeys("right", "l", "pgdown", "f", "d")),
		goToStart:        key.NewBinding(key.WithKeys("home", "g")),
		goToEnd:          key.NewBinding(key.WithKeys("end", "G")),
		run:              key.NewBinding(key.WithKeys("enter")),
		parent:           key.NewBinding(key.WithKeys("backspace", "ctrl+h")),
		startFilter:      key.NewBinding(key.WithKeys("/")),
		clearFilter:      key.NewBinding(key.WithKeys("esc")),
		cancelFilter:     key.NewBinding(key.WithKeys("esc")),
		confirmFilter:    key.NewBinding(key.WithKeys("enter")),
		toggleFilterType: key.NewBinding(key.WithKeys("ctrl+x")),
		switchView:       key.NewBinding(key.WithKeys("tab", "shift+tab")),
		showHelp:         key.NewBinding(key.WithKeys("?")),
		closeHelp:        key.NewBinding(key.WithKeys("?", "backspace", "ctrl+h")),
		quit:             key.NewBinding(key.WithKeys("q", "esc")),
		forceQuit:        key.NewBinding(key.WithKeys("ctrl+c")),
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
