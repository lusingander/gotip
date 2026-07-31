package tip

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

func defaultKeybindingsConfig() KeybindingsConfig {
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
