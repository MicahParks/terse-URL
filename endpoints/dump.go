package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
)

func HandleDump(logger *zap.SugaredLogger) operations.TerseDumpHandlerFunc {
	return func(params operations.TerseDumpParams, _ *models.JWTInfo) middleware.Responder {
		// TODO
	}
}
