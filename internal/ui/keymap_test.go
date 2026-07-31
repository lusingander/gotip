package ui

import (
	"reflect"
	"testing"

	"github.com/lusingander/gotip/internal/theme"
	"github.com/lusingander/gotip/internal/tip"
)

func TestHelpItemsUseModelKeybindings(t *testing.T) {
	keybindings := tip.DefaultKeybindingsConfig()
	keybindings.Run = []string{"r", "ctrl+r"}
	m := newModel(
		nil,
		nil,
		allView,
		fuzzyMatchFilterType,
		legacyFuzzyMatchFilter,
		keybindings,
		theme.DefaultColorTheme(),
	)

	var got []string
	for _, item := range m.helpItems() {
		if item.desc == "Run the selected test" {
			got = item.keys
			break
		}
	}
	want := []string{"r", "ctrl+r"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("run help keys = %v, want %v", got, want)
	}

	m.keys.run.Unbind()
	for _, item := range m.helpItems() {
		if item.desc == "Run the selected test" {
			t.Error("unbound run action is shown in help")
		}
	}
}

func TestKeyLabel(t *testing.T) {
	tests := []struct {
		key  string
		want string
	}{
		{key: "ctrl+c", want: "Ctrl-c"},
		{key: "shift+tab", want: "Shift-Tab"},
		{key: "pgdown", want: "PgDown"},
		{key: "down", want: "Down"},
		{key: "G", want: "G"},
		{key: " ", want: "Space"},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := keyLabel(tt.key); got != tt.want {
				t.Errorf("keyLabel(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}
