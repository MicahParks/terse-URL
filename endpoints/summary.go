package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/restapi/operations/api"
	"github.com/MicahParks/terseurl/storage"
)

// HandleShortenedSummary creates a /api/summary endpoint handler via a closure. It can provide Summary data for the
// requested shortened URLs.
func HandleShortenedSummary(logger *zap.SugaredLogger, manager storage.StoreManager) api.ShortenedSummaryHandlerFunc {
	return func(params api.ShortenedSummaryParams, principal *models.Principal) middleware.Responder {

		// Debug info.
		logger.Debugw("Requested summary data.",
			"shortenedURLs", params.ShortenedURLs,
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Gather the summary information for the requested shortened URLs.
		summaries, err := manager.Summary(ctx, params.ShortenedURLs)
		if err != nil {

			// Log at the appropriate level.
			message := "Failed to gather summary information for requested shortened URLs."
			logger.Infow(message,
				"error", err.Error(),
			)

			// Report the error to the client.
			return ErrorResponse(500, message, &api.ShortenedSummaryDefault{})
		}

		return &api.ShortenedSummaryOK{
			Payload: summaries,
		}
	}
}
