package endpoints

import (
	"errors"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
	"github.com/MicahParks/terse-URL/storage"
)

func HandleRead(logger *zap.SugaredLogger, terseStore storage.TerseStore) operations.TerseReadHandlerFunc {
	return func(params operations.TerseReadParams, _ *models.JWTInfo) middleware.Responder {

		// Debug info.
		logger.Debugw("Parameters",
			"shortened", params.Shortened,
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the Terse from the TerseStore.
		terse, err := terseStore.GetTerse(ctx, params.Shortened, nil)
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
			return &operations.TerseReadDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		return &operations.TerseReadOK{
			Payload: terse,
		}
	}
}
