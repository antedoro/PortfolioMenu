package server

import (
	"bytes"
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/antedoro/PortfolioMenu/internal/assets"
	"github.com/antedoro/PortfolioMenu/internal/history"
	"github.com/antedoro/PortfolioMenu/internal/models"
	"github.com/antedoro/PortfolioMenu/internal/reports"
	"github.com/antedoro/PortfolioMenu/internal/services"
	"github.com/antedoro/PortfolioMenu/internal/taxes"
)

// ViewModel dati passati ai template.
type ViewModel struct {
	Portfolio   models.Portfolio
	History     []history.Snapshot
	Tax         taxes.Summary
	DarkMode    bool
	DataFile    string
	Now         time.Time
	CurrentPage string
}

type Server struct {
	svc       *services.Service
	templates *template.Template
}

func New(
	svc *services.Service,
	templates embed.FS,
) *Server {

	tmpl :=
		template.Must(
			template.ParseFS(
				templates,
				"templates/*.html",
			),
		)

	return &Server{
		svc:       svc,
		templates: tmpl,
	}
}

func (s *Server) view() ViewModel {

	p := s.svc.Get()

	return ViewModel{
		Portfolio: p,
		History:   s.svc.History(),
		Tax:       s.svc.TaxSummary(),
		DarkMode:  s.svc.DarkMode(),
		DataFile:  s.svc.StorePath(),
		Now:       time.Now(),
	}
}

// Start avvia il server HTTP.
func (s *Server) Start(
	address string,
) {

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/",
		s.index,
	)

	mux.HandleFunc(
		"/partial/",
		s.partial,
	)

	mux.HandleFunc(
		"/api/charts",
		s.charts,
	)

	mux.HandleFunc(
		"/api/history",
		s.apiHistory,
	)

	mux.HandleFunc(
		"/api/taxes",
		s.apiTaxes,
	)

	mux.HandleFunc(
		"/api/refresh",
		s.apiRefresh,
	)

	mux.HandleFunc(
		"/api/settings",
		s.apiSettings,
	)

	mux.HandleFunc(
		"/api/export/csv",
		s.apiExportCSV,
	)

	mux.HandleFunc(
		"/api/assets",
		s.apiAssets,
	)

	mux.HandleFunc(
		"/api/assets/",
		s.apiAssetByID,
	)

	mux.HandleFunc(
		"/about",
		s.about,
	)

	mux.HandleFunc(
		"/config",
		s.configPage,
	)

	mux.HandleFunc(
		"/api/config/menubar",
		s.apiConfigMenubar,
	)

	// Servire l'icona
	mux.HandleFunc(
		"/icon.png",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/png")
			w.Write(assets.Icon)
		},
	)

	log.Printf(
		"Dashboard: http://%s",
		address,
	)

	go func() {

		if err := http.ListenAndServe(
			address,
			mux,
		); err != nil {
			log.Fatal(err)
		}

	}()

}

func (s *Server) about(
	w http.ResponseWriter,
	r *http.Request,
) {
	v := s.view()
	v.CurrentPage = "about"
	s.render(w, "about.html", v)
}

func (s *Server) configPage(
	w http.ResponseWriter,
	r *http.Request,
) {
	v := s.view()
	v.CurrentPage = "config"
	// Convert array to json for the frontend
	formatJSON, _ := json.Marshal(s.svc.MenubarFormat())
	
	data := struct {
		ViewModel
		MenubarFormatJSON template.JS
	}{
		ViewModel:         v,
		MenubarFormatJSON: template.JS(formatJSON),
	}
	s.render(w, "config.html", data)
}

func (s *Server) apiConfigMenubar(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Format []string `json:"format"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.svc.UpdateMenubarFormat(req.Format); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (s *Server) index(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	v := s.view()
	v.CurrentPage = "dashboard"
	s.render(w, "index.html", v)
}

// partial gestisce /partial/header, /partial/summary, ecc.
func (s *Server) partial(
	w http.ResponseWriter,
	r *http.Request,
) {

	name := strings.TrimPrefix(
		r.URL.Path,
		"/partial/",
	)

	allowed := map[string]bool{
		"header":  true,
		"summary": true,
		"charts":  true,
		"assets":  true,
		"history": true,
		"taxes":   true,
		"all":     true,
	}

	if !allowed[name] {
		http.NotFound(w, r)
		return
	}

	// refresh=1 forza l'aggiornamento prezzi prima del render.
	if r.URL.Query().Get("refresh") == "1" {
		s.svc.Refresh()
	}

	s.render(w, name, s.view())
}

func (s *Server) render(
	w http.ResponseWriter,
	name string,
	data interface{},
) {

	var buf bytes.Buffer

	err :=
		s.templates.ExecuteTemplate(
			&buf,
			name,
			data,
		)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	_, _ = w.Write(buf.Bytes())

}

func (s *Server) charts(
	w http.ResponseWriter,
	r *http.Request,
) {

	p := s.svc.Get()

	allocation :=
		make(map[string]float64)

	brokerAlloc :=
		make(map[string]float64)

	var names []string
	var gains []float64
	var values []float64
	var invested []float64

	for _, a := range p.Assets {

		allocation[string(a.Type)] +=
			a.MarketValue

		broker := a.Broker
		if broker == "" {
			broker = "—"
		}
		brokerAlloc[broker] +=
			a.MarketValue

		names = append(names, a.Ticker)
		gains = append(gains, a.GainLoss)
		values = append(values, a.MarketValue)
		invested = append(invested, a.Invested)

	}

	hist := s.svc.History()
	var hLabels []string
	var hValues []float64
	for _, h := range hist {
		hLabels = append(
			hLabels,
			dateLabel(h.Timestamp),
		)
		hValues = append(
			hValues,
			h.TotalValue,
		)
	}

	resp := struct {
		Allocation  map[string]float64 `json:"allocation"`
		BrokerAlloc map[string]float64 `json:"broker_alloc"`
		Names       []string           `json:"names"`
		Gains       []float64          `json:"gains"`
		Values      []float64          `json:"values"`
		Invested    []float64          `json:"invested"`
		HistLabels  []string           `json:"hist_labels"`
		HistValues  []float64          `json:"hist_values"`
	}{
		Allocation:  allocation,
		BrokerAlloc: brokerAlloc,
		Names:       names,
		Gains:       gains,
		Values:      values,
		Invested:    invested,
		HistLabels:  hLabels,
		HistValues:  hValues,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(resp)
}

func (s *Server) apiHistory(
	w http.ResponseWriter,
	r *http.Request,
) {

	snaps := s.svc.History()

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(snaps)
}

func (s *Server) apiTaxes(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(s.svc.TaxSummary())
}

func (s *Server) apiRefresh(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	s.svc.Refresh()

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]bool{"ok": true},
	)
}

func (s *Server) apiSettings(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	dark :=
		r.FormValue("dark_mode") == "true" ||
			r.FormValue("dark_mode") == "on"

	if err := s.svc.UpdateSettings(dark); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]bool{"ok": true},
	)
}

func (s *Server) apiExportCSV(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"text/csv",
	)

	w.Header().Set(
		"Content-Disposition",
		"attachment; filename=portfolio.csv",
	)

	_ = writeCSV(w, s.svc.Get())
}

func (s *Server) apiAssets(
	w http.ResponseWriter,
	r *http.Request,
) {

	switch r.Method {

	case http.MethodGet:

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		json.NewEncoder(w).
			Encode(s.svc.Get().Assets)

	case http.MethodPost:

		a := assetFromForm(r)
		if _, err := s.svc.AddAsset(a); err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}
		s.renderAll(w)

	default:
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

	}

}

func (s *Server) apiAssetByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	idStr := strings.TrimPrefix(
		r.URL.Path,
		"/api/assets/",
	)

	// Gestisce /api/assets/{id}/duplicate
	if idx := strings.Index(idStr, "/"); idx >= 0 {
		suffix := idStr[idx+1:]
		idStr = idStr[:idx]
		if suffix != "duplicate" {
			http.NotFound(w, r)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(
				w,
				"method not allowed",
				http.StatusMethodNotAllowed,
			)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(
				w,
				"id non valido",
				http.StatusBadRequest,
			)
			return
		}

		if _, err := s.svc.DuplicateAsset(id); err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusNotFound,
			)
			return
		}

		s.renderAll(w)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(
			w,
			"id non valido",
			http.StatusBadRequest,
		)
		return
	}

	switch r.Method {

	case http.MethodPut:

		a := assetFromForm(r)
		a.ID = id
		if err := s.svc.UpdateAsset(a); err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusNotFound,
			)
			return
		}
		s.renderAll(w)

	case http.MethodDelete:

		if err := s.svc.DeleteAsset(id); err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusNotFound,
			)
			return
		}
		s.renderAll(w)

	default:
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

	}

}

func (s *Server) partialAssets(
	w http.ResponseWriter,
) {

	s.render(w, "assets", s.view())
}

func (s *Server) renderAll(
	w http.ResponseWriter,
) {

	s.render(w, "all", s.view())
}

func dateLabel(timestamp string) string {
	if len(timestamp) >= 10 {
		return timestamp[:10]
	}
	return timestamp
}

func writeCSV(
	w http.ResponseWriter,
	p models.Portfolio,
) error {

	return reports.WriteAssetsCSV(w, p)
}

func assetFromForm(
	r *http.Request,
) models.Asset {

	_ = r.ParseForm()

	f64 := func(k string) float64 {
		v, _ := strconv.ParseFloat(
			r.FormValue(k),
			64,
		)
		return v
	}

	t := models.AssetType(
		r.FormValue("type"),
	)
	if t == "" {
		t = models.ETF
	}

	return models.Asset{
		Name:        r.FormValue("name"),
		Ticker:      r.FormValue("ticker"),
		Type:        t,
		Broker:      r.FormValue("broker"),
		Symbol:      r.FormValue("symbol"),
		YahooSymbol: r.FormValue("yahoo_symbol"),
		ISIN:        r.FormValue("isin"),
		GovBond: r.FormValue("gov_bond") == "on" ||
			r.FormValue("gov_bond") == "true",
		Quantity:     f64("quantity"),
		AvgCost:      f64("avg_cost"),
		Fees:         f64("fees"),
		PurchaseDate: r.FormValue("purchase_date"),
		Currency:     r.FormValue("currency"),
		ManualPrice:  f64("manual_price"),
	}
}
