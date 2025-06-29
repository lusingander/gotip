package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/gotip/internal/tip"
)

var outputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("34"))

func Test(target *tip.Target) error {
	if target == nil {
		return nil
	}

	packageName := relativePathToPackageName(target.Path)
	testNameRegex := testNameToTestRunRegex(target.Name, target.IsUnresolved)

	cmd := exec.Command("go", "test", "-run", testNameRegex, packageName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println(outputStyle.Render(cmd.String()))
	return cmd.Run()
}

func relativePathToPackageName(path string) string {
	name := filepath.Dir(path)
	name = filepath.ToSlash(name)
	if !strings.HasPrefix(name, "./") {
		name = "./" + name
	}
	return name
}

func testNameToTestRunRegex(name string, isUnresolved bool) string {
	if isUnresolved {
		return fmt.Sprintf("^%s", strings.TrimSuffix(name, tip.UnresolvedTestCaseName))
	}
	return fmt.Sprintf("^%s$", name)
}
