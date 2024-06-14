package main

import (
	"FullText_SearchEngine/handlers"
	"FullText_SearchEngine/views"
	"context"
	"flag"
	"fmt"
	"github.com/a-h/templ"
	"net/http"
	"strconv"
	"strings"
)

var SearchEngine *handlers.SearchEngine

var searchFilePath string

// init initializes the searchFilePath variable by parsing the command-line flags.
func init() {
	flag.StringVar(&searchFilePath, "file", "", "Path to the XML file for search engine initialization")
	flag.Parse()
}

// main is the entry point of the application.
func main() {
	// Check if the searchFilePath is provided as a command-line flag.
	if searchFilePath == "" {
		fmt.Println("Usage: ./yourApp -file <path_to_xml_file>")
		return
	}

	// Initialize the SearchEngine using the provided searchFilePath.
	SearchEngine = handlers.NewSearchEngine(searchFilePath)

	// Create an index page view with the number of documents in the SearchEngine.
	indexPage := views.Index(strconv.Itoa(len(SearchEngine.Documents)))

	// Handle the root path with the indexPage view.
	http.Handle("/", templ.Handler(indexPage))

	// Handle the "/search" path with the SearchHandler function.
	http.Handle("/search", http.HandlerFunc(SearchHandler))

	// Handle the "/doc" path with the DocHandler function.
	http.Handle("/doc", http.HandlerFunc(DocHandler))

	// Print a message indicating that the application is listening on port 3000.
	fmt.Println("Listening on :3000")

	// Start the HTTP server and listen on port 3000.
	http.ListenAndServe(":3000", nil)
}

// SearchHandler handles the "/search" path and processes the search query.
// It takes a http.ResponseWriter and a pointer to a http.Request as parameters.
// It retrieves the search query from the request and performs the search using the SearchEngine.
// If the query is empty, it writes "No query provided" to the response.
// If the search yields no results, it writes "No results found" to the response.
// Otherwise, it iterates through the search results, retrieves the corresponding documents from SearchEngine.Documents,
// creates a view for each document, and renders the view to the response.
func SearchHandler(writer http.ResponseWriter, request *http.Request) {
	query := request.FormValue("q")
	if len(query) == 0 {
		fmt.Fprintf(writer, "No query provided")
		return
	}
	var docs []int
	if strings.Contains(query, "\"") {
		docs = SearchEngine.SearchPhrase(query)
	} else if strings.Contains(query, "*") {
		docs = SearchEngine.FindWildcardMatches(query)
	} else {
		docs = SearchEngine.Search(query)
	}

	if docs == nil || len(docs) <= 0 {
		fmt.Fprintf(writer, "No results found")
	}
	for i, _ := range docs {
		doc := SearchEngine.Documents[docs[i]]
		item := views.Item(doc.Title, fmt.Sprint(doc.ID))
		item.Render(context.Background(), writer)
	}
}

// DocHandler handles the "/doc" path and processes the document request.
// It takes a http.ResponseWriter and a pointer to a http.Request as parameters.
// It retrieves the document ID from the request and fetches the corresponding document from SearchEngine.Documents.
// It then creates a view for the document using its title, text, and ID, and renders the view to the response.
func DocHandler(writer http.ResponseWriter, request *http.Request) {
	query := request.FormValue("id")                                    // Retrieve the document ID from the request.
	id, _ := strconv.Atoi(query)                                        // Convert the document ID to an integer.
	doc := SearchEngine.Documents[id]                                   // Fetch the document from SearchEngine.Documents using the ID.
	page := views.DocumentPage(doc.Title, doc.Text, fmt.Sprint(doc.ID)) // Create a view for the document with its title, text, and ID.
	page.Render(context.Background(), writer)                           // Render the view to the response.
}
