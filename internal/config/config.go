package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	RefreshMinutes  int    `toml:"refresh_minutes"`
	Theme           string `toml:"theme"`
	DefaultBroker   string `toml:"default_broker"`
	DefaultProvider string `toml:"default_provider"`
	Proxy           string `toml:"proxy"`
	CacheMinutes    int    `toml:"cache_minutes"`
	FilePath        string `toml:"-"`
	// DataFile percorso del file portfolio (TOML/JSON).
	DataFile string `toml:"data_file"`

	// Impostazioni dell'applicazione (non sono asset).
	Settings AppSettings `toml:"settings"`
}

type AppSettings struct {
	DarkMode      bool     `toml:"dark_mode"`
	Language      string   `toml:"language"`
	BaseCurrency  string   `toml:"base_currency"`
	MenubarFormat []string `toml:"menubar_format"`
}

func Load(filename string) (*Config, error) {

	var cfg Config

	if _, err := toml.DecodeFile(filename, &cfg); err != nil {
		return nil, err
	}

	cfg.FilePath = filename

	if cfg.RefreshMinutes <= 0 {
		cfg.RefreshMinutes = 15
	}

	if cfg.Theme == "" {
		cfg.Theme = "light"
	}

	if cfg.DataFile == "" {
		cfg.DataFile = "configs/portfolio.toml"
	}

	if cfg.CacheMinutes <= 0 {
		cfg.CacheMinutes = 60
	}

	if cfg.Settings.Language == "" {
		cfg.Settings.Language = "it"
	}

	if cfg.Settings.BaseCurrency == "" {
		cfg.Settings.BaseCurrency = "EUR"
	}

	if len(cfg.Settings.MenubarFormat) == 0 {
		cfg.Settings.MenubarFormat = []string{"ticker", "value", "percent", "gainloss"}
	}

	return &cfg, nil

}

// SaveFile persiste la configurazione su disco.
func (c *Config) SaveFile(filename string) error {

	if filename == "" {
		return fmt.Errorf("nome file vuoto")
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(c)
}
