package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/antedoro/PortfolioMenu/internal/models"
)

type YahooProvider struct {
	Client *http.Client
}

func NewYahooProvider() *YahooProvider {

	return &YahooProvider{

		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

}

type yahooResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`

				PreviousClose float64 `json:"previousClose"`

				Currency string `json:"currency"`
			} `json:"meta"`
		} `json:"result"`
	} `json:"chart"`
}

func (y *YahooProvider) GetPrice(asset *models.Asset) error {

	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v8/finance/chart/%s",
		asset.YahooSymbol,
	)

	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)

	if err != nil {
		return err
	}

	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0",
	)

	resp, err := y.Client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var data yahooResponse

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return err
	}

	if len(data.Chart.Result) == 0 {

		return fmt.Errorf(
			"no data for %s",
			asset.YahooSymbol,
		)

	}

	meta := data.Chart.Result[0].Meta

	asset.LastPrice = meta.RegularMarketPrice

	asset.PreviousClose = meta.PreviousClose

	asset.Currency = meta.Currency

	switch meta.Currency {

	case "USD":
		asset.CurrencySymbol = "$"

	case "EUR":
		asset.CurrencySymbol = "€"

	default:
		asset.CurrencySymbol = meta.Currency

	}

	asset.LastUpdate = time.Now()

	return nil
}
