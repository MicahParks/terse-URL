package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
)

func HandleDumpShortened(logger *zap.SugaredLogger) operations.TerseDumpShortenedHandlerFunc {
	return func(params operations.TerseDumpShortenedParams, _ *models.JWTInfo) middleware.Responder {
		// TODO
	}
}
