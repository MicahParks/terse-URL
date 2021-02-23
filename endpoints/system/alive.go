package system

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/MicahParks/terseurl/restapi/operations/system"
)

// HandleAlive creates and /api/alive endpoint handler via a closure.
func HandleAlive() system.SystemAliveHandlerFunc {
	return func(params system.SystemAliveParams) middleware.Responder {
		return &system.SystemAliveOK{}
	}
}
