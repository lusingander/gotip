package tip

import (
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
	Path            string
	PackageName     string
	TestNamePattern string
	IsPrefix        bool
}

func NewTarget(path, name string, isUnresolved bool) *Target {
	if isUnresolved {
		name = strings.TrimSuffix(name, UnresolvedTestCaseName)
	}
	return &Target{
		Path:            path,
		PackageName:     relativePathToPackageName(path),
		TestNamePattern: name,
		IsPrefix:        isUnresolved,
	}
}

func relativePathToPackageName(path string) string {
	name := filepath.Dir(path)
	name = filepath.ToSlash(name)
	if !strings.HasPrefix(name, "./") {
		name = "./" + name
	}
	return name
}

func (t *Target) DropLastSegment() {
	pattern := strings.TrimSuffix(t.TestNamePattern, "/")
	if lastSlash := strings.LastIndex(pattern, "/"); lastSlash != -1 {
		t.TestNamePattern = pattern[:lastSlash+1]
	} else {
		t.TestNamePattern = ""
	}
	t.IsPrefix = true
}
