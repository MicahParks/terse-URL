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

// HandleExportOne creates and /api/export/{shortened} endpoint handler via a closure. It can perform exports of a
// single shortened URL's Terse and Visits data.
func HandleExportOne(logger *zap.SugaredLogger, terseStore storage.TerseStore) api.TerseExportOneHandlerFunc {
	return func(params api.TerseExportOneParams) middleware.Responder {

		// Log the event.
		logger.Infow("Exporting a shortened URL's Terse and Visits data.",
			"shortened", params.Shortened,
		)

		// Create a request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the data dump.
		dump, err := terseStore.ExportOne(ctx, params.Shortened)
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
			resp := &api.TerseExportOneDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
			resp.SetStatusCode(int(code))
			return resp
		}

		return &api.TerseExportOneOK{
			Payload: &dump,
		}
	}
}
