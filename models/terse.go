// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Terse terse
//
// swagger:model Terse
type Terse struct {

	// javascript tracking
	JavascriptTracking bool `json:"javascriptTracking,omitempty"`

	// media preview
	MediaPreview *MediaPreview `json:"mediaPreview,omitempty"`

	// original URL
	// Required: true
	OriginalURL *string `json:"originalURL"`

	// redirect type
	RedirectType RedirectType `json:"redirectType,omitempty"`

	// shortened URL
	// Required: true
	ShortenedURL *string `json:"shortenedURL"`

	// visit count
	VisitCount int64 `json:"visitCount,omitempty"`
}

// Validate validates this terse
func (m *Terse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMediaPreview(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOriginalURL(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRedirectType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateShortenedURL(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Terse) validateMediaPreview(formats strfmt.Registry) error {

	if swag.IsZero(m.MediaPreview) { // not required
		return nil
	}

	if m.MediaPreview != nil {
		if err := m.MediaPreview.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("mediaPreview")
			}
			return err
		}
	}

	return nil
}

func (m *Terse) validateOriginalURL(formats strfmt.Registry) error {

	if err := validate.Required("originalURL", "body", m.OriginalURL); err != nil {
		return err
	}

	return nil
}

func (m *Terse) validateRedirectType(formats strfmt.Registry) error {

	if swag.IsZero(m.RedirectType) { // not required
		return nil
	}

	if err := m.RedirectType.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("redirectType")
		}
		return err
	}

	return nil
}

func (m *Terse) validateShortenedURL(formats strfmt.Registry) error {

	if err := validate.Required("shortenedURL", "body", m.ShortenedURL); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Terse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Terse) UnmarshalBinary(b []byte) error {
	var res Terse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
