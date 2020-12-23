package account

import (
	"strconv"
)

type Account struct {
	// The account's available funds
	// Required: true
	AvailableFunds float64
	// The account's available to withdrawal funds
	// Required: true
	AvailableWithdrawalFunds float64
	// The account's balance
	// Required: true
	Balance float64

	// The selected currency
	// Required: true
	Currency string // *string

	// The sum of position deltas (currently bugged : miscomputed by Deribit)
	// Required: true
	DeltaTotal float64

	// The deposit address for the account (if available)
	DepositAddress string

	// User email (available when parameter `extended` = `true`)
	// Required: true
	Email string //*string

	// The account's current equity
	// Required: true
	Equity float64 // *float64 `json:"equity"`

	// Futures profit and Loss
	// Required: true
	FuturesPl float64 // *float64

	// Futures session realized profit and Loss
	// Required: true
	FuturesSessionRpl float64 // *float64

	// Futures session unrealized profit and Loss
	// Required: true
	FuturesSessionUpl float64 // *float64

	// Account id (available when parameter `extended` = `true`)
	// Required: true
	ID int64 // *int64

	// The account's initial margin
	// Required: true
	InitialMargin float64 // *float64

	// The maintenance margin.
	// Required: true
	MaintenanceMargin float64 //*float64

	// The account's margin balance
	MarginBalance float64

	// Options summary delta
	// Required: true
	OptionsDelta float64 // *float64

	// Options summary gamma
	// Required: true
	OptionsGamma float64 // *float64

	// Options profit and Loss
	// Required: true
	OptionsPl float64 // *float64

	// Options session realized profit and Loss
	// Required: true
	OptionsSessionRpl float64 //*float64

	// Options session unrealized profit and Loss
	// Required: true
	OptionsSessionUpl float64 //*float64

	// Options summary theta
	// Required: true
	OptionsTheta float64 // *float64

	// Options summary vega
	// Required: true
	OptionsVega float64 // *float64 `

	// `true` when portfolio margining is enabled for user
	PortfolioMarginingEnabled bool

	// Projected initial margin (for portfolio margining users)
	ProjectedInitialMargin float64

	// Projected maintenance margin (for portfolio margining users)
	ProjectedMaintenanceMargin float64

	// Session funding
	// Required: true
	SessionFunding float64 // *float64

	// Session realized profit and loss
	// Required: true
	SessionRpl float64 // *float64

	// Session unrealized profit and loss
	// Required: true
	SessionUpl float64 //*float64

	// System generated user nickname (available when parameter `extended` = `true`)
	// Required: true
	SystemName string // *string

	// Whether two factor authentication is enabled (available when parameter `extended` = `true`)
	// Required: true
	TfaEnabled bool // *bool

	// Profit and loss
	// Required: true
	TotalPl float64 // *float64

	// Account type (available when parameter `extended` = `true`)
	// Required: true
	// Enum: [main subaccount]
	Type string // *string

	// Account name (given by user) (available when parameter `extended` = `true`)
	// Required: true
	Username string // *string
}

func (pAccount Account) Sprintf() string {
	var output string
	output = "Currency: " + pAccount.Currency + "\n"
	output += "Available funds:" + strconv.FormatFloat(pAccount.AvailableFunds, 'f', 2, 32) + "\n"
	output += "Balance: " + strconv.FormatFloat(pAccount.Balance, 'f', 2, 32) + "\n"
	output += "Equity : " + strconv.FormatFloat(pAccount.Equity, 'f', 2, 32) + "\n"
	output += "Delta Total: " + strconv.FormatFloat(pAccount.DeltaTotal, 'f', 2, 32) + "\n"
	output += "Options Delta: " + strconv.FormatFloat(pAccount.OptionsDelta, 'f', 2, 32) + "\n"
	output += "Options Gamma: " + strconv.FormatFloat(pAccount.OptionsGamma, 'f', 2, 32) + "\n"
	output += "Options Vega: " + strconv.FormatFloat(pAccount.OptionsVega, 'f', 2, 32) + "\n"
	output += "Options Theta: " + strconv.FormatFloat(pAccount.OptionsTheta, 'f', 2, 32) + "\n"
	// fmt.Printf("Session funding: %f\n", (*pAccount).SessionFunding)
	output += "Futures PnL: " + strconv.FormatFloat(pAccount.FuturesPl, 'f', 2, 32) + "\n"
	output += "Options PnL: " + strconv.FormatFloat(pAccount.OptionsPl, 'f', 2, 32) + "\n"
	output += "Total PnL: " + strconv.FormatFloat(pAccount.TotalPl, 'f', 2, 32) + "\n"
	return output
}
