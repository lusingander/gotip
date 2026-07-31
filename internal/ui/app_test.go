package ui

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/gotip/internal/theme"
	"github.com/lusingander/gotip/internal/tip"
)

func TestToggleMatchFilterRestoresConfiguredFuzzyMatcher(t *testing.T) {
	m := newModel(
		nil,
		nil,
		allView,
		exactMatchFilterType,
		legacyFuzzyMatchFilter,
		theme.DefaultColorTheme(),
	)

	m.toggleMatchFilter()

	ranks := m.allList.Filter("abc", []string{"axxbxxcxxabc"})
	if len(ranks) != 1 {
		t.Fatalf("want 1 rank, got %d", len(ranks))
	}
	want := []int{0, 3, 6}
	for i, wantIndex := range want {
		if got := ranks[0].MatchedIndexes[i]; got != wantIndex {
			t.Errorf("matched index %d = %d, want %d", i, got, wantIndex)
		}
	}
}

func TestNewListUsesExistingKeybindings(t *testing.T) {
	l := newList(
		nil,
		testCaseItemDelegate{},
		fuzzyMatchFilterType,
		legacyFuzzyMatchFilter,
		newAppStyles(theme.DefaultColorTheme()),
	)

	tests := []struct {
		name string
		got  []string
		want []string
	}{
		{name: "cursor up", got: l.KeyMap.CursorUp.Keys(), want: []string{"up", "k"}},
		{name: "cursor down", got: l.KeyMap.CursorDown.Keys(), want: []string{"down", "j"}},
		{name: "previous page", got: l.KeyMap.PrevPage.Keys(), want: []string{"left", "h", "pgup", "b", "u"}},
		{name: "next page", got: l.KeyMap.NextPage.Keys(), want: []string{"right", "l", "pgdown", "f", "d"}},
		{name: "go to start", got: l.KeyMap.GoToStart.Keys(), want: []string{"home", "g"}},
		{name: "go to end", got: l.KeyMap.GoToEnd.Keys(), want: []string{"end", "G"}},
		{name: "filter", got: l.KeyMap.Filter.Keys(), want: []string{"/"}},
		{name: "clear filter", got: l.KeyMap.ClearFilter.Keys(), want: []string{"esc"}},
		{name: "cancel filter", got: l.KeyMap.CancelWhileFiltering.Keys(), want: []string{"esc"}},
		{name: "confirm filter", got: l.KeyMap.AcceptWhileFiltering.Keys(), want: []string{"enter"}},
		{name: "quit", got: l.KeyMap.Quit.Keys(), want: []string{"q", "esc"}},
		{name: "force quit", got: l.KeyMap.ForceQuit.Keys(), want: []string{"ctrl+c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("keys = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestExistingApplicationKeybindings(t *testing.T) {
	tests := []struct {
		name   string
		keyMsg tea.KeyMsg
		check  func(t *testing.T, got model)
	}{
		{
			name:   "run selected test",
			keyMsg: tea.KeyMsg{Type: tea.KeyEnter},
			check: func(t *testing.T, got model) {
				t.Helper()
				if got.retTarget == nil || got.retTarget.TestNamePattern != "TestFoo/Bar" {
					t.Errorf("retTarget = %#v, want TestFoo/Bar", got.retTarget)
				}
			},
		},
		{
			name:   "select parent with backspace",
			keyMsg: tea.KeyMsg{Type: tea.KeyBackspace},
			check: func(t *testing.T, got model) {
				t.Helper()
				if got.tmpTarget.TestNamePattern != "TestFoo/" {
					t.Errorf("target pattern = %q, want %q", got.tmpTarget.TestNamePattern, "TestFoo/")
				}
			},
		},
		{
			name:   "select parent with ctrl+h",
			keyMsg: tea.KeyMsg{Type: tea.KeyCtrlH},
			check: func(t *testing.T, got model) {
				t.Helper()
				if got.tmpTarget.TestNamePattern != "TestFoo/" {
					t.Errorf("target pattern = %q, want %q", got.tmpTarget.TestNamePattern, "TestFoo/")
				}
			},
		},
		{
			name:   "switch view with tab",
			keyMsg: tea.KeyMsg{Type: tea.KeyTab},
			check: func(t *testing.T, got model) {
				t.Helper()
				if got.currentView != historyView {
					t.Errorf("currentView = %v, want historyView", got.currentView)
				}
			},
		},
		{
			name:   "switch view with shift+tab",
			keyMsg: tea.KeyMsg{Type: tea.KeyShiftTab},
			check: func(t *testing.T, got model) {
				t.Helper()
				if got.currentView != historyView {
					t.Errorf("currentView = %v, want historyView", got.currentView)
				}
			},
		},
		{
			name:   "toggle filter type",
			keyMsg: tea.KeyMsg{Type: tea.KeyCtrlX},
			check: func(t *testing.T, got model) {
				t.Helper()
				if got.matchFilterType != exactMatchFilterType {
					t.Errorf("matchFilterType = %v, want exactMatchFilterType", got.matchFilterType)
				}
			},
		},
		{
			name:   "show help",
			keyMsg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}},
			check: func(t *testing.T, got model) {
				t.Helper()
				if !got.showHelp {
					t.Error("showHelp = false, want true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newModel(
				[]list.Item{&testCaseItem{path: "./foo/foo_test.go", name: "TestFoo/Bar"}},
				nil,
				allView,
				fuzzyMatchFilterType,
				legacyFuzzyMatchFilter,
				theme.DefaultColorTheme(),
			)
			m.tmpTarget = tip.NewTarget("./foo/foo_test.go", "TestFoo/Bar", false)
			m.allBeforeSelected = m.allList.GlobalIndex()

			updated, _ := m.Update(tt.keyMsg)
			tt.check(t, updated.(model))
		})
	}
}

func TestExistingHelpKeybindings(t *testing.T) {
	tests := []struct {
		name   string
		keyMsg tea.KeyMsg
	}{
		{name: "question mark", keyMsg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}},
		{name: "backspace", keyMsg: tea.KeyMsg{Type: tea.KeyBackspace}},
		{name: "ctrl+h", keyMsg: tea.KeyMsg{Type: tea.KeyCtrlH}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newModel(
				nil,
				nil,
				allView,
				fuzzyMatchFilterType,
				legacyFuzzyMatchFilter,
				theme.DefaultColorTheme(),
			)
			m.openHelp()

			updated, _ := m.Update(tt.keyMsg)
			if updated.(model).showHelp {
				t.Error("showHelp = true, want false")
			}
		})
	}
}
