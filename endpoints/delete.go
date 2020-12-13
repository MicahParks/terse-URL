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

func HandleDelete(logger *zap.SugaredLogger, terseStore storage.TerseStore) api.TerseDeleteHandlerFunc {
	return func(params api.TerseDeleteParams) middleware.Responder {

		// Do not have debug level logging on in production, as it will log clog up the logs.
		logger.Debugw("Parameters",
			"delete", fmt.Sprintf("%+v", params.Delete),
			"shortened", params.Shortened,
		)

		// TODO Non-debug level log?

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Delete the shortened URL's Terse pair from storage.
		var err error
		if params.Delete.Terse == nil || *params.Delete.Terse { // TODO Does the default value populate?
			if err = terseStore.DeleteTerse(ctx, params.Shortened); err != nil {

				// Log at the appropriate level.
				message := "Failed to delete Terse."
				logger.Warnw(message,
					"shortened", params.Shortened,
					"error", err.Error(),
				)

				// Report the error to the client.
				code := int64(500)
				return &api.TerseDeleteDefault{Payload: &models.Error{
					Code:    &code,
					Message: &message,
				}}
			}
		}

		// Delete the visits for the shortened URL.
		if params.Delete.Visits == nil || *params.Delete.Visits { // TODO Does the default value populate?
			if err = terseStore.VisitsStore().DeleteVisits(ctx, params.Shortened); err != nil {

				// Log with the appropriate level.
				message := "Failed to delete the visits for this shortened URL. Terse deleted."
				logger.Warnw(message,
					"shortened", params.Shortened,
					"error", err.Error(),
				)

				// Report the error to the client.
				code := int64(500)
				return &api.TerseDeleteDefault{Payload: &models.Error{
					Code:    &code,
					Message: &message,
				}}
			}
		}

		return &api.TerseDeleteOK{}
	}
}
