package portfolio

import (
	"github.com/antedoro/PortfolioMenu/internal/models"
)

// Calculate ricalcola valori e totali del portafoglio.
//
// Tutti i valori sono convertiti nella valuta base (EUR). Gli asset
// quotati in valuta estera (es. crypto in USD) vengono convertiti
// usando il cambio EUR/USD fornito dal provider.
func Calculate(
	p *models.Portfolio,
) {

	var invested float64

	var value float64

	for i := range p.Assets {

		asset :=
			&p.Assets[i]

		// capitale investito (PMC + commissioni) in valuta nativa
		investedNative :=
			asset.Quantity*
				asset.AvgCost +
				asset.Fees

		// valore di mercato in valuta nativa
		marketNative :=
			asset.Quantity *
				asset.LastPrice

		// conversione a valuta base (EUR)
		toBase := func(v float64) float64 {

			if asset.Currency != "EUR" &&
				p.ExchangeRate > 0 {

				return v / p.ExchangeRate

			}

			return v

		}

		assetInvested :=
			toBase(investedNative)

		assetMarket :=
			toBase(marketNative)

		asset.Invested =
			assetInvested

		asset.MarketValue =
			assetMarket

		asset.GainLoss =
			assetMarket -
				assetInvested

		if asset.Invested != 0 {

			asset.GainPercent =
				(asset.GainLoss /
					asset.Invested) * 100
		} else {
			asset.GainPercent = 0

		}

		invested +=
			assetInvested

		value +=
			assetMarket

	}

	p.TotalInvested =
		invested

	p.TotalValue =
		value

	p.TotalGain =
		value - invested

	if invested != 0 {

		p.GainPercent =
			(p.TotalGain /
				invested) * 100
	} else {
		p.GainPercent = 0

	}

}
