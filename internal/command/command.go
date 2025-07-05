package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/gotip/internal/tip"
)

var outputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00A29C"))

func Test(target *tip.Target, extraArgs []string) error {
	if target == nil {
		return nil
	}

	nameRegex := testNameToTestRunRegex(target.TestNamePattern, target.IsPrefix)

	args := []string{"test"}
	if target.TestNamePattern != "" {
		args = append(args, "-run", nameRegex)
	}
	args = append(args, target.PackageName)

	cmd := exec.Command("go", append(args, extraArgs...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println(outputStyle.Render(cmd.String()))
	return cmd.Run()
}

func testNameToTestRunRegex(pattern string, isPrefix bool) string {
	if isPrefix {
		return fmt.Sprintf("^%s", pattern)
	}
	return fmt.Sprintf("^%s$", pattern)
}
