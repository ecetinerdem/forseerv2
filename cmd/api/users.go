package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ecetinerdem/forseerv2/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	userURLParam := chi.URLParam(r, "userID")

	userID, err := strconv.ParseInt(userURLParam, 10, 64)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	user, err := app.store.Users.GetUserByID(ctx, userID)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = app.writeJsonResponse(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	ctx := r.Context()
	err := app.store.Users.Activate(ctx, token)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = app.writeJsonResponse(w, http.StatusNoContent, "")
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
