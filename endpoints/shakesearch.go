package endpoints

import (
	"strings"
	"unicode"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/shakesearch"
	"github.com/MicahParks/shakesearch/models"
	"github.com/MicahParks/shakesearch/restapi/operations/public"
)

var (

	// defaultMaxMatches is the largest number of matches to return if none was given.
	defaultMaxMatches = int64(20)
)

// HandleSearch creates a /api/search endpoint handler via a closure. It will do some basic cleaning of the query and
// respond with an array of results that contain some metadata.
func HandleSearch(logger *zap.SugaredLogger, shakeSearcher shakesearch.ShakeSearcher) public.ShakeSearchHandlerFunc {
	return func(params public.ShakeSearchParams) middleware.Responder {

		// Clean the query and confirm the cleaned query is not an empty string.
		if params.Q = cleanQuery(params.Q); params.Q == "" {

			// Debug info.
			message := "Client query failed validation."
			logger.Debugw(message,
				"query", params.Q,
			)

			// Report the error back to the client.
			code := int64(400)
			return &public.ShakeSearchDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		// Debug info.
		logger.Debugw("Performing search.",
			"query", params.Q,
		)

		// Determine the maximum number of matches to return.
		if params.MaxResults == nil {
			params.MaxResults = &defaultMaxMatches
		}

		// Find up to the maximum number of matches.
		matches := shakeSearcher.Search(int(*params.MaxResults), params.Q)

		// Perform the search on Shakespeare's works and return the info to the client.
		return &public.ShakeSearchOK{
			Payload: matches,
		}
	}
}

// cleanQuery trims spaces from the query, gets rid of non-alphanumeric characters, and keeps at most once space between
// them.
func cleanQuery(query string) (clean string) {

	// Trim whitespaces.
	query = strings.TrimSpace(query)

	// Only keep alphanumeric characters and at most one space between them.
	wasSpace := true
	for _, r := range query {

		// If the rune is a space and the previous rune was not the start of the string or another space, keep it.
		if r == ' ' {
			if !wasSpace {
				clean += string(r)
			}
			wasSpace = true
			continue
		}

		// If the rune is a letter or a digit, keep it.
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			clean += string(r)
		}

		wasSpace = false
	}

	return clean
}
