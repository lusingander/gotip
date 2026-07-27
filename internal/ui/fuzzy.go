package ui

import (
	"cmp"
	"slices"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
)

func fuzzyMatchFilterFromStr(s string) list.FilterFunc {
	switch s {
	case "gotip":
		return gotipFuzzyMatchFilter
	case "legacy":
		return legacyFuzzyMatchFilter
	default:
		panic("unknown fuzzy matcher: " + s)
	}
}

func legacyFuzzyMatchFilter(term string, targets []string) []list.Rank {
	ranks := list.DefaultFilter(term, targets)
	return convertRanks(ranks, targets)
}

// gotipFuzzyMatchFilter implements fuzzy matching tailored to Go test names.
//
// The legacy sahilm/fuzzy matcher may commit an early, scattered alignment
// while scanning even when a later contiguous match exists. Its ranking also
// mixes byte-based length penalties with rapidly increasing adjacency bonuses,
// which can make Unicode targets and longer queries rank unexpectedly.
//
// This matcher instead works in runes and uses dynamic programming to select
// the globally best case-insensitive subsequence alignment for each target.
// Candidates are compared lexicographically by shorter match span, more word
// boundary matches, fewer gaps, earlier start, and shorter target length.
// Exact priority rules keep ranking predictable without balancing unrelated
// criteria through a single weighted score; complete ties retain input order.
func gotipFuzzyMatchFilter(term string, targets []string) []list.Rank {
	pattern := []rune(term)
	if len(pattern) == 0 {
		return nil
	}

	matches := make([]rankedFuzzyMatch, 0, len(targets))
	for i, target := range targets {
		targetRunes := []rune(target)
		quality, matchedIndexes, ok := findFuzzyMatch(pattern, targetRunes)
		if !ok {
			continue
		}
		matches = append(matches, rankedFuzzyMatch{
			index:          i,
			matchedIndexes: matchedIndexes,
			quality:        quality,
			targetLength:   len(targetRunes),
		})
	}

	slices.SortStableFunc(matches, compareRankedFuzzyMatches)

	ranks := make([]list.Rank, len(matches))
	for i, match := range matches {
		ranks[i] = list.Rank{
			Index:          match.index,
			MatchedIndexes: match.matchedIndexes,
		}
	}
	return ranks
}

type rankedFuzzyMatch struct {
	index          int
	matchedIndexes []int
	quality        fuzzyMatchQuality
	targetLength   int
}

// fuzzyMatchQuality fields are compared in declaration order, except that more
// boundary matches are better. A contiguous match has the smallest possible
// span, so it does not need a separate priority flag.
type fuzzyMatchQuality struct {
	span       int
	boundaries int
	gaps       int
	start      int
}

type fuzzyMatchState struct {
	valid      bool
	start      int
	boundaries int
	gaps       int
	previous   int
}

// findFuzzyMatch uses dynamic programming to find the globally best alignment
// in O(len(pattern) * len(target)) time. States retain predecessor indexes so
// the selected runes can be highlighted by the list UI.
func findFuzzyMatch(pattern, target []rune) (fuzzyMatchQuality, []int, bool) {
	if len(pattern) == 0 || len(pattern) > len(target) {
		return fuzzyMatchQuality{}, nil, false
	}

	targetLength := len(target)
	states := make([]fuzzyMatchState, len(pattern)*targetLength)
	for targetIndex, targetRune := range target {
		if equalFoldRune(pattern[0], targetRune) {
			states[targetIndex] = fuzzyMatchState{
				valid:      true,
				start:      targetIndex,
				boundaries: fuzzyBoundaryValue(target, targetIndex),
				previous:   -1,
			}
		}
	}

	for patternIndex := 1; patternIndex < len(pattern); patternIndex++ {
		previousRow := (patternIndex - 1) * targetLength
		currentRow := patternIndex * targetLength
		bestNonAdjacent := -1

		for targetIndex, targetRune := range target {
			nonAdjacentIndex := targetIndex - 2
			if nonAdjacentIndex >= 0 {
				candidate := states[previousRow+nonAdjacentIndex]
				if candidate.valid && (bestNonAdjacent == -1 ||
					betterFuzzyPredecessor(
						candidate,
						nonAdjacentIndex,
						states[previousRow+bestNonAdjacent],
						bestNonAdjacent,
					)) {
					bestNonAdjacent = nonAdjacentIndex
				}
			}

			if !equalFoldRune(pattern[patternIndex], targetRune) {
				continue
			}

			boundary := fuzzyBoundaryValue(target, targetIndex)
			var best fuzzyMatchState

			if targetIndex > 0 {
				previous := states[previousRow+targetIndex-1]
				if previous.valid {
					best = fuzzyMatchState{
						valid:      true,
						start:      previous.start,
						boundaries: previous.boundaries + boundary,
						gaps:       previous.gaps,
						previous:   targetIndex - 1,
					}
				}
			}

			if bestNonAdjacent != -1 {
				previous := states[previousRow+bestNonAdjacent]
				candidate := fuzzyMatchState{
					valid:      true,
					start:      previous.start,
					boundaries: previous.boundaries + boundary,
					gaps:       previous.gaps + 1,
					previous:   bestNonAdjacent,
				}
				if !best.valid || betterFuzzyStateAtSameEnd(candidate, best) {
					best = candidate
				}
			}

			states[currentRow+targetIndex] = best
		}
	}

	lastRow := (len(pattern) - 1) * targetLength
	bestEnd := -1
	var bestQuality fuzzyMatchQuality
	for targetIndex := range target {
		state := states[lastRow+targetIndex]
		if !state.valid {
			continue
		}
		quality := fuzzyMatchQuality{
			span:       targetIndex - state.start + 1,
			boundaries: state.boundaries,
			gaps:       state.gaps,
			start:      state.start,
		}
		if bestEnd == -1 || compareFuzzyMatchQualities(quality, bestQuality) < 0 {
			bestEnd = targetIndex
			bestQuality = quality
		}
	}
	if bestEnd == -1 {
		return fuzzyMatchQuality{}, nil, false
	}

	matchedIndexes := make([]int, len(pattern))
	targetIndex := bestEnd
	for patternIndex := len(pattern) - 1; patternIndex >= 0; patternIndex-- {
		matchedIndexes[patternIndex] = targetIndex
		targetIndex = states[patternIndex*targetLength+targetIndex].previous
	}
	return bestQuality, matchedIndexes, true
}

func compareRankedFuzzyMatches(a, b rankedFuzzyMatch) int {
	if c := compareFuzzyMatchQualities(a.quality, b.quality); c != 0 {
		return c
	}
	return cmp.Compare(a.targetLength, b.targetLength)
}

func compareFuzzyMatchQualities(a, b fuzzyMatchQuality) int {
	if c := cmp.Compare(a.span, b.span); c != 0 {
		return c
	}
	if c := cmp.Compare(b.boundaries, a.boundaries); c != 0 {
		return c
	}
	if c := cmp.Compare(a.gaps, b.gaps); c != 0 {
		return c
	}
	return cmp.Compare(a.start, b.start)
}

func betterFuzzyStateAtSameEnd(a, b fuzzyMatchState) bool {
	if a.start != b.start {
		return a.start > b.start
	}
	if a.boundaries != b.boundaries {
		return a.boundaries > b.boundaries
	}
	if a.gaps != b.gaps {
		return a.gaps < b.gaps
	}
	return a.previous > b.previous
}

func betterFuzzyPredecessor(a fuzzyMatchState, aEnd int, b fuzzyMatchState, bEnd int) bool {
	if a.start != b.start {
		return a.start > b.start
	}
	if a.boundaries != b.boundaries {
		return a.boundaries > b.boundaries
	}
	if a.gaps != b.gaps {
		return a.gaps < b.gaps
	}
	return aEnd > bEnd
}

func fuzzyBoundaryValue(target []rune, index int) int {
	if index == 0 {
		return 1
	}
	previous := target[index-1]
	if isFuzzySeparator(previous) ||
		unicode.IsLower(previous) && unicode.IsUpper(target[index]) {
		return 1
	}
	return 0
}

func isFuzzySeparator(r rune) bool {
	return strings.ContainsRune("/-_ .\\", r)
}

func equalFoldRune(a, b rune) bool {
	if a == b {
		return true
	}
	for folded := unicode.SimpleFold(a); folded != a; folded = unicode.SimpleFold(folded) {
		if folded == b {
			return true
		}
	}
	return false
}
