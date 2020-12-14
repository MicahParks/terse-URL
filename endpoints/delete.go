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
func HandleDelete(logger *zap.SugaredLogger, terseStore storage.TerseStore) api.TerseDeleteHandlerFunc {
	return func(params api.TerseDeleteParams) middleware.Responder {

		// Log the event.
		logger.Infow("Deleting a shortened URL's assets.",
			"deleteTerse", params.Delete.Terse,
			"deleteVisits", params.Delete.Visits,
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Delete the shortened URL's Terse pair from storage.
		var err error
		if params.Delete.Terse == nil || *params.Delete.Terse {
			if err = terseStore.DeleteTerse(ctx, params.Shortened); err != nil {

				// Log at the appropriate level.
				message := "Failed to delete Terse."
				logger.Warnw(message,
					"shortened", params.Shortened,
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
		}

		// Delete the visits for the shortened URL.
		if params.Delete.Visits == nil || *params.Delete.Visits && terseStore.VisitsStore() != nil {
			if err = terseStore.VisitsStore().DeleteVisits(ctx, params.Shortened); err != nil {

				// Log with the appropriate level.
				message := "Failed to delete the visits for this shortened URL. Terse deleted."
				logger.Warnw(message,
					"shortened", params.Shortened,
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
		}

		return &api.TerseDeleteOK{}
	}
}
