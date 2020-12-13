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

func HandleDumpShortened(logger *zap.SugaredLogger, terseStore storage.TerseStore) operations.TerseDumpShortenedHandlerFunc {
	return func(params operations.TerseDumpShortenedParams) middleware.Responder {

		// Debug info.
		logger.Debugw("Performing data dump for shortened URL.",
			"shortened", params.Shortened,
		)

		// Create a request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the data dump.
		dump, err := terseStore.Dump(ctx, params.Shortened)
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
				message = "Failed to dump data for shortened URL."
				logger.Errorw(message,
					"shortened", params.Shortened,
					"error", err.Error(),
				)
			}

			// Report the error to the client.
			return &operations.TerseDumpShortenedDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		return &operations.TerseDumpShortenedOK{
			Payload: &dump,
		}
	}
}
