package providers

import (
	"time"

	"github.com/antedoro/PortfolioMenu/internal/models"
)

type ManualProvider struct{}

func NewManualProvider() *ManualProvider {

	return &ManualProvider{}

}

func (m *ManualProvider) GetPrice(asset *models.Asset) error {

	asset.LastPrice = asset.ManualPrice

	asset.Currency = "EUR"

	asset.CurrencySymbol = "€"

	asset.LastUpdate = time.Now()

	return nil
}
