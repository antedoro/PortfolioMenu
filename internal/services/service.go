package services

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/antedoro/PortfolioMenu/internal/config"
	"github.com/antedoro/PortfolioMenu/internal/history"
	"github.com/antedoro/PortfolioMenu/internal/models"
	"github.com/antedoro/PortfolioMenu/internal/portfolio"
	"github.com/antedoro/PortfolioMenu/internal/providers"
	"github.com/antedoro/PortfolioMenu/internal/storage"
	"github.com/antedoro/PortfolioMenu/internal/taxes"
)

// Service è il coordinatore centrale dell'applicazione:
// detiene lo stato del portafoglio, lo persiste, lo aggiorna
// e ne salva lo storico.
type Service struct {
	mu        sync.RWMutex
	portfolio models.Portfolio
	cfg       *config.Config

	store   *storage.Store
	history *history.Store

	stop chan bool
}

func New(
	cfg *config.Config,
) *Service {

	s := &Service{
		cfg:   cfg,
		store: storage.New(cfg.DataFile),
		history: history.New(
			filepath.Join(
				filepath.Dir(cfg.DataFile),
				"history.json",
			),
		),
		stop: make(chan bool),
	}

	return s
}

// Load carica il portafoglio dal file dati.
func (s *Service) Load() error {

	p, err := s.store.Load()
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.portfolio = p
	s.mu.Unlock()

	return nil
}

// Get restituisce una copia sicura del portafoglio.
func (s *Service) Get() models.Portfolio {

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.portfolio
}

// History restituisce gli snapshot salvati.
func (s *Service) History() []history.Snapshot {

	snaps, err := s.history.Load()
	if err != nil {
		return nil
	}
	return snaps
}

// StorePath restituisce il percorso del file dati.
func (s *Service) StorePath() string {
	return s.store.Path()
}

// DarkMode restituisce lo stato del tema scuro.
func (s *Service) DarkMode() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg.Settings.DarkMode
}

// MenubarFormat restituisce il formato della menubar.
func (s *Service) MenubarFormat() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg.Settings.MenubarFormat
}

// UpdateMenubarFormat aggiorna e salva il formato della menubar.
func (s *Service) UpdateMenubarFormat(format []string) error {
	s.mu.Lock()
	s.cfg.Settings.MenubarFormat = format
	s.mu.Unlock()

	if s.cfg.FilePath == "" {
		return fmt.Errorf("percorso configurazione non impostato")
	}

	return s.cfg.SaveFile(s.cfg.FilePath)
}

// TaxSummary restituisce il riepilogo fiscale.
func (s *Service) TaxSummary() taxes.Summary {

	return taxes.OfPortfolio(s.Get())
}

// Start avvia l'aggiornamento automatico periodico.
func (s *Service) Start(
	refreshMinutes int,
) {

	go func() {

		s.Refresh()

		if refreshMinutes <= 0 {
			refreshMinutes = 15
		}

		ticker :=
			time.NewTicker(
				time.Duration(refreshMinutes) *
					time.Minute,
			)

		defer ticker.Stop()

		for {

			select {

			case <-ticker.C:
				s.Refresh()

			case <-s.stop:
				return

			}

		}

	}()

}

// Stop ferma l'aggiornamento.
func (s *Service) Stop() {

	select {

	case s.stop <- true:

	default:

	}

}

// Refresh aggiorna prezzi, cambio, ricalcola e salva.
func (s *Service) Refresh() {

	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println()
	fmt.Println("Aggiornamento portfolio...")

	// Cambio EUR/USD
	currency :=
		providers.NewCurrencyProvider()

	rate, err := currency.GetEURUSD()
	if err == nil {
		s.portfolio.ExchangeRate = rate
	}

	for i := range s.portfolio.Assets {

		asset := &s.portfolio.Assets[i]

		switch {

		case asset.IsManual():
			providers.NewManualProvider().
				GetPrice(asset)

		case asset.Type == models.Bond &&
			asset.ISIN != "":
			providers.NewBorsaProvider().
				GetPrice(asset)

		case asset.YahooSymbol != "":
			providers.NewYahooProvider().
				GetPrice(asset)

		}

	}

	portfolio.Calculate(&s.portfolio)

	s.portfolio.LastUpdate = time.Now()

	// Salva prezzi aggiornati
	_ = s.store.Save(s.portfolio)

	// Snapshot storico
	_ = s.history.Append(s.portfolio)

	fmt.Println(
		"Portfolio aggiornato:",
		s.portfolio.LastUpdate.Format(
			"02/01/2006 15:04",
		),
	)

}

func (s *Service) nextID() int {

	max := 0
	for _, a := range s.portfolio.Assets {
		if a.ID > max {
			max = a.ID
		}
	}
	return max + 1
}

// AddAsset aggiunge un nuovo asset e lo persiste.
func (s *Service) AddAsset(
	a models.Asset,
) (models.Asset, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if a.ID <= 0 {
		a.ID = s.nextID()
	}

	a.LastUpdate = time.Now()

	s.portfolio.Assets =
		append(s.portfolio.Assets, a)

	portfolio.Calculate(&s.portfolio)

	if err := s.store.Save(s.portfolio); err != nil {
		return a, err
	}

	return a, nil
}

// UpdateAsset sostituisce un asset esistente.
func (s *Service) UpdateAsset(
	a models.Asset,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	found := false
	for i := range s.portfolio.Assets {
		if s.portfolio.Assets[i].ID == a.ID {
			a.LastUpdate = time.Now()
			s.portfolio.Assets[i] = a
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf(
			"asset %d non trovato",
			a.ID,
		)
	}

	portfolio.Calculate(&s.portfolio)

	if err := s.store.Save(s.portfolio); err != nil {
		return err
	}

	return nil
}

// DeleteAsset rimuove un asset per ID.
func (s *Service) DeleteAsset(
	id int,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	filtered := s.portfolio.Assets[:0]
	found := false
	for _, a := range s.portfolio.Assets {
		if a.ID == id {
			found = true
			continue
		}
		filtered = append(filtered, a)
	}
	s.portfolio.Assets = filtered

	if !found {
		return fmt.Errorf(
			"asset %d non trovato",
			id,
		)
	}

	portfolio.Calculate(&s.portfolio)

	if err := s.store.Save(s.portfolio); err != nil {
		return err
	}

	return nil
}

// DuplicateAsset crea una copia di un asset esistente.
func (s *Service) DuplicateAsset(
	id int,
) (models.Asset, error) {

	s.mu.RLock()
	var src models.Asset
	found := false
	for _, a := range s.portfolio.Assets {
		if a.ID == id {
			src = a
			found = true
			break
		}
	}
	s.mu.RUnlock()

	if !found {
		return models.Asset{},
			fmt.Errorf(
				"asset %d non trovato",
				id,
			)
	}

	src.ID = 0
	src.Name = src.Name + " (copia)"
	src.Ticker = src.Ticker + "_2"

	return s.AddAsset(src)
}

// UpdateSettings aggiorna le impostazioni applicative.
func (s *Service) UpdateSettings(
	darkMode bool,
) error {

	s.mu.Lock()
	s.cfg.Settings.DarkMode = darkMode
	s.mu.Unlock()

	if s.cfg.FilePath == "" {
		return fmt.Errorf("percorso configurazione non impostato")
	}

	return s.cfg.SaveFile(s.cfg.FilePath)
}
