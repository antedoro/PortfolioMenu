package models

import "time"

type AssetType string

const (
	Stock  AssetType = "Stock"
	ETF    AssetType = "ETF"
	Bond   AssetType = "Bond"
	Crypto AssetType = "Crypto"
)

type Asset struct {
	ID int

	Name string
	Type AssetType

	Broker string

	Symbol string

	YahooSymbol string

	ISINBond string

	Quantity float64

	AvgCost float64

	ManualPrice float64

	Currency string

	CurrencySymbol string

	PriceSource string

	PreviousClose float64

	LastPrice float64

	MarketValue float64

	Invested float64

	GainLoss float64

	GainPercent float64

	LastUpdate time.Time
}
