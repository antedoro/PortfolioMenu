package portfolio

import (
	"fmt"
	"sync"
	"time"

	"github.com/antedoro/PortfolioMenu/internal/models"
	"github.com/antedoro/PortfolioMenu/internal/providers"
)

type Updater struct {
	Portfolio *models.Portfolio

	RefreshMinutes int

	mu sync.RWMutex

	stop chan bool
}

func NewUpdater(
	p *models.Portfolio,
	refreshMinutes int,
) *Updater {

	return &Updater{

		Portfolio: p,

		RefreshMinutes: refreshMinutes,

		stop: make(chan bool),
	}

}

// Start avvia l'aggiornamento automatico

func (u *Updater) Start() {

	go func() {

		// primo aggiornamento immediato

		u.Update()

		ticker := time.NewTicker(
			time.Duration(u.RefreshMinutes) *
				time.Minute,
		)

		defer ticker.Stop()

		for {

			select {

			case <-ticker.C:

				u.Update()

			case <-u.stop:

				return

			}

		}

	}()

}

// Stop ferma il servizio

func (u *Updater) Stop() {

	u.stop <- true

}

// Update aggiorna prezzi e cambio

func (u *Updater) Update() {

	u.mu.Lock()

	defer u.mu.Unlock()

	fmt.Println()

	fmt.Println(
		"Aggiornamento portfolio...",
	)

	// Cambio EUR/USD

	currency := providers.NewCurrencyProvider()

	rate, err := currency.GetEURUSD()

	if err == nil {

		u.Portfolio.ExchangeRate = rate

	}

	// Aggiornamento asset

	for i := range u.Portfolio.Assets {

		asset := &u.Portfolio.Assets[i]

		switch {

		case asset.ManualPrice > 0:

			provider :=
				providers.NewManualProvider()

			err :=
				provider.GetPrice(asset)

			if err != nil {

				fmt.Println(
					"Errore manuale:",
					asset.Name,
					err,
				)

			}

		case asset.Type == models.Bond &&
			asset.ISINBond != "":

			provider :=
				providers.NewBorsaProvider()

			err :=
				provider.GetPrice(asset)

			if err != nil {

				fmt.Println(
					"Errore Borsa:",
					asset.Name,
					err,
				)

			}

		case asset.YahooSymbol != "":

			provider :=
				providers.NewYahooProvider()

			err :=
				provider.GetPrice(asset)

			if err != nil {

				fmt.Println(
					"Errore Yahoo:",
					asset.Name,
					err,
				)

			}

		}

	}

	Calculate(u.Portfolio)

	u.Portfolio.LastUpdate = time.Now()

	fmt.Println(
		"Portfolio aggiornato:",
		u.Portfolio.LastUpdate.Format(
			"02/01/2006 15:04",
		),
	)

}

// Get restituisce una copia sicura del portfolio

func (u *Updater) Get() models.Portfolio {

	u.mu.RLock()

	defer u.mu.RUnlock()

	return *u.Portfolio

}
