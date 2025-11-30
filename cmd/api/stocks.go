package main

import (
	"errors"
	"net/http"

	"github.com/ecetinerdem/forseerv2/internal/store"
)

type AddStockPayload struct {
	Symbol       string  `json:"symbol" validate:"required,min=1,max=4"`
	Shares       float64 `json:"shares" validate:"required,gt=0"`
	AveragePrice float64 `json:"average_price" validate:"requiredgt=0"`
}

// AddStockToPortfolio godoc
//
//	@Summary		Add a stock to a portfolio
//	@Description	Adds a new stock to the specified portfolio
//	@Tags			portfolios
//	@Accept			json
//	@Produce		json
//	@Param			portfolioID	path		int				true	"Portfolio ID"
//	@Param			payload		body		AddStockPayload	true	"Stock payload"
//	@Success		201			{object}	store.Stock
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		404			{object}	error
//	@Failure		409			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/portfolios/{portfolioID}/stocks [post]
func (app *application) addStockToPortfolioHandler(w http.ResponseWriter, r *http.Request) {
	portfolio := getPortfolioFromCtx(r)

	//Get from the auth later
	UserID := 1

	ctx := r.Context()

	var addStockPayload AddStockPayload
	err := readJson(w, r, &addStockPayload)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err = Validate.Struct(&addStockPayload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	stock := &store.Stock{
		Symbol:       addStockPayload.Symbol,
		Shares:       addStockPayload.Shares,
		AveragePrice: addStockPayload.AveragePrice,
	}
	err = app.store.Portfolio.AddStockToPortfolio(ctx, portfolio.ID, int64(UserID), stock)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		case errors.Is(err, store.ErrDuplicateStock):
			app.duplicateError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = app.writeJsonResponse(w, http.StatusOK, &stock)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
