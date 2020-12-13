package endpoints

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/MicahParks/terse-URL/models"
	"github.com/MicahParks/terse-URL/restapi/operations"
)

func HandleWrite(logger *zap.SugaredLogger) operations.TerseWriteHandlerFunc {
	return func(params operations.TerseWriteParams, _ *models.JWTInfo) middleware.Responder {
		// TODO
	}
}
