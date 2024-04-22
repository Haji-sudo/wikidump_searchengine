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

func init() {
	flag.StringVar(&searchFilePath, "file", "", "Path to the XML file for search engine initialization")
	flag.Parse()
}
func main() {
	if searchFilePath == "" {
		fmt.Println("Usage: ./yourapp -file <path_to_xml_file>")
		return
	}
	SearchEngine = handlers.NewSearchEngine(searchFilePath)

	indexPage := views.Index(strconv.Itoa(len(SearchEngine.Documents)))

	http.Handle("/", templ.Handler(indexPage))
	http.Handle("/search", http.HandlerFunc(SearchHandler))
	http.Handle("/doc", http.HandlerFunc(DocHandler))
	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}

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

func DocHandler(writer http.ResponseWriter, request *http.Request) {
	query := request.FormValue("id")
	id, _ := strconv.Atoi(query)
	doc := SearchEngine.Documents[id]
	page := views.DocumentPage(doc.Title, doc.Text, fmt.Sprint(doc.ID))
	page.Render(context.Background(), writer)
}
