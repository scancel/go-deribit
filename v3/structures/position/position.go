package position

import (
	"strconv"
	"strings"

	"github.com/tuanito/go-deribit/v3/models"
)

// Position is almost the same struct as in v3/models/position.go
type Position struct {

	// Average price of trades that built this position
	// Required: true
	AveragePrice float64

	// Only for options, average price in USD
	AveragePriceUsd float64

	// Delta parameter
	// Required: true
	Delta float64

	// direction
	// Required: true
	Direction models.Direction

	// Only for futures, estimated liquidation price
	EstimatedLiquidationPrice float64

	// Floating profit or loss
	// Required: true
	FloatingProfitLoss float64

	// Only for options, floating profit or loss in USD
	FloatingProfitLossUsd float64

	// Current index price
	// Required: true
	IndexPrice float64

	// Initial margin
	// Required: true
	InitialMargin float64

	// instrument name
	// Required: true
	InstrumentName models.InstrumentName

	// kind
	// Required: true
	Kind models.Kind

	// Maintenance margin
	// Required: true
	MaintenanceMargin float64

	// Current mark price for position's instrument
	// Required: true
	MarkPrice float64

	// Open orders margin
	// Required: true
	OpenOrdersMargin float64

	// Realized profit or loss
	// Required: true
	RealizedProfitLoss float64

	// Last settlement price for position's instrument 0 if instrument wasn't settled yet
	// Required: true
	SettlementPrice float64

	// Position size for futures size in quote currency (e.g. USD), for options size is in base currency (e.g. BTC)
	// Required: true
	Size float64

	// Only for futures, position size in base currency
	SizeCurrency float64

	// Profit or loss from position
	// Required: true
	TotalProfitLoss float64
}

// Sprintf exports a string
func (pPosition Position) Sprintf() string {
	var output string
	output += "AveragePrice: " + strconv.FormatFloat(pPosition.AveragePrice, 'f', 6, 32) + "\n"
	output += "AveragePriceUsd: " + strconv.FormatFloat(pPosition.AveragePriceUsd, 'f', 6, 32) + "\n"
	output += "Delta: " + strconv.FormatFloat(pPosition.Delta, 'f', 2, 32) + "\n"
	output += "Direction:" + strings.ToUpper(string(pPosition.Direction)) + "\n"
	output += "EstimatedLiquidationPrice: " + strconv.FormatFloat(pPosition.EstimatedLiquidationPrice, 'f', 6, 32) + "\n"
	output += "FloatingProfitLoss:" + strconv.FormatFloat(pPosition.FloatingProfitLoss, 'f', 2, 32) + "\n"
	output += "FloatingProfitLossUsd: " + strconv.FormatFloat(pPosition.FloatingProfitLossUsd, 'f', 2, 32) + "\n"
	output += "IndexPrice:" + strconv.FormatFloat(pPosition.IndexPrice, 'f', 6, 32) + "\n"
	output += "InitialMargin: " + strconv.FormatFloat(pPosition.InitialMargin, 'f', 6, 32) + "\n"
	output += "InstrumentName: " + strings.ToUpper(string(pPosition.InstrumentName)) + "\n"
	output += "Kind: " + string(pPosition.Kind) + "\n"
	output += "MaintenanceMargin: " + strconv.FormatFloat(pPosition.MaintenanceMargin, 'f', 6, 32) + "\n"
	output += "MarkPrice: " + strconv.FormatFloat(pPosition.MarkPrice, 'f', 6, 32) + "\n"
	output += "OpenOrdersMargin: " + strconv.FormatFloat(pPosition.OpenOrdersMargin, 'f', 6, 32) + "\n"
	output += "RealizedProfitLoss: " + strconv.FormatFloat(pPosition.RealizedProfitLoss, 'f', 6, 32) + "\n"
	output += "SettlementPrice: " + strconv.FormatFloat(pPosition.SettlementPrice, 'f', 6, 32) + "\n"
	output += "Size: " + strconv.FormatFloat(pPosition.Size, 'f', 6, 32) + "\n"
	output += "SizeCurrency: " + strconv.FormatFloat(pPosition.SizeCurrency, 'f', 6, 32) + "\n"
	output += "TotalProfitLoss: " + strconv.FormatFloat(pPosition.TotalProfitLoss, 'f', 6, 32) + "\n"

	return output
}
