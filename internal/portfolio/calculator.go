package portfolio

import (
	"github.com/antedoro/PortfolioMenu/internal/models"
)

func Calculate(
	p *models.Portfolio,
) {

	var invested float64

	var value float64

	for i := range p.Assets {

		asset :=
			&p.Assets[i]

		// capitale investito
		asset.Invested =
			asset.Quantity *
				asset.AvgCost

		asset.CapitalInvested =
			asset.Invested

		// valore attuale base

		marketValue :=
			asset.Quantity *
				asset.LastPrice

		// Conversione crypto USD -> EUR

		if asset.Currency == "USD" &&
			p.ExchangeRate > 0 {

			marketValue =
				marketValue /
					p.ExchangeRate

		}

		asset.MarketValue =
			marketValue

		asset.GainLoss =
			asset.MarketValue -
				asset.Invested

		if asset.Invested != 0 {

			asset.GainPercent =
				(asset.GainLoss /
					asset.Invested) * 100

		}

		invested +=
			asset.Invested

		value +=
			asset.MarketValue

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

	}

}
