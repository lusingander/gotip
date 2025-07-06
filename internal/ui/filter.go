package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

type matchFilterType int

const (
	fuzzyMatchFilterType matchFilterType = iota
	exactMatchFilterType
)

func fuzzyMatchFilter(term string, targets []string) []list.Rank {
	// todo: consider multi-byte characters
	return list.DefaultFilter(term, targets)
}

func exactMatchFilter(term string, targets []string) []list.Rank {
	// todo: consider multi-byte characters
	ranks := make([]list.Rank, 0, len(targets))
	termLower := strings.ToLower(term)
	for i, target := range targets {
		targetLower := strings.ToLower(target)
		if idx := strings.Index(targetLower, termLower); idx != -1 {
			matchedIndexes := make([]int, 0)
			for j := range len(termLower) {
				matchedIndexes = append(matchedIndexes, idx+j)
			}
			ranks = append(ranks, list.Rank{Index: i, MatchedIndexes: matchedIndexes})
		}
	}
	return ranks
}
