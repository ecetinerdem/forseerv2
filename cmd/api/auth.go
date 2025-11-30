package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/ecetinerdem/forseerv2/internal/store"
	"github.com/google/uuid"
)

type RegisterUserpayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
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

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashCode := hex.EncodeToString(hash[:])
	err = app.store.Users.CreateAndInvite(ctx, user, hashCode, app.config.mail.expiry)
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
	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	if err := app.writeJsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
