package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(str string) []string {
	limit := 10

	freqMap := make(map[string]int)

	for _, word := range strings.Fields(str) {
		if normalized := normalizeWord(word); len(normalized) > 0 {
			freqMap[normalized]++
		}
	}

	list := make([]string, 0, len(freqMap))

	for word := range freqMap {
		list = append(list, word)
	}

	sort.Slice(list, func(i, j int) bool {
		diff := freqMap[list[i]] - freqMap[list[j]]
		switch {
		case diff < 0:
			return false
		case diff > 0:
			return true
		default:
			return list[i] < list[j]
		}
	})

	result := make([]string, 0, limit)
	for i := 0; i < limit && i < len(list); i++ {
		result = append(result, list[i])
	}
	return result
}

var normalizingRegexp = regexp.MustCompile(`^\p{P}?(.*?)\p{P}?$`)

func normalizeWord(s string) string {
	if matches := normalizingRegexp.FindStringSubmatch(strings.ToLower(s)); len(matches) > 1 {
		return matches[1]
	}
	return ""
}
