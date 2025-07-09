package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/gotip/internal/tip"
)

var outputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00A29C"))

func Test(target *tip.Target, extraArgs []string) (int, error) {
	if target == nil {
		return 0, nil
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

	fmt.Fprintln(os.Stderr, outputStyle.Render(cmd.String()))
	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return 1, err
		}
	}
	return cmd.ProcessState.ExitCode(), nil
}

func testNameToTestRunRegex(pattern string, isPrefix bool) string {
	if isPrefix {
		return fmt.Sprintf("^%s", pattern)
	}
	return fmt.Sprintf("^%s$", pattern)
}
