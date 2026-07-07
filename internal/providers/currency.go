package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CurrencyProvider struct {
	Client *http.Client
}

func NewCurrencyProvider() *CurrencyProvider {

	return &CurrencyProvider{

		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

}

type yahooCurrencyResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`

				Currency string `json:"currency"`
			} `json:"meta"`
		} `json:"result"`
	} `json:"chart"`
}

func (c *CurrencyProvider) GetEURUSD() (float64, error) {

	url :=
		"https://query1.finance.yahoo.com/v8/finance/chart/EURUSD=X"

	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)

	if err != nil {
		return 0, err
	}

	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0",
	)

	resp, err := c.Client.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var data yahooCurrencyResponse

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return 0, err
	}

	if len(data.Chart.Result) == 0 {

		return 0,
			fmt.Errorf("cambio EURUSD non trovato")

	}

	rate :=
		data.Chart.Result[0].
			Meta.RegularMarketPrice

	return rate, nil

}
