package endpoints

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
	"github.com/MicahParks/terse-URL/storage"
)

func HandleDelete(logger *zap.SugaredLogger, terseStore storage.TerseStore) operations.URLDeleteHandlerFunc {
	return func(params operations.URLDeleteParams, _ *models.JWTInfo) middleware.Responder {

		// Do not have debug level logging on in production, as it will log clog up the logs.
		logger.Debugw("Parameters",
			"shortened", fmt.Sprintf("%+v", params.Shortened),
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Delete the shortened URL's Terse pair from storage.
		var err error
		if err = terseStore.DeleteTerse(ctx, params.Shortened); err != nil {

			// Log with the appropriate level.
			message := "Failed to delete Terse pair."
			logger.Warnw(message,
				"shortened", params.Shortened,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			return &operations.URLDeleteDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		// Delete the visits for the shortened URL.
		if err = terseStore.VisitsStore().DeleteVisits(ctx, params.Shortened); err != nil {

			// Log with the appropriate level.
			message := "Failed to delete the visits for this shortened URL. Terse pair deleted."
			logger.Warnw(message,
				"shortened", params.Shortened,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			return &operations.URLDeleteDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		return &operations.URLDeleteOK{}
	}
}
