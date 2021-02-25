// Code generated by go-swagger; DO NOT EDIT.

package api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

// NewVisitsReadParams creates a new VisitsReadParams object
//
// There are no default values defined in the spec.
func NewVisitsReadParams() VisitsReadParams {

	return VisitsReadParams{}
}

// VisitsReadParams contains all the bound params for the visits read operation
// typically these are obtained from a http.Request
//
// swagger:parameters visitsRead
type VisitsReadParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The shortened URLs to read the Visits data for.
	  Required: true
	  In: body
	*/
	ShortenedURLs []string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewVisitsReadParams() beforehand.
func (o *VisitsReadParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body []string
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("shortenedURLs", "body", ""))
			} else {
				res = append(res, errors.NewParseError("shortenedURLs", "body", "", err))
			}
		} else {
			// no validation required on inline body
			o.ShortenedURLs = body
		}
	} else {
		res = append(res, errors.Required("shortenedURLs", "body", ""))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
