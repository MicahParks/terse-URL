// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/MicahParks/terse-URL/models"
)

// TerseDumpShortenedHandlerFunc turns a function with the right signature into a terse dump shortened handler
type TerseDumpShortenedHandlerFunc func(TerseDumpShortenedParams, *models.JWTInfo) middleware.Responder

// Handle executing the request and returning a response
func (fn TerseDumpShortenedHandlerFunc) Handle(params TerseDumpShortenedParams, principal *models.JWTInfo) middleware.Responder {
	return fn(params, principal)
}

// TerseDumpShortenedHandler interface for that can handle valid terse dump shortened params
type TerseDumpShortenedHandler interface {
	Handle(TerseDumpShortenedParams, *models.JWTInfo) middleware.Responder
}

// NewTerseDumpShortened creates a new http.Handler for the terse dump shortened operation
func NewTerseDumpShortened(ctx *middleware.Context, handler TerseDumpShortenedHandler) *TerseDumpShortened {
	return &TerseDumpShortened{Context: ctx, Handler: handler}
}

/*TerseDumpShortened swagger:route GET /api/dump/{shortened} terseDumpShortened

TerseDumpShortened terse dump shortened API

*/
type TerseDumpShortened struct {
	Context *middleware.Context
	Handler TerseDumpShortenedHandler
}

func (o *TerseDumpShortened) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewTerseDumpShortenedParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.JWTInfo
	if uprinc != nil {
		principal = uprinc.(*models.JWTInfo) // this is really a models.JWTInfo, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}