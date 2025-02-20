// Code generated by go-swagger; DO NOT EDIT.

package wallet

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/scancel/go-deribit/v3/models"
)

// GetPrivateCancelWithdrawalReader is a Reader for the GetPrivateCancelWithdrawal structure.
type GetPrivateCancelWithdrawalReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPrivateCancelWithdrawalReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPrivateCancelWithdrawalOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetPrivateCancelWithdrawalOK creates a GetPrivateCancelWithdrawalOK with default headers values
func NewGetPrivateCancelWithdrawalOK() *GetPrivateCancelWithdrawalOK {
	return &GetPrivateCancelWithdrawalOK{}
}

/*GetPrivateCancelWithdrawalOK handles this case with default header values.

GetPrivateCancelWithdrawalOK get private cancel withdrawal o k
*/
type GetPrivateCancelWithdrawalOK struct {
	Payload *models.PrivateWithdrawResponse
}

func (o *GetPrivateCancelWithdrawalOK) Error() string {
	return fmt.Sprintf("[GET /private/cancel_withdrawal][%d] getPrivateCancelWithdrawalOK  %+v", 200, o.Payload)
}

func (o *GetPrivateCancelWithdrawalOK) GetPayload() *models.PrivateWithdrawResponse {
	return o.Payload
}

func (o *GetPrivateCancelWithdrawalOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PrivateWithdrawResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
