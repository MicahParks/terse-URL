package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/models"
	"github.com/MicahParks/terseurl/restapi/operations/api"
)

// HandleShortenedPrefix creates an /api/prefix endpoint handler via a closure. It let's the frontend client know the
// HTTP prefix for all shortened URLs.
func HandleShortenedPrefix(logger *zap.SugaredLogger, prefix string) api.ShortenedPrefixHandlerFunc {
	return func(params api.ShortenedPrefixParams, principal *models.Principal) middleware.Responder {

		// TODO Use principal.

		// Debug info.
		logger.Debug("Requested.")

		return &api.ShortenedPrefixOK{
			Payload: prefix,
		}
	}
}
