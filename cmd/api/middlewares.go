package main

import (
	"net/http"
	"strings"

	"github.com/kubil6y/dukkan-go/internal/data"
	"github.com/kubil6y/dukkan-go/internal/validator"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = app.setUserContext(r, data.AnonUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] == "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()
		if validateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		r = app.setUserContext(r, user)

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.getUserContext(r)
		if user.IsAnon() {
			app.notPermittedResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireActivation(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.getUserContext(r)
		if !user.IsActivated {
			app.inactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
	return app.requireAuthentication(fn)
}

func (app *application) requireRole(roleName string, next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.getUserContext(r)
		if user.Role.Name != roleName {
			app.notPermittedResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})

	return app.requireActivation(fn)
}
