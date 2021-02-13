package endpoints

import (
	"errors"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/restapi/operations/api"
	"github.com/MicahParks/terseurl/storage"
)

// HandleExportSome creates and /api/export/some endpoint handler via a closure. It can perform exports of a
// single shortened URL's Terse and Visits data.
func HandleExportSome(logger *zap.SugaredLogger, terseStore storage.TerseStore) api.TerseExportSomeHandlerFunc {
	return func(params api.TerseExportSomeParams) middleware.Responder {

		// Log the event.
		logger.Infow("Exporting a shortened URL's Terse and Visits data.",
			"shortened", fmt.Sprintf("%v", params.ShortenedURLs),
		)

		// Create a request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the data dump.
		dump, err := terseStore.ExportSome(ctx, params.ShortenedURLs)
		if err != nil {

			// Log at the appropriate level. Assign the response code and message.
			var code int
			var message string
			if errors.Is(err, storage.ErrShortenedNotFound) {
				code = 400
				message = "Shortened URL not found."
				logger.Infow(message,
					"shortened", fmt.Sprintf("%v", params.ShortenedURLs),
					"error", err.Error(),
				)
			} else {
				code = 500
				message = "Failed to dump data for shortened URL."
				logger.Errorw(message,
					"shortened", fmt.Sprintf("%v", params.ShortenedURLs),
					"error", err.Error(),
				)
			}

			// Report the error to the client.
			resp := &api.TerseExportSomeDefault{Payload: &models.Error{
				Code:    int64(code),
				Message: message,
			}}
			resp.SetStatusCode(code)
			return resp
		}

		return &api.TerseExportSomeOK{
			Payload: dump,
		}
	}
}
