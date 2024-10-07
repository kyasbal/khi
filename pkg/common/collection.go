package common

import (
	"slices"
	"sort"
	"strings"
)

func DedupeStringArray(arr []string) []string {
	arrByMap := map[string]struct{}{}
	for _, v := range arr {
		arrByMap[v] = struct{}{}
	}
	deduped := []string{}
	for v := range arrByMap {
		deduped = append(deduped, v)
	}
	slices.SortFunc(deduped, func(a, b string) int { return strings.Compare(a, b) })
	return deduped
}

func levenshteinDistance(s, t string) int {
	m := len(s)
	n := len(t)

	d := make([][]int, m+1)
	for i := range d {
		d[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		d[i][0] = i
	}
	for j := 1; j <= n; j++ {
		d[0][j] = j
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			m := min(d[i-1][j]+1, d[i][j-1]+1)
			if s[i-1] == t[j-1] {
				d[i][j] = min(m, d[i-1][j-1])
			} else {
				d[i][j] = min(m, d[i-1][j-1]+1)
			}
		}
	}

	return d[m][n]
}

func SortForAutocomplete(input string, elements []string) []string {
	// Sort criteria
	// Priority group 1: completely matches the name
	// Priority group 2: has the query prefix and levenshtein distance
	// Priority 3: levenshtein distance
	result := []string{}
	prefixMatches := []string{}
	prefixNonMatches := []string{}
	for _, element := range elements {
		if element == input {
			result = append(result, element)
			continue
		}
		if strings.HasPrefix(element, input) {
			prefixMatches = append(prefixMatches, element)
		} else {
			prefixNonMatches = append(prefixNonMatches, element)
		}
	}

	distances := map[string]int{}
	for _, s := range elements {
		distances[s] = levenshteinDistance(input, s)
	}
	sort.Slice(prefixMatches, func(i, j int) bool {
		if distances[prefixMatches[i]] == distances[prefixMatches[j]] {
			return strings.Compare(prefixMatches[i], prefixMatches[j]) > 0
		}
		return distances[prefixMatches[i]] < distances[prefixMatches[j]]
	})
	sort.Slice(prefixNonMatches, func(i, j int) bool {
		if distances[prefixNonMatches[i]] == distances[prefixNonMatches[j]] {
			return strings.Compare(prefixNonMatches[i], prefixNonMatches[j]) > 0
		}
		return distances[prefixNonMatches[i]] < distances[prefixNonMatches[j]]
	})

	return append(append(result, prefixMatches...), prefixNonMatches...)
}
