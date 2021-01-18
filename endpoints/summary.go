package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terseurl/restapi/operations/api"
)

// TODO
func HandleSummary(logger *zap.SugaredLogger) api.TerseSummaryHandlerFunc {
	return func(params api.TerseSummaryParams) middleware.Responder {
		// TODO
	}
}
