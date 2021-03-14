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

// HandleTerseRead creates and /api/terse/{shortened} endpoint handler via a closure. It can perform exports of a single
// shortened URL's Terse data.
func HandleTerseRead(logger *zap.SugaredLogger, manager storage.StoreManager) api.TerseReadHandlerFunc {
	return func(params api.TerseReadParams, principal *models.Principal) middleware.Responder {

		// Log the event.
		logger.Infow("Reading a shortened URL's Terse data.",
			"shortenedURLs", params.ShortenedURLs,
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the Terse.
		terse, err := manager.Terse(ctx, principal, params.ShortenedURLs)
		if err != nil {

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
				message = "Failed to get Terse from shortened URL."
				logger.Errorw(message,
					"error", err.Error(),
				)
			}

			// Report the error to the client.
			return ErrorResponse(code, message, &api.TerseReadDefault{})
		}

		return &api.TerseReadOK{
			Payload: terse,
		}
	}
}
