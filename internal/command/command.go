package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/gotip/internal/tip"
)

const (
	commandTestNameMarker = "${name}"
	commandPackageMarker  = "${package}"
)

var outputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00A29C"))

func Test(target *tip.Target, extraArgs []string, conf *tip.Config) (int, error) {
	if target == nil {
		return 0, nil
	}

	nameRegex := testNameToTestRunRegex(target.TestNamePattern, target.IsPrefix)

	cmd := buildTestExecCommand(target, nameRegex, extraArgs, conf.Command)
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

func buildTestExecCommand(target *tip.Target, nameRegex string, extraArgs []string, command []string) *exec.Cmd {
	if len(command) == 0 {
		// default Go test command
		args := []string{"test"}
		if target.TestNamePattern != "" {
			args = append(args, "-run", nameRegex)
		}
		args = append(args, target.PackageName)

		return exec.Command("go", append(args, extraArgs...)...)
	}

	// custom command from configuration
	args := make([]string, 0)
	if len(command) > 1 {
		for _, arg := range command[1:] {
			switch arg {
			case commandTestNameMarker:
				args = append(args, nameRegex)
			case commandPackageMarker:
				args = append(args, target.PackageName)
			default:
				args = append(args, arg)
			}
		}
	}
	return exec.Command(command[0], append(args, extraArgs...)...)
}

func testNameToTestRunRegex(pattern string, isPrefix bool) string {
	if isPrefix {
		return fmt.Sprintf("^%s", pattern)
	}
	return fmt.Sprintf("^%s$", pattern)
}
