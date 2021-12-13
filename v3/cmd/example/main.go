package main

import (
	"encoding/csv"
	"fmt"
	"github.com/scancel/go-deribit/v3"
	"github.com/scancel/go-deribit/v3/client/public"
	flag "github.com/spf13/pflag"
	"log"
	"math"
	"os"
	"time"
)
func saveCSV(records map[string]float64){
	f, err := os.Create("benchmarks.csv")
	defer f.Close()

	if err != nil {

		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()
	for action, duration := range records{
		w.Write([]string{action, fmt.Sprintf("%f", duration)})
	}
}
func main() {
	key := flag.String("key", "", "Access key")
	secret := flag.String("secret", "", "Secret access key")
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
	var updatedPrice float64 = 20000

	NewOrder, errNewOrder := deribit.Buy(client, 10, "BTC-PERPETUAL", 20000, "limit")
	if errNewOrder != nil {
		log.Fatalf("Error creating new order: %s", errNewOrder)
	}

	var benchmarks = make(map[string]float64)
	depth := "20"
	interval := "100ms"
	book, err := exchange.SubscribeBookGroup("BTC-PERPETUAL", "none", depth, interval)
	if err != nil {
		log.Fatalf("Error subscribing to the book: %s", err)
	}
	for b := range book {
		fmt.Printf("Top bid: %f Top ask: %f\n", b.Bids[0][0], b.Asks[0][0])
		updatedPrice = math.Round(b.Bids[0][0])
		if errNewOrder == nil {
			start := time.Now()
			_, errEditOrder := deribit.Edit(client, string(NewOrder.Payload.Result.Order.OrderID), 10, updatedPrice/2)
			if errEditOrder != nil {
				log.Fatalf("Error editing order: %s", errEditOrder)
			}
			dur := time.Since(start).Microseconds()
			benchmarks["editOrder_" + start.String()] = float64(dur)/1000
			saveCSV(benchmarks)
			fmt.Printf("Edit order took %0d microseconds to execute", dur)
			//print(EditedOrder)
		}
	}

	exchange.Close()
}
