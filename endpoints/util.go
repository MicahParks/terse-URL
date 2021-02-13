package endpoints

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/MicahParks/terseurl/models"
)

// defaultResponse is an interface used to pass different types of default responses and return an error responder.
type defaultResponse interface {
	SetStatusCode(code int)
	SetPayload(payload *models.Error)
	WriteResponse(rw http.ResponseWriter, producer runtime.Producer)
}

// ErrorResponse creates a response given the required assets.
func ErrorResponse(code int, message string, resp defaultResponse) middleware.Responder {

	// Set the payload for the response.
	resp.SetPayload(&models.Error{
		Code:    int64(code),
		Message: message,
	})

	// Set the status code for the response.
	resp.SetStatusCode(code)

	return resp
}
