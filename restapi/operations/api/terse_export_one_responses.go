// Code generated by go-swagger; DO NOT EDIT.

package api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/MicahParks/terse-URL/models"
)

// TerseExportOneOKCode is the HTTP code returned for type TerseExportOneOK
const TerseExportOneOKCode int = 200

/*TerseExportOneOK The export was successfully retrieved.

swagger:response terseExportOneOK
*/
type TerseExportOneOK struct {

	/*
	  In: Body
	*/
	Payload *models.Export `json:"body,omitempty"`
}

// NewTerseExportOneOK creates TerseExportOneOK with default headers values
func NewTerseExportOneOK() *TerseExportOneOK {

	return &TerseExportOneOK{}
}

// WithPayload adds the payload to the terse export one o k response
func (o *TerseExportOneOK) WithPayload(payload *models.Export) *TerseExportOneOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the terse export one o k response
func (o *TerseExportOneOK) SetPayload(payload *models.Export) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *TerseExportOneOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*TerseExportOneDefault Unexpected error.

swagger:response terseExportOneDefault
*/
type TerseExportOneDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewTerseExportOneDefault creates TerseExportOneDefault with default headers values
func NewTerseExportOneDefault(code int) *TerseExportOneDefault {
	if code <= 0 {
		code = 500
	}

	return &TerseExportOneDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the terse export one default response
func (o *TerseExportOneDefault) WithStatusCode(code int) *TerseExportOneDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the terse export one default response
func (o *TerseExportOneDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the terse export one default response
func (o *TerseExportOneDefault) WithPayload(payload *models.Error) *TerseExportOneDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the terse export one default response
func (o *TerseExportOneDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *TerseExportOneDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}