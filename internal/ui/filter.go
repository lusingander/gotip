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

func matchFilterTypeFromStr(s string) matchFilterType {
	switch s {
	case "fuzzy":
		return fuzzyMatchFilterType
	case "exact":
		return exactMatchFilterType
	default:
		panic("unknown match filter type: " + s)
	}
}

func fuzzyMatchFilter(term string, targets []string) []list.Rank {
	ranks := list.DefaultFilter(term, targets)
	return convertRanks(ranks, targets)
}

func exactMatchFilter(term string, targets []string) []list.Rank {
	ranks := make([]list.Rank, 0, len(targets))
	termLower := strings.ToLower(term)
	for i, target := range targets {
		targetLower := strings.ToLower(target)
		if idx := strings.Index(targetLower, termLower); idx != -1 {
			matchedIndexes := make([]int, 0)
			for j := range termLower {
				matchedIndexes = append(matchedIndexes, idx+j)
			}
			ranks = append(ranks, list.Rank{Index: i, MatchedIndexes: matchedIndexes})
		}
	}
	return convertRanks(ranks, targets)
}

func convertRanks(ranks []list.Rank, targets []string) []list.Rank {
	ret := make([]list.Rank, len(ranks))
	for i, rank := range ranks {
		target := targets[rank.Index]
		ret[i] = list.Rank{
			Index:          rank.Index,
			MatchedIndexes: byteOffsetsToRuneIndices(target, rank.MatchedIndexes),
		}
	}
	return ret
}

func byteOffsetsToRuneIndices(s string, offsets []int) []int {
	m := make(map[int]int)
	byteOffset := 0
	runeIndex := 0
	for _, r := range s {
		m[byteOffset] = runeIndex
		byteOffset += len(string(r))
		runeIndex++
	}
	runeIndices := make([]int, 0, len(offsets))
	for _, offset := range offsets {
		runeIndices = append(runeIndices, m[offset])
	}
	return runeIndices
}
