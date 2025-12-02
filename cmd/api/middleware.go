package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ecetinerdem/forseerv2/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

type userKey string

const userCtx userKey = "user"

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unAuthorizedBasicError(w, r, fmt.Errorf("missing authorization header"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unAuthorizedBasicError(w, r, fmt.Errorf("malformed authorization header"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unAuthorizedBasicError(w, r, err)
				return
			}

			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass

			creds := strings.SplitN(string(decoded), ":", 2)

			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unAuthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
				return

			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) TokenAuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unAuthorizedError(w, r, fmt.Errorf("missing authorization header"))
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unAuthorizedError(w, r, fmt.Errorf("malformed authorization header"))
			return
		}

		token := parts[1]

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unAuthorizedError(w, r, err)
			return
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok || claims == nil {
			app.unAuthorizedError(w, r, fmt.Errorf("invalid token claims"))
			return
		}

		subRaw, ok := claims["sub"]
		if !ok {
			app.unAuthorizedError(w, r, fmt.Errorf("missing sub claim"))
			return
		}

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", subRaw), 10, 64)
		if err != nil {
			app.unAuthorizedBasicError(w, r, err)
			return
		}

		ctx := r.Context()

		// Check Cache or go db
		user, err := app.getUser(ctx, userID)
		if err != nil {
			app.unAuthorizedBasicError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// for endpoints
func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)

	return user
}

// Check Cache or go db
func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {

	if !app.config.redisCfg.enabled {
		return app.store.Users.GetUserByID(ctx, userID)
	}

	user, err := app.cacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = app.store.Users.GetUserByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		err = app.cacheStorage.Users.Set(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (app *application) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.rateLimiter.Enabled {
			if allow, retryAfter := app.rateLimiter.Allow(r.RemoteAddr); !allow {
				app.rateLimitExceedResponse(w, r, retryAfter.String())
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
