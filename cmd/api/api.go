package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ecetinerdem/forseerv2/docs" //Required for generating swagger docs
	"github.com/ecetinerdem/forseerv2/internal/auth"
	"github.com/ecetinerdem/forseerv2/internal/env"
	"github.com/ecetinerdem/forseerv2/internal/mailer"
	"github.com/ecetinerdem/forseerv2/internal/store"
	"github.com/ecetinerdem/forseerv2/internal/store/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2" // http-swagger middleware
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         *store.Storage
	cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	version     string
	mail        mailConfig
	apiURL      string
	frontEndURL string
	auth        authConfig
	redisCfg    redisConfig
}

type dbConfig struct {
	addr        string
	maxOpenConn int
	maxIdleConn int
	maxIdleTime string
}

type mailConfig struct {
	sendGrid  sendGridConfig
	expiry    time.Duration
	fromEmail string
}

type sendGridConfig struct {
	apiKey string
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}
type basicConfig struct {
	user string
	pass string
}

type tokenConfig struct {
	secret string
	expiry time.Duration
	iss    string
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

func (app *application) mount() *chi.Mux {

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5174")},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	//r.Use(app.RateLimiterMiddleware)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.With(app.BasicAuthMiddleware()).Get("/healthz", app.healthzCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.TokenAuthMiddleware)
				r.Get("/", app.getUserHandler)
			})
		})
		//Only Public Route
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})

		r.Route("/portfolios", func(r chi.Router) {
			r.Use(app.TokenAuthMiddleware)
			r.Post("/", app.createPortfolioHandler)
			r.Get("/", app.getPortfoliosHandler)
			r.Get("/search", app.searchPortfoliosHandler)
			r.Route("/{portfolioID}", func(r chi.Router) {
				r.Use(app.portfoliosContextMiddleware)
				r.Get("/", app.getPortfolioHandler)
				r.Patch("/", app.updatePortfolioHandler)
				r.Delete("/", app.deletePortfolioHandler)

				r.Route("/stocks", func(r chi.Router) {
					r.Post("/", app.addStockHandler)
					r.Put("/{symbol}", app.updateStockHandler)
					r.Delete("/{symbol}", app.deleteStockHandler)
				})

			})
		})

	})

	return r
}

func (app *application) run(mux *chi.Mux) error {

	//Docs
	docs.SwaggerInfo.Version = app.config.version
	docs.SwaggerInfo.Host = "localhost:3000"
	docs.SwaggerInfo.BasePath = "/v1"
	srvr := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		app.logger.Infow("signal caught", "addr", "signal", s.String())
		shutdown <- srvr.Shutdown(ctx)
	}()

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	err := srvr.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
