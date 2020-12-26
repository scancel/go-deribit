package main

import (
	"fmt"
	"log"
	"strings"

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
	"github.com/tuanito/go-deribit/v3/structures/utils"
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
			Type:           utils.StrPointer("market"),
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
		// SessionFunding: utils.GetNumber(accountSummary.Payload.Result.SessionFunding), // <-- for some reasons, nil pointer crash
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
		fmt.Println("utils.HumanReadable : CreationTimestamp : ", myBookSummary.CreationTime, ":", utils.HumanReadableTime(int64(myBookSummary.CreationTimestamp)))
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
	return myPosition
}

// GetAccountPositions retrieves positions by derivatives types currency might be "BTC" or "ETH", Kind might be "future" or "option"
func GetAccountPositions(pClient *client.Deribit, pCurrency string, pKind string) []position.Position {
	// see v3/models/position.go && v3/client/private/private_client.go
	var outputPositions []position.Position
	positionsParams := &private.GetPrivateGetPositionsParams{
		Currency: pCurrency,
		Kind:     utils.StrPointer(strings.ToLower(pKind)), // utils.StrPointer("future"),
	}
	positions, errPos := pClient.Private.GetPrivateGetPositions(positionsParams)
	if errPos != nil {
		log.Fatalf("Error getting positions : %s", errPos)
	}
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
	return outputPositions
}
