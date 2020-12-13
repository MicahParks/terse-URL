package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
)

func HandleVisits(logger *zap.SugaredLogger) operations.TerseVisitsHandlerFunc {
	return func(params operations.TerseVisitsParams, _ *models.JWTInfo) middleware.Responder {
		// TODO
	}
}
