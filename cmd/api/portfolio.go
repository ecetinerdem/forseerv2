package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/ecetinerdem/forseerv2/internal/store"
	"github.com/go-chi/chi/v5"
)

type portfolioKey string

const portfolioCtx portfolioKey = "portfolio"

type CreatePortfolioPayload struct {
	Name   string        `json:"name" validate:"required,max=50"`
	Stocks []store.Stock `json:"stocks,omitempty"`
}

type UpdatePortfolioPayload struct {
	Name string `json:"name" validate:"required,max=50"`
}

func (app *application) createPortfolioHandler(w http.ResponseWriter, r *http.Request) {

	userID := int64(1) // TODO: get from auth

	var createPortfolio CreatePortfolioPayload

	err := readJson(w, r, &createPortfolio)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err = Validate.Struct(createPortfolio)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	portfolio := &store.Portfolio{
		//Change after auth
		UserID: userID,
		Name:   createPortfolio.Name,
		Stocks: createPortfolio.Stocks,
	}

	ctx := r.Context()

	err = app.store.Portfolio.CreatePortfolioWithStocks(ctx, portfolio)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.writeJsonResponse(w, http.StatusCreated, portfolio)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPortfoliosHandler(w http.ResponseWriter, r *http.Request) {

	userID := int64(1) // TODO: get from auth

	ctx := r.Context()

	portfolios, err := app.store.Portfolio.GetPortfolios(ctx, userID)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.writeJsonResponse(w, http.StatusOK, portfolios)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPortfolioHandler(w http.ResponseWriter, r *http.Request) {

	userID := int64(1) // TODO: get from auth

	URLPortfolioID := chi.URLParam(r, "portfolioID")
	portfolioID, err := strconv.ParseInt(URLPortfolioID, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	portfolio, err := app.store.Portfolio.GetPortfolioByID(ctx, portfolioID, userID)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = app.writeJsonResponse(w, http.StatusOK, portfolio)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) searchPortfoliosHandler(w http.ResponseWriter, r *http.Request) {
	userID := int64(1) // TODO: get from auth

	searchParam := r.URL.Query().Get("name")

	if searchParam == "" {
		app.badRequestError(w, r, errors.New("search query parameter is required"))
		return
	}

	ctx := r.Context()

	portfolios, err := app.store.Portfolio.SearchPortfoliosByName(ctx, userID, searchParam)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if len(portfolios) == 0 {
		err = app.writeJsonResponse(w, http.StatusOK, []interface{}{})
		if err != nil {
			app.internalServerError(w, r, err)
		}
		return
	}

	err = app.writeJsonResponse(w, http.StatusOK, portfolios)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) updatePortfolioHandler(w http.ResponseWriter, r *http.Request) {
	userID := int64(1) // TODO: get from auth

	URLPortfolioID := chi.URLParam(r, "portfolioID")
	portfolioID, err := strconv.ParseInt(URLPortfolioID, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	var updatePortfolioPayload UpdatePortfolioPayload

	err = readJson(w, r, &updatePortfolioPayload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err = Validate.Struct(&updatePortfolioPayload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	portfolio, err := app.store.Portfolio.UpdatePortfolio(ctx, portfolioID, userID, updatePortfolioPayload.Name)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = app.writeJsonResponse(w, http.StatusOK, portfolio)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePortfolioHandler(w http.ResponseWriter, r *http.Request) {
	userID := int64(1) // TODO: get from auth

	URLPortfolioID := chi.URLParam(r, "portfolioID")
	portfolioID, err := strconv.ParseInt(URLPortfolioID, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	err = app.store.Portfolio.DeletePortfolio(ctx, portfolioID, userID)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) portfoliosContextMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := int64(1) // TODO: get from auth
		URLPortfolioID := chi.URLParam(r, "portfolioID")
		portfolioID, err := strconv.ParseInt(URLPortfolioID, 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()

		portfolio, err := app.store.Portfolio.GetPortfolioByID(ctx, portfolioID, userID)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
		ctx = context.WithValue(ctx, portfolioCtx, portfolio)
		next.ServeHTTP(w, r)
	})
}

func getPortfolioFromCtx(r *http.Request) *store.Portfolio {
	portfolio, _ := r.Context().Value(portfolioCtx).(*store.Portfolio)

	return portfolio
}
