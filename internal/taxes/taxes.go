package taxes

import (
	"github.com/antedoro/PortfolioMenu/internal/models"
)

// Result contiene il calcolo fiscale di un asset o del portafoglio.
type Result struct {
	GrossGain    float64
	TaxRate      float64
	Tax          float64
	NetGain      float64
	TaxableBase  float64
}

// OfAsset calcola imposta e gain netto per un singolo asset.
//
// L'imposta italiana colpisce il gain (plusvalenza). Le minusvalenze
// (gain negativo) non generano imposta ma compensano i gain futuri
// nello stesso anno fiscale (report).
func OfAsset(a models.Asset) Result {

	// Base imponibile = gain lordo (positivo tassato, negativo no).
	gross := a.GainLoss
	rate := a.TaxRate()

	var tax float64
	if gross > 0 {
		tax = gross * rate
	}

	return Result{
		GrossGain:   gross,
		TaxRate:     rate,
		Tax:         tax,
		NetGain:     gross - tax,
		TaxableBase: gross,
	}
}

// Summary aggrega i risultati dell'intero portafoglio.
type Summary struct {
	GrossGain  float64
	TotalTax   float64
	NetGain    float64
	ByRate     map[string]float64 // aliquota -> imposta
	Plusvalenze  float64
	Minusvalenze float64
}

// OfPortfolio calcola il riepilogo fiscale del portafoglio.
func OfPortfolio(p models.Portfolio) Summary {

	s := Summary{
		ByRate: map[string]float64{},
	}

	for _, a := range p.Assets {

		r := OfAsset(a)

		s.GrossGain += r.GrossGain
		s.TotalTax += r.Tax
		s.NetGain += r.NetGain

		key := formatRate(r.TaxRate)
		s.ByRate[key] += r.Tax

		if r.GrossGain >= 0 {
			s.Plusvalenze += r.GrossGain
		} else {
			s.Minusvalenze += r.GrossGain
		}
	}

	return s
}

func formatRate(r float64) string {
	switch r {
	case 0.125:
		return "12,5%"
	case 0.26:
		return "26%"
	default:
		return "26%"
	}
}
