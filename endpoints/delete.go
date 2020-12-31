package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/restapi/operations/api"
	"github.com/MicahParks/terseurl/storage"
)

// HandleDelete creates a /api/delete/{shortened} endpoint handler via a closure. It can delete Terse and Visits data
// given the associated shortened URL.
func HandleDelete(logger *zap.SugaredLogger, terseStore storage.TerseStore) api.TerseDeleteHandlerFunc {
	return func(params api.TerseDeleteParams) middleware.Responder {

		// Log the event.
		logger.Infow("Deleting Terse and or Visits data.") // TODO Log delete info?

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Delete the shortened URL's Terse pair from storage.
		if err := terseStore.Delete(ctx, *params.Delete); err != nil {

			// Log at the appropriate level.
			message := "Failed to delete Terse or Visits data."
			logger.Warnw(message,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			resp := &api.TerseDeleteDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
			resp.SetStatusCode(int(code))
			return resp
		}

		return &api.TerseDeleteOK{}
	}
}
