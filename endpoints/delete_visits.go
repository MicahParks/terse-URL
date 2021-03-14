package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/restapi/operations/api"
	"github.com/MicahParks/terseurl/storage"
)

// HandlerVisitsDelete TODO
func HandlerVisitsDelete(logger *zap.SugaredLogger, manager storage.StoreManager) api.VisitsDeleteHandlerFunc {
	return func(params api.VisitsDeleteParams, principal *models.Principal) middleware.Responder {

		// Debug info.
		logger.Debugw("Deleting Visits data.",
			"shortenedURLs", params.ShortenedURLs,
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Delete Visits data for the requested shortened URLs.
		if err := manager.DeleteShortened(ctx, principal, params.ShortenedURLs); err != nil {

			// Log at the appropriate level.
			message := "Failed to delete Visits data for the requested shortened URLs."
			logger.Infow(message,
				"error", err.Error(),
			)

			// Report the error to the client.
			return ErrorResponse(500, message, &api.VisitsDeleteDefault{})
		}

		return &api.VisitsDeleteOK{}
	}
}
