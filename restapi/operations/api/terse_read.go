// Code generated by go-swagger; DO NOT EDIT.

package api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// TerseReadHandlerFunc turns a function with the right signature into a terse read handler
type TerseReadHandlerFunc func(TerseReadParams) middleware.Responder

// Handle executing the request and returning a response
func (fn TerseReadHandlerFunc) Handle(params TerseReadParams) middleware.Responder {
	return fn(params)
}

// TerseReadHandler interface for that can handle valid terse read params
type TerseReadHandler interface {
	Handle(TerseReadParams) middleware.Responder
}

// NewTerseRead creates a new http.Handler for the terse read operation
func NewTerseRead(ctx *middleware.Context, handler TerseReadHandler) *TerseRead {
	return &TerseRead{Context: ctx, Handler: handler}
}

/* TerseRead swagger:route POST /api/terse api terseRead

Read the Terse data for the given shortened URL.

Read the Terse data for the given shortened URL.

*/
type TerseRead struct {
	Context *middleware.Context
	Handler TerseReadHandler
}

func (o *TerseRead) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewTerseReadParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}