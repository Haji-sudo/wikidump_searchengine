package handlers

import (
	"regexp"
	"strings"
)

// FindWildcardMatches finds matches for wildcard token in the Index
func (s *SearchEngine) FindWildcardMatches(wildcardToken string) []int {
	var wildcardMatches []int
	wildcardPattern := strings.ReplaceAll(wildcardToken, "*", ".*")
	wildcardRegex := regexp.MustCompile("^" + wildcardPattern + "$")
	for token := range s.Index {
		if wildcardRegex.MatchString(token) {
			wildcardMatches = append(wildcardMatches, s.Index[token]...)
		}
	}
	return wildcardMatches
}
