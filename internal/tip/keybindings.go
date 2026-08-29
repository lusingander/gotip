package tip

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type KeybindingsConfig struct {
	SelectPrevious   []string `toml:"select_previous"`
	SelectNext       []string `toml:"select_next"`
	PreviousPage     []string `toml:"previous_page"`
	NextPage         []string `toml:"next_page"`
	GoToStart        []string `toml:"go_to_start"`
	GoToEnd          []string `toml:"go_to_end"`
	Run              []string `toml:"run"`
	Parent           []string `toml:"parent"`
	StartFilter      []string `toml:"start_filter"`
	ClearFilter      []string `toml:"clear_filter"`
	CancelFilter     []string `toml:"cancel_filter"`
	ConfirmFilter    []string `toml:"confirm_filter"`
	ToggleFilterType []string `toml:"toggle_filter_type"`
	SwitchView       []string `toml:"switch_view"`
	ShowHelp         []string `toml:"show_help"`
	CloseHelp        []string `toml:"close_help"`
	Quit             []string `toml:"quit"`
	ForceQuit        []string `toml:"force_quit"`
}

// DefaultKeybindingsConfig returns the built-in UI keybindings.
func DefaultKeybindingsConfig() KeybindingsConfig {
	return KeybindingsConfig{
		SelectPrevious:   []string{"up", "k"},
		SelectNext:       []string{"down", "j"},
		PreviousPage:     []string{"left", "h", "pgup", "b", "u"},
		NextPage:         []string{"right", "l", "pgdown", "f", "d"},
		GoToStart:        []string{"home", "g"},
		GoToEnd:          []string{"end", "G"},
		Run:              []string{"enter"},
		Parent:           []string{"backspace", "ctrl+h"},
		StartFilter:      []string{"/"},
		ClearFilter:      []string{"esc"},
		CancelFilter:     []string{"esc"},
		ConfirmFilter:    []string{"enter"},
		ToggleFilterType: []string{"ctrl+x"},
		SwitchView:       []string{"tab", "shift+tab"},
		ShowHelp:         []string{"?"},
		CloseHelp:        []string{"?", "backspace", "ctrl+h"},
		Quit:             []string{"q", "esc"},
		ForceQuit:        []string{"ctrl+c"},
	}
}

type namedKeybinding struct {
	name string
	keys []string
}

func (c KeybindingsConfig) validate() error {
	bindings := c.namedBindings()
	for _, binding := range bindings {
		seen := make(map[string]struct{}, len(binding.keys))
		for _, key := range binding.keys {
			if !validKeyName(key) {
				return fmt.Errorf("invalid keybindings.%s key %q", binding.name, key)
			}
			if _, ok := seen[key]; ok {
				return fmt.Errorf("keybindings.%s contains duplicate key %q", binding.name, key)
			}
			seen[key] = struct{}{}
		}
	}
	if len(c.ForceQuit) == 0 {
		return fmt.Errorf("keybindings.force_quit must contain at least one key")
	}

	contexts := []struct {
		name     string
		bindings []namedKeybinding
	}{
		{name: "browsing", bindings: c.browsingBindings(c.binding("quit", c.Quit))},
		{name: "filtered browsing", bindings: c.browsingBindings(c.binding("clear_filter", c.ClearFilter))},
		{
			name: "filtering",
			bindings: []namedKeybinding{
				c.binding("cancel_filter", c.CancelFilter),
				c.binding("confirm_filter", c.ConfirmFilter),
			},
		},
		{
			name: "help",
			bindings: []namedKeybinding{
				c.binding("select_previous", c.SelectPrevious),
				c.binding("select_next", c.SelectNext),
				c.binding("close_help", c.CloseHelp),
			},
		},
	}
	for _, context := range contexts {
		if err := validateKeybindingContext(context.name, context.bindings); err != nil {
			return err
		}
	}

	for _, forceKey := range c.ForceQuit {
		for _, binding := range bindings {
			if binding.name == "force_quit" {
				continue
			}
			for _, key := range binding.keys {
				if forceKey == key {
					return fmt.Errorf(
						"key %q is assigned to both keybindings.force_quit and keybindings.%s",
						key,
						binding.name,
					)
				}
			}
		}
	}
	return nil
}

func (c KeybindingsConfig) namedBindings() []namedKeybinding {
	return []namedKeybinding{
		c.binding("select_previous", c.SelectPrevious),
		c.binding("select_next", c.SelectNext),
		c.binding("previous_page", c.PreviousPage),
		c.binding("next_page", c.NextPage),
		c.binding("go_to_start", c.GoToStart),
		c.binding("go_to_end", c.GoToEnd),
		c.binding("run", c.Run),
		c.binding("parent", c.Parent),
		c.binding("start_filter", c.StartFilter),
		c.binding("clear_filter", c.ClearFilter),
		c.binding("cancel_filter", c.CancelFilter),
		c.binding("confirm_filter", c.ConfirmFilter),
		c.binding("toggle_filter_type", c.ToggleFilterType),
		c.binding("switch_view", c.SwitchView),
		c.binding("show_help", c.ShowHelp),
		c.binding("close_help", c.CloseHelp),
		c.binding("quit", c.Quit),
		c.binding("force_quit", c.ForceQuit),
	}
}

func (c KeybindingsConfig) browsingBindings(stateBinding namedKeybinding) []namedKeybinding {
	return []namedKeybinding{
		c.binding("select_previous", c.SelectPrevious),
		c.binding("select_next", c.SelectNext),
		c.binding("previous_page", c.PreviousPage),
		c.binding("next_page", c.NextPage),
		c.binding("go_to_start", c.GoToStart),
		c.binding("go_to_end", c.GoToEnd),
		c.binding("run", c.Run),
		c.binding("parent", c.Parent),
		c.binding("start_filter", c.StartFilter),
		c.binding("toggle_filter_type", c.ToggleFilterType),
		c.binding("switch_view", c.SwitchView),
		c.binding("show_help", c.ShowHelp),
		stateBinding,
	}
}

func (KeybindingsConfig) binding(name string, keys []string) namedKeybinding {
	return namedKeybinding{name: name, keys: keys}
}

func validateKeybindingContext(context string, bindings []namedKeybinding) error {
	assigned := make(map[string]string)
	for _, binding := range bindings {
		for _, key := range binding.keys {
			if previous, ok := assigned[key]; ok {
				return fmt.Errorf(
					"key %q is assigned to both keybindings.%s and keybindings.%s in the %s context",
					key,
					previous,
					binding.name,
					context,
				)
			}
			assigned[key] = binding.name
		}
	}
	return nil
}

func validKeyName(key string) bool {
	if remaining, ok := strings.CutPrefix(key, "alt+"); ok {
		if remaining == "" || strings.HasPrefix(remaining, "alt+") {
			return false
		}
		key = remaining
	}

	if utf8.RuneCountInString(key) == 1 {
		r, _ := utf8.DecodeRuneInString(key)
		return unicode.IsPrint(r)
	}

	if _, ok := namedKeys[key]; ok {
		return true
	}
	if after, ok := strings.CutPrefix(key, "f"); ok {
		number, err := strconv.Atoi(after)
		return err == nil && number >= 1 && number <= 20
	}
	return false
}

var namedKeys = map[string]struct{}{
	"ctrl+@":           {},
	"ctrl+a":           {},
	"ctrl+b":           {},
	"ctrl+c":           {},
	"ctrl+d":           {},
	"ctrl+e":           {},
	"ctrl+f":           {},
	"ctrl+g":           {},
	"ctrl+h":           {},
	"ctrl+j":           {},
	"ctrl+k":           {},
	"ctrl+l":           {},
	"ctrl+n":           {},
	"ctrl+o":           {},
	"ctrl+p":           {},
	"ctrl+q":           {},
	"ctrl+r":           {},
	"ctrl+s":           {},
	"ctrl+t":           {},
	"ctrl+u":           {},
	"ctrl+v":           {},
	"ctrl+w":           {},
	"ctrl+x":           {},
	"ctrl+y":           {},
	"ctrl+z":           {},
	"ctrl+\\":          {},
	"ctrl+]":           {},
	"ctrl+^":           {},
	"ctrl+_":           {},
	"enter":            {},
	"backspace":        {},
	"tab":              {},
	"esc":              {},
	"up":               {},
	"down":             {},
	"right":            {},
	"left":             {},
	"shift+tab":        {},
	"home":             {},
	"end":              {},
	"pgup":             {},
	"pgdown":           {},
	"ctrl+pgup":        {},
	"ctrl+pgdown":      {},
	"delete":           {},
	"insert":           {},
	"ctrl+up":          {},
	"ctrl+down":        {},
	"ctrl+right":       {},
	"ctrl+left":        {},
	"ctrl+home":        {},
	"ctrl+end":         {},
	"shift+up":         {},
	"shift+down":       {},
	"shift+right":      {},
	"shift+left":       {},
	"shift+home":       {},
	"shift+end":        {},
	"ctrl+shift+up":    {},
	"ctrl+shift+down":  {},
	"ctrl+shift+left":  {},
	"ctrl+shift+right": {},
	"ctrl+shift+home":  {},
	"ctrl+shift+end":   {},
}
