package models

import "time"

type AssetType string

const (
	Stock  AssetType = "Stock"
	ETF    AssetType = "ETF"
	Bond   AssetType = "Bond"
	Crypto AssetType = "Crypto"
	Cash   AssetType = "Cash"
)

func (t AssetType) String() string { return string(t) }

// AllTypes lista ordinata dei tipi supportati.
func AllTypes() []AssetType {
	return []AssetType{ETF, Stock, Bond, Crypto, Cash}
}

// Asset rappresenta uno strumento finanziario posseduto.
type Asset struct {
	ID   int
	Name string
	// Ticker breve usato in menubar / tabelle.
	Ticker string
	Type   AssetType
	Broker string
	// Simbolo su Borsa Italiana / mercato di riferimento.
	Symbol string
	// Simbolo Yahoo Finance (vuoto se non usato).
	YahooSymbol string
	// Codice ISIN (obbligazioni o generico).
	ISIN string
	// GovBond indica titoli di Stato agevolati (12,5%).
	GovBond bool

	// Posizione.
	Quantity     float64
	AvgCost      float64 // Prezzo Medio di Carico (PMC).
	Fees         float64 // Commissioni di acquisto.
	PurchaseDate string  // Data acquisto (YYYY-MM-DD).

	// Valuta di quotazione (EUR, USD, ...).
	Currency string

	// Prezzo manuale (ignorato se 0).
	ManualPrice float64

	// Dati di mercato aggiornati dai provider.
	LastPrice     float64
	PreviousClose float64
	CurrencySymbol string
	PriceSource   string
	LastUpdate    time.Time

	// Valori calcolati.
	Invested     float64
	MarketValue  float64
	GainLoss     float64
	GainPercent  float64
}

// TaxRate restituisce l'aliquota fiscale italiana applicabile.
//
// 12,5% per titoli di Stato agevolati (BTP, BOT, CCT, BEI, UE),
// 26% per tutto il resto (ETF, azioni, corporate bond, crypto).
func (a Asset) TaxRate() float64 {
	if a.Type == Bond && a.GovBond {
		return 0.125
	}
	return 0.26
}

// IsManual indica se il prezzo è inserito manualmente.
func (a Asset) IsManual() bool {
	return a.ManualPrice > 0
}
