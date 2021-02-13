package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/restapi/operations/api"
	"github.com/MicahParks/terseurl/storage"
)

// HandleExport creates and /api/export endpoint handler via a closure. It can perform exports of all Terse and Visits
// data.
func HandleExport(logger *zap.SugaredLogger, terseStore storage.TerseStore) api.TerseExportHandlerFunc {
	return func(params api.TerseExportParams) middleware.Responder {

		// Log the event.
		logger.Info("Exporting all Terse and Visits data.")

		// Create a request context.
		//
		// Maybe make a longer context if timing out.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the data dump.
		dump, err := terseStore.Export(ctx)
		if err != nil {

			// Log at the appropriate level.
			logger.Errorw("Failed to perform data dump.",
				"error", err.Error(),
			)

			// Report the error to the client.
			code := 500
			message := "Failed to perform data dump."
			resp := &api.TerseExportDefault{Payload: &models.Error{
				Code:    int64(code),
				Message: message,
			}}
			resp.SetStatusCode(code)
			return resp
		}

		return &api.TerseExportOK{
			Payload: dump,
		}
	}
}
