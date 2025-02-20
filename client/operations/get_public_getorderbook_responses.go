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

// GetPublicGetorderbookReader is a Reader for the GetPublicGetorderbook structure.
type GetPublicGetorderbookReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPublicGetorderbookReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetPublicGetorderbookOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetPublicGetorderbookOK creates a GetPublicGetorderbookOK with default headers values
func NewGetPublicGetorderbookOK() *GetPublicGetorderbookOK {
	return &GetPublicGetorderbookOK{}
}

/*GetPublicGetorderbookOK handles this case with default header values.

foo
*/
type GetPublicGetorderbookOK struct {
	Payload *models.PublicGetorderbookResponse
}

func (o *GetPublicGetorderbookOK) Error() string {
	return fmt.Sprintf("[GET /public/getorderbook][%d] getPublicGetorderbookOK  %+v", 200, o.Payload)
}

func (o *GetPublicGetorderbookOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PublicGetorderbookResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
