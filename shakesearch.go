package shakesearch

import (
	"io/ioutil"
	"strings"

	"github.com/sahilm/fuzzy"
)

// ShakeSearcher is the data structure that holds the unqiue lines of Shakespeare's complete works in a convent format
// to search. Its methods allow for loading a searching the data.
type ShakeSearcher struct {
	uniqueLines []string
}

// Load loads the give file as a slice of string with trimmed space.
func NewShakeSearcher(filePath string) (shakeSearcher *ShakeSearcher, err error) {

	// Create the ShakeSearcher.
	shakeSearcher = &ShakeSearcher{}

	// Read the file into memory.
	var fileData []byte
	if fileData, err = ioutil.ReadFile(filePath); err != nil {
		return nil, err
	}

	// Split the file data by newlines.
	split := strings.Split(string(fileData), "\n")

	// Create a slice to store the file data.
	shakeSearcher.uniqueLines = make([]string, 0)

	// Create a map to behave as a set of strings.
	lineSet := make(map[string]bool)

	// Trim the space for the line and add it to the resulting slice. Do not add empty strings or non-unique strings.
	for _, line := range split {
		if line = strings.TrimSpace(line); line != "" {
			lineSet[line] = true
		}
	}

	// Allocate the required memory for the return slice so it's faster to add all the strings.
	shakeSearcher.uniqueLines = make([]string, len(lineSet))

	// Add every unique line to the list of complete works.
	index := 0
	for line := range lineSet {
		shakeSearcher.uniqueLines[index] = line
		index++
	}

	return shakeSearcher, nil
}

func (s *ShakeSearcher) Search(query string) (results []string) {

	// Find all matching strings. Make that external dependency do all the work.
	matches := fuzzy.Find(query, s.uniqueLines)

	// Create a map to behave as a set of strings. The value is the original index.
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
