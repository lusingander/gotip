package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func TestTrimRightSpace_StyledSpaces(t *testing.T) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#DDDDDD"))
	input := style.Render("Filtering: ") + style.Render("query   ")

	got := trimRightSpace(input)
	if ansi.Strip(got) != "Filtering: query" {
		t.Fatalf("got %q", ansi.Strip(got))
	}
	if lipgloss.Width(got) != len("Filtering: query") {
		t.Fatalf("got width %d", lipgloss.Width(got))
	}
}
