package tip

import "strings"

func FilterTestsByPackages(tests map[string][]*TestFunction, packages []string) map[string][]*TestFunction {
	if len(packages) == 0 {
		return tests
	}

	packageSet := make(map[string]struct{}, len(packages))
	for _, pkg := range packages {
		packageSet[normalizePackageName(pkg)] = struct{}{}
	}

	filtered := make(map[string][]*TestFunction)
	for path, testFunctions := range tests {
		if _, ok := packageSet[relativePathToPackageName(path)]; ok {
			filtered[path] = testFunctions
		}
	}
	return filtered
}

func normalizePackageName(name string) string {
	name = strings.TrimSuffix(name, "/")
	if name == "" || name == "." {
		return "."
	}
	if !strings.HasPrefix(name, "./") {
		return "./" + name
	}
	return name
}
