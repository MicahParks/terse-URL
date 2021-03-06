// Code generated by go-swagger; DO NOT EDIT.

package api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/MicahParks/terseurl/models"
)

// FrontendMetaHandlerFunc turns a function with the right signature into a frontend meta handler
type FrontendMetaHandlerFunc func(FrontendMetaParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn FrontendMetaHandlerFunc) Handle(params FrontendMetaParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// FrontendMetaHandler interface for that can handle valid frontend meta params
type FrontendMetaHandler interface {
	Handle(FrontendMetaParams, *models.Principal) middleware.Responder
}

// NewFrontendMeta creates a new http.Handler for the frontend meta operation
func NewFrontendMeta(ctx *middleware.Context, handler FrontendMetaHandler) *FrontendMeta {
	return &FrontendMeta{Context: ctx, Handler: handler}
}

/* FrontendMeta swagger:route POST /api/frontend/meta api frontendMeta

Used by the frontend to inherit HTML meta for social media link previews.

This endpoint is intended for use only by the frontend. It will perform an HTTP GET request to the originalURL, extract all meta tag information for social media link previews, then return it.

*/
type FrontendMeta struct {
	Context *middleware.Context
	Handler FrontendMetaHandler
}

func (o *FrontendMeta) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewFrontendMetaParams()
	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// FrontendMetaOKBody frontend meta o k body
//
// swagger:model FrontendMetaOKBody
type FrontendMetaOKBody struct {

	// og
	Og models.OpenGraph `json:"og,omitempty"`

	// twitter
	Twitter models.Twitter `json:"twitter,omitempty"`
}

// Validate validates this frontend meta o k body
func (o *FrontendMetaOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateOg(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateTwitter(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *FrontendMetaOKBody) validateOg(formats strfmt.Registry) error {
	if swag.IsZero(o.Og) { // not required
		return nil
	}

	if o.Og != nil {
		if err := o.Og.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("frontendMetaOK" + "." + "og")
			}
			return err
		}
	}

	return nil
}

func (o *FrontendMetaOKBody) validateTwitter(formats strfmt.Registry) error {
	if swag.IsZero(o.Twitter) { // not required
		return nil
	}

	if o.Twitter != nil {
		if err := o.Twitter.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("frontendMetaOK" + "." + "twitter")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this frontend meta o k body based on the context it is used
func (o *FrontendMetaOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateOg(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := o.contextValidateTwitter(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *FrontendMetaOKBody) contextValidateOg(ctx context.Context, formats strfmt.Registry) error {

	if err := o.Og.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("frontendMetaOK" + "." + "og")
		}
		return err
	}

	return nil
}

func (o *FrontendMetaOKBody) contextValidateTwitter(ctx context.Context, formats strfmt.Registry) error {

	if err := o.Twitter.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("frontendMetaOK" + "." + "twitter")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *FrontendMetaOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *FrontendMetaOKBody) UnmarshalBinary(b []byte) error {
	var res FrontendMetaOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
