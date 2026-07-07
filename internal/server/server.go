package server

import (
	"embed"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/antedoro/PortfolioMenu/internal/portfolio"
)

type Server struct {
	Updater *portfolio.Updater

	Templates *template.Template
}

func New(
	updater *portfolio.Updater,
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

		Updater: updater,

		Templates: tmpl,
	}

}

func (s *Server) Start(
	address string,
) {

	mux :=
		http.NewServeMux()

	mux.HandleFunc(
		"/",
		s.index,
	)

	mux.HandleFunc(
		"/dashboard",
		s.dashboard,
	)

	mux.HandleFunc(
		"/api/charts",
		s.charts,
	)

	println(
		"Dashboard running:",
		"http://"+address,
	)

	go func() {

		err :=
			http.ListenAndServe(
				address,
				mux,
			)

		if err != nil {

			panic(err)

		}

	}()

}

func (s *Server) index(
	w http.ResponseWriter,
	r *http.Request,
) {

	p :=
		s.Updater.Get()

	err :=
		s.Templates.ExecuteTemplate(
			w,
			"index.html",
			p,
		)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

	}

}

// Endpoint HTMX
// aggiorna solo il blocco dashboard

func (s *Server) dashboard(
	w http.ResponseWriter,
	r *http.Request,
) {

	p :=
		s.Updater.Get()

	err :=
		s.Templates.ExecuteTemplate(
			w,
			"dashboard.html",
			p,
		)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

	}

}

func (s *Server) charts(
	w http.ResponseWriter,
	r *http.Request,
) {

	p :=
		s.Updater.Get()

	allocation :=
		make(map[string]float64)

	var names []string

	var gains []float64

	for _, asset := range p.Assets {

		allocation[string(asset.Type)] +=
			asset.MarketValue

		names =
			append(
				names,
				asset.Ticker,
			)

		gains =
			append(
				gains,
				asset.GainLoss,
			)

	}

	response :=
		struct {
			Allocation map[string]float64 `json:"allocation"`

			Names []string `json:"names"`

			Gains []float64 `json:"gains"`
		}{

			Allocation: allocation,

			Names: names,

			Gains: gains,
		}

	w.Header().
		Set(
			"Content-Type",
			"application/json",
		)

	json.NewEncoder(w).
		Encode(response)

}
