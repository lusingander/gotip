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

	packageName := target.RelativePathToPackageName()
	testNameRegex := target.TestNameToTestRunRegex()

	args := []string{"test", "-run", testNameRegex, packageName}
	cmd := exec.Command("go", append(args, extraArgs...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println(outputStyle.Render(cmd.String()))
	return cmd.Run()
}
