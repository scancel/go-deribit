package main

import (
	"fmt"
	"log"
	"math"

	flag "github.com/spf13/pflag"
	"github.com/tuanito/go-deribit/v3"
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
	myAccount := account.Account{
		Currency:       *accountSummary.Payload.Result.Currency,
		AvailableFunds: *accountSummary.Payload.Result.AvailableFunds,
		Balance:        *accountSummary.Payload.Result.Balance,
		Equity:         *accountSummary.Payload.Result.Equity,
		DeltaTotal:     *accountSummary.Payload.Result.DeltaTotal,
		OptionsDelta:   *accountSummary.Payload.Result.OptionsDelta,
		OptionsGamma:   *accountSummary.Payload.Result.OptionsGamma,
		OptionsVega:    *accountSummary.Payload.Result.OptionsVega,
		OptionsTheta:   *accountSummary.Payload.Result.OptionsTheta,
		// *accountSummary.Payload.Result.SessionFunding)
		FuturesPl: *accountSummary.Payload.Result.FuturesPl,
		OptionsPl: *accountSummary.Payload.Result.OptionsPl,
		TotalPl:   *accountSummary.Payload.Result.TotalPl,
	}

	fmt.Println(myAccount.Sprintf())
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

		myInstrument := instrument.Instrument{
			BaseCurrency:        *value.BaseCurrency,
			InstrumentName:      *&value.InstrumentName,
			ContractSize:        *value.ContractSize,
			ExpirationTimestamp: *value.ExpirationTimestamp,
			IsActive:            *value.IsActive,
			Kind:                value.Kind,
			MinTradeAmount:      *value.MinTradeAmount,
			OptionType:          value.OptionType,
			QuoteCurrency:       *value.QuoteCurrency,
			SettlementPeriod:    *value.SettlementPeriod,
			Strike:              value.Strike,
			TickSize:            *value.TickSize,
		}

		fmt.Println("Referential : \n", key, ":", myInstrument.Sprintf())
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

	// See v3/models/book_summary.go
	for key, value := range bookSummary.Payload.Result {
		// fmt.Printf("Instrument raw price request : # %v \n %+v \n", key, value)
		fmt.Printf("Instrument price request : # %v \n", key)
		myBookSummary := summary.BookSummary{
			AskPrice:               getNumber(value.AskPrice),
			BaseCurrency:           *value.BaseCurrency,
			BidPrice:               getNumber(value.BidPrice),
			CreationTimestamp:      value.CreationTimestamp,
			CurrentFunding:         value.CurrentFunding,
			EstimatedDeliveryPrice: value.EstimatedDeliveryPrice,
			Funding8h:              value.Funding8h,
			High:                   getNumber(value.High),
			InstrumentName:         value.InstrumentName,

			InterestRate:    value.InterestRate,
			Last:            getNumber(value.Last),
			Low:             getNumber(value.Low),
			MarkPrice:       getNumber(value.MarkPrice),
			MidPrice:        getNumber(value.MidPrice),
			OpenInterest:    *value.OpenInterest,
			QuoteCurrency:   *value.QuoteCurrency,
			UnderlyingIndex: value.UnderlyingIndex,
			// underlying price for implied volatility calculations (options only)
			UnderlyingPrice: value.UnderlyingPrice,
			Volume:          *value.Volume,
			VolumeUsd:       value.VolumeUsd,
		}
		fmt.Println(myBookSummary.Sprintf())
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

	privatePosition, errPos := client.Private.GetPrivateGetPosition(positionParams)
	if errPos != nil {
		log.Fatalf("Error getting position : %s", errPos)
	}

	myPosition := position.Position{
		AveragePrice:              *privatePosition.Payload.Result.AveragePrice,
		AveragePriceUsd:           privatePosition.Payload.Result.AveragePriceUsd,
		Delta:                     *privatePosition.Payload.Result.Delta,
		Direction:                 privatePosition.Payload.Result.Direction,
		EstimatedLiquidationPrice: privatePosition.Payload.Result.EstimatedLiquidationPrice,
		FloatingProfitLoss:        *privatePosition.Payload.Result.FloatingProfitLoss,
		FloatingProfitLossUsd:     privatePosition.Payload.Result.FloatingProfitLossUsd,
		IndexPrice:                *privatePosition.Payload.Result.IndexPrice,
		InitialMargin:             *privatePosition.Payload.Result.InitialMargin,
		InstrumentName:            privatePosition.Payload.Result.InstrumentName,
		Kind:                      privatePosition.Payload.Result.Kind,
		MaintenanceMargin:         *privatePosition.Payload.Result.MaintenanceMargin,
		MarkPrice:                 *privatePosition.Payload.Result.MarkPrice,
		OpenOrdersMargin:          *privatePosition.Payload.Result.OpenOrdersMargin,
		RealizedProfitLoss:        *privatePosition.Payload.Result.RealizedProfitLoss,
		SettlementPrice:           *privatePosition.Payload.Result.SettlementPrice,
		Size:                      *privatePosition.Payload.Result.Size,
		SizeCurrency:              privatePosition.Payload.Result.SizeCurrency,
		// TotalProfitLoss:           *privatePosition.TotalProfitLoss, // <-- should be here though BUG
	}

	fmt.Println("Position : ", myPosition.Sprintf())

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

// getNumber will return NaN or the number, depending on the pointer
func getNumber(pNumber *float64) float64 {
	var output float64
	output = math.NaN()
	if pNumber != nil {
		output = *pNumber
	}
	return output
}
