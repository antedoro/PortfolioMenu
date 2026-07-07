package main

import (
	"fmt"
	"log"
	"time"

	"github.com/antedoro/PortfolioMenu/internal/config"
	"github.com/antedoro/PortfolioMenu/internal/models"
	"github.com/antedoro/PortfolioMenu/internal/portfolio"
	"github.com/antedoro/PortfolioMenu/internal/providers"
)

func main() {

	cfg, err := config.Load("configs/portfoliomenu.toml")

	if err != nil {

		log.Fatal(err)

	}

	fmt.Println()

	fmt.Println("PortfolioMenu")

	fmt.Println("========================================")

	fmt.Printf(
		"Refresh ogni %d minuti\n\n",
		cfg.RefreshMinutes,
	)

	p := buildPortfolio(cfg)

	// Primo calcolo iniziale

	portfolio.Calculate(&p)

	// Avvio updater automatico

	updater := portfolio.NewUpdater(
		&p,
		cfg.RefreshMinutes,
	)

	updater.Start()

	// Stampa iniziale

	printPortfolio(
		updater.Get(),
	)

	/*
		Il processo resta attivo.

		In futuro qui verranno inseriti:

		- systray macOS
		- web server dashboard
		- notifiche
	*/

	select {}

}

func buildPortfolio(
	cfg *config.Config,
) models.Portfolio {

	var p models.Portfolio

	for _, a := range cfg.Assets {

		asset := models.Asset{

			ID: a.ID,

			Name: a.Name,

			Type: models.AssetType(a.Type),

			Broker: a.Broker,

			Symbol: a.Symbol,

			YahooSymbol: a.YahooSymbol,

			ISINBond: a.ISINBond,

			Quantity: a.Quantity,

			AvgCost: a.AvgCost,

			ManualPrice: a.ManualPrice,

			LastUpdate: time.Now(),
		}

		switch {

		case asset.ManualPrice > 0:

			provider :=
				providers.NewManualProvider()

			err :=
				provider.GetPrice(&asset)

			if err != nil {

				fmt.Println(
					"Errore prezzo manuale:",
					asset.Name,
					err,
				)

			}

		case asset.Type == models.Bond &&
			asset.ISINBond != "":

			provider :=
				providers.NewBorsaProvider()

			err :=
				provider.GetPrice(&asset)

			if err != nil {

				fmt.Println(
					"Errore Borsa Italiana:",
					asset.Name,
					err,
				)

			}

		case asset.YahooSymbol != "":

			provider :=
				providers.NewYahooProvider()

			err :=
				provider.GetPrice(&asset)

			if err != nil {

				fmt.Println(
					"Errore Yahoo:",
					asset.Name,
					err,
				)

			}

		default:

			fmt.Println(
				"Nessun provider per:",
				asset.Name,
			)

		}

		p.Assets =
			append(
				p.Assets,
				asset,
			)

	}

	// cambio iniziale EUR/USD

	currency :=
		providers.NewCurrencyProvider()

	rate, err :=
		currency.GetEURUSD()

	if err == nil {

		p.ExchangeRate = rate

	}

	p.LastUpdate = time.Now()

	return p

}

func printPortfolio(
	p models.Portfolio,
) {

	fmt.Println("ASSETS")

	fmt.Println(
		"--------------------------------------------------------------------------------",
	)

	fmt.Printf(
		"%-3s %-15s %-10s %12s %5s %12s %12s %-20s\n",
		"ID",
		"Nome",
		"Tipo",
		"Prezzo",
		"Val.",
		"Valore",
		"Gain",
		"Fonte",
	)

	for _, a := range p.Assets {

		fmt.Printf(
			"%-3d %-15s %-10s %12.2f %5s %12.2f %12.2f %-20s\n",

			a.ID,

			a.Name,

			a.Type,

			a.LastPrice,

			a.CurrencySymbol,

			a.MarketValue,

			a.GainLoss,

			a.PriceSource,
		)

	}

	fmt.Println()

	fmt.Println("PORTFOLIO")

	fmt.Println("----------------------------------------")

	fmt.Printf(
		"EUR/USD           : %.4f\n",
		p.ExchangeRate,
	)

	fmt.Printf(
		"Capitale investito : %12.2f €\n",
		p.TotalCost,
	)

	fmt.Printf(
		"Valore attuale     : %12.2f €\n",
		p.TotalMarketValue,
	)

	fmt.Printf(
		"Gain/Loss          : %12.2f €\n",
		p.TotalGain,
	)

	fmt.Printf(
		"Gain/Loss %%        : %12.2f %%\n",
		p.TotalGainPercent,
	)

}
