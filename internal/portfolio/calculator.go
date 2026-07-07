package portfolio

import (
	"github.com/antedoro/PortfolioMenu/internal/models"
)

// Calculate aggiorna tutti i valori del portafoglio
func Calculate(p *models.Portfolio) {

	var totalCost float64
	var totalValue float64
	var totalGain float64

	for i := range p.Assets {

		asset := &p.Assets[i]

		// Capitale investito
		asset.Invested = calculateInvested(*asset)

		// Valore corrente
		asset.MarketValue = calculateMarketValue(*asset)

		// Gain/Loss lordo
		asset.GainLoss = asset.MarketValue - asset.Invested

		// Percentuale Gain/Loss
		if asset.Invested > 0 {

			asset.GainPercent =
				(asset.GainLoss / asset.Invested) * 100

		} else {

			asset.GainPercent = 0

		}

		totalCost += asset.Invested

		totalValue += asset.MarketValue

		totalGain += asset.GainLoss

	}

	p.TotalCost = totalCost

	p.TotalMarketValue = totalValue

	p.TotalGain = totalGain

	if totalCost > 0 {

		p.TotalGainPercent =
			(totalGain / totalCost) * 100

	} else {

		p.TotalGainPercent = 0

	}

}

// calculateInvested calcola il capitale iniziale investito
//
// ETF / Stock / Crypto:
// quantità × prezzo medio
//
// Bond:
// valore nominale × prezzo medio / 100
func calculateInvested(a models.Asset) float64 {

	switch a.Type {

	case models.Bond:

		return a.Quantity * a.AvgCost / 100

	default:

		return a.Quantity * a.AvgCost

	}

}

// calculateMarketValue calcola il valore attuale.
//
// ETF / Stock / Crypto:
// quantità × prezzo attuale
//
// Bond:
// valore nominale × prezzo mercato / 100
func calculateMarketValue(a models.Asset) float64 {

	switch a.Type {

	case models.Bond:

		return a.Quantity * a.LastPrice / 100

	default:

		return a.Quantity * a.LastPrice

	}

}
