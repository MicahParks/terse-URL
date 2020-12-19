// Code generated by go-swagger; DO NOT EDIT.

package api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/MicahParks/terse-URL/models"
)

// TerseVisitsOKCode is the HTTP code returned for type TerseVisitsOK
const TerseVisitsOKCode int = 200

/*TerseVisitsOK The Visits data was successfully retrieved.

swagger:response terseVisitsOK
*/
type TerseVisitsOK struct {

	/*The visit data for a single shortened URL.
	  In: Body
	*/
	Payload []*models.Visit `json:"body,omitempty"`
}

// NewTerseVisitsOK creates TerseVisitsOK with default headers values
func NewTerseVisitsOK() *TerseVisitsOK {

	return &TerseVisitsOK{}
}

// WithPayload adds the payload to the terse visits o k response
func (o *TerseVisitsOK) WithPayload(payload []*models.Visit) *TerseVisitsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the terse visits o k response
func (o *TerseVisitsOK) SetPayload(payload []*models.Visit) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *TerseVisitsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.Visit, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*TerseVisitsDefault Unexpected error.

swagger:response terseVisitsDefault
*/
type TerseVisitsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewTerseVisitsDefault creates TerseVisitsDefault with default headers values
func NewTerseVisitsDefault(code int) *TerseVisitsDefault {
	if code <= 0 {
		code = 500
	}

	return &TerseVisitsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the terse visits default response
func (o *TerseVisitsDefault) WithStatusCode(code int) *TerseVisitsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the terse visits default response
func (o *TerseVisitsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the terse visits default response
func (o *TerseVisitsDefault) WithPayload(payload *models.Error) *TerseVisitsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the terse visits default response
func (o *TerseVisitsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *TerseVisitsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
