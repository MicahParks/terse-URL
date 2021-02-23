package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/restapi/operations/api"
)

// HandleShortenedURLPrefix creates an /api/prefix endpoint handler via a closure. It let's the frontend client know the
// HTTP prefix for all shortened URLs.
func HandleShortenedURLPrefix(logger *zap.SugaredLogger, prefix string) api.ShortenedURLPrefixHandlerFunc {
	return func(params api.ShortenedURLPrefixParams) middleware.Responder {

		// Debug info.
		logger.Debug("Requested.")

		return &api.ShortenedURLPrefixOK{
			Payload: prefix,
		}
	}
}
