package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations/api"
	"github.com/MicahParks/terse-URL/storage"
)

// HandleDelete creates a /api/delete/{shortened} endpoint handler via a closure. It can delete Terse and Visits data
// given the associated shortened URL.
func HandleDeleteOne(logger *zap.SugaredLogger, terseStore storage.TerseStore) api.TerseDeleteOneHandlerFunc {
	return func(params api.TerseDeleteOneParams) middleware.Responder {

		// Log the event.
		logger.Infow("Deleting a shortened URL's assets.", // TODO Log delete info?
			"shortened", params.Shortened,
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Delete the shortened URL's Terse pair from storage.
		if err := terseStore.DeleteOne(ctx, *params.Delete, params.Shortened); err != nil {

			// Log at the appropriate level.
			message := "Failed to delete Terse or Visits data."
			logger.Warnw(message,
				"shortened", params.Shortened,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			resp := &api.TerseDeleteOneDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
			resp.SetStatusCode(int(code))
			return resp
		}

		return &api.TerseDeleteOneOK{}
	}
}
