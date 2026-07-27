package ui

import (
	"math/rand"
	"testing"
)

func TestFuzzyMatchFilter_MatchedIndexes(t *testing.T) {
	tests := []struct {
		target string
		term   string
		want   []int
	}{
		{"abcdeあいうえおxyzわをん", "abc", []int{0, 1, 2}},
		{"abcdeあいうえおxyzわをん", "deあい", []int{3, 4, 5, 6}},
		{"abcdeあいうえおxyzわをん", "うえお", []int{7, 8, 9}},
		{"abcdeあいうえおxyzわをん", "xyz", []int{10, 11, 12}},
		{"abcdeあいうえおxyzわをん", "adz", []int{0, 3, 12}},
		{"abcdeあいうえおxyzわをん", "いうお", []int{6, 7, 9}},
		{"abcdeあいうえおxyzわをん", "eあyん", []int{4, 5, 11, 15}},
		{"abcdeあいうえおxyzわをん", "fgh", nil},
		{"abcdeあいうえおxyzわをん", "かきくけこ", nil},
		{"abcde", "", nil},
		{"axxbxxcxxabc", "abc", []int{9, 10, 11}},
		{"TestÄpfel", "äPF", []int{4, 5, 6}},
		{"TestFoo/日本語", "日本", []int{8, 9}},
	}

	for _, tt := range tests {
		t.Run(tt.term, func(t *testing.T) {
			targets := []string{tt.target}
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

func TestFuzzyMatchFilter_Ranking(t *testing.T) {
	tests := []struct {
		name    string
		term    string
		targets []string
		want    []int
	}{
		{
			name:    "prefers contiguous matches and then shorter spans",
			term:    "abc",
			targets: []string{"TestAxxBxxC", "TestABC", "TestAbxxC"},
			want:    []int{1, 2, 0},
		},
		{
			name:    "prefers matches at word boundaries when spans are equal",
			term:    "ab",
			targets: []string{"TestXaxb", "Test/A-b"},
			want:    []int{1, 0},
		},
		{
			name:    "prefers fewer gaps when spans are equal",
			term:    "abcd",
			targets: []string{"TestAbxcxd", "TestAbcxxd"},
			want:    []int{1, 0},
		},
		{
			name:    "prefers earlier matches after match quality",
			term:    "abc",
			targets: []string{"TestXXabcZ", "TestXabcZZ"},
			want:    []int{1, 0},
		},
		{
			name:    "prefers shorter targets after match quality",
			term:    "abc",
			targets: []string{"TestXabcZZ", "TestXabc"},
			want:    []int{1, 0},
		},
		{
			name:    "preserves input order for ties",
			term:    "abc",
			targets: []string{"TestXabc", "TestXabc"},
			want:    []int{0, 1},
		},
		{
			name:    "excludes targets that are not subsequence matches",
			term:    "abc",
			targets: []string{"TestAC", "TestABC", "TestBAC"},
			want:    []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ranks := fuzzyMatchFilter(tt.term, tt.targets)
			if len(ranks) != len(tt.want) {
				t.Fatalf("want %d ranks, got %d", len(tt.want), len(ranks))
			}
			for i, want := range tt.want {
				if ranks[i].Index != want {
					t.Errorf("rank %d: want target index %d, got %d", i, want, ranks[i].Index)
				}
			}
		})
	}
}

func TestFuzzyMatcherAgainstExhaustiveSearch(t *testing.T) {
	random := rand.New(rand.NewSource(1))
	alphabet := []rune("aAbB_/")
	for range 10_000 {
		target := make([]rune, random.Intn(8)+1)
		for i := range target {
			target[i] = alphabet[random.Intn(len(alphabet))]
		}
		pattern := make([]rune, random.Intn(4)+1)
		for i := range pattern {
			pattern[i] = alphabet[random.Intn(len(alphabet))]
		}

		gotQuality, gotIndexes, gotOK := findFuzzyMatch(pattern, target)
		wantQuality, wantOK := exhaustiveFuzzyMatch(pattern, target)
		if gotOK != wantOK || gotOK && compareFuzzyMatchQualities(gotQuality, wantQuality) != 0 {
			t.Fatalf(
				"pattern %q, target %q: got (%+v, %v), want (%+v, %v)",
				string(pattern), string(target), gotQuality, gotOK, wantQuality, wantOK,
			)
		}
		if gotOK {
			for i, index := range gotIndexes {
				if !equalFoldRune(pattern[i], target[index]) {
					t.Fatalf("pattern %q, target %q: invalid indexes %v", string(pattern), string(target), gotIndexes)
				}
				if i > 0 && gotIndexes[i-1] >= index {
					t.Fatalf("pattern %q, target %q: unordered indexes %v", string(pattern), string(target), gotIndexes)
				}
			}
		}
	}
}

func exhaustiveFuzzyMatch(pattern, target []rune) (fuzzyMatchQuality, bool) {
	indexes := make([]int, len(pattern))
	var best fuzzyMatchQuality
	found := false
	var search func(int, int)
	search = func(patternIndex, targetStart int) {
		if patternIndex == len(pattern) {
			quality := fuzzyMatchQuality{
				span:  indexes[len(indexes)-1] - indexes[0] + 1,
				start: indexes[0],
			}
			for i, index := range indexes {
				quality.boundaries += fuzzyBoundaryValue(target, index)
				if i > 0 && indexes[i-1]+1 != index {
					quality.gaps++
				}
			}
			if !found || compareFuzzyMatchQualities(quality, best) < 0 {
				best = quality
				found = true
			}
			return
		}
		for targetIndex := targetStart; targetIndex < len(target); targetIndex++ {
			if equalFoldRune(pattern[patternIndex], target[targetIndex]) {
				indexes[patternIndex] = targetIndex
				search(patternIndex+1, targetIndex+1)
			}
		}
	}
	search(0, 0)
	return best, found
}
