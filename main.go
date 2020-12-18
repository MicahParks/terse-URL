package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/sahilm/fuzzy"
)

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}
	searcher.Search("Hamlet")

	//fs := http.FileServer(http.Dir("./static"))
	//http.Handle("/", fs)
	//
	//http.HandleFunc("/search", handleSearch(searcher))
	//
	//port := os.Getenv("PORT")
	//if port == "" {
	//	port = "3001"
	//}
	//
	//fmt.Printf("Listening on port %s...", port)
	//err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

type Searcher struct {
	CompleteWorks []string
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		results := searcher.Search(query[0])
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

// Load loads the give file as a slice of string with trimmed space.
func (s *Searcher) Load(filePath string) (err error) {

	// Read the file into memory.
	var fileData []byte
	if fileData, err = ioutil.ReadFile(filePath); err != nil {
		return err
	}

	// Split the file data by newlines.
	split := strings.Split(string(fileData), "\n")

	// Create a slice to store the file data.
	s.CompleteWorks = make([]string, 0)

	// Trim the space for the line and add it to the resulting slice. Do not add empty strings.
	for _, line := range split {
		if line = strings.TrimSpace(line); line != "" {
			s.CompleteWorks = append(s.CompleteWorks, line) // TODO Only add unique strings?
		}
	}

	return nil
}

func (s *Searcher) Search(query string) (results []string) {

	// Find all matching strings. Make that external dependency do all the work.
	matches := fuzzy.Find(query, s.CompleteWorks)

	// Create a map to behave as set of strings. The value is the original index.
	strSet := make(map[string]uint)

	// Add each match to the set of strings only once. Keep track of the indexes.
	index := uint(0)
	for _, match := range matches {
		if _, ok := strSet[match.Str]; !ok {
			strSet[match.Str] = index
			index++
		}
	}

	// Allocate the required memory for the return slice so it's faster to add all the strings.
	results = make([]string, len(strSet))

	// Turn the set of strings back into a slice. Insert the matches to their proper index.
	for str, index := range strSet {
		results[index] = str
	}

	return results
}
