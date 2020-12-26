package main

import (
	"fmt"
	"log"

	flag "github.com/spf13/pflag"
	"github.com/tuanito/go-deribit/v3"
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

	fmt.Println("---------------")
	fmt.Println("Account summary BTC")
	fmt.Println("---------------")
	BtcAccount, errAcctBtc := deribit.GetAccountSummary(client, "BTC")
	if errAcctBtc != nil {
		log.Fatalf("Error getting account summary: %s", errAcctBtc)
	}
	BtcAccount.Sprintf()
	fmt.Println("---------------")
	fmt.Println("Account summary ETH")
	fmt.Println("---------------")
	EthAccount, errAcctEth := deribit.GetAccountSummary(client, "ETH")
	if errAcctEth != nil {
		log.Fatalf("Error getting account summary: %s", errAcctEth)
	}
	EthAccount.Sprintf()

	// --------------------
	// 2. Referential
	// --------------------
	btcReferential, _ := deribit.GetReferential(client, "BTC")
	ethReferential, _ := deribit.GetReferential(client, "ETH")

	fmt.Println("---------------")
	fmt.Println("Referential BTC")
	fmt.Println("---------------")
	for key, value := range btcReferential {
		fmt.Println("#", key, ":\n", value.Sprintf())
	}
	fmt.Println("---------------")
	fmt.Println("Referential ETH")
	fmt.Println("---------------")
	for key, value := range ethReferential {
		fmt.Println("# ", key, ":\n", value.Sprintf())
	}

	// --------------------
	// 3. Book Summary
	// Makes a snapshot of prices on all derivatives
	// --------------------

	currency := "BTC"
	fmt.Println("Book Summary ", currency)
	fmt.Println("------------------")
	bookSummaries, _ := deribit.GetBookSummary(client, currency)

	for key, value := range bookSummaries {
		fmt.Println("Book summary :", key, " ", value.Sprintf())
	}

	// --------------------
	// 4. Account position for a specific underlying
	// --------------------
	// see v3/models/position.go && v3/client/private/private_client.go
	fmt.Println("Account Position")
	fmt.Println("------------------")
	pInstrumentName := "BTC-PERPETUAL"
	accountPosition, _ := deribit.GetAccountPosition(client, pInstrumentName)

	fmt.Println("Position of ", pInstrumentName)
	accountPosition.Sprintf()

	// --------------------
	// 5. Account positions for a kind of underlyings
	// --------------------
	fmt.Println("Account Positions")
	fmt.Println("------------------")
	futureType := "FUTURE" // option
	optionType := "OPTION" // option
	futuresPositions, _ := deribit.GetAccountPositions(client, currency, futureType)
	optionPositions, _ := deribit.GetAccountPositions(client, currency, optionType)

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

	// --------------------
	// 7. Transaction
	// --------------------
	// This has not been tried yet
	// amount := 0.0001
	// instrumentName := "BTC-PERPETUAL"
	// price := 12532.0
	// orderType := "limit" // "market"
	// buy, err := deribit.Buy(client, amount, instrumentName, price, orderType)

	exchange.Close()
}
