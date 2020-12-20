package endpoints

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations/api"
	"github.com/MicahParks/terse-URL/storage"
)

// HandleImport creates and /api/import endpoint handler via a closure. It can import Terse and or Visits data. It will
// delete existing data before importing, if told to do so.
func HandleImport(logger *zap.SugaredLogger, terseStore storage.TerseStore) api.TerseImportHandlerFunc {
	return func(params api.TerseImportParams) middleware.Responder {

		// Log the event.
		logger.Infow("Importing data.",
			"delete", fmt.Sprintf("%+v", params.ImportDelete.Delete),
		)

		// Create a request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Import the given data.
		if err := terseStore.Import(ctx, params.ImportDelete.Delete, params.ImportDelete.Import); err != nil {

			// Log at the appropriate level.
			message := "Failed to import data. Clean up may be necessary."
			logger.Warnw(message,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			resp := &api.TerseImportDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
			resp.SetStatusCode(int(code))
			return resp
		}

		return &api.TerseImportOK{}
	}
}
