package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ecetinerdem/forseerv2/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  *store.Storage
}

type config struct {
	addr    string
	db      dbConfig
	env     string
	version string
}

type dbConfig struct {
	addr        string
	maxOpenConn int
	maxIdleConn int
	maxIdleTime string
}

func (app *application) mount() *chi.Mux {

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthz", app.healthzCheckHandler)
		r.Route("/portfolio", func(r chi.Router) {
			r.Post("/", app.createPortfolioHandler)
			r.Get("/", app.getPortfoliosHandler)
			r.Route("/{portfolioID}", func(r chi.Router) {
				r.Get("/", app.getPortfolioHandler)
			})
		})

	})

	return r
}

func (app *application) run(mux *chi.Mux) error {

	srvr := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server started at%v", srvr.Addr)

	return srvr.ListenAndServe()
}
