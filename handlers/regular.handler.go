package handlers

import (
	"math"
	"sort"
)

func (s *SearchEngine) Search(text string) []int {
	// Tokenize the search query
	queryTokens := analyze(text)
	if len(queryTokens) == 0 {
		return []int{}
	}
	// Initialize the result set with document IDs from the first token
	var resultSet []int
	if ids, ok := s.Index[queryTokens[0]]; ok {
		resultSet = append(resultSet, ids...)
	} else {
		return nil
	}
	// Calculate intersection of document IDs for subsequent tokens
	for _, token := range queryTokens[1:] {
		if ids, ok := s.Index[token]; ok {
			resultSet = Intersection(resultSet, ids)
		} else {
			// Token doesn't exist in Index.
			return nil
		}
	}
	// Calculate TF-IDF score for each document in the result set
	docScores := make(map[int]float64)
	for _, docID := range resultSet {
		// Calculate TF-IDF score for document text
		for _, token := range queryTokens {
			idf := math.Log(float64(len(s.Documents)) / float64(len(s.Index[token])))
			tf := calculateTF(token, s.Documents[docID].Text)
			tfidf := tf * idf
			docScores[docID] += tfidf
			// Check if token is also in document title
			titleTokens := analyze(s.Documents[docID].Title)
			for _, titleToken := range titleTokens {
				if token == titleToken {
					// Boost score if token is in title
					docScores[docID] += tfidf * 0.5 // Adjust the boosting factor as needed
				}
			}
		}
	}
	// Convert map to slice for sorting
	var results []struct {
		DocID int
		Score float64
	}
	for docID, score := range docScores {
		results = append(results, struct {
			DocID int
			Score float64
		}{docID, score})
	}

	// Sort results by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Extract document IDs from sorted results
	var rankedDocs []int
	for _, result := range results {
		rankedDocs = append(rankedDocs, result.DocID)
	}

	return rankedDocs
}

func calculateTF(term string, documentText string) float64 {
	tokens := analyze(documentText)
	termCount := 0
	for _, token := range tokens {
		if token == term {
			termCount++
		}
	}
	return float64(termCount) / float64(len(tokens))
}
