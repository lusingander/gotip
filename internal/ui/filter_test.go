package ui

import "testing"

func TestFuzzyMatchFilter_MatchedIndexes(t *testing.T) {
	target := "abcdeあいうえおxyzわをん"
	tests := []struct {
		term string
		want []int
	}{
		{"abc", []int{0, 1, 2}},
		{"deあい", []int{3, 4, 5, 6}},
		{"うえお", []int{7, 8, 9}},
		{"xyz", []int{10, 11, 12}},
		{"adz", []int{0, 3, 12}},
		{"いうお", []int{6, 7, 9}},
		{"eあyん", []int{4, 5, 11, 15}},
		{"fgh", nil},
		{"かきくけこ", nil},
	}

	for _, tt := range tests {
		t.Run(tt.term, func(t *testing.T) {
			targets := []string{target}
			ranks := fuzzyMatchFilter(tt.term, targets)
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

func TestExactMatchFilter_MatchedIndexes(t *testing.T) {
	target := "abcdeあいうえおxyzわをん"
	tests := []struct {
		term string
		want []int
	}{
		{"abc", []int{0, 1, 2}},
		{"deあい", []int{3, 4, 5, 6}},
		{"うえお", []int{7, 8, 9}},
		{"xyz", []int{10, 11, 12}},
		{"adz", nil},
		{"いうお", nil},
		{"eあyん", nil},
		{"fgh", nil},
		{"かきくけこ", nil},
	}

	for _, tt := range tests {
		t.Run(tt.term, func(t *testing.T) {
			targets := []string{target}
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
