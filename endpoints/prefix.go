package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/restapi/operations/api"
)

// HandlePrefix creates an /api/prefix endpoint handler via a closure. It let's the frontend client know the HTTP prefix
// for all shortened URLs.
func HandlePrefix(logger *zap.SugaredLogger, prefix string) api.TersePrefixHandlerFunc {
	return func(params api.TersePrefixParams) middleware.Responder {

		// Debug info.
		logger.Debug("Requested.")

		return &api.TersePrefixOK{
			Payload: prefix,
		}
	}
}
