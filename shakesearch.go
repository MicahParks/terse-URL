package shakesearch

import (
	"io/ioutil"
	"strings"

	"github.com/sahilm/fuzzy"

	"github.com/MicahParks/shakesearch/models"
)

// ShakeSearcher is the data structure that holds the unique lines of Shakespeare's complete works in a convent format
// to search. The format is a map of all lines and a slice of the line numbers they belong to.
type ShakeSearcher map[string]*models.Result

// NewShakeSearcher loads the give file as a slice of string with trimmed space.
func NewShakeSearcher(filePath string) (shakeSearcher ShakeSearcher, err error) {

	// Create the ShakeSearcher.
	shakeSearcher = make(map[string]*models.Result)

	// Read the file into memory.
	var fileData []byte
	if fileData, err = ioutil.ReadFile(filePath); err != nil {
		return nil, err
	}

	// Split the file data by newlines.
	split := strings.Split(string(fileData), "\n")

	// Add every unique line to the list of complete works.
	for index, line := range split {

		// Trim the space of the line. If it was only whitespace, ignore the line.
		if line = strings.TrimSpace(line); line != "" {

			// Get the existing UniqueLine data structure's pointer from the map.
			uniqueLine, ok := shakeSearcher[line]

			// If this is the first time this line has been seen, create the data structure and assign its pointer to
			// the map.
			if !ok {
				uniqueLine = &models.Result{
					Line:        line,
					LineNumbers: make([]int64, 0),
				}
				shakeSearcher[line] = uniqueLine
			}

			// Add the line number for this line to the UniqueLine data structure in the map.
			uniqueLine.LineNumbers = append(uniqueLine.LineNumbers, int64(index+1))
		}
	}

	return shakeSearcher, nil
}

// Search does a fuzzy search of the complete works of Shakespeare and returns a slice of results with a length less
// than or equal to the maxMatches. A maxMatches with a value of -1 will return all matches.
func (s ShakeSearcher) Search(maxMatches int, query string) (results []*models.Result) {

	// Create a slice that will hold all the unique lines. Allocate the memory up front so it inserts faster.
	uniqueLines := make([]string, len(s))

	// Insert all of the unique lines to the slice.
	index := uint(0)
	for line := range s {
		uniqueLines[index] = line
		index++
	}

	// Find all fuzzy matching strings.
	matches := fuzzy.Find(query, uniqueLines)

	// Figure out how many results to return.
	var matchCount int
	if maxMatches <= len(matches) {
		matchCount = maxMatches
	} else {
		matchCount = len(matches)
	}

	// Allocate the required memory for the return slice so it's faster to add all the strings.
	results = make([]*models.Result, matchCount)

	// Add each match to the set of strings only once. Keep track of the indexes.
	for index, match := range matches {

		// Make sure to only return the maximum number of matches.
		if maxMatches <= 0 || index == maxMatches {
			break
		}

		// Get a copy of the result to manipulate and return.
		result := *s[match.Str]

		// Allocate the memory for the matching indexes slice so inserting them is faster.
		result.MatchedIndexes = make([]int64, len(match.MatchedIndexes))

		// Turn the matched indexes into the correct integer format.
		for i, matchedIndex := range match.MatchedIndexes {
			result.MatchedIndexes[i] = int64(matchedIndex)
		}

		// Use the modified copy in the result to return.
		results[index] = &result
	}

	return results
}
