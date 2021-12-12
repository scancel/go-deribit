package deribit

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/scancel/go-deribit/v3/client"
	"github.com/scancel/go-deribit/v3/client/account_management"
	"github.com/scancel/go-deribit/v3/client/market_data"
	"github.com/scancel/go-deribit/v3/client/private"
	"github.com/scancel/go-deribit/v3/client/public"
	"github.com/scancel/go-deribit/v3/client/trading"
	"github.com/scancel/go-deribit/v3/structures/account"
	"github.com/scancel/go-deribit/v3/structures/book"
	"github.com/scancel/go-deribit/v3/structures/instrument"
	"github.com/scancel/go-deribit/v3/structures/position"
	"github.com/scancel/go-deribit/v3/utils"
)

// Renew login 10 minutes before we have to
const renewBefore int64 = 600

// Keep clientID and clientSecret for re-connecting automatically
var (
	clientID     string
	clientSecret string
)

// Authenticate : 1st param: API Key, 2nd param: API Secret
func (e *Exchange) Authenticate(keys ...string) error {
	if len(keys) == 2 {
		clientID = keys[0]
		clientSecret = keys[1]
	} else {
		if clientID == "" || clientSecret == "" {
			return fmt.Errorf("API Key and Secret must be provided")
		}
	}
	auth, err := e.Client().Public.GetPublicAuth(&public.GetPublicAuthParams{ClientID: clientID, ClientSecret: clientSecret, GrantType: "client_credentials"})
	if err != nil {
		return fmt.Errorf("error authenticating: %s", err)
	}
	e.auth = auth.Payload
	d, err := time.ParseDuration(fmt.Sprintf("%ds", *(e.auth.Result.ExpiresIn)-renewBefore))
	if err != nil {
		return fmt.Errorf("unable to parse %ds as a duration: %s", *(e.auth.Result.ExpiresIn)-renewBefore, err)
	}
	go e.refreshAuth(d)
	return nil
}

func (e *Exchange) refreshAuth(d time.Duration) {
	time.Sleep(d)
	auth, err := e.Client().Public.GetPublicAuth(&public.GetPublicAuthParams{RefreshToken: *e.auth.Result.RefreshToken, GrantType: "refresh_token"})
	if err != nil {
		e.errors <- fmt.Errorf("error authenticating: %s", err)
	}
	e.auth = auth.Payload
}

// GetReferential takes the derivatives referential on each currency
func GetReferential(pClient *client.Deribit, pCurrency string) ([]instrument.Instrument, error) {

	var instruments []instrument.Instrument

	referentialParams := &public.GetPublicGetInstrumentsParams{
		Currency: pCurrency,
	}
	referential, errRef := pClient.Public.GetPublicGetInstruments(referentialParams)
	// if errRef != nil {
	//	log.Fatalf("Error loading referential: %s", errRef)
	//}
	for _, value := range referential.Payload.Result {
		myInstrument := instrument.Instrument{
			BaseCurrency:        *value.BaseCurrency,
			InstrumentName:      *&value.InstrumentName,
			ContractSize:        *value.ContractSize,
			ExpirationTimestamp: *value.ExpirationTimestamp,
			ExpirationTime:      utils.HumanReadableTime(*value.ExpirationTimestamp),
			IsActive:            *value.IsActive,
			Kind:                value.Kind,
			MinTradeAmount:      *value.MinTradeAmount,
			OptionType:          value.OptionType,
			QuoteCurrency:       *value.QuoteCurrency,
			SettlementPeriod:    *value.SettlementPeriod,
			Strike:              value.Strike,
			TickSize:            *value.TickSize,
		}
		instruments = append(instruments, myInstrument)
	}
	return instruments, errRef
}

// GetAccountSummary gets all details from the trading account
// see v3/models/account_management.go
func GetAccountSummary(pClient *client.Deribit, pCurrency string) (account.Account, error) {
	var myAccount account.Account

	accountSummaryParams := &account_management.GetPrivateGetAccountSummaryParams{Currency: pCurrency}
	accountSummary, err := pClient.AccountManagement.GetPrivateGetAccountSummary(accountSummaryParams)

	myAccount = account.Account{
		Currency:       *accountSummary.Payload.Result.Currency,
		AvailableFunds: *accountSummary.Payload.Result.AvailableFunds,
		Balance:        *accountSummary.Payload.Result.Balance,
		Equity:         *accountSummary.Payload.Result.Equity,
		DeltaTotal:     *accountSummary.Payload.Result.DeltaTotal,
		OptionsDelta:   *accountSummary.Payload.Result.OptionsDelta,
		OptionsGamma:   *accountSummary.Payload.Result.OptionsGamma,
		OptionsVega:    *accountSummary.Payload.Result.OptionsVega,
		OptionsTheta:   *accountSummary.Payload.Result.OptionsTheta,
		// SessionFunding: utils.GetNumber(accountSummary.Payload.Result.SessionFunding), // <-- for some reasons, nil pointer crash
		FuturesPl: *accountSummary.Payload.Result.FuturesPl,
		OptionsPl: *accountSummary.Payload.Result.OptionsPl,
		TotalPl:   *accountSummary.Payload.Result.TotalPl,
	}

	return myAccount, err
}

// GetBookSummary gets the
func GetBookSummary(pClient *client.Deribit, pCurrency string) ([]book.BookSummary, error) {
	var bookSummaries []book.BookSummary
	bookSummaryParams := &market_data.GetPublicGetBookSummaryByCurrencyParams{
		Currency: pCurrency,
	}

	bookSummary, errBook := pClient.MarketData.GetPublicGetBookSummaryByCurrency(bookSummaryParams)
	// if errBook != nil {
	//	log.Fatalf("Error loading book summary : %s", errBook)
	//}

	// See v3/models/book_summary.go
	for _, value := range bookSummary.Payload.Result {
		// fmt.Printf("Instrument raw price request : # %v \n %+v \n", key, value)
		// fmt.Printf("Instrument price request : # %v \n", key)
		myBookSummary := book.BookSummary{
			AskPrice:               utils.GetNumber(value.AskPrice),
			BaseCurrency:           *value.BaseCurrency,
			BidPrice:               utils.GetNumber(value.BidPrice),
			CreationTimestamp:      value.CreationTimestamp,
			CreationTime:           utils.HumanReadableTime(int64(value.CreationTimestamp)),
			CurrentFunding:         value.CurrentFunding,
			EstimatedDeliveryPrice: value.EstimatedDeliveryPrice,
			Funding8h:              value.Funding8h,
			High:                   utils.GetNumber(value.High),
			InstrumentName:         value.InstrumentName,

			InterestRate:    value.InterestRate,
			Last:            utils.GetNumber(value.Last),
			Low:             utils.GetNumber(value.Low),
			MarkPrice:       utils.GetNumber(value.MarkPrice),
			MidPrice:        utils.GetNumber(value.MidPrice),
			OpenInterest:    *value.OpenInterest,
			QuoteCurrency:   *value.QuoteCurrency,
			UnderlyingIndex: value.UnderlyingIndex,
			// underlying price for implied volatility calculations (options only)
			UnderlyingPrice: value.UnderlyingPrice,
			Volume:          *value.Volume,
			VolumeUsd:       value.VolumeUsd,
		}
		// fmt.Println("utils.HumanReadable : CreationTimestamp : ", myBookSummary.CreationTime, ":", utils.HumanReadableTime(int64(myBookSummary.CreationTimestamp)))
		bookSummaries = append(bookSummaries, myBookSummary)
	}

	return bookSummaries, errBook
}

// GetAccountPosition retrives the position for a given underlying
func GetAccountPosition(pClient *client.Deribit, pInstrumentName string) (position.Position, error) {
	positionParams := &private.GetPrivateGetPositionParams{
		InstrumentName: pInstrumentName, // "BTC-PERPETUAL",
	}

	privatePosition, errPos := pClient.Private.GetPrivateGetPosition(positionParams)
	if errPos != nil {
		log.Fatalf("Error getting position : %s", errPos)
	}

	myPosition := position.Position{
		AveragePrice:              utils.GetNumber(privatePosition.Payload.Result.AveragePrice),
		AveragePriceUsd:           privatePosition.Payload.Result.AveragePriceUsd,
		Delta:                     utils.GetNumber(privatePosition.Payload.Result.Delta),
		Direction:                 privatePosition.Payload.Result.Direction,
		EstimatedLiquidationPrice: privatePosition.Payload.Result.EstimatedLiquidationPrice,
		FloatingProfitLoss:        utils.GetNumber(privatePosition.Payload.Result.FloatingProfitLoss),
		FloatingProfitLossUsd:     privatePosition.Payload.Result.FloatingProfitLossUsd,
		IndexPrice:                utils.GetNumber(privatePosition.Payload.Result.IndexPrice),
		InitialMargin:             utils.GetNumber(privatePosition.Payload.Result.InitialMargin),
		InstrumentName:            privatePosition.Payload.Result.InstrumentName,
		Kind:                      privatePosition.Payload.Result.Kind,
		MaintenanceMargin:         utils.GetNumber(privatePosition.Payload.Result.MaintenanceMargin),
		MarkPrice:                 utils.GetNumber(privatePosition.Payload.Result.MarkPrice),
		OpenOrdersMargin:          utils.GetNumber(privatePosition.Payload.Result.OpenOrdersMargin),
		RealizedProfitLoss:        utils.GetNumber(privatePosition.Payload.Result.RealizedProfitLoss),
		SettlementPrice:           utils.GetNumber(privatePosition.Payload.Result.SettlementPrice),
		Size:                      utils.GetNumber(privatePosition.Payload.Result.Size),
		SizeCurrency:              privatePosition.Payload.Result.SizeCurrency,
		// TotalProfitLoss:           *privatePosition.TotalProfitLoss, // <-- should be here though BUG
	}
	return myPosition, errPos
}

// GetAccountPositions retrieves positions by derivatives types currency might be "BTC" or "ETH", Kind might be "future" or "option"
func GetAccountPositions(pClient *client.Deribit, pCurrency string, pKind string) ([]position.Position, error) {
	// see v3/models/position.go && v3/client/private/private_client.go
	var outputPositions []position.Position
	positionsParams := &private.GetPrivateGetPositionsParams{
		Currency: pCurrency,
		Kind:     utils.StrPointer(strings.ToLower(pKind)), // utils.StrPointer("future"),
	}
	positions, errPos := pClient.Private.GetPrivateGetPositions(positionsParams)
	// if errPos != nil {
	//	log.Fatalf("Error getting positions : %s", errPos)
	//}
	for _, value := range positions.Payload.Result {
		// fmt.Printf("#%v Instrument :  %s %s %+v \n", key, strings.ToUpper(pKind), value.InstrumentName, value)

		myPosition := position.Position{
			AveragePrice:              utils.GetNumber(value.AveragePrice),
			AveragePriceUsd:           value.AveragePriceUsd,
			Delta:                     utils.GetNumber(value.Delta),
			Direction:                 value.Direction,
			EstimatedLiquidationPrice: value.EstimatedLiquidationPrice,
			FloatingProfitLoss:        utils.GetNumber(value.FloatingProfitLoss),
			FloatingProfitLossUsd:     value.FloatingProfitLossUsd,
			IndexPrice:                utils.GetNumber(value.IndexPrice),
			InitialMargin:             utils.GetNumber(value.InitialMargin),
			InstrumentName:            value.InstrumentName,
			Kind:                      value.Kind,
			MaintenanceMargin:         utils.GetNumber(value.MaintenanceMargin),
			MarkPrice:                 utils.GetNumber(value.MarkPrice),
			OpenOrdersMargin:          utils.GetNumber(value.OpenOrdersMargin),
			RealizedProfitLoss:        utils.GetNumber(value.RealizedProfitLoss),
			SettlementPrice:           utils.GetNumber(value.SettlementPrice),
			Size:                      utils.GetNumber(value.Size),
			SizeCurrency:              value.SizeCurrency,
			// TotalProfitLoss:           *privatePosition.TotalProfitLoss, // <-- should be here though BUG
		}
		outputPositions = append(outputPositions, myPosition)
	}
	return outputPositions, errPos
}

func Edit(pClient *client.Deribit, pOrderID string, pAmount float64, pPrice float64) (*trading.GetPrivateEditOK, error) {
	editParams := &trading.GetPrivateEditParams{
		OrderID: 		pOrderID,
		Amount:         pAmount,
		Price:          pPrice,
	}
	edit, err := pClient.Trading.GetPrivateEdit(editParams)
	if err != nil {
		log.Fatalf("Error submitting edit order: %s", err)
	}
	fmt.Printf("BOT %s %f %f at %f\n", pOrderID, pAmount, pPrice, edit.Payload.Result.Order.AveragePrice)
	return edit, err
}

// Buy a derivatives : example : Buy(10, "BTC-PERPETUAL", "market")
// This has not been tried yet.
// pOrderType : "limit" or "market"
// Responses from the Deribit API :
// order_state : "open", "filled", "rejected", "cancelled", "untriggered"
// Please note that any partial fill is a "filled" for the relevant quantity
// Note that the answer from Deribit contains fields that can probably be implemented as an input :
// -"advanced" : for options orders, can be in currency ("usd") or in volatility level ("implv")
// -"time_in_force" : "good_til_cancelled", "fill_or_kill", "immediate_or_cancel"
// -"post_only" : bool true/false
// -"filled_amount" : returns the filled amount
// These answers from Deribit cannot be inputs but are still relevant :
// -"is_liquidation" : bool true/false
func Buy(pClient *client.Deribit, pAmount float64, pInstrumentName string, pPrice float64, pOrderType string) (*private.GetPrivateBuyOK, error) {
	buyParams := &private.GetPrivateBuyParams{
		Amount:         pAmount,
		InstrumentName: pInstrumentName,
		Price:          utils.FloatPointer(pPrice),
		Type:           utils.StrPointer(strings.ToLower(pOrderType)),
		Advanced:       utils.StrPointer("usd"),
	}
	buy, err := pClient.Private.GetPrivateBuy(buyParams)
	if err != nil {
		log.Fatalf("Error submitting buy order: %s", err)
	}
	fmt.Printf("BOT %f %s at %f\n", pAmount, pInstrumentName, buy.Payload.Result.Order.Price)
	return buy, err
}

// Sell is untested
func Sell(pClient *client.Deribit, pAmount float64, pInstrumentName string, pPrice float64, pOrderType string) (*private.GetPrivateSellOK, error) {
	sellParams := &private.GetPrivateSellParams{
		Amount:         pAmount,
		InstrumentName: pInstrumentName,
		Price:          utils.FloatPointer(pPrice),
		Type:           utils.StrPointer(strings.ToLower(pOrderType)),
		Advanced:       utils.StrPointer("usd"),
	}
	sell, err := pClient.Private.GetPrivateSell(sellParams)
	if err != nil {
		log.Fatalf("Error submitting sell order: %s", err)
	}
	fmt.Printf("SOLD %f %s at %f\n", pAmount, pInstrumentName, sell.Payload.Result.Order.Price)
	return sell, err
}
