package ui

import "testing"

func TestExactMatchFilter_MatchedIndexes(t *testing.T) {
	tests := []struct {
		target string
		term   string
		want   []int
	}{
		{"abcdeあいうえおxyzわをん", "abc", []int{0, 1, 2}},
		{"abcdeあいうえおxyzわをん", "deあい", []int{3, 4, 5, 6}},
		{"abcdeあいうえおxyzわをん", "うえお", []int{7, 8, 9}},
		{"abcdeあいうえおxyzわをん", "xyz", []int{10, 11, 12}},
		{"abcdeあいうえおxyzわをん", "adz", nil},
		{"abcdeあいうえおxyzわをん", "いうお", nil},
		{"abcdeあいうえおxyzわをん", "eあyん", nil},
		{"abcdeあいうえおxyzわをん", "fgh", nil},
		{"abcdeあいうえおxyzわをん", "かきくけこ", nil},
		{"axxbxxcxxabc", "abc", []int{9, 10, 11}},
	}

	for _, tt := range tests {
		t.Run(tt.term, func(t *testing.T) {
			targets := []string{tt.target}
			ranks := exactMatchFilter(tt.term, targets)
			if tt.want == nil {
				if len(ranks) != 0 {
					t.Errorf("want no ranks, got %d", len(ranks))
					return
				}
			} else {
				if len(ranks) != 1 {
					t.Errorf("want 1 rank, got %d", len(ranks))
					return
				}
				if len(ranks[0].MatchedIndexes) != len(tt.want) {
					t.Errorf("want %d matched indexes, got %d", len(tt.want), len(ranks[0].MatchedIndexes))
					return
				}
				for i, idx := range tt.want {
					if ranks[0].MatchedIndexes[i] != idx {
						t.Errorf("want matched index %d at position %d, got %d", idx, i, ranks[0].MatchedIndexes[i])
					}
				}
			}
		})
	}
}
