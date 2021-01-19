package endpoints

import (
	"errors"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/restapi/operations/api"
	"github.com/MicahParks/terseurl/storage"
)

// HandleVisits creates and /api/visits/{shortened} endpoint handler via a closure. It can perform exports of a single
// shortened URL's Visits data.
func HandleVisits(logger *zap.SugaredLogger, visitsStore storage.VisitsStore) api.TerseVisitsHandlerFunc {
	return func(params api.TerseVisitsParams) middleware.Responder {

		// Log the event.
		logger.Infow("Reading a shortened URL's Visits data.",
			"shortened", params.Shortened,
		)

		// Create a request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the visits from storage.
		var err error
		visits := make([]*models.Visit, 0)
		if visitsStore != nil {
			if visits, err = visitsStore.ExportSome(ctx, params.Shortened); err != nil {

				// Log at the appropriate level. Assign the response code and message.
				var code int64
				var message string
				if errors.Is(err, storage.ErrShortenedNotFound) {
					code = 400
					message = "Shortened URL not found."
					logger.Infow(message,
						"shortened", params.Shortened,
						"error", err.Error(),
					)
				} else {
					code = 500
					message = "Failed to find the visits for the shortened URL."
					logger.Errorw(message,
						"shortened", params.Shortened,
						"error", err.Error(),
					)
				}

				// Report the error to the client.
				resp := &api.TerseVisitsDefault{Payload: &models.Error{
					Code:    &code,
					Message: &message,
				}}
				resp.SetStatusCode(int(code))
				return resp
			}
		}

		return &api.TerseVisitsOK{
			Payload: visits,
		}
	}
}
