package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ecetinerdem/forseerv2/internal/mailer"
	"github.com/ecetinerdem/forseerv2/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  *store.Storage
	mailer mailer.Client
}

type config struct {
	addr    string
	db      dbConfig
	env     string
	version string
	mail    mailConfig
}

type dbConfig struct {
	addr        string
	maxOpenConn int
	maxIdleConn int
	maxIdleTime string
}

type mailConfig struct {
	expiry time.Duration
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

		r.Route("/users", func(r chi.Router) {
			r.Route("/activate/{token}", app.activateUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Get("/", app.getUserHandler)
			})
		})
		//Only Public Route
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
		})

		r.Route("/portfolios", func(r chi.Router) {
			r.Post("/", app.createPortfolioHandler)
			r.Get("/", app.getPortfoliosHandler)
			r.Get("/search", app.searchPortfoliosHandler)
			r.Route("/{portfolioID}", func(r chi.Router) {
				r.Use(app.portfoliosContextMiddleware)
				r.Get("/", app.getPortfolioHandler)
				r.Patch("/", app.updatePortfolioHandler)
				r.Delete("/", app.deletePortfolioHandler)
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
