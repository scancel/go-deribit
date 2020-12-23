package book

type BookSummary struct {

	// The current best ask price, `null` if there aren't any asks
	// Required: true
	AskPrice float64

	// Base currency
	// Required: true
	BaseCurrency string

	// The current best bid price, `null` if there aren't any bids
	// Required: true
	BidPrice float64

	// creation timestamp
	// Required: true
	CreationTimestamp Timestamp

	// Current funding (perpetual only)
	CurrentFunding float64

	// Estimated delivery price, in USD. For more details, see Documentation > General > Expiration Price
	EstimatedDeliveryPrice float64

	// Funding 8h (perpetual only)
	Funding8h float64

	// Price of the 24h highest trade
	// Required: true
	High float64

	// instrument name
	// Required: true
	InstrumentName InstrumentName

	// Interest rate used in implied volatility calculations (options only)
	InterestRate float64

	// The price of the latest trade, `null` if there weren't any trades
	// Required: true
	Last float64

	// Price of the 24h lowest trade, `null` if there weren't any trades
	// Required: true
	Low float64

	// The current instrument market price
	// Required: true
	MarkPrice float64

	// The average of the best bid and ask, `null` if there aren't any asks or bids
	// Required: true
	MidPrice float64

	// The total amount of outstanding contracts in the corresponding amount units. For perpetual and futures the amount is in USD units, for options it is amount of corresponding cryptocurrency contracts, e.g., BTC or ETH.
	// Required: true
	OpenInterest float64

	// Quote currency
	// Required: true
	QuoteCurrency string

	// Name of the underlying future, or `'index_price'` (options only)
	UnderlyingIndex string

	// underlying price for implied volatility calculations (options only)
	UnderlyingPrice float64

	// The total 24h traded volume (in base currency)
	// Required: true
	Volume float64

	// Volume in usd (futures only)
	VolumeUsd float64
}
