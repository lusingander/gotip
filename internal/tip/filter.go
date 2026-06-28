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

func FilterHistoriesByPackages(histories *Histories, packages []string) *Histories {
	if len(packages) == 0 {
		return histories
	}

	packageSet := make(map[string]struct{}, len(packages))
	for _, pkg := range packages {
		packageSet[normalizePackageName(pkg)] = struct{}{}
	}

	filtered := &Histories{
		ProjectDir: histories.ProjectDir,
		Histories:  make([]*History, 0, len(histories.Histories)),
	}
	for _, history := range histories.Histories {
		if _, ok := packageSet[historyPackageName(history)]; ok {
			filtered.Histories = append(filtered.Histories, history)
		}
	}
	return filtered
}

func historyPackageName(history *History) string {
	if history.PackageName != "" {
		return normalizePackageName(history.PackageName)
	}
	return relativePathToPackageName(history.Path)
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
