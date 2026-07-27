package ui

import (
	"math/rand"
	"testing"
)

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
