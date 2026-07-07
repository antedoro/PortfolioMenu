package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/antedoro/PortfolioMenu/internal/portfolio"
)

type Server struct {
	Portfolio *portfolio.Updater

	Templates *template.Template
}

func NewServer(
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

		Portfolio: updater,

		Templates: tmpl,
	}

}

func (s *Server) Start(
	addr string,
) {

	http.HandleFunc(
		"/",
		s.dashboard,
	)

	fmt.Println(
		"Dashboard:",
		"http://"+addr,
	)

	go func() {

		err :=
			http.ListenAndServe(
				addr,
				nil,
			)

		if err != nil {

			panic(err)

		}

	}()

}

func (s *Server) dashboard(
	w http.ResponseWriter,
	r *http.Request,
) {

	data :=
		s.Portfolio.Get()

	err :=
		s.Templates.ExecuteTemplate(
			w,
			"index.html",
			data,
		)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			500,
		)

	}

}
