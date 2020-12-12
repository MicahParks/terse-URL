// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/MicahParks/terse-URL/models"
)

// URLDumpOKCode is the HTTP code returned for type URLDumpOK
const URLDumpOKCode int = 200

/*URLDumpOK url dump o k

swagger:response urlDumpOK
*/
type URLDumpOK struct {

	/*
	  In: Body
	*/
	Payload []*models.Dump `json:"body,omitempty"`
}

// NewURLDumpOK creates URLDumpOK with default headers values
func NewURLDumpOK() *URLDumpOK {

	return &URLDumpOK{}
}

// WithPayload adds the payload to the url dump o k response
func (o *URLDumpOK) WithPayload(payload []*models.Dump) *URLDumpOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the url dump o k response
func (o *URLDumpOK) SetPayload(payload []*models.Dump) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *URLDumpOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.Dump, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*URLDumpDefault Unexpected error.

swagger:response urlDumpDefault
*/
type URLDumpDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewURLDumpDefault creates URLDumpDefault with default headers values
func NewURLDumpDefault(code int) *URLDumpDefault {
	if code <= 0 {
		code = 500
	}

	return &URLDumpDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the url dump default response
func (o *URLDumpDefault) WithStatusCode(code int) *URLDumpDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the url dump default response
func (o *URLDumpDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the url dump default response
func (o *URLDumpDefault) WithPayload(payload *models.Error) *URLDumpDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the url dump default response
func (o *URLDumpDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *URLDumpDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}