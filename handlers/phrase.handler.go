package handlers

import "strings"

// SearchPhrase performs a phrase search and returns matching document IDs.
func (s *SearchEngine) SearchPhrase(query string) []int {
	query = strings.ReplaceAll(query, "\"", "")
	queryTokens := strings.Fields(query)

	var resultSet []int
	if len(queryTokens) == 0 {
		return nil
	}

	firstToken := queryTokens[0]
	if ids, ok := s.Index[firstToken]; ok {
		resultSet = append(resultSet, ids...)
	} else {
		return nil
	}

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
			// Remove docID from resultSet if phrase not found
			resultSet = removeElement(resultSet, docID)
		}
	}

	return resultSet
}

// removeElement removes an element from an int slice.
func removeElement(slice []int, elem int) []int {
	index := -1
	for i, v := range slice {
		if v == elem {
			index = i
			break
		}
	}
	if index != -1 {
		return append(slice[:index], slice[index+1:]...)
	}
	return slice
}
