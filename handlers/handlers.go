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

func NewSearchEngine(path string) *SearchEngine {
	s := &SearchEngine{
		Index: make(map[string][]int),
	}
	err := s.LoadDocuments(path)
	if err != nil {
		panic(err)
	}
	s.IndexDoc()
	return s
}
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

// add adds documents to the Index.
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

// tokenize returns a slice of tokens.
func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// Split on any character that is not a letter or a number.
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// analyze returns a slice of tokens.
func analyze(text string) []string {
	tokens := tokenize(text)
	tokens = lowercaseFilter(tokens)
	tokens = stopWordFilter(tokens)
	tokens = stemmerFilter(tokens)
	return tokens
}

func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}

// stopWordFilter returns a slice of tokens with stop words removed.
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
		if _, ok := stopWords[token]; !ok {
			filteredTokens = append(filteredTokens, token)
		}
	}
	return filteredTokens
}

// stemmerFilter returns a slice of stemmed tokens.
func stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false)
	}
	return r
}
