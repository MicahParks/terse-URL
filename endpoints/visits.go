package endpoints

import (
	"errors"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/configure"
	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
	"github.com/MicahParks/terse-URL/storage"
)

func HandleVisits(logger *zap.SugaredLogger, visitsStore storage.VisitsStore) operations.TerseVisitsHandlerFunc {
	return func(params operations.TerseVisitsParams, _ *models.JWTInfo) middleware.Responder {

		// Debug info.
		logger.Debugw("Parameters",
			"shortened", params.Shortened,
		)

		// Create a request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Get the visits from storage.
		visits, err := visitsStore.ReadVisits(ctx, params.Shortened)
		if err != nil {

			// Log at the appropriate level.
			if errors.Is(err, storage.ErrShortenedNotFound) {

			}
			logger.Infow("Failed to find the visits for the requested shortened URL.",
				"shortened", params.Shortened,
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			message := "Failed to retrieve visits from storage."
			return &operations.TerseVisitsDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
		}

		return &operations.TerseVisitsOK{
			Payload: visits,
		}
	}
}
