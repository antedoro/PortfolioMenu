package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	RefreshMinutes int           `toml:"refresh_minutes"`
	Assets         []AssetConfig `toml:"assets"`
}

type AssetConfig struct {
	ID          int     `toml:"id"`
	Name        string  `toml:"name"`
	Type        string  `toml:"type"`
	Broker      string  `toml:"broker"`
	Symbol      string  `toml:"symbol"`
	YahooSymbol string  `toml:"yahoo_symbol"`
	ISINBond    string  `toml:"isin_bond"`
	Quantity    float64 `toml:"quantity"`
	AvgCost     float64 `toml:"avg_cost"`
	ManualPrice float64 `toml:"manual_price"`
}

func Load(filename string) (*Config, error) {

	var cfg Config

	if _, err := toml.DecodeFile(filename, &cfg); err != nil {
		return nil, err
	}

	if cfg.RefreshMinutes <= 0 {
		cfg.RefreshMinutes = 15
	}

	for _, a := range cfg.Assets {

		if a.Name == "" {
			return nil, fmt.Errorf("asset without name")
		}

		if a.Quantity <= 0 {
			return nil, fmt.Errorf("%s: quantity must be > 0", a.Name)
		}

	}

	return &cfg, nil

}
