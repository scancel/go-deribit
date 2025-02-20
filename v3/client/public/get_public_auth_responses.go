// Code generated by go-swagger; DO NOT EDIT.

package public

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/scancel/go-deribit/v3/models"
)

// GetPublicAuthReader is a Reader for the GetPublicAuth structure.
type GetPublicAuthReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPublicAuthReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPublicAuthOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewGetPublicAuthTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetPublicAuthOK creates a GetPublicAuthOK with default headers values
func NewGetPublicAuthOK() *GetPublicAuthOK {
	return &GetPublicAuthOK{}
}

/*GetPublicAuthOK handles this case with default header values.

ok response
*/
type GetPublicAuthOK struct {
	Payload *models.PublicAuthResponse
}

func (o *GetPublicAuthOK) Error() string {
	return fmt.Sprintf("[GET /public/auth][%d] getPublicAuthOK  %+v", 200, o.Payload)
}

func (o *GetPublicAuthOK) GetPayload() *models.PublicAuthResponse {
	return o.Payload
}

func (o *GetPublicAuthOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PublicAuthResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPublicAuthTooManyRequests creates a GetPublicAuthTooManyRequests with default headers values
func NewGetPublicAuthTooManyRequests() *GetPublicAuthTooManyRequests {
	return &GetPublicAuthTooManyRequests{}
}

/*GetPublicAuthTooManyRequests handles this case with default header values.

over limit
*/
type GetPublicAuthTooManyRequests struct {
	Payload *models.ErrorMessage
}

func (o *GetPublicAuthTooManyRequests) Error() string {
	return fmt.Sprintf("[GET /public/auth][%d] getPublicAuthTooManyRequests  %+v", 429, o.Payload)
}

func (o *GetPublicAuthTooManyRequests) GetPayload() *models.ErrorMessage {
	return o.Payload
}

func (o *GetPublicAuthTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorMessage)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
