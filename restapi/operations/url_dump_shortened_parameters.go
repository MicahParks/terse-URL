// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// NewURLDumpShortenedParams creates a new URLDumpShortenedParams object
// no default values defined in spec.
func NewURLDumpShortenedParams() URLDumpShortenedParams {

	return URLDumpShortenedParams{}
}

// URLDumpShortenedParams contains all the bound params for the url dump shortened operation
// typically these are obtained from a http.Request
//
// swagger:parameters urlDumpShortened
type URLDumpShortenedParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	Shortened string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewURLDumpShortenedParams() beforehand.
func (o *URLDumpShortenedParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rShortened, rhkShortened, _ := route.Params.GetOK("shortened")
	if err := o.bindShortened(rShortened, rhkShortened, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindShortened binds and validates parameter Shortened from path.
func (o *URLDumpShortenedParams) bindShortened(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.Shortened = raw

	return nil
}
