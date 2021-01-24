// Code generated by go-swagger; DO NOT EDIT.

package api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

// NewTerseSummaryParams creates a new TerseSummaryParams object
//
// There are no default values defined in the spec.
func NewTerseSummaryParams() TerseSummaryParams {

	return TerseSummaryParams{}
}

// TerseSummaryParams contains all the bound params for the terse summary operation
// typically these are obtained from a http.Request
//
// swagger:parameters terseSummary
type TerseSummaryParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The array of shortened URLs to get Terse summary data for. If none is provided, all will summaries will be returned.
	  In: body
	*/
	Shortened []string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewTerseSummaryParams() beforehand.
func (o *TerseSummaryParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body []string
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			res = append(res, errors.NewParseError("shortened", "body", "", err))
		} else {
			// no validation required on inline body
			o.Shortened = body
		}
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
