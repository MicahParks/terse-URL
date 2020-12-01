package endpoints

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/storage"
)

func HandleTrack(logger *zap.SugaredLogger, visitsStorage storage.VisitsStore) operations.URLTrackHandlerFunc {
	return func(params operations.URLTrackParams, _ *models.JWTInfo) middleware.Responder {

		// Do not have debug level logging on in production, as it will log clog up the logs.
		logger.Debugw("Parameters",
			"shortened", fmt.Sprintf("%+v", params.Shortened),
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the visits from storage.
		var err error
		var visits []*models.Visit
		if visits, err = visitsStorage.GetVisits(ctx, params.Shortened); err != nil {

			// Log with the appropriate level.
			logger.Warnw("Failed to find the visits for the given shortened URL.",
				"shortened", params.Shortened,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			message := "Failed to retrieve visits from storage."
			return &operations.URLTrackDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		return &operations.URLTrackOK{
			Payload: visits,
		}
	}
}
