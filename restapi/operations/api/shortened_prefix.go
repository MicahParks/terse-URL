// Code generated by go-swagger; DO NOT EDIT.

package api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ShortenedPrefixHandlerFunc turns a function with the right signature into a shortened prefix handler
type ShortenedPrefixHandlerFunc func(ShortenedPrefixParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ShortenedPrefixHandlerFunc) Handle(params ShortenedPrefixParams) middleware.Responder {
	return fn(params)
}

// ShortenedPrefixHandler interface for that can handle valid shortened prefix params
type ShortenedPrefixHandler interface {
	Handle(ShortenedPrefixParams) middleware.Responder
}

// NewShortenedPrefix creates a new http.Handler for the shortened prefix operation
func NewShortenedPrefix(ctx *middleware.Context, handler ShortenedPrefixHandler) *ShortenedPrefix {
	return &ShortenedPrefix{Context: ctx, Handler: handler}
}

/* ShortenedPrefix swagger:route GET /api/prefix api shortenedPrefix

Client's web browser is requesting what HTTP prefix all shortened URLs have.

Provides the HTTP prefix all shortened URLs have.

*/
type ShortenedPrefix struct {
	Context *middleware.Context
	Handler ShortenedPrefixHandler
}

func (o *ShortenedPrefix) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewShortenedPrefixParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}