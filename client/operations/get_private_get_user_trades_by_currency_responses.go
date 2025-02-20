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

// GetPrivateGetUserTradesByCurrencyReader is a Reader for the GetPrivateGetUserTradesByCurrency structure.
type GetPrivateGetUserTradesByCurrencyReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPrivateGetUserTradesByCurrencyReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetPrivateGetUserTradesByCurrencyOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetPrivateGetUserTradesByCurrencyOK creates a GetPrivateGetUserTradesByCurrencyOK with default headers values
func NewGetPrivateGetUserTradesByCurrencyOK() *GetPrivateGetUserTradesByCurrencyOK {
	return &GetPrivateGetUserTradesByCurrencyOK{}
}

/*GetPrivateGetUserTradesByCurrencyOK handles this case with default header values.

foo
*/
type GetPrivateGetUserTradesByCurrencyOK struct {
	Payload *models.UserTradesHistoryResponse
}

func (o *GetPrivateGetUserTradesByCurrencyOK) Error() string {
	return fmt.Sprintf("[GET /private/get_user_trades_by_currency][%d] getPrivateGetUserTradesByCurrencyOK  %+v", 200, o.Payload)
}

func (o *GetPrivateGetUserTradesByCurrencyOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.UserTradesHistoryResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
