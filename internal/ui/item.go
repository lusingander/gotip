package ui

import (
	"cmp"
	"slices"

	"github.com/charmbracelet/bubbles/list"
	"github.com/lusingander/gotip/internal/tip"
)

type testCaseItem struct {
	path         string
	name         string
	isUnresolved bool
}

var _ list.Item = (*testCaseItem)(nil)

func toTestCaseItems(tests map[string][]*tip.TestFunction) []list.Item {
	items := make([]list.Item, 0)
	for path, tfs := range tests {
		for _, tf := range tfs {
			if len(tf.Subs) == 0 {
				item := &testCaseItem{
					path:         path,
					name:         tf.Name,
					isUnresolved: false,
				}
				items = append(items, item)
			} else {
				items = append(items, toTestCaseItemsFromSubTests(tf.Subs, path, tf.Name)...)
			}
		}
	}
	slices.SortStableFunc(items, func(a, b list.Item) int {
		return cmp.Compare(a.(*testCaseItem).path, b.(*testCaseItem).path)
	})
	return items
}

func toTestCaseItemsFromSubTests(ss []*tip.SubTest, path, base string) []list.Item {
	items := make([]list.Item, 0)
	for _, s := range ss {
		name := base + "/" + s.Name
		if len(s.Subs) == 0 {
			item := &testCaseItem{
				path:         path,
				name:         name,
				isUnresolved: s.IsUnresolvedSubTests,
			}
			items = append(items, item)
		} else {
			items = append(items, toTestCaseItemsFromSubTests(s.Subs, path, name)...)
		}
	}
	return items
}

func (i *testCaseItem) FilterValue() string {
	return i.name
}
