package main

import (
	"errors"
	"net/http"

	"github.com/ecetinerdem/forseerv2/internal/store"
)

type RegisterUserpayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserpayload

	err := readJson(w, r, &payload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err = Validate.Struct(&payload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	//hash password
	err = user.Password.Set(payload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	err = app.store.Users.CreateAndInvite(ctx, user, "", app.config.mail.expiry)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			app.badRequestError(w, r, err)
		case errors.Is(err, store.ErrDuplicateUsername):
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.writeJsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
