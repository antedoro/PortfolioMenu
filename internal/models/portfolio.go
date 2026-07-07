package models

import "time"

type Portfolio struct {
	Assets []Asset

	TotalCost float64

	TotalMarketValue float64

	TotalGain float64

	TotalGainPercent float64

	EURUSD float64

	ExchangeRate float64

	LastUpdate time.Time
}
