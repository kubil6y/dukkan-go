package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/v1/register", app.registerHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodPost, "/v1/admin/roles", app.createRoleHandler)
	router.HandlerFunc(http.MethodGet, "/v1/admin/roles", app.getAllRolesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/admin/roles/:id", app.getRoleHandler)
	router.HandlerFunc(http.MethodPut, "/v1/admin/roles/:id", app.updateRoleHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/admin/roles/:id", app.deleteRoleHandler)

	router.HandlerFunc(http.MethodPost, "/v1/admin/users/:id/role", app.updateUserRoleHandler)

	return router
}
