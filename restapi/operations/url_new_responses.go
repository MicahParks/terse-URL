// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/MicahParks/terse-URL/models"
)

// URLNewOKCode is the HTTP code returned for type URLNewOK
const URLNewOKCode int = 200

/*URLNewOK The shortened URL to visit that will redirect to the given full URL.

swagger:response urlNewOK
*/
type URLNewOK struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewURLNewOK creates URLNewOK with default headers values
func NewURLNewOK() *URLNewOK {

	return &URLNewOK{}
}

// WithPayload adds the payload to the url new o k response
func (o *URLNewOK) WithPayload(payload string) *URLNewOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the url new o k response
func (o *URLNewOK) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *URLNewOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*URLNewDefault Unexpected error.

swagger:response urlNewDefault
*/
type URLNewDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewURLNewDefault creates URLNewDefault with default headers values
func NewURLNewDefault(code int) *URLNewDefault {
	if code <= 0 {
		code = 500
	}

	return &URLNewDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the url new default response
func (o *URLNewDefault) WithStatusCode(code int) *URLNewDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the url new default response
func (o *URLNewDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the url new default response
func (o *URLNewDefault) WithPayload(payload *models.Error) *URLNewDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the url new default response
func (o *URLNewDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *URLNewDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}