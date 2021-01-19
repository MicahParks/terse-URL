package endpoints

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/configure"
	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/restapi/operations/api"
	"github.com/MicahParks/terseurl/storage"
)

// HandleDeleteSome creates a /api/delete/some endpoint handler via a closure. It can delete Terse and Visits data
// given the associated shortened URL.
func HandleDeleteSome(logger *zap.SugaredLogger, terseStore storage.TerseStore) api.TerseDeleteSomeHandlerFunc {
	return func(params api.TerseDeleteSomeParams) middleware.Responder {

		// Log the event.
		logger.Infow("Deleting a shortened URL's assets.", // TODO Log delete info?
			"shortened", fmt.Sprintf("%v", params.Info.ShortenedURLs),
		)

		// Create a new request context.
		ctx, cancel := configure.DefaultCtx()
		defer cancel()

		// Delete the shortened URL's Terse pair from storage.
		if err := terseStore.DeleteSome(ctx, *params.Info.Delete, params.Info.ShortenedURLs); err != nil {

			// Log at the appropriate level.
			message := "Failed to delete Terse or Visits data."
			logger.Warnw(message,
				"shortened", fmt.Sprintf("%v", params.Info.ShortenedURLs),
				"error", err.Error(),
			)

			// Report the error to the client.
			code := int64(500)
			resp := &api.TerseDeleteSomeDefault{Payload: &models.Error{
				Code:    &code,
				Message: &message,
			}}
			resp.SetStatusCode(int(code))
			return resp
		}

		return &api.TerseDeleteSomeOK{}
	}
}
