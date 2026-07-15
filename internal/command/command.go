package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/gotip/internal/theme"
	"github.com/lusingander/gotip/internal/tip"
)

const (
	commandTestNameMarker = "${name}"
	commandPackageMarker  = "${package}"
)

type styles struct {
	output lipgloss.Style
}

func newStyles(colorTheme theme.ColorTheme) styles {
	return styles{
		output: lipgloss.NewStyle().Foreground(colorTheme.Command),
	}
}

func Test(target *tip.Target, extraArgs []string, conf *tip.Config, colorTheme theme.ColorTheme) (int, error) {
	if target == nil {
		return 0, nil
	}
	styles := newStyles(colorTheme)

	nameRegex := testNameToTestRunRegex(target.TestNamePattern, target.IsPrefix)

	cmd := buildTestExecCommand(target, nameRegex, extraArgs, conf.Command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Fprintln(os.Stderr, styles.output.Render(cmd.String()))
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
	segments := strings.Split(pattern, "/")
	for i, segment := range segments {
		if segment == "" {
			continue
		}
		if isPrefix && i == len(segments)-1 {
			segments[i] = "^" + segment
		} else {
			segments[i] = "^" + segment + "$"
		}
	}
	return strings.Join(segments, "/")
}
