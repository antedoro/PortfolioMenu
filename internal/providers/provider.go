package providers

import "github.com/antedoro/PortfolioMenu/internal/models"

type PriceProvider interface {
	GetPrice(asset *models.Asset) error
}
