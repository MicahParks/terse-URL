package endpoints

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/MicahParks/terse-URL/restapi/operations/system"
)

func HandleAlive() system.AliveHandlerFunc {
	return func(params system.AliveParams) middleware.Responder {
		return &system.AliveOK{}
	}
}
