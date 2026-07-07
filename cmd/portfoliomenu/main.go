package main

import (
	"fmt"
	"log"
	"time"

	"github.com/antedoro/PortfolioMenu/internal/config"
	"github.com/antedoro/PortfolioMenu/internal/models"
	"github.com/antedoro/PortfolioMenu/internal/portfolio"
	"github.com/antedoro/PortfolioMenu/internal/providers"
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

	appTray :=
		tray.New(
			updater,
		)

	// Avvio menubar macOS
	appTray.Run()

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

			provider.GetPrice(
				&asset,
			)

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
					"Errore Borsa:",
					err,
				)

			}

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
