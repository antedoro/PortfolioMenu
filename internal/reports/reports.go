package reports

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/antedoro/PortfolioMenu/internal/models"
	"github.com/antedoro/PortfolioMenu/internal/taxes"
)

func f(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}

// WriteAssetsCSV scrive la tabella degli asset in formato CSV.
func WriteAssetsCSV(
	w io.Writer,
	p models.Portfolio,
) error {

	cw := csv.NewWriter(w)

	header := []string{
		"ID", "Nome", "Ticker", "Tipo", "Broker",
		"ISIN", "Quantita", "PMC", "Commissioni",
		"Prezzo", "Valore", "Investito",
		"Gain", "Perc", "Valuta", "Aliquota",
	}

	if err := cw.Write(header); err != nil {
		return err
	}

	for _, a := range p.Assets {

		row := []string{
			strconv.Itoa(a.ID),
			a.Name,
			a.Ticker,
			string(a.Type),
			a.Broker,
			a.ISIN,
			f(a.Quantity),
			f(a.AvgCost),
			f(a.Fees),
			f(a.LastPrice),
			f(a.MarketValue),
			f(a.Invested),
			f(a.GainLoss),
			f(a.GainPercent),
			a.Currency,
			strconv.FormatFloat(
				a.TaxRate()*100,
				'f', 1, 64,
			) + "%",
		}

		if err := cw.Write(row); err != nil {
			return err
		}
	}

	cw.Flush()

	return cw.Error()
}

// WriteTaxCSV scrive il report fiscale per anno.
func WriteTaxCSV(
	w io.Writer,
	p models.Portfolio,
	year string,
) error {

	cw := csv.NewWriter(w)

	if err := cw.Write([]string{
		"Anno Fiscale", year,
	}); err != nil {
		return err
	}

	if err := cw.Write([]string{
		"Tipo", "Plusvalenze",
		"Minusvalenze", "Imposta",
	}); err != nil {
		return err
	}

	sum := taxes.OfPortfolio(p)

	if err := cw.Write([]string{
		"Portafoglio",
		f(sum.Plusvalenze),
		f(sum.Minusvalenze),
		f(sum.TotalTax),
	}); err != nil {
		return err
	}

	cw.Flush()

	return cw.Error()
}

// BuildAssetCSV compatta in stringa (utility test/CLI).
func BuildAssetCSV(p models.Portfolio) (string, error) {

	var sb strings.Builder
	if err := WriteAssetsCSV(&sb, p); err != nil {
		return "", fmt.Errorf("csv: %w", err)
	}
	return sb.String(), nil
}
