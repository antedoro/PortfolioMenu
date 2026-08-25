package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/antedoro/PortfolioMenu/internal/models"
)

// Snapshot è un punto salvato nel tempo del portafoglio.
type Snapshot struct {
	Timestamp     string  `json:"timestamp"`
	TotalValue    float64 `json:"total_value"`
	TotalInvested float64 `json:"total_invested"`
	TotalGain     float64 `json:"total_gain"`
	GainPercent   float64 `json:"gain_percent"`
	ExchangeRate  float64 `json:"exchange_rate"`
}

// Store salva gli snapshot su file JSON.
type Store struct {
	path string
}

func New(path string) *Store {
	if path == "" {
		path = "configs/history.json"
	}
	return &Store{path: path}
}

func (s *Store) Path() string { return s.path }

// Load legge gli snapshot esistenti.
func (s *Store) Load() ([]Snapshot, error) {

	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		return nil, nil
	}

	raw, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}

	var snaps []Snapshot
	if err := json.Unmarshal(raw, &snaps); err != nil {
		return nil, err
	}
	return snaps, nil
}

// Append aggiunge uno snapshot e lo persiste.
func (s *Store) Append(p models.Portfolio) error {

	snaps, err := s.Load()
	if err != nil {
		return err
	}

	snap := Snapshot{
		Timestamp:     p.LastUpdate.Format(time.RFC3339),
		TotalValue:    p.TotalValue,
		TotalInvested: p.TotalInvested,
		TotalGain:     p.TotalGain,
		GainPercent:   p.GainPercent,
		ExchangeRate:  p.ExchangeRate,
	}

	snaps = append(snaps, snap)

	// Mantiene al massimo 3650 punti (~10 anni giornalieri).
	if len(snaps) > 3650 {
		snaps = snaps[len(snaps)-3650:]
	}

	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].Timestamp < snaps[j].Timestamp
	})

	if err := os.MkdirAll(
		filepath.Dir(s.path),
		os.ModePerm,
	); err != nil {
		return err
	}

	data, err := json.MarshalIndent(snaps, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}
