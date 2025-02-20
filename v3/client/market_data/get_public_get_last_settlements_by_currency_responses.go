// Code generated by go-swagger; DO NOT EDIT.

package market_data

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/scancel/go-deribit/v3/models"
)

// GetPublicGetLastSettlementsByCurrencyReader is a Reader for the GetPublicGetLastSettlementsByCurrency structure.
type GetPublicGetLastSettlementsByCurrencyReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPublicGetLastSettlementsByCurrencyReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPublicGetLastSettlementsByCurrencyOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetPublicGetLastSettlementsByCurrencyOK creates a GetPublicGetLastSettlementsByCurrencyOK with default headers values
func NewGetPublicGetLastSettlementsByCurrencyOK() *GetPublicGetLastSettlementsByCurrencyOK {
	return &GetPublicGetLastSettlementsByCurrencyOK{}
}

/*GetPublicGetLastSettlementsByCurrencyOK handles this case with default header values.

GetPublicGetLastSettlementsByCurrencyOK get public get last settlements by currency o k
*/
type GetPublicGetLastSettlementsByCurrencyOK struct {
	Payload *models.PublicSettlementResponse
}

func (o *GetPublicGetLastSettlementsByCurrencyOK) Error() string {
	return fmt.Sprintf("[GET /public/get_last_settlements_by_currency][%d] getPublicGetLastSettlementsByCurrencyOK  %+v", 200, o.Payload)
}

func (o *GetPublicGetLastSettlementsByCurrencyOK) GetPayload() *models.PublicSettlementResponse {
	return o.Payload
}

func (o *GetPublicGetLastSettlementsByCurrencyOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PublicSettlementResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
