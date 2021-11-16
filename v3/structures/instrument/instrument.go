package instrument

import (
	"strconv"

	"github.com/scancel/go-deribit/v3/models"
)

// Instrument is a slightly modified version of models.Instrument in v3/models/instrument.go
type Instrument struct {

	// The underlying currency being traded.
	// Required: true
	// Enum: [BTC ETH]
	BaseCurrency string

	// Contract size for instrument
	// Required: true
	ContractSize float64

	// The time when the instrument was first created (milliseconds)
	// Required: true
	CreationTimestamp int64
	CreationTime      string

	// The time when the instrument will expire (milliseconds)
	// Required: true
	ExpirationTimestamp int64
	ExpirationTime      string

	// instrument name
	// Required: true
	InstrumentName models.InstrumentName

	// Indicates if the instrument can currently be traded.
	// Required: true
	IsActive bool

	// kind
	// Required: true
	Kind models.Kind

	// Minimum amount for trading. For perpetual and futures - in USD units, for options it is amount of corresponding cryptocurrency contracts, e.g., BTC or ETH.
	// Required: true
	MinTradeAmount float64

	// The option type (only for options)
	// Enum: [call put]
	OptionType string

	// The currency in which the instrument prices are quoted.
	// Required: true
	// Enum: [USD]
	QuoteCurrency string

	// The settlement period.
	// Required: true
	// Enum: [month week perpetual]
	SettlementPeriod string

	// The strike value. (only for options)
	Strike float64

	// specifies minimal price change and, as follows, the number of decimal places for instrument prices
	// Required: true
	TickSize float64
}

// Sprintf exports a string
func (pInstrument Instrument) Sprintf() string {
	var output string
	output += "BaseCurrency: " + pInstrument.BaseCurrency + "\n"
	output += "InstrumentName: " + string(pInstrument.InstrumentName) + "\n"
	output += "ContractSize: " + strconv.FormatFloat(pInstrument.ContractSize, 'f', 2, 32) + "\n"
	output += "IsActive:" + strconv.FormatBool(pInstrument.IsActive) + "\n"
	output += "Kind: " + string(pInstrument.Kind) + "\n"
	output += "ExpirationTimestamp:" + strconv.FormatInt(pInstrument.ExpirationTimestamp, 10) + "\n"
	output += "ExpirationTime:" + pInstrument.ExpirationTime + "\n"
	output += "MinTradeAmount: " + strconv.FormatFloat(pInstrument.MinTradeAmount, 'f', 6, 32) + "\n"
	output += "OptionType: " + pInstrument.OptionType + "\n"
	output += "QuoteCurrency: " + pInstrument.QuoteCurrency + "\n"
	output += "SettlementPeriod: " + pInstrument.SettlementPeriod + "\n"
	output += "Strike: " + strconv.FormatFloat(pInstrument.Strike, 'f', 6, 32) + "\n"
	output += "TickSize: " + strconv.FormatFloat(pInstrument.TickSize, 'f', 6, 32) + "\n"

	return output
}
