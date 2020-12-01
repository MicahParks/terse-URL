package endpoints

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/MicahParks/terse-URL/restapi/operations"
)

func HandleAlive() operations.AliveHandlerFunc {
	return func(params operations.AliveParams) middleware.Responder {
		return &operations.AliveOK{}
	}
}
