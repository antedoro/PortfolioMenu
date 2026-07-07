package models

import "time"

type Portfolio struct {
	Assets []Asset

	// Cambio EUR/USD
	ExchangeRate float64

	// Totali portfolio

	TotalValue float64

	TotalInvested float64

	TotalGain float64

	GainPercent float64

	LastUpdate time.Time
}
