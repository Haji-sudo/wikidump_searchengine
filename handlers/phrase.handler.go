package handlers

// SearchPhrase performs a phrase search and returns matching document IDs.
// It takes a query string as input, removes any double quotes from the query,
// tokenizes the query into individual words, and then searches for the presence
// of the tokens in the indexed documents. It returns a slice of document IDs
// that match the entire phrase query.
func (s *SearchEngine) SearchPhrase(query string) []int {
	queryTokens := analyze(query) // Use the analyze function to process the query

	if len(queryTokens) == 0 {
		return nil // Return nil if the query contains no tokens
	}

	// Get the list of documents containing the first token
	firstToken := queryTokens[0]
	resultSet, ok := s.Index[firstToken]
	if !ok {
		return nil // Return nil if no document IDs are associated with the first token
	}

	// Iterate through the result set to check for the presence of the phrase in the documents
	finalResults := []int{}
	for _, docID := range resultSet {
		doc := s.Documents[docID]
		docTokens := analyze(doc.Text) // Use the analyze function to process the document text

		if containsExactPhrase(docTokens, queryTokens) {
			finalResults = append(finalResults, docID)
		}
	}

	return finalResults // Return the resulting document IDs that match the entire phrase query
}

// containsPhrase checks if the document tokens contain the query tokens in the same order
func containsExactPhrase(docTokens, queryTokens []string) bool {
	for i := 0; i <= len(docTokens)-len(queryTokens); i++ {
		match := true
		for j := 0; j < len(queryTokens); j++ {
			if docTokens[i+j] != queryTokens[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
