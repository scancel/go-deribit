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

// GetPrivateGetCurrentDepositAddressReader is a Reader for the GetPrivateGetCurrentDepositAddress structure.
type GetPrivateGetCurrentDepositAddressReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPrivateGetCurrentDepositAddressReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPrivateGetCurrentDepositAddressOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetPrivateGetCurrentDepositAddressOK creates a GetPrivateGetCurrentDepositAddressOK with default headers values
func NewGetPrivateGetCurrentDepositAddressOK() *GetPrivateGetCurrentDepositAddressOK {
	return &GetPrivateGetCurrentDepositAddressOK{}
}

/*GetPrivateGetCurrentDepositAddressOK handles this case with default header values.

GetPrivateGetCurrentDepositAddressOK get private get current deposit address o k
*/
type GetPrivateGetCurrentDepositAddressOK struct {
	Payload *models.PrivateDepositAddressResponse
}

func (o *GetPrivateGetCurrentDepositAddressOK) Error() string {
	return fmt.Sprintf("[GET /private/get_current_deposit_address][%d] getPrivateGetCurrentDepositAddressOK  %+v", 200, o.Payload)
}

func (o *GetPrivateGetCurrentDepositAddressOK) GetPayload() *models.PrivateDepositAddressResponse {
	return o.Payload
}

func (o *GetPrivateGetCurrentDepositAddressOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PrivateDepositAddressResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
