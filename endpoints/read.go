package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
)

func HandleRead(logger *zap.SugaredLogger) operations.TerseReadHandlerFunc {
	return func(params operations.TerseReadParams, _ *models.JWTInfo) middleware.Responder {
		// TODO
	}
}
