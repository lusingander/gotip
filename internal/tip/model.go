package tip

import (
	"fmt"
	"path/filepath"
	"strings"
)

type TestFunction struct {
	Name string
	Subs []*SubTest
}

type SubTest struct {
	Name                 string
	Subs                 []*SubTest
	IsUnresolvedSubTests bool
}

type Target struct {
	Path         string
	Name         string
	IsUnresolved bool
}

func (t *Target) RelativePathToPackageName() string {
	name := filepath.Dir(t.Path)
	name = filepath.ToSlash(name)
	if !strings.HasPrefix(name, "./") {
		name = "./" + name
	}
	return name
}

func (t *Target) TestNameToTestRunRegex() string {
	if t.IsUnresolved {
		return fmt.Sprintf("^%s", strings.TrimSuffix(t.Name, UnresolvedTestCaseName))
	}
	return fmt.Sprintf("^%s$", t.Name)
}
