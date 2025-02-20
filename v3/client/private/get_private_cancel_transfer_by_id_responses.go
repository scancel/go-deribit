// Code generated by go-swagger; DO NOT EDIT.

package private

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/scancel/go-deribit/v3/models"
)

// GetPrivateCancelTransferByIDReader is a Reader for the GetPrivateCancelTransferByID structure.
type GetPrivateCancelTransferByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPrivateCancelTransferByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPrivateCancelTransferByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetPrivateCancelTransferByIDOK creates a GetPrivateCancelTransferByIDOK with default headers values
func NewGetPrivateCancelTransferByIDOK() *GetPrivateCancelTransferByIDOK {
	return &GetPrivateCancelTransferByIDOK{}
}

/*GetPrivateCancelTransferByIDOK handles this case with default header values.

ok response
*/
type GetPrivateCancelTransferByIDOK struct {
	Payload *models.PrivateSubmitTransferResponse
}

func (o *GetPrivateCancelTransferByIDOK) Error() string {
	return fmt.Sprintf("[GET /private/cancel_transfer_by_id][%d] getPrivateCancelTransferByIdOK  %+v", 200, o.Payload)
}

func (o *GetPrivateCancelTransferByIDOK) GetPayload() *models.PrivateSubmitTransferResponse {
	return o.Payload
}

func (o *GetPrivateCancelTransferByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PrivateSubmitTransferResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
