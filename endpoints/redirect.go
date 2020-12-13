package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/restapi/operations"
)

func HandleRedirect(logger *zap.SugaredLogger) operations.TerseRedirectHandlerFunc {
	return func(params operations.TerseRedirectParams) middleware.Responder {
		// TODO
	}
}
