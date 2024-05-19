package handlers

import (
	"compress/gzip"
	"encoding/xml"
	"os"
	"strings"
	"unicode"

	snowballeng "github.com/kljensen/snowball/english"
)

type Document struct {
	Title string `xml:"title"`
	URL   string `xml:"url"`
	Text  string `xml:"abstract"`
	ID    int
}

type SearchEngine struct {
	Documents []Document
	Index     map[string][]int
}

// NewSearchEngine creates a new SearchEngine instance and initializes it with the documents loaded from the specified path.
// It returns a pointer to the newly created SearchEngine.
func NewSearchEngine(path string) *SearchEngine {
	s := &SearchEngine{
		Index: make(map[string][]int), // Initialize the Index map.
	}
	err := s.LoadDocuments(path) // Load documents from the specified path.
	if err != nil {
		panic(err) // Panic if an error occurs while loading documents.
	}
	s.IndexDoc() // Index the loaded documents.
	return s     // Return the pointer to the newly created SearchEngine.
}

// LoadDocuments loads and processes the documents from the specified path.
// It opens the file, creates a gzip reader, and decodes the XML data into a temporary slice of documents.
// Then, it converts the temporary documents to actual documents and assigns IDs to them.
// Parameters:
//
//	path: a string representing the path to the file containing the documents in gzip-compressed XML format.
//
// Return values:
//
//	error: an error if any occurred during the process, or nil if the operation was successful.
func (s *SearchEngine) LoadDocuments(path string) error {
	// Open the file
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a gzip reader
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	// Create an XML decoder
	dec := xml.NewDecoder(gz)

	// Define a temporary struct for decoding
	type tempDocument struct {
		Title    string `xml:"title"`
		URL      string `xml:"url"`
		Abstract string `xml:"abstract"`
	}

	// Decode XML data into temporary slice of documents
	dump := struct {
		Documents []tempDocument `xml:"doc"`
	}{}
	if err := dec.Decode(&dump); err != nil {
		return err
	}

	// Convert temporary documents to actual documents
	var documents []Document
	for i, tempDoc := range dump.Documents {
		document := Document{
			ID:    i,
			Title: tempDoc.Title,
			URL:   tempDoc.URL,
			Text:  tempDoc.Abstract,
		}
		documents = append(documents, document)
	}

	// Assign IDs to documents
	s.Documents = documents

	return nil
}

// IndexDoc indexes the documents in the SearchEngine by tokenizing and adding them to the Index map.
func (s *SearchEngine) IndexDoc() {
	for _, doc := range s.Documents {
		for _, token := range analyze(doc.Text) {
			ids := s.Index[token]
			if ids != nil && ids[len(ids)-1] == doc.ID {
				// Token already exists in Index.
				continue
			}
			s.Index[token] = append(ids, doc.ID)
		}
	}
}

// Intersection returns the intersection of two slices.
// It takes two integer slices a and b as input and returns a new slice containing the common elements between the two input slices.
// Parameters:
//
//	a: an integer slice representing the first input slice.
//	b: an integer slice representing the second input slice.
//
// Return values:
//
//	[]int: a new slice containing the common elements between the input slices a and b.
func Intersection(a []int, b []int) []int {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	r := make([]int, 0, maxLen)
	var i, j int
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		} else {
			r = append(r, a[i])
			i++
			j++
		}
	}
	return r
}

// tokenize returns a slice of tokens by splitting the input text based on non-letter and non-number characters.
// It takes a string representing the input text and returns a slice of strings representing the tokens obtained after splitting the text.
// Parameters:
//
//	text: a string representing the input text to be tokenized.
//
// Return values:
//
//	[]string: a slice of strings representing the tokens obtained after splitting the input text.
func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// Split on any character that is not a letter or a number.
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// analyze returns a slice of tokens after processing the input text through tokenization, lowercase filtering, stop word filtering, and stemming.
// It takes a string representing the input text and returns a slice of strings representing the processed tokens.
// Parameters:
//
//	text: a string representing the input text to be analyzed.
//
// Return values:
//
//	[]string: a slice of strings representing the processed tokens obtained after tokenization, lowercase filtering, stop word filtering, and stemming.
func analyze(text string) []string {
	tokens := tokenize(text)         // Tokenize the input text.
	tokens = lowercaseFilter(tokens) // Apply lowercase filtering to the tokens.
	tokens = stopWordFilter(tokens)  // Apply stop word filtering to the tokens.
	tokens = stemmerFilter(tokens)   // Apply stemming to the tokens.
	return tokens                    // Return the processed tokens.
}

// lowercaseFilter applies lowercase filtering to the input tokens and returns a new slice of strings with all tokens converted to lowercase.
// It takes a slice of strings representing the input tokens and returns a new slice of strings with all tokens converted to lowercase.
// Parameters:
//
//	tokens: a slice of strings representing the input tokens to be converted to lowercase.
//
// Return values:
//
//	[]string: a new slice of strings containing the input tokens converted to lowercase.
func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}

// stopWordFilter returns a slice of tokens with stop words removed.
// It takes a slice of strings representing the input tokens and removes any stop words from the input.
// Stop words are common words that are often filtered out from text data because they do not carry significant meaning or are too common to be useful for searching or indexing.
// Parameters:
//   tokens: a slice of strings representing the input tokens from which stop words are to be removed.
// Return values:
//   []string: a new slice of strings containing the input tokens with stop words removed.

func stopWordFilter(tokens []string) []string {
	stopWords := map[string]struct{}{
		"a": {}, "about": {}, "above": {}, "after": {}, "again": {},
		"against": {}, "all": {}, "am": {}, "an": {}, "and": {},
		"any": {}, "are": {}, "aren't": {}, "as": {}, "at": {},
		"be": {}, "because": {}, "been": {}, "before": {}, "being": {},
		"below": {}, "between": {}, "both": {}, "but": {}, "by": {},
		"can't": {}, "cannot": {}, "could": {}, "couldn't": {}, "did": {},
		"didn't": {}, "do": {}, "does": {}, "doesn't": {}, "doing": {},
		"don't": {}, "down": {}, "during": {}, "each": {}, "few": {},
		"for": {}, "from": {}, "further": {}, "had": {}, "hadn't": {},
		"has": {}, "hasn't": {}, "have": {}, "haven't": {}, "having": {},
		"he": {}, "he'd": {}, "he'll": {}, "he's": {}, "her": {},
		"here": {}, "here's": {}, "hers": {}, "herself": {}, "him": {},
		"himself": {}, "his": {}, "how": {}, "how's": {}, "i": {},
		"i'd": {}, "i'll": {}, "i'm": {}, "i've": {}, "if": {},
		"in": {}, "into": {}, "is": {}, "isn't": {}, "it": {},
		"it's": {}, "its": {}, "itself": {}, "let's": {}, "me": {},
		"more": {}, "most": {}, "mustn't": {}, "my": {}, "myself": {},
		"no": {}, "nor": {}, "not": {}, "of": {}, "off": {},
		"on": {}, "once": {}, "only": {}, "or": {}, "other": {},
		"ought": {}, "our": {}, "ours": {}, "ourselves": {}, "out": {},
		"over": {}, "own": {}, "same": {}, "shan't": {}, "she": {},
		"she'd": {}, "she'll": {}, "she's": {}, "should": {}, "shouldn't": {},
		"so": {}, "some": {}, "such": {}, "than": {}, "that": {},
		"that's": {}, "the": {}, "their": {}, "theirs": {}, "them": {},
		"themselves": {}, "then": {}, "there": {}, "there's": {}, "these": {},
		"they": {}, "they'd": {}, "they'll": {}, "they're": {}, "they've": {},
		"this": {}, "those": {}, "through": {}, "to": {}, "too": {},
		"under": {}, "until": {}, "up": {}, "very": {}, "was": {},
		"wasn't": {}, "we": {}, "we'd": {}, "we'll": {}, "we're": {},
		"we've": {}, "were": {}, "weren't": {}, "what": {}, "what's": {},
		"when": {}, "when's": {}, "where": {}, "where's": {}, "which": {},
		"while": {}, "who": {}, "who's": {}, "whom": {}, "why": {},
		"why's": {}, "with": {}, "won't": {}, "would": {}, "wouldn't": {},
		"you": {}, "you'd": {}, "you'll": {}, "you're": {}, "you've": {},
		"your": {}, "yours": {}, "yourself": {}, "yourselves": {},
	}

	var filteredTokens []string
	for _, token := range tokens {
		// Check if the token is not a stop word
		if _, ok := stopWords[token]; !ok {
			filteredTokens = append(filteredTokens, token)
		}
	}
	return filteredTokens
}

// stemmerFilter returns a slice of stemmed tokens after applying the Snowball stemming algorithm.
// It takes a slice of strings representing the input tokens and returns a new slice of strings containing the stemmed tokens.
// Snowball stemming algorithm is used to reduce words to their root form, which helps in improving search and indexing capabilities by treating different forms of the same word as equivalent.
// Parameters:
//
//	tokens: a slice of strings representing the input tokens to be stemmed.
//
// Return values:
//
//	[]string: a new slice of strings containing the stemmed tokens obtained after applying the Snowball stemming algorithm to the input tokens.
func stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false) // Apply Snowball stemming algorithm to the token.
	}
	return r
}
