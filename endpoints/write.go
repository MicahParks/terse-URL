package endpoints

import (
	"errors"

	"github.com/go-openapi/runtime/middleware"
	"github.com/teris-io/shortid"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/restapi/operations/api"
	"github.com/MicahParks/terseurl/storage"
)

// HandleWrite creates and /api/write/{operation} endpoint handler via a closure. It can perform write operations on a
// single shortened URL's Terse data.
func HandleWrite(logger *zap.SugaredLogger, shortID *shortid.Shortid, manager storage.StoreManager) api.TerseWriteHandlerFunc {
	return func(params api.TerseWriteParams) middleware.Responder {

		// Debug info.
		logger.Debugw("Terse data",
			"terseData", params.Terse,
		)

		// Log the event.
		logger.Infow("Writing terse data.",
			"operation", params.Operation,
			"shortened", params.Terse.ShortenedURL,
		)

		// Create a request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// If no redirect type was give, use 302 as the default.
		if params.Terse.RedirectType == "" {
			params.Terse.RedirectType = models.RedirectTypeNr302
		}

		// Create the Terse data structure.
		terse := &models.Terse{
			JavascriptTracking: params.Terse.JavascriptTracking,
			MediaPreview:       params.Terse.MediaPreview,
			OriginalURL:        params.Terse.OriginalURL,
			RedirectType:       params.Terse.RedirectType,
			ShortenedURL:       params.Terse.ShortenedURL,
		}

		// If no shortened URL was given, create one.
		var err error
		if params.Terse.ShortenedURL == "" {
			if terse.ShortenedURL, err = shortID.Generate(); err != nil { // TODO Loop this in paranoid mode?

				// Log at the appropriate level.
				message := "Failed to create random shortened URL."
				logger.Errorw(message,
					"error", err.Error(),
				)

				// Report the error to the client.
				return ErrorResponse(500, message, &api.TerseWriteDefault{})
			}
		}

		// Decide which operation to do.
		switch terseMap := map[string]*models.Terse{terse.ShortenedURL: terse}; params.Operation {
		case "insert":
			err = manager.WriteTerse(ctx, terseMap, storage.Insert)
		case "update":
			err = manager.WriteTerse(ctx, terseMap, storage.Update)
		case "upsert":
			err = manager.WriteTerse(ctx, terseMap, storage.Upsert)
		}

		// Check for an error when writing to the TerseStore.
		if err != nil {

			// Log at the appropriate level. Assign the response code and message.
			var code int
			var message string
			if errors.Is(err, storage.ErrShortenedExists) {
				code = 400
				message = "Not going to overwrite existing shortened URL."
				logger.Infow(message,
					"shortened", terse.ShortenedURL,
					"error", err.Error(),
				)
			} else if errors.Is(err, storage.ErrShortenedNotFound) {
				code = 400
				message = "Shortened URL not found."
				logger.Infow(message,
					"shortened", terse.ShortenedURL,
					"error", err.Error(),
				)
			} else {
				code = 500
				message = "Failed to write Terse."
				logger.Errorw(message,
					"shortened", terse.ShortenedURL,
					"error", err.Error(),
				)
			}

			// Report the error to the client.
			return ErrorResponse(code, message, &api.TerseWriteDefault{})
		}

		return &api.TerseWriteOK{
			Payload: terse.ShortenedURL,
		}
	}
}
