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

	// TODO
	router.HandlerFunc(http.MethodPost, "/v1/addresses", app.createAddressHandler)
	router.HandlerFunc(http.MethodGet, "/v1/addresses/:id", app.getAddressHandler)
	router.HandlerFunc(http.MethodPut, "/v1/addresses/:id", app.updateAddressHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/addresses/:id", app.deleteAddressHandler)

	router.HandlerFunc(http.MethodPost, "/v1/admin/roles", app.createRoleHandler)
	router.HandlerFunc(http.MethodGet, "/v1/admin/roles", app.getAllRolesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/admin/roles/:id", app.getRoleHandler)
	router.HandlerFunc(http.MethodPut, "/v1/admin/roles/:id", app.updateRoleHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/admin/roles/:id", app.deleteRoleHandler)

	router.HandlerFunc(http.MethodPost, "/v1/admin/users/:id/role", app.updateUserRoleHandler)

	return router
}
