package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/restapi/operations/api"
	"github.com/MicahParks/terseurl/storage"
)

// HandleSummary creates a /api/summary endpoint handler via a closure. It can provide Terse summary data for the
// requested shortened URLs.
func HandleSummary(logger *zap.SugaredLogger, summaryStore storage.SummaryStore) api.TerseSummaryHandlerFunc {
	return func(params api.TerseSummaryParams) middleware.Responder {

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Gather the summary information for the requested shortened URLs.
		summaries, err := summaryStore.Summarize(ctx, params.Shortened)
		if err != nil {

			// Log at the appropriate level.
			message := "Failed to gather summary information for requested shortened URLs."
			logger.Infow(message,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			resp := &api.TerseSummaryDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
			resp.SetStatusCode(int(code))
			return resp
		}

		return &api.TerseSummaryOK{
			Payload: summaries,
		}
	}
}
