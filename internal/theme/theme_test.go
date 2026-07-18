package theme

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestColorThemeValidate(t *testing.T) {
	tests := []struct {
		name  string
		color lipgloss.Color
	}{
		{name: "short hex", color: "#0aF"},
		{name: "long hex", color: "#00aF12"},
		{name: "lowest ANSI", color: "0"},
		{name: "highest ANSI", color: "255"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colorTheme := themeWithAllColors(tt.color)
			if err := colorTheme.Validate(); err != nil {
				t.Fatalf("Validate() error = %v", err)
			}
		})
	}
}

func TestColorThemeValidateRejectsInvalidColor(t *testing.T) {
	tests := []struct {
		name  string
		color lipgloss.Color
	}{
		{name: "empty", color: ""},
		{name: "invalid hex digit", color: "#12xz56"},
		{name: "invalid hex length", color: "#12345"},
		{name: "negative ANSI", color: "-1"},
		{name: "ANSI out of range", color: "256"},
		{name: "color name", color: "red"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colorTheme := themeWithAllColors("#123456")
			colorTheme.Accent = tt.color
			err := colorTheme.Validate()
			if err == nil {
				t.Fatal("Validate() error = nil, want error")
			}
			if !strings.Contains(err.Error(), "theme.accent") {
				t.Fatalf("Validate() error = %q, want theme.accent context", err)
			}
		})
	}
}

func themeWithAllColors(color lipgloss.Color) ColorTheme {
	return ColorTheme{
		Text:      color,
		Accent:    color,
		Highlight: color,
		Muted:     color,
		Dimmed:    color,
		Border:    color,
		Match:     color,
		Command:   color,
	}
}
