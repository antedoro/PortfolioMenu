package providers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antedoro/PortfolioMenu/internal/models"
)

type BorsaProvider struct {
	Client *http.Client
}

func NewBorsaProvider() *BorsaProvider {

	return &BorsaProvider{
		Client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (b *BorsaProvider) GetPrice(asset *models.Asset) error {

	if asset.ISINBond == "" {
		return fmt.Errorf("ISIN bond mancante")
	}

	url := fmt.Sprintf(
		"https://www.borsaitaliana.it/borsa/obbligazioni/mot/btp/scheda/%s-MOTX.html?lang=it",
		asset.ISINBond,
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
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
	)

	resp, err := b.Client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	html := string(body)

	/*
	   Cerchiamo:

	   Prezzo ufficiale

	   e successivamente
	   il primo numero decimale italiano
	*/

	re := regexp.MustCompile(
		`Prezzo ufficiale[\s\S]{0,500}?([0-9]{1,3},[0-9]{2,5})`,
	)

	match := re.FindStringSubmatch(html)

	if len(match) < 2 {

		return fmt.Errorf(
			"prezzo non trovato per %s",
			asset.ISINBond,
		)

	}

	priceString := strings.Replace(
		match[1],
		",",
		".",
		1,
	)

	price, err := strconv.ParseFloat(
		priceString,
		64,
	)

	if err != nil {
		return err
	}

	asset.LastPrice = price

	asset.Currency = "EUR"

	asset.CurrencySymbol = "€"

	asset.PriceSource = "Borsa Italiana"

	asset.LastUpdate = time.Now()

	return nil
}
