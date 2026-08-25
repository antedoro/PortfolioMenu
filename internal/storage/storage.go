package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/antedoro/PortfolioMenu/internal/models"
)

// Store gestisce la persistenza del portafoglio su file
// TOML o JSON (rilevato dall'estensione).
type Store struct {
	path string
}

// New crea uno store per il percorso indicato.
func New(path string) *Store {
	return &Store{path: path}
}

// Path restituisce il percorso del file dati.
func (s *Store) Path() string { return s.path }

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func parseTime(v string) time.Time {
	if v == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return time.Time{}
	}
	return t
}

// assetFile è la forma serializzata di un asset.
type assetFile struct {
	ID           int     `toml:"id" json:"id"`
	Name         string  `toml:"name" json:"name"`
	Ticker       string  `toml:"ticker" json:"ticker"`
	Type         string  `toml:"type" json:"type"`
	Broker       string  `toml:"broker" json:"broker"`
	Symbol       string  `toml:"symbol" json:"symbol"`
	YahooSymbol  string  `toml:"yahoo_symbol" json:"yahoo_symbol"`
	ISIN         string  `toml:"isin" json:"isin"`
	GovBond      bool    `toml:"gov_bond" json:"gov_bond"`
	Quantity     float64 `toml:"quantity" json:"quantity"`
	AvgCost      float64 `toml:"avg_cost" json:"avg_cost"`
	Fees         float64 `toml:"fees" json:"fees"`
	PurchaseDate string  `toml:"purchase_date" json:"purchase_date"`
	Currency     string  `toml:"currency" json:"currency"`
	ManualPrice  float64 `toml:"manual_price" json:"manual_price"`

	LastPrice      float64 `toml:"last_price" json:"last_price"`
	PreviousClose  float64 `toml:"previous_close" json:"previous_close"`
	CurrencySymbol string  `toml:"currency_symbol" json:"currency_symbol"`
	PriceSource    string  `toml:"price_source" json:"price_source"`
	LastUpdate     string  `toml:"last_update" json:"last_update"`
}

type portfolioFile struct {
	Version int         `toml:"version" json:"version"`
	Assets  []assetFile `toml:"assets" json:"assets"`
}

func toFile(a models.Asset) assetFile {
	return assetFile{
		ID:             a.ID,
		Name:           a.Name,
		Ticker:         a.Ticker,
		Type:           string(a.Type),
		Broker:         a.Broker,
		Symbol:         a.Symbol,
		YahooSymbol:    a.YahooSymbol,
		ISIN:           a.ISIN,
		GovBond:        a.GovBond,
		Quantity:       a.Quantity,
		AvgCost:        a.AvgCost,
		Fees:           a.Fees,
		PurchaseDate:   a.PurchaseDate,
		Currency:       a.Currency,
		ManualPrice:    a.ManualPrice,
		LastPrice:      a.LastPrice,
		PreviousClose:  a.PreviousClose,
		CurrencySymbol: a.CurrencySymbol,
		PriceSource:    a.PriceSource,
		LastUpdate:     formatTime(a.LastUpdate),
	}
}

func fromFile(a assetFile) models.Asset {
	t := models.AssetType(a.Type)
	if t == "" {
		t = models.ETF
	}
	return models.Asset{
		ID:             a.ID,
		Name:           a.Name,
		Ticker:         a.Ticker,
		Type:           t,
		Broker:         a.Broker,
		Symbol:         a.Symbol,
		YahooSymbol:    a.YahooSymbol,
		ISIN:           a.ISIN,
		GovBond:        a.GovBond,
		Quantity:       a.Quantity,
		AvgCost:        a.AvgCost,
		Fees:           a.Fees,
		PurchaseDate:   a.PurchaseDate,
		Currency:       a.Currency,
		ManualPrice:    a.ManualPrice,
		LastPrice:      a.LastPrice,
		PreviousClose:  a.PreviousClose,
		CurrencySymbol: a.CurrencySymbol,
		PriceSource:    a.PriceSource,
		LastUpdate:     parseTime(a.LastUpdate),
	}
}

// Load carica il portafoglio dal file (TOML o JSON).
// Se il file non esiste, ne crea uno con i valori di default.
func (s *Store) Load() (models.Portfolio, error) {

	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		p := Default()
		if serr := s.Save(p); serr != nil {
			return p, serr
		}
		return p, nil
	}

	raw, err := os.ReadFile(s.path)
	if err != nil {
		return models.Portfolio{}, err
	}

	var pf portfolioFile

	switch strings.ToLower(filepath.Ext(s.path)) {
	case ".json":
		if err := json.Unmarshal(raw, &pf); err != nil {
			return models.Portfolio{}, err
		}
	default:
		if err := toml.Unmarshal(raw, &pf); err != nil {
			return models.Portfolio{}, err
		}
	}

	var p models.Portfolio
	for _, a := range pf.Assets {
		p.Assets = append(p.Assets, fromFile(a))
	}
	return p, nil
}

// Save persiste il portafoglio sul file.
func (s *Store) Save(p models.Portfolio) error {

	pf := portfolioFile{Version: 1}
	for _, a := range p.Assets {
		pf.Assets = append(pf.Assets, toFile(a))
	}

	if err := os.MkdirAll(
		filepath.Dir(s.path),
		os.ModePerm,
	); err != nil {
		return err
	}

	var data []byte
	var err error

	switch strings.ToLower(filepath.Ext(s.path)) {
	case ".json":
		data, err = json.MarshalIndent(pf, "", "  ")
	default:
		var sb strings.Builder
		err = toml.NewEncoder(&sb).Encode(pf)
		data = []byte(sb.String())
	}

	if err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}

// Default restituisce un portafoglio di esempio.
func Default() models.Portfolio {
	now := time.Now()
	seed := []models.Asset{
		{
			ID: 1, Name: "VWCE", Ticker: "VWCE",
			Type: models.ETF, Broker: "Directa",
			Symbol: "VWCE.MI", YahooSymbol: "VWCE.MI",
			Quantity: 162, AvgCost: 165, Currency: "EUR",
			LastUpdate: now,
		},
		{
			ID: 2, Name: "XEON", Ticker: "XEON",
			Type: models.ETF, Broker: "Fineco",
			Symbol: "XEON.DE", YahooSymbol: "XEON.DE",
			Quantity: 43, AvgCost: 149.4508, Currency: "EUR",
			LastUpdate: now,
		},
		{
			ID: 3, Name: "OAT 2027", Ticker: "OAT27",
			Type: models.Bond, Broker: "Directa",
			Symbol: "FR0013250560", ISIN: "FR0013250560",
			GovBond:  true,
			Quantity: 250, AvgCost: 98.612, Currency: "EUR",
			LastUpdate: now,
		},
		{
			ID: 4, Name: "BTP 2028", Ticker: "BTP28",
			Type: models.Bond, Broker: "Directa",
			Symbol: "IT0005641029", ISIN: "IT0005641029",
			GovBond:  true,
			Quantity: 250, AvgCost: 99.91, Currency: "EUR",
			LastUpdate: now,
		},
		{
			ID: 5, Name: "BTP 2030", Ticker: "BTP30",
			Type: models.Bond, Broker: "Directa",
			Symbol: "IT0005383309", ISIN: "IT0005383309",
			GovBond:  true,
			Quantity: 250, AvgCost: 94.73, Currency: "EUR",
			LastUpdate: now,
		},
		{
			ID: 6, Name: "Bitcoin", Ticker: "BTC",
			Type: models.Crypto, Broker: "Binance",
			Symbol: "BTC-USD", YahooSymbol: "BTC-USD",
			Quantity: 0.0239652, AvgCost: 54000,
			Currency: "USD", LastUpdate: now,
		},
	}
	return models.Portfolio{Assets: seed, LastUpdate: now}
}
