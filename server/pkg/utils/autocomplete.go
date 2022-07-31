package utils

import (
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// Autocomplete tries to find a string in a list and returns
// the sorted and ranked matches
func Autocomplete(value string, list []string) []string {
	value = strings.TrimSpace(strings.ToLower(value))

	// compute matches
	matches := fuzzy.RankFindNormalizedFold(value, list)
	sort.Sort(matches)

	// split between those who match the prefix
	prefixed := []string{}
	rest := []string{}
	for _, match := range matches {
		matchString := strings.TrimSpace(strings.ToLower(match.Target))
		if strings.HasPrefix(matchString, value) {
			prefixed = append(prefixed, match.Target)
		} else {
			rest = append(rest, match.Target)
		}
	}

	// return the prefixed first then the rest
	return append(prefixed, rest...)
}
