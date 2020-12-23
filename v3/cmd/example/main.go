package main

import (
	"fmt"
	"log"

	flag "github.com/spf13/pflag"
	"github.com/tuanito/go-deribit/v3"
	"github.com/tuanito/go-deribit/v3/client/account_management"
	"github.com/tuanito/go-deribit/v3/client/market_data"
	"github.com/tuanito/go-deribit/v3/client/private"
	"github.com/tuanito/go-deribit/v3/client/public"
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
	// Hit the test RPC endpoint
	res, err := client.Public.GetPublicTest(&public.GetPublicTestParams{})
	if err != nil {
		log.Fatalf("Error testing connection: %s", err)
	}
	log.Printf("Connected to Deribit API v%s", *res.Payload.Result.Version)
	if err := exchange.Authenticate(*key, *secret); err != nil {
		log.Fatalf("Error authenticating: %s", err)
	}

	// --------------------
	// 1. Account Summary
	// --------------------
	// see v3/models/account_management.go
	accountSummary, err := client.AccountManagement.GetPrivateGetAccountSummary(&account_management.GetPrivateGetAccountSummaryParams{Currency: "BTC"})

	if err != nil {
		log.Fatalf("Error getting account accountSummary: %s", err)
	}
	fmt.Printf("Currency: %s\n", *accountSummary.Payload.Result.Currency)
	fmt.Printf("Available funds: %f\n", *accountSummary.Payload.Result.AvailableFunds)
	fmt.Printf("Balance : %f\n", *accountSummary.Payload.Result.Balance)
	fmt.Printf("Equity : %f\n", *accountSummary.Payload.Result.Equity)
	fmt.Printf("Delta Total: %f\n", *accountSummary.Payload.Result.DeltaTotal)
	fmt.Printf("Options Delta: %f\n", *accountSummary.Payload.Result.OptionsDelta)
	fmt.Printf("Options Gamma: %f\n", *accountSummary.Payload.Result.OptionsGamma)
	fmt.Printf("Options Vega: %f\n", *accountSummary.Payload.Result.OptionsVega)
	fmt.Printf("Options Theta: %f\n", *accountSummary.Payload.Result.OptionsTheta)
	// fmt.Printf("Session funding: %f\n", *accountSummary.Payload.Result.SessionFunding)
	fmt.Printf("Futures PnL: %f\n", *accountSummary.Payload.Result.FuturesPl)
	fmt.Printf("Options PnL: %f\n", *accountSummary.Payload.Result.OptionsPl)
	fmt.Printf("Total PnL: %f\n", *accountSummary.Payload.Result.TotalPl)
	fmt.Printf("Equity : %f\n", *accountSummary.Payload.Result.Equity)

	// --------------------
	// 2. Referential
	// --------------------

	referentialParams := &public.GetPublicGetInstrumentsParams{
		Currency: "BTC",
	}
	referential, errRef := client.Public.GetPublicGetInstruments(referentialParams)
	if errRef != nil {
		log.Fatalf("Error loading referential: %s", errRef)
	}
	for key, value := range referential.Payload.Result {
		fmt.Printf("Instrument :  %d %s %+v \n", key, value.InstrumentName, value)
	}

	// --------------------
	// 3. Book Summary
	// Makes a snapshot of prices on all derivatives
	// --------------------

	bookSummaryParams := &market_data.GetPublicGetBookSummaryByCurrencyParams{
		Currency: "BTC",
	}

	bookSummary, errBook := client.MarketData.GetPublicGetBookSummaryByCurrency(bookSummaryParams)
	if errBook != nil {
		log.Fatalf("Error loading book summary : %s", errBook)
	}
	for key, value := range bookSummary.Payload.Result {
		fmt.Printf("Instrument price request :  %d %s %+v \n", key, value.InstrumentName, value)
	}
	/*
		// Buy
		buyParams := &private.GetPrivateBuyParams{
			Amount:         10,
			InstrumentName: "BTC-PERPETUAL",
			Type:           strPointer("market"),
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
	positionParams := &private.GetPrivateGetPositionParams{
		InstrumentName: "BTC-PERPETUAL",
	}

	position, errPos := client.Private.GetPrivateGetPosition(positionParams)
	if errPos != nil {
		log.Fatalf("Error getting position : %s", errPos)
	}

	fmt.Printf("Position: %+v\n", position.Payload.Result /*positions.Payload.Result*/)

	// --------------------
	// 5. Account positions
	// --------------------
	// see v3/models/position.go && v3/client/private/private_client.go
	positionsParams := &private.GetPrivateGetPositionsParams{
		Currency: "BTC",
		Kind:     strPointer("future"),
	}
	positions, errPos2 := client.Private.GetPrivateGetPositions(positionsParams)
	if errPos != nil {
		log.Fatalf("Error getting positions : %s", errPos2)
	}
	for key, value := range positions.Payload.Result {
		fmt.Printf("FUTURES : Instrument :  %d %s %+v \n", key, value.InstrumentName, value)
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

func strPointer(str string) *string {
	return &str
}
