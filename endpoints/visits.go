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

// HandleVisitsRead creates and /api/visits/{shortened} endpoint handler via a closure. It can perform exports of a
// single shortened URL's Visits data.
func HandleVisitsRead(logger *zap.SugaredLogger, manager storage.StoreManager) api.VisitsReadHandlerFunc {
	return func(params api.VisitsReadParams) middleware.Responder {

		// Log the event.
		logger.Debugw("Reading shortened URL Visits data.",
			"shortened", params.ShortenedURLs,
		)

		// Create a request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the visits from storage.
		var err error
		var visits map[string][]models.Visit
		if visits, err = manager.Visits(ctx, params.ShortenedURLs); err != nil {

			// Log at the appropriate level. Assign the response code and message.
			var code int
			var message string
			if errors.Is(err, storage.ErrShortenedNotFound) {
				code = 400
				message = "Shortened URL not found."
				logger.Infow(message,
					"error", err.Error(),
				)
			} else {
				code = 500
				message = "Failed to read Visits data for the shortened URLs."
				logger.Errorw(message,
					"error", err.Error(),
				)
			}

			// Report the error to the client.
			return ErrorResponse(code, message, &api.VisitsReadDefault{})
		}

		return &api.VisitsReadOK{
			Payload: visits,
		}
	}
}
