[![Build Status](https://travis-ci.com/adampointer/go-deribit.svg?branch=master)](https://travis-ci.com/adampointer/go-deribit)  [![Go Report Card](https://goreportcard.com/badge/github.com/adampointer/go-deribit)](https://goreportcard.com/report/github.com/adampointer/go-deribit)  [![codebeat badge](https://codebeat.co/badges/5bf32114-b7e1-4e70-91bf-fae2449fe2cb)](https://codebeat.co/projects/github-com-adampointer-go-deribit-master)

# go-deribit

## Status

This code has been forked from adampointer/go-deribit original code. Maintenance is now taken over on this fork (as of 20201216). This code adds encapsulated functions which allow a more user-friendly use of Adampointer's fantastic original API.

Here is a sample code:

``package main

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
}``

## V3 

This project is now using go1.13 with Go Modules, but should remain compatible with `dep`. Also, as there are some breaking changes introduced by the latest schema changes from the remote API, I have decided to carry on development in the new `v3` namespace with the project root containing the code tagged `v2.x`.

`import "github.com/adampointer/go-deribit/v3"`

We now have the latest API methods which were recently released such as `public/get_tradingview_chart_data`.

I recommend using the `v3` project in your projects as all onward development will now be within this project.`

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
