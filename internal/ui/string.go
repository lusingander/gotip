package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// trimRightSpace trims styled trailing spaces without breaking ANSI escape sequences.
func trimRightSpace(s string) string {
	width := lipgloss.Width(strings.TrimRight(ansi.Strip(s), " "))
	return ansi.Truncate(s, width, "")
}
