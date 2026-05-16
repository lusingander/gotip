package listfmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/lusingander/gotip/internal/tip"
)

type document struct {
	Files []file `json:"files"`
}

type file struct {
	Path  string `json:"path"`
	Tests []test `json:"tests"`
}

type test struct {
	Name     string    `json:"name"`
	Subtests []subtest `json:"subtests"`
}

type subtest struct {
	Name     *string   `json:"name"`
	Resolved bool      `json:"resolved"`
	Subtests []subtest `json:"subtests"`
}

func WriteText(w io.Writer, tests map[string][]*tip.TestFunction) error {
	paths := sortedPaths(tests)
	for i, path := range paths {
		if i > 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "# %s\n", path); err != nil {
			return err
		}
		for _, tf := range tests[path] {
			if _, err := fmt.Fprintf(w, "- %s\n", tf.Name); err != nil {
				return err
			}
			if err := writeTextSubtests(w, tf.Subs, 1); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeTextSubtests(w io.Writer, subs []*tip.SubTest, depth int) error {
	for _, sub := range subs {
		name := sub.Name
		if !sub.Resolved {
			name = tip.UnresolvedTestCaseName + " [unresolved]"
		}
		if _, err := fmt.Fprintf(w, "%s- %s\n", strings.Repeat("  ", depth), name); err != nil {
			return err
		}
		if err := writeTextSubtests(w, sub.Subs, depth+1); err != nil {
			return err
		}
	}
	return nil
}

func WriteJSON(w io.Writer, tests map[string][]*tip.TestFunction) error {
	doc := newDocument(tests)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(doc)
}

func FormatText(tests map[string][]*tip.TestFunction) (string, error) {
	var buf bytes.Buffer
	if err := WriteText(&buf, tests); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func FormatJSON(tests map[string][]*tip.TestFunction) (string, error) {
	var buf bytes.Buffer
	if err := WriteJSON(&buf, tests); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func newDocument(tests map[string][]*tip.TestFunction) document {
	paths := sortedPaths(tests)
	files := make([]file, 0, len(paths))
	for _, path := range paths {
		functions := tests[path]
		outTests := make([]test, 0, len(functions))
		for _, tf := range functions {
			outTests = append(outTests, test{
				Name:     tf.Name,
				Subtests: newSubtests(tf.Subs),
			})
		}
		files = append(files, file{
			Path:  path,
			Tests: outTests,
		})
	}
	return document{Files: files}
}

func newSubtests(subs []*tip.SubTest) []subtest {
	out := make([]subtest, 0, len(subs))
	for _, sub := range subs {
		var name *string
		if sub.Resolved {
			value := sub.Name
			name = &value
		}
		out = append(out, subtest{
			Name:     name,
			Resolved: sub.Resolved,
			Subtests: newSubtests(sub.Subs),
		})
	}
	return out
}

func sortedPaths(tests map[string][]*tip.TestFunction) []string {
	paths := make([]string, 0, len(tests))
	for path := range tests {
		paths = append(paths, path)
	}
	slices.Sort(paths)
	return paths
}
