package tip

import "strings"

func FilterTestsByPackages(tests map[string][]*TestFunction, packages []string) map[string][]*TestFunction {
	if len(packages) == 0 {
		return tests
	}

	matcher := newPackageMatcher(packages)

	filtered := make(map[string][]*TestFunction)
	for path, testFunctions := range tests {
		if matcher.match(relativePathToPackageName(path)) {
			filtered[path] = testFunctions
		}
	}
	return filtered
}

func FilterHistoriesByPackages(histories *Histories, packages []string) *Histories {
	if len(packages) == 0 {
		return histories
	}

	matcher := newPackageMatcher(packages)

	filtered := &Histories{
		ProjectDir: histories.ProjectDir,
		Histories:  make([]*History, 0, len(histories.Histories)),
	}
	for _, history := range histories.Histories {
		if matcher.match(historyPackageName(history)) {
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

type packageMatcher struct {
	patterns []string
}

func newPackageMatcher(packages []string) packageMatcher {
	patterns := make([]string, 0, len(packages))
	for _, pkg := range packages {
		patterns = append(patterns, normalizePackagePattern(pkg))
	}
	return packageMatcher{patterns: patterns}
}

func (m packageMatcher) match(pkg string) bool {
	pkg = normalizePackageName(pkg)
	for _, pattern := range m.patterns {
		if pattern == "./..." {
			return true
		}

		if before, ok := strings.CutSuffix(pattern, "/..."); ok {
			base := before
			if packagePathContains(base, pkg) {
				return true
			}
			continue
		}

		if packagePathContains(pattern, pkg) {
			return true
		}
	}
	return false
}

func packagePathContains(base, pkg string) bool {
	return pkg == base || strings.HasPrefix(pkg, base+"/")
}

func normalizePackagePattern(name string) string {
	name = strings.TrimSuffix(name, "/")
	if before, ok := strings.CutSuffix(name, "/..."); ok {
		base := normalizePackageName(before)
		if base == "." {
			return "./..."
		}
		return base + "/..."
	}
	return normalizePackageName(name)
}

func normalizePackageName(name string) string {
	name = strings.TrimSuffix(name, "/")
	if name == "" || name == "." || name == "./." {
		return "."
	}
	if !strings.HasPrefix(name, "./") {
		return "./" + name
	}
	return name
}
