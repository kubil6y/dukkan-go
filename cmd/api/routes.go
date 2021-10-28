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

	router.HandlerFunc(http.MethodPost, "/v1/review-product/:id", app.requireAuthentication(app.createReviewHandler))   // authenticated
	router.HandlerFunc(http.MethodPut, "/v1/review-product/:id", app.requireAuthentication(app.updateReviewHandler))    // authenticated
	router.HandlerFunc(http.MethodDelete, "/v1/review-product/:id", app.requireAuthentication(app.deleteReviewHandler)) // authenticated

	router.HandlerFunc(http.MethodGet, "/v1/admin/users", app.requireRole("admin", app.getAllUsersHandler)) // admin
	router.HandlerFunc(http.MethodGet, "/v1/admin/users/:id", app.requireRole("admin", app.getUserHandler)) // admin

	router.HandlerFunc(http.MethodPost, "/v1/admin/roles", app.requireRole("admin", app.createRoleHandler))             // admin
	router.HandlerFunc(http.MethodGet, "/v1/admin/roles", app.requireRole("admin", app.getAllRolesHandler))             // admin
	router.HandlerFunc(http.MethodGet, "/v1/admin/roles/:id", app.requireRole("admin", app.getRoleHandler))             // admin
	router.HandlerFunc(http.MethodPut, "/v1/admin/roles/:id", app.requireRole("admin", app.updateRoleHandler))          // admin
	router.HandlerFunc(http.MethodDelete, "/v1/admin/roles/:id", app.requireRole("admin", app.deleteRoleHandler))       // admin
	router.HandlerFunc(http.MethodPut, "/v1/admin/users/:id/role", app.requireRole("admin", app.updateUserRoleHandler)) // admin

	router.HandlerFunc(http.MethodGet, "/v1/products", app.getAllProductsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/products/:id", app.getProductHandler)
	router.HandlerFunc(http.MethodPost, "/v1/admin/products", app.createProductHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/admin/products/:id", app.updateProductHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/admin/products/:id", app.deleteProductHandler)

	return app.authenticate(router)
}
