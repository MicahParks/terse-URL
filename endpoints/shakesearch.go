package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/shakesearch"
	"github.com/MicahParks/shakesearch/models"
	"github.com/MicahParks/shakesearch/restapi/operations/public"
)

func HandleSearch(logger *zap.SugaredLogger, shakeSearcher *shakesearch.ShakeSearcher) public.ShakeSearchHandlerFunc {
	return func(params public.ShakeSearchParams) middleware.Responder {

		// Do some basic input validation.
		if params.Q == "" {

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

		// Perform the search on Shakespeare's works and return the info to the client.
		return &public.ShakeSearchOK{
			Matches: shakeSearcher.Search(params.Q),
		}
	}
}
