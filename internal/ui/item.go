package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/lusingander/gotip/internal/parse"
)

type testCaseItem struct {
	path string
	name string
}

var _ list.Item = (*testCaseItem)(nil)

func toTestCaseItems(tests map[string][]*parse.TestFunction) []list.Item {
	items := make([]list.Item, 0)
	for path, tfs := range tests {
		for _, tf := range tfs {
			if len(tf.Subs) == 0 {
				item := &testCaseItem{
					path: path,
					name: tf.Name,
				}
				items = append(items, item)
			} else {
				items = append(items, toTestCaseItemsFromSubTests(tf.Subs, path, tf.Name)...)
			}
		}
	}
	return items
}

func toTestCaseItemsFromSubTests(ss []*parse.SubTest, path, base string) []list.Item {
	items := make([]list.Item, 0)
	for _, s := range ss {
		name := base + "/" + s.Name
		if len(s.Subs) == 0 {
			item := &testCaseItem{
				path: path,
				name: name,
			}
			items = append(items, item)
		} else {
			items = append(items, toTestCaseItemsFromSubTests(s.Subs, path, name)...)
		}
	}
	return items
}

func (i *testCaseItem) Title() string {
	return i.name
}

func (i *testCaseItem) Description() string {
	return i.path
}

func (i *testCaseItem) FilterValue() string {
	return i.name
}
