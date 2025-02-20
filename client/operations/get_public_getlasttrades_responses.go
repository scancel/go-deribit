// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/scancel/go-deribit/models"
)

// GetPublicGetlasttradesReader is a Reader for the GetPublicGetlasttrades structure.
type GetPublicGetlasttradesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPublicGetlasttradesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetPublicGetlasttradesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetPublicGetlasttradesOK creates a GetPublicGetlasttradesOK with default headers values
func NewGetPublicGetlasttradesOK() *GetPublicGetlasttradesOK {
	return &GetPublicGetlasttradesOK{}
}

/*GetPublicGetlasttradesOK handles this case with default header values.

foo
*/
type GetPublicGetlasttradesOK struct {
	Payload *models.PublicGetlasttradesResponse
}

func (o *GetPublicGetlasttradesOK) Error() string {
	return fmt.Sprintf("[GET /public/getlasttrades][%d] getPublicGetlasttradesOK  %+v", 200, o.Payload)
}

func (o *GetPublicGetlasttradesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PublicGetlasttradesResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
