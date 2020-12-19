package endpoints

import (
	"errors"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations/api"
	"github.com/MicahParks/terse-URL/storage"
)

// HandleTerse creates and /api/terse/{shortened} endpoint handler via a closure. It can perform exports of a single
// shortened URL's Terse data.
func HandleTerse(logger *zap.SugaredLogger, terseStore storage.TerseStore) api.TerseTerseHandlerFunc {
	return func(params api.TerseTerseParams) middleware.Responder {

		// Log the event.
		logger.Infow("Reading a shortened URL's Terse data.",
			"shortened", params.Shortened,
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the Terse from the TerseStore.
		terse, err := terseStore.Read(ctx, params.Shortened, nil)
		if err != nil {

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
				message = "Failed to get Terse from shortened URL."
				logger.Errorw(message,
					"shortened", params.Shortened,
					"error", err.Error(),
				)
			}

			// Report the error to the client.
			resp := &api.TerseTerseDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
			resp.SetStatusCode(int(code))
			return resp
		}

		return &api.TerseTerseOK{
			Payload: terse,
		}
	}
}