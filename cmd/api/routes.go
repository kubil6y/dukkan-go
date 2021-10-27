package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/v1/register", app.registerHandler)                                  // public
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)    // public
	router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.activateAccountHandler)                  // public
	router.HandlerFunc(http.MethodPost, "/v1/tokens/generate-activation", app.generateActivationTokenHandler) // public

	router.HandlerFunc(http.MethodGet, "/v1/profile", app.requireAuthentication(app.getProfileHandler))         // authenticated
	router.HandlerFunc(http.MethodPatch, "/v1/profile/edit", app.requireAuthentication(app.editProfileHandler)) // authenticated

	router.HandlerFunc(http.MethodGet, "/v1/admin/users", app.getAllUsersHandler) // admin
	router.HandlerFunc(http.MethodGet, "/v1/admin/users/:id", app.getUserHandler) // admin

	router.HandlerFunc(http.MethodPost, "/v1/admin/roles", app.createRoleHandler)             // admin
	router.HandlerFunc(http.MethodGet, "/v1/admin/roles", app.getAllRolesHandler)             // admin
	router.HandlerFunc(http.MethodGet, "/v1/admin/roles/:id", app.getRoleHandler)             // admin
	router.HandlerFunc(http.MethodPut, "/v1/admin/roles/:id", app.updateRoleHandler)          // admin
	router.HandlerFunc(http.MethodDelete, "/v1/admin/roles/:id", app.deleteRoleHandler)       // admin
	router.HandlerFunc(http.MethodPut, "/v1/admin/users/:id/role", app.updateUserRoleHandler) // admin

	return app.authenticate(router)
}
