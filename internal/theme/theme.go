package theme

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ColorTheme struct {
	Text      lipgloss.Color `toml:"text"`
	Accent    lipgloss.Color `toml:"accent"`
	Highlight lipgloss.Color `toml:"highlight"`
	Muted     lipgloss.Color `toml:"muted"`
	Dimmed    lipgloss.Color `toml:"dimmed"`
	Border    lipgloss.Color `toml:"border"`
	Match     lipgloss.Color `toml:"match"`
	Command   lipgloss.Color `toml:"command"`
}

func DefaultColorTheme() ColorTheme {
	return ColorTheme{
		Text:      lipgloss.Color("#DDDDDD"),
		Accent:    lipgloss.Color("#00ADD8"),
		Highlight: lipgloss.Color("#5DC9E2"),
		Muted:     lipgloss.Color("#777777"),
		Dimmed:    lipgloss.Color("#4D4D4D"),
		Border:    lipgloss.Color("240"),
		Match:     lipgloss.Color("#CE3262"),
		Command:   lipgloss.Color("#00A29C"),
	}
}

func (t ColorTheme) Validate() error {
	colors := []struct {
		name  string
		value lipgloss.Color
	}{
		{name: "text", value: t.Text},
		{name: "accent", value: t.Accent},
		{name: "highlight", value: t.Highlight},
		{name: "muted", value: t.Muted},
		{name: "dimmed", value: t.Dimmed},
		{name: "border", value: t.Border},
		{name: "match", value: t.Match},
		{name: "command", value: t.Command},
	}
	for _, color := range colors {
		if !isValidColor(string(color.value)) {
			return fmt.Errorf("invalid theme.%s color %q: expected #RGB, #RRGGBB, or an ANSI color number from 0 to 255", color.name, color.value)
		}
	}
	return nil
}

func isValidColor(value string) bool {
	if strings.HasPrefix(value, "#") {
		hex := value[1:]
		if len(hex) != 3 && len(hex) != 6 {
			return false
		}
		_, err := strconv.ParseUint(hex, 16, 24)
		return err == nil
	}

	n, err := strconv.Atoi(value)
	return err == nil && n >= 0 && n <= 255
}
