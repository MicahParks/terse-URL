package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/restapi/operations/api"
	"github.com/MicahParks/terseurl/storage"
)

// HandleShortenedDelete TODO
func HandleShortenedDelete(logger *zap.SugaredLogger, manager storage.StoreManager) api.ShortenedDeleteHandlerFunc {
	return func(params api.ShortenedDeleteParams, principal *models.Principal) middleware.Responder {

		// Debug info.
		logger.Debugw("Deleting data.",
			"shortenedURLs", params.ShortenedURLs,
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Delete all data for the requested shortened URLs.
		if err := manager.DeleteShortened(ctx, params.ShortenedURLs); err != nil {

			// Log at the appropriate level.
			message := "Failed to delete data for the requested shortened URLs."
			logger.Infow(message,
				"error", err.Error(),
			)

			// Report the error to the client.
			return ErrorResponse(500, message, &api.ShortenedDeleteDefault{})
		}

		return &api.ShortenedDeleteOK{}
	}
}
