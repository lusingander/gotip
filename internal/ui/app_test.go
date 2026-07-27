package ui

import (
	"testing"

	"github.com/lusingander/gotip/internal/theme"
)

func TestToggleMatchFilterRestoresConfiguredFuzzyMatcher(t *testing.T) {
	m := newModel(
		nil,
		nil,
		allView,
		exactMatchFilterType,
		legacyFuzzyMatchFilter,
		theme.DefaultColorTheme(),
	)

	m.toggleMatchFilter()

	ranks := m.allList.Filter("abc", []string{"axxbxxcxxabc"})
	if len(ranks) != 1 {
		t.Fatalf("want 1 rank, got %d", len(ranks))
	}
	want := []int{0, 3, 6}
	for i, wantIndex := range want {
		if got := ranks[0].MatchedIndexes[i]; got != wantIndex {
			t.Errorf("matched index %d = %d, want %d", i, got, wantIndex)
		}
	}
}
