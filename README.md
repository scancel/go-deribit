[![Build Status](https://travis-ci.com/adampointer/go-deribit.svg?branch=master)](https://travis-ci.com/adampointer/go-deribit)  [![Go Report Card](https://goreportcard.com/badge/github.com/adampointer/go-deribit)](https://goreportcard.com/report/github.com/adampointer/go-deribit)  [![codebeat badge](https://codebeat.co/badges/5bf32114-b7e1-4e70-91bf-fae2449fe2cb)](https://codebeat.co/projects/github-com-adampointer-go-deribit-master)

# go-deribit

## Forked code

This code has been forked from adampointer/go-deribit original code. Maintenance is now done on this fork (as of 20201216). This code adds encapsulated functions which allow a more user-friendly use of Adampointer's fantastic original API.

Here is a sample code:
`
package main

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/tuanito/go-deribit/v3"
	"github.com/tuanito/go-deribit/v3/client"
	"github.com/tuanito/go-deribit/v3/client/account_management"
	"github.com/tuanito/go-deribit/v3/client/market_data"
	"github.com/tuanito/go-deribit/v3/client/private"
	"github.com/tuanito/go-deribit/v3/client/public"
	"github.com/tuanito/go-deribit/v3/structures/account"
	summary "github.com/tuanito/go-deribit/v3/structures/book"
	"github.com/tuanito/go-deribit/v3/structures/instrument"
	"github.com/tuanito/go-deribit/v3/structures/position"
)

func main() {
	key := flag.String("access-key", "", "Access key")
	secret := flag.String("secret-key", "", "Secret access key")
	flag.Parse()
	errs := make(chan error)
	stop := make(chan bool)
	realExchange := false
	exchange, err := deribit.NewExchange(realExchange, errs, stop)

	if err != nil {
		log.Fatalf("Error creating connection: %s", err)
	}
	if err := exchange.Connect(); err != nil {
		log.Fatalf("Error connecting to exchange: %s", err)
	}
	go func() {
		err := <-errs
		stop <- true
		log.Fatalf("RPC error: %s", err)
	}()
	client := exchange.Client()
	// Step #1 : Hit the test public RPC endpoint
	res, err := client.Public.GetPublicTest(&public.GetPublicTestParams{})
	if err != nil {
		log.Fatalf("Error testing connection: %s", err)
	}
	log.Printf("Connected to Deribit API v%s", *res.Payload.Result.Version)
	// Step #2 : Actually authenticate on private account
	if err := exchange.Authenticate(*key, *secret); err != nil {
		log.Fatalf("Error authenticating: %s", err)
	}

	// --------------------
	// 1. Account Summary
	// --------------------

	BtcAccount := GetAccountSummary(client, "BTC")
	BtcAccount.Sprintf()
	EthAccount := GetAccountSummary(client, "ETH")
	EthAccount.Sprintf()

	// --------------------
	// 2. Referential
	// --------------------
	btcReferential := GetReferential(client, "BTC")
	ethReferential := GetReferential(client, "ETH")

	fmt.Println("---------------")
	fmt.Println("BTC Referential")
	fmt.Println("---------------")
	for key, value := range btcReferential {
		fmt.Println("Referential BTC:", key, ":\n", value.Sprintf())
	}
	fmt.Println("---------------")
	fmt.Println("ETH Referential")
	fmt.Println("---------------")
	for key, value := range ethReferential {
		fmt.Println("Referential ETH: ", key, ":\n", value.Sprintf())
	}

	// --------------------
	// 3. Book Summary
	// Makes a snapshot of prices on all derivatives
	// --------------------

	pCurrency := "BTC"
	bookSummaries := GetBookSummary(client, pCurrency)

	fmt.Println("Book summaries ", pCurrency)
	for key, value := range bookSummaries {
		fmt.Println("Book summary :", key, " ", value.Sprintf())
	}

	// This has not been tried yet
	/*
		// Buy
		buyParams := &private.GetPrivateBuyParams{
			Amount:         10,
			InstrumentName: "BTC-PERPETUAL",
			Type:           StrPointer("market"),
		}
		buy, err := client.Private.GetPrivateBuy(buyParams)
		if err != nil {
			log.Fatalf("Error submitting buy order: %s", err)
		}
		fmt.Printf("Bought at %f\n", buy.Payload.Result.Order.AveragePrice)
	*/

	// --------------------
	// 4. Account position
	// --------------------
	// see v3/models/position.go && v3/client/private/private_client.go
	pInstrumentName := "BTC-PERPETUAL"
	accountPosition := GetAccountPosition(client, pInstrumentName)

	fmt.Println("Position of ", pInstrumentName)
	accountPosition.Sprintf()

	// --------------------
	// 5. Account positions
	// --------------------
	futureType := "FUTURE" // option
	optionType := "OPTION" // option
	futuresPositions := GetAccountPositions(client, pCurrency, futureType)
	optionPositions := GetAccountPositions(client, pCurrency, optionType)

	for key, value := range futuresPositions {
		fmt.Println("FUTURES Position: #", key, " ", value.Sprintf())
	}

	for key, value := range optionPositions {
		fmt.Println("OPTION Position: #", key, " ", value.Sprintf())
	}

	// --------------------
	// 6. Order book subscription (market data)
	// --------------------
	depth := "1"
	interval := "100ms"
	book, err := exchange.SubscribeBookGroup("BTC-PERPETUAL", "none", depth, interval)
	if err != nil {
		log.Fatalf("Error subscribing to the book: %s", err)
	}
	for b := range book {
		fmt.Printf("Top bid: %f Top ask: %f\n", b.Bids[0][0], b.Asks[0][0])
	}

	exchange.Close()
}

// GetReferential takes the derivatives referential on each currency
func GetReferential(pClient *client.Deribit, pCurrency string) []instrument.Instrument {

	var instruments []instrument.Instrument

	referentialParams := &public.GetPublicGetInstrumentsParams{
		Currency: pCurrency,
	}
	referential, errRef := pClient.Public.GetPublicGetInstruments(referentialParams)
	if errRef != nil {
		log.Fatalf("Error loading referential: %s", errRef)
	}
	for _, value := range referential.Payload.Result {
		myInstrument := instrument.Instrument{
			BaseCurrency:        *value.BaseCurrency,
			InstrumentName:      *&value.InstrumentName,
			ContractSize:        *value.ContractSize,
			ExpirationTimestamp: *value.ExpirationTimestamp,
			ExpirationTime:      HumanReadableTime(*value.ExpirationTimestamp),
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
	return instruments
}

// GetAccountSummary gets all details from the trading account
// see v3/models/account_management.go
func GetAccountSummary(pClient *client.Deribit, pCurrency string) account.Account {
	var myAccount account.Account

	accountSummaryParams := &account_management.GetPrivateGetAccountSummaryParams{Currency: pCurrency}
	accountSummary, err := pClient.AccountManagement.GetPrivateGetAccountSummary(accountSummaryParams)

	if err != nil {
		log.Fatalf("Error getting account accountSummary: %s", err)
	}
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
		// SessionFunding: GetNumber(accountSummary.Payload.Result.SessionFunding), // <-- for some reasons, nil pointer crash
		FuturesPl: *accountSummary.Payload.Result.FuturesPl,
		OptionsPl: *accountSummary.Payload.Result.OptionsPl,
		TotalPl:   *accountSummary.Payload.Result.TotalPl,
	}

	return myAccount
}

// GetBookSummary gets the
func GetBookSummary(pClient *client.Deribit, pCurrency string) []summary.BookSummary {
	var bookSummaries []summary.BookSummary
	bookSummaryParams := &market_data.GetPublicGetBookSummaryByCurrencyParams{
		Currency: pCurrency,
	}

	bookSummary, errBook := pClient.MarketData.GetPublicGetBookSummaryByCurrency(bookSummaryParams)
	if errBook != nil {
		log.Fatalf("Error loading book summary : %s", errBook)
	}

	// See v3/models/book_summary.go
	for key, value := range bookSummary.Payload.Result {
		// fmt.Printf("Instrument raw price request : # %v \n %+v \n", key, value)
		fmt.Printf("Instrument price request : # %v \n", key)
		myBookSummary := summary.BookSummary{
			AskPrice:               GetNumber(value.AskPrice),
			BaseCurrency:           *value.BaseCurrency,
			BidPrice:               GetNumber(value.BidPrice),
			CreationTimestamp:      value.CreationTimestamp,
			CreationTime:           HumanReadableTime(int64(value.CreationTimestamp)),
			CurrentFunding:         value.CurrentFunding,
			EstimatedDeliveryPrice: value.EstimatedDeliveryPrice,
			Funding8h:              value.Funding8h,
			High:                   GetNumber(value.High),
			InstrumentName:         value.InstrumentName,

			InterestRate:    value.InterestRate,
			Last:            GetNumber(value.Last),
			Low:             GetNumber(value.Low),
			MarkPrice:       GetNumber(value.MarkPrice),
			MidPrice:        GetNumber(value.MidPrice),
			OpenInterest:    *value.OpenInterest,
			QuoteCurrency:   *value.QuoteCurrency,
			UnderlyingIndex: value.UnderlyingIndex,
			// underlying price for implied volatility calculations (options only)
			UnderlyingPrice: value.UnderlyingPrice,
			Volume:          *value.Volume,
			VolumeUsd:       value.VolumeUsd,
		}
		fmt.Println("HumanReadable : CreationTimestamp : ", myBookSummary.CreationTime, ":", HumanReadableTime(int64(myBookSummary.CreationTimestamp)))
		bookSummaries = append(bookSummaries, myBookSummary)
	}

	return bookSummaries
}

// GetAccountPosition retrives the position for a given underlying

func GetAccountPosition(pClient *client.Deribit, pInstrumentName string) position.Position {
	positionParams := &private.GetPrivateGetPositionParams{
		InstrumentName: pInstrumentName, // "BTC-PERPETUAL",
	}

	privatePosition, errPos := pClient.Private.GetPrivateGetPosition(positionParams)
	if errPos != nil {
		log.Fatalf("Error getting position : %s", errPos)
	}

	myPosition := position.Position{
		AveragePrice:              GetNumber(privatePosition.Payload.Result.AveragePrice),
		AveragePriceUsd:           privatePosition.Payload.Result.AveragePriceUsd,
		Delta:                     GetNumber(privatePosition.Payload.Result.Delta),
		Direction:                 privatePosition.Payload.Result.Direction,
		EstimatedLiquidationPrice: privatePosition.Payload.Result.EstimatedLiquidationPrice,
		FloatingProfitLoss:        GetNumber(privatePosition.Payload.Result.FloatingProfitLoss),
		FloatingProfitLossUsd:     privatePosition.Payload.Result.FloatingProfitLossUsd,
		IndexPrice:                GetNumber(privatePosition.Payload.Result.IndexPrice),
		InitialMargin:             GetNumber(privatePosition.Payload.Result.InitialMargin),
		InstrumentName:            privatePosition.Payload.Result.InstrumentName,
		Kind:                      privatePosition.Payload.Result.Kind,
		MaintenanceMargin:         GetNumber(privatePosition.Payload.Result.MaintenanceMargin),
		MarkPrice:                 GetNumber(privatePosition.Payload.Result.MarkPrice),
		OpenOrdersMargin:          GetNumber(privatePosition.Payload.Result.OpenOrdersMargin),
		RealizedProfitLoss:        GetNumber(privatePosition.Payload.Result.RealizedProfitLoss),
		SettlementPrice:           GetNumber(privatePosition.Payload.Result.SettlementPrice),
		Size:                      GetNumber(privatePosition.Payload.Result.Size),
		SizeCurrency:              privatePosition.Payload.Result.SizeCurrency,
		// TotalProfitLoss:           *privatePosition.TotalProfitLoss, // <-- should be here though BUG
	}
	return myPosition
}

// GetAccountPositions retrieves positions by derivatives types currency might be "BTC" or "ETH", Kind might be "future" or "option"
func GetAccountPositions(pClient *client.Deribit, pCurrency string, pKind string) []position.Position {
	// see v3/models/position.go && v3/client/private/private_client.go
	var outputPositions []position.Position
	positionsParams := &private.GetPrivateGetPositionsParams{
		Currency: pCurrency,
		Kind:     StrPointer(strings.ToLower(pKind)), // StrPointer("future"),
	}
	positions, errPos := pClient.Private.GetPrivateGetPositions(positionsParams)
	if errPos != nil {
		log.Fatalf("Error getting positions : %s", errPos)
	}
	for _, value := range positions.Payload.Result {
		// fmt.Printf("#%v Instrument :  %s %s %+v \n", key, strings.ToUpper(pKind), value.InstrumentName, value)

		myPosition := position.Position{
			AveragePrice:              GetNumber(value.AveragePrice),
			AveragePriceUsd:           value.AveragePriceUsd,
			Delta:                     GetNumber(value.Delta),
			Direction:                 value.Direction,
			EstimatedLiquidationPrice: value.EstimatedLiquidationPrice,
			FloatingProfitLoss:        GetNumber(value.FloatingProfitLoss),
			FloatingProfitLossUsd:     value.FloatingProfitLossUsd,
			IndexPrice:                GetNumber(value.IndexPrice),
			InitialMargin:             GetNumber(value.InitialMargin),
			InstrumentName:            value.InstrumentName,
			Kind:                      value.Kind,
			MaintenanceMargin:         GetNumber(value.MaintenanceMargin),
			MarkPrice:                 GetNumber(value.MarkPrice),
			OpenOrdersMargin:          GetNumber(value.OpenOrdersMargin),
			RealizedProfitLoss:        GetNumber(value.RealizedProfitLoss),
			SettlementPrice:           GetNumber(value.SettlementPrice),
			Size:                      GetNumber(value.Size),
			SizeCurrency:              value.SizeCurrency,
			// TotalProfitLoss:           *privatePosition.TotalProfitLoss, // <-- should be here though BUG
		}
		outputPositions = append(outputPositions, myPosition)
	}
	return outputPositions
}

// ----------
// Utilities
// ----------

func StrPointer(str string) *string {
	return &str
}

// GetNumber will return NaN or the number, depending on the pointer
func GetNumber(pNumber *float64) float64 {
	var output float64
	output = math.NaN()
	if pNumber != nil {
		output = *pNumber
	}
	return output
}

// HumanReadableTime takes a timestamp in milliseconds and converts it a human-readable time.
func HumanReadableTime(pTimestamp int64) string {
	tm := time.Unix(pTimestamp/1000, 0)
	output := tm.Format("2006-01-02 15:04:05")
	return output
}

// HumanReadableTime takes a timestamp in milliseconds and converts it a human-readable time.
// func HumanReadableTime2(pTimestamp Timestamp) string {
//	output := pTimestamp.Format("2006-01-02 15:04:05")
//	return output
//}
`

## V3 

This project is now using go1.13 with Go Modules, but should remain compatible with `dep`. Also, as there are some breaking changes introduced by the latest schema changes from the remote API, I have decided to carry on development in the new `v3` namespace with the project root containing the code tagged `v2.x`.

`import "github.com/adampointer/go-deribit/v3"`

We now have the latest API methods which were recently released such as `public/get_tradingview_chart_data`.

I recommend using the `v3` project in your projects as all onward development will now be within this project.

[GoDoc API Documentation](https://godoc.org/github.com/adampointer/go-deribit/v3)

## Overview

Go library for using the [Deribit's](https://www.deribit.com/reg-3027.8327) **v2** Websocket API. 

Deribit is a modern, fast BitCoin derivatives exchange. 

This library is a port of the [official wrapper libraries](https://github.com/deribit) to Go.

If you wish to try it out, be kind and use my affiliate link [https://www.deribit.com/reg-3027.8327](https://www.deribit.com/reg-3027.8327)

Or tip me!

btc: 3HmLfHJvrJuM48zHFY6HstUCxbuwV3dvxd

eth: 0x9Dc9129185E79211534D0039Af1C6f1ff585F5e3

ltc: MEpFjCdR3uXd6QjuJTSu3coLtcSWY3S2Hg


*p.s.* If you want a good BitMEX client library then try [go-bitmex](https://github.com/adampointer/go-bitmex)

[V2 API Documentation](http://docs.deribit.com/v2/?javascript#deribit-api-v2-0-0)

## Example Usage

Look at `v3/cmd/example/main.go`

```
go build main.go
./main --access-key XXX --secret-key YYYYYY
```

## Reconnect Behaviour

There is a heartbeat which triggers every 10 seconds to keep the websocket connection alive. In the event of a connection error, the library will automatically attempt to reconnect, re-authenticate and reestablish subscriptions.

This behaviour is overrideable with the `SetDisconnectHandler` method.

```
// Example reconnect code
exchange.SetDisconnectHandler(func (core *deribit.RPCCore) {
    exg := &deribit.NewExchangeFromCore(true, core)
	log.Warn("Disconnected from exchange. Attempting reconnection...")
	if err := exg.Connect(); err != nil {
		log.Fatalf("Error re-connecting to exchange: %s", err)
	}
	log.Info("Reconnected")
})
```

## Logging

The standard logger has been used within the library. You can plug this into your own application's logger by overriding the output io.Writer.

```
logger := logrus.New()
exchange.SetLogOutput(logger.Writer())
```

## Development

The `models` and `client` directories are where all the requests and responses are stored. The contents is automatically generated from the `schema` directory by `go-swagger`.

If you need to rebuild these use `make generate-models`.

The RPC subscriptions are also auto-generated. Use `make generate-methods` to rebuild these. They are in `rpc_subscriptions.go`.
