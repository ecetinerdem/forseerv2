package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ecetinerdem/forseerv2/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePortfolioPayload struct {
	Name   string        `json:"name" validate:"required,max=50"`
	Stocks []store.Stock `json:"stocks,omitempty"`
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

	err = writeJson(w, http.StatusCreated, portfolio)
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

	err = writeJson(w, http.StatusOK, portfolios)
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

	portfolio, err := app.store.Portfolio.GetPortfolioByID(ctx, userID, portfolioID)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = writeJson(w, http.StatusOK, portfolio)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) updatePortfolioHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deletePortfolioHandler(w http.ResponseWriter, r *http.Request) {

}
