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

// CreatePost godoc
//
//	@Summary		Creates a portfolio
//	@Description	Creates a portfolio
//	@Tags			portfolios
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePortfolioPayload	true	"Portfolio payload"
//	@Success		201		{object}	store.Portfolio
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/portfolios [post]
func (app *application) createPortfolioHandler(w http.ResponseWriter, r *http.Request) {

	user := getUserFromCtx(r)

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
		UserID: user.ID,
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

// GetPortfolios godoc
//
//	@Summary		Fetches all portfolios
//	@Description	Fetches all portfolios belonging to the authenticated user
//	@Tags			portfolios
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		store.Portfolio
//	@Failure		401	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/portfolios [get]
func (app *application) getPortfoliosHandler(w http.ResponseWriter, r *http.Request) {

	user := getUserFromCtx(r)

	pfq := &store.PaginatedFeedQuery{
		Limit:  5,
		Offset: 0,
		Sort:   "desc",
	}

	pfq, err := pfq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err = Validate.Struct(pfq)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	portfolios, err := app.store.Portfolio.GetPortfolios(ctx, user.ID, pfq)

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

// GetPost godoc
//
//	@Summary		Fetches a portfolio
//	@Description	Fetches a portfolio by ID
//	@Tags			portfolios
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Portfolio ID"
//	@Success		204	{object}	store.Portfolio
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/portfolios/{id} [get]
func (app *application) getPortfolioHandler(w http.ResponseWriter, r *http.Request) {

	user := getUserFromCtx(r)
	portfolioID, err := strconv.ParseInt(chi.URLParam(r, "portfolioID"), 10, 64)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	portfolio, err := app.getPortfolio(ctx, portfolioID, user.ID)

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

// SearchPortfolios godoc
//
//	@Summary		Search portfolios by name
//	@Description	Searches the user's portfolios whose names partially match the provided query
//	@Tags			portfolios
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	true	"Search term for portfolio name"
//	@Success		200		{array}		store.Portfolio
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/portfolios/search [get]
func (app *application) searchPortfoliosHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	searchParam := r.URL.Query().Get("name")

	if searchParam == "" {
		app.badRequestError(w, r, errors.New("search query parameter is required"))
		return
	}

	ctx := r.Context()

	portfolios, err := app.store.Portfolio.SearchPortfoliosByName(ctx, user.ID, searchParam)

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

// UpdatePortfolio godoc
//
//	@Summary		Updates an existing portfolio
//	@Description	Updates portfolio fields (currently only the name). Requires version check.
//	@Tags			portfolios
//	@Accept			json
//	@Produce		json
//	@Param			portfolioID	path		int						true	"Portfolio ID"
//	@Param			payload		body		UpdatePortfolioPayload	true	"Update payload"
//	@Success		200			{object}	store.Portfolio
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		404			{object}	error
//	@Failure		409			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/portfolios/{portfolioID} [patch]
func (app *application) updatePortfolioHandler(w http.ResponseWriter, r *http.Request) {

	user := getUserFromCtx(r)
	portfolioID, err := strconv.ParseInt(chi.URLParam(r, "portfolioID"), 10, 64)

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

	portfolio, err := app.getPortfolio(ctx, portfolioID, user.ID)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if updatePortfolioPayload.Name != "" {
		portfolio.Name = updatePortfolioPayload.Name
	}

	updatedPortfolio, err := app.store.Portfolio.UpdatePortfolio(ctx, portfolio, user.ID)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		case errors.Is(err, store.ErrVersionConflict):
			app.conflictError(w, r, err)
		default:
			app.internalServerError(w, r, err)

		}
		return
	}

	err = app.writeJsonResponse(w, http.StatusOK, updatedPortfolio)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// DeletePortfolio godoc
//
//	@Summary		Deletes a portfolio
//	@Description	Removes a portfolio permanently by its ID
//	@Tags			portfolios
//	@Accept			json
//	@Produce		json
//	@Param			portfolioID	path	int	true	"Portfolio ID"
//	@Success		204			"Portfolio deleted successfully"
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/portfolios/{portfolioID} [delete]
func (app *application) deletePortfolioHandler(w http.ResponseWriter, r *http.Request) {

	user := getUserFromCtx(r)
	portfolio := getPortfolioFromCtx(r)

	ctx := r.Context()

	err := app.store.Portfolio.DeletePortfolio(ctx, portfolio.ID, user.ID)

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
		user := getUserFromCtx(r)
		URLPortfolioID := chi.URLParam(r, "portfolioID")
		portfolioID, err := strconv.ParseInt(URLPortfolioID, 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()
		// Check Cache or go db
		portfolio, err := app.getPortfolio(ctx, portfolioID, user.ID)

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
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPortfolioFromCtx(r *http.Request) *store.Portfolio {
	portfolio, _ := r.Context().Value(portfolioCtx).(*store.Portfolio)

	return portfolio
}

// Check Cache or go db
func (app *application) getPortfolio(ctx context.Context, portfolioID int64, userID int64) (*store.Portfolio, error) {
	if !app.config.redisCfg.enabled {
		return app.store.Portfolio.GetPortfolioByID(ctx, portfolioID, userID)
	}

	portfolio, err := app.cacheStorage.Portfolio.Get(ctx, portfolioID)
	if err != nil {
		return nil, err
	}

	if portfolio == nil {
		portfolio, err = app.store.Portfolio.GetPortfolioByID(ctx, portfolioID, userID)
		if err != nil {
			return nil, err
		}
		err = app.cacheStorage.Portfolio.Set(ctx, portfolio)
		if err != nil {
			return nil, err
		}
	}

	return portfolio, nil
}
