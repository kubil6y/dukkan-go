package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/register", app.registerHandler)
	router.HandlerFunc(http.MethodGet, "/v1/login", app.loginHandler)

	return router
}
