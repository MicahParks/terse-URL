package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/restapi/operations/api"
	"github.com/MicahParks/terseurl/storage"
)

// HandleExport creates and /api/export endpoint handler via a closure. It can perform exports of all Terse and Visits
// data.
func HandleExport(logger *zap.SugaredLogger, manager storage.StoreManager) api.ExportHandlerFunc {
	return func(params api.ExportParams) middleware.Responder {

		// Log the event.
		logger.Info("Exporting data.")

		// Create a request context.
		//
		// Maybe make a longer context if timing out.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the data dump.
		dump, err := manager.Export(ctx, params.ShortenedURLs)
		if err != nil {

			// Log at the appropriate level.
			message := "Failed to perform data dump."
			logger.Warnw(message,
				"error", err.Error(),
			)

			// Report the error to the client.
			return ErrorResponse(500, message, &api.ExportDefault{})
		}

		return &api.ExportOK{
			Payload: dump,
		}
	}
}
