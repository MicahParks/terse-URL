package endpoints

import (
	"errors"

	"github.com/go-openapi/runtime/middleware"
	"github.com/teris-io/shortid"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
	"github.com/MicahParks/terse-URL/storage"
)

func HandleWrite(logger *zap.SugaredLogger, shortID *shortid.Shortid, terseStore storage.TerseStore) operations.TerseWriteHandlerFunc {
	return func(params operations.TerseWriteParams) middleware.Responder {

		// Debug info.
		logger.Debugw("Performing Terse write.",
			"operation", params.Operation,
			"shortened", params.Terse.ShortenedURL,
		)

		// Create a request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Create the Terse data structure.
		terse := &models.Terse{
			DeleteAt:     params.Terse.DeleteAt, // TODO How to get zero value?
			MediaPreview: params.Terse.MediaPreview,
			OriginalURL:  params.Terse.OriginalURL,
			ShortenedURL: &params.Terse.ShortenedURL,
		}

		// If no shortened URL was given, create one.
		var err error
		if *terse.ShortenedURL, err = shortID.Generate(); err != nil { // TODO Loop this in paranoid mode?

			// Log at the appropriate level.
			message := "Failed to create random shortened URL."
			logger.Errorw(message,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			return &operations.TerseWriteDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		// Decide which operation to do.
		switch params.Operation {
		case "create":
			err = terseStore.CreateTerse(ctx, terse)
		case "update":
			err = terseStore.UpdateTerse(ctx, terse)
		case "upsert":
			err = terseStore.UpsertTerse(ctx, terse)
		}

		// Check for an error when writing to the TerseStore.
		if err != nil {

			// Log at the appropriate level. Assign the response code and message.
			var code int64
			var message string
			if errors.Is(err, storage.ErrShortenedExists) {
				code = 400
				message = "Not going to existing shortened URL."
				logger.Infow(message,
					"shortened", terse.ShortenedURL,
					"error", err.Error(),
				)
			} else {
				code = 500
				message = "Failed to write Terse."
				logger.Errorw(message,
					"shortened", *terse.ShortenedURL,
					"error", err.Error(),
				)
			}

			// Report the error to the client.
			return &operations.TerseWriteDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		return &operations.TerseWriteOK{
			Payload: *terse.ShortenedURL,
		}
	}
}
