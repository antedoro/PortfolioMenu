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

	// Nome completo asset
	Name string

	// Simbolo breve per menubar
	Ticker string

	Type AssetType

	Broker string

	Symbol string

	YahooSymbol string

	ISINBond string

	Quantity float64

	AvgCost float64

	ManualPrice float64

	// Prezzi
	LastPrice float64

	PreviousClose float64

	// Valute
	CurrencySymbol string

	// Compatibilità provider
	Currency string

	// Origine prezzo
	PriceSource string

	// Valori economici
	// campo storico usato dal calculator
	Invested float64

	// nuovo nome più descrittivo
	CapitalInvested float64

	MarketValue float64

	GainLoss float64

	GainPercent float64

	LastUpdate time.Time
}
