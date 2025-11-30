package main

import "net/http"

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

	UserID := 1

	ctx := r.Context()
}
