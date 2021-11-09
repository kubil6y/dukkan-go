package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/kubil6y/dukkan-go/internal/data"
	"github.com/kubil6y/dukkan-go/internal/validator"
	"golang.org/x/time/rate"
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

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
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
			app.authenticationRequiredResponse(w, r)
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

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	// Define a client struct to hold the rate limiter and last seen time
	// for each client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Launch a background goroutine which removes old entries from the clients map
	// once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			// Lock the mutex to prevent any rate limiter checks from happening
			// while the cleanup is taking place.
			mu.Lock()

			// Loop through all clients. If they haven't been seen within the last
			// three minutes, delete the corresponding entry from the map.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

// NOTE If your code makes a decision about what to return based on the content of a request header,
// you should include that header name in your Vary response header — even if the request
// didn’t include that header.
// we can only set one header with Access-Control-Allow-Origin,
// so we will overcome decide it on runtime.
//go run ./cmd/api -cors-trusted-origins="http://localhost:9000 http://localhost:9001"
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		// Only run this if there's an Origin request header present AND at least one
		// trusted origin is configured.
		if origin != "" && len(app.config.cors.trustedOrigins) != 0 {
			// Loop through the list of trusted origins, checking to see if the request
			// origin exactly matches one of them.
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] {
					// If there is a match, then set a "Access-Control-Allow-Origin"
					// response header with the request origin as the value.
					w.Header().Set("Access-Control-Allow-Origin", origin)
					// Check if the request has the HTTP method OPTIONS and contains the
					// "Access-Control-Request-Method" header. If it does, then we treat
					// it as a preflight request.
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						// Set the necessary preflight response headers, as discussed
						// previously.
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						// Write the headers along with a 200 OK status and return from
						// the middleware with no further action.
						w.WriteHeader(http.StatusOK)
						return
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
