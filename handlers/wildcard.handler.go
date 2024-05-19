package handlers

import (
	"regexp"
	"strings"
)

// FindWildcardMatches finds matches for wildcard token in the Index.
// It replaces the wildcard character '*' with '.*' to create a wildcard pattern,
// compiles the pattern into a regular expression, and then iterates through the Index
// to find tokens that match the wildcard pattern.
// Parameters:
//
//	wildcardToken: the wildcard token to be matched in the Index
//
// Return:
//
//	[]int: a slice of integers representing the matches found in the Index
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
