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

// GetPrivateToggleDepositAddressCreationReader is a Reader for the GetPrivateToggleDepositAddressCreation structure.
type GetPrivateToggleDepositAddressCreationReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPrivateToggleDepositAddressCreationReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetPrivateToggleDepositAddressCreationOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetPrivateToggleDepositAddressCreationOK creates a GetPrivateToggleDepositAddressCreationOK with default headers values
func NewGetPrivateToggleDepositAddressCreationOK() *GetPrivateToggleDepositAddressCreationOK {
	return &GetPrivateToggleDepositAddressCreationOK{}
}

/*GetPrivateToggleDepositAddressCreationOK handles this case with default header values.

foo
*/
type GetPrivateToggleDepositAddressCreationOK struct {
	Payload *models.OkResponse
}

func (o *GetPrivateToggleDepositAddressCreationOK) Error() string {
	return fmt.Sprintf("[GET /private/toggle_deposit_address_creation][%d] getPrivateToggleDepositAddressCreationOK  %+v", 200, o.Payload)
}

func (o *GetPrivateToggleDepositAddressCreationOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.OkResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
