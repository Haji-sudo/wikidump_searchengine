package handlers

import "strings"

// SearchPhrase performs a phrase search and returns matching document IDs.
// It takes a query string as input, removes any double quotes from the query,
// tokenizes the query into individual words, and then searches for the presence
// of the tokens in the indexed documents. It returns a slice of document IDs
// that match the entire phrase query.
func (s *SearchEngine) SearchPhrase(query string) []int {
	query = strings.ReplaceAll(query, "\"", "") // Remove double quotes from the query
	queryTokens := strings.Fields(query)        // Tokenize the query into individual words

	var resultSet []int // Initialize an empty slice to store the resulting document IDs
	if len(queryTokens) == 0 {
		return nil // Return nil if the query contains no tokens
	}

	firstToken := queryTokens[0] // Get the first token from the query
	if ids, ok := s.Index[firstToken]; ok {
		resultSet = append(resultSet, ids...) // Add the document IDs associated with the first token to the result set
	} else {
		return nil // Return nil if no document IDs are associated with the first token
	}

	// Iterate through the result set to check for the presence of the remaining tokens in the indexed documents
	for _, docID := range resultSet {
		doc := s.Documents[docID]
		docText := strings.ToLower(doc.Text)
		phraseFound := true
		for i := 1; i < len(queryTokens); i++ {
			if !strings.Contains(docText, queryTokens[i]) {
				phraseFound = false
				break
			}
		}
		if !phraseFound {
			// Remove docID from resultSet if the entire phrase is not found in the document
			resultSet = removeElement(resultSet, docID)
		}
	}

	return resultSet // Return the resulting document IDs that match the entire phrase query
}

// removeElement removes an element from an int slice.
// It takes a slice of integers and an integer element as input, and removes the
// first occurrence of the element from the slice. It returns the modified slice
// with the element removed, or the original slice if the element is not found.
func removeElement(slice []int, elem int) []int {
	index := -1
	for i, v := range slice {
		if v == elem {
			index = i
			break
		}
	}
	if index != -1 {
		return append(slice[:index], slice[index+1:]...) // Remove the element from the slice and return the modified slice
	}
	return slice // Return the original slice if the element is not found
}
