package ui

import "testing"

func TestExactMatchFilter_MatchedIndexes(t *testing.T) {
	tests := []struct {
		name   string
		target string
		term   string
		want   []int
	}{
		{"contiguous ASCII at start", "abcdeあいうえおxyzわをん", "abc", []int{0, 1, 2}},
		{"contiguous match across ASCII and Unicode", "abcdeあいうえおxyzわをん", "deあい", []int{3, 4, 5, 6}},
		{"contiguous Unicode", "abcdeあいうえおxyzわをん", "うえお", []int{7, 8, 9}},
		{"ASCII match after Unicode", "abcdeあいうえおxyzわをん", "xyz", []int{10, 11, 12}},
		{"rejects sparse ASCII subsequence", "abcdeあいうえおxyzわをん", "adz", nil},
		{"rejects sparse Unicode subsequence", "abcdeあいうえおxyzわをん", "いうお", nil},
		{"rejects sparse mixed subsequence", "abcdeあいうえおxyzわをん", "eあyん", nil},
		{"missing ASCII pattern", "abcdeあいうえおxyzわをん", "fgh", nil},
		{"missing Unicode pattern", "abcdeあいうえおxyzわをん", "かきくけこ", nil},
		{"later contiguous substring", "axxbxxcxxabc", "abc", []int{9, 10, 11}},
		{"first repeated substring", "abc---abc", "abc", []int{0, 1, 2}},
		{"case-insensitive ASCII", "TestFoo", "foo", []int{4, 5, 6}},
		{"case-insensitive Unicode", "TestÄpfel", "äPF", []int{4, 5, 6}},
		{"pattern longer than target", "abc", "abcd", nil},
		{"no Unicode normalization", "TestCafe\u0301", "Café", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestExactMatchFilter_PreservesInputOrder(t *testing.T) {
	targets := []string{"zzabc", "abc", "xabc"}
	ranks := exactMatchFilter("abc", targets)
	if len(ranks) != len(targets) {
		t.Fatalf("want %d ranks, got %d", len(targets), len(ranks))
	}
	for i := range targets {
		if ranks[i].Index != i {
			t.Errorf("rank %d: want target index %d, got %d", i, i, ranks[i].Index)
		}
	}
}
