package main

import (
	"fmt"
	"log"
	"time"

	"github.com/antedoro/PortfolioMenu/internal/assets"
	"github.com/antedoro/PortfolioMenu/internal/config"
	"github.com/antedoro/PortfolioMenu/internal/models"
	"github.com/antedoro/PortfolioMenu/internal/portfolio"
	"github.com/antedoro/PortfolioMenu/internal/providers"
	"github.com/antedoro/PortfolioMenu/internal/server"
	"github.com/antedoro/PortfolioMenu/internal/tray"
)

func main() {

	cfg, err := config.Load(
		"configs/portfoliomenu.toml",
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("PortfolioMenu")
	fmt.Println("========================================")

	fmt.Printf(
		"Refresh ogni %d minuti\n",
		cfg.RefreshMinutes,
	)

	p := buildPortfolio(cfg)

	portfolio.Calculate(&p)

	updater :=
		portfolio.NewUpdater(
			&p,
			cfg.RefreshMinutes,
		)

	updater.Start()

	fmt.Println(
		"Portfolio updater avviato",
	)

	webServer :=
		server.New(
			updater,
			assets.Templates,
		)

	webServer.Start(
		"localhost:8080",
	)

	fmt.Println(
		"Dashboard disponibile su http://localhost:8080",
	)

	appTray :=
		tray.New(
			updater,
		)

	appTray.Run()

}

func buildPortfolio(
	cfg *config.Config,
) models.Portfolio {

	var p models.Portfolio

	for _, a := range cfg.Assets {

		asset :=
			models.Asset{

				ID: a.ID,

				Name: a.Name,

				Ticker: a.Ticker,

				Type: models.AssetType(a.Type),

				Broker: a.Broker,

				Symbol: a.Symbol,

				YahooSymbol: a.YahooSymbol,

				ISINBond: a.ISINBond,

				Quantity: a.Quantity,

				AvgCost: a.AvgCost,

				ManualPrice: a.ManualPrice,

				Currency: "EUR",

				CurrencySymbol: "€",

				LastUpdate: time.Now(),
			}

		// Crypto quotate in USD

		if asset.Type == models.Crypto {

			asset.Currency = "USD"

			asset.CurrencySymbol = "$"

		}

		switch {

		// Prezzo inserito manualmente

		case asset.ManualPrice > 0:

			provider :=
				providers.NewManualProvider()

			err :=
				provider.GetPrice(
					&asset,
				)

			if err != nil {

				fmt.Println(
					"Errore prezzo manuale:",
					err,
				)

			}

		// Bond Borsa Italiana

		case asset.Type == models.Bond &&
			asset.ISINBond != "":

			provider :=
				providers.NewBorsaProvider()

			err :=
				provider.GetPrice(
					&asset,
				)

			if err != nil {

				fmt.Println(
					"Errore Borsa Italiana:",
					asset.Name,
					err,
				)

			}

		// Yahoo Finance

		case asset.YahooSymbol != "":

			provider :=
				providers.NewYahooProvider()

			err :=
				provider.GetPrice(
					&asset,
				)

			if err != nil {

				fmt.Println(
					"Errore Yahoo:",
					asset.Name,
					err,
				)

			}

		}

		p.Assets =
			append(
				p.Assets,
				asset,
			)

	}

	// Cambio EUR/USD

	currency :=
		providers.NewCurrencyProvider()

	rate, err :=
		currency.GetEURUSD()

	if err == nil {

		p.ExchangeRate = rate

	}

	p.LastUpdate =
		time.Now()

	return p

}
