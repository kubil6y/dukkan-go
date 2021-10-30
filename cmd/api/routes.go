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
	router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.activateAccountHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/generate-activation", app.generateActivationTokenHandler)

	router.HandlerFunc(http.MethodGet, "/v1/profile", app.requireAuthentication(app.getProfileHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/profile/edit", app.requireAuthentication(app.editProfileHandler))

	router.HandlerFunc(http.MethodPost, "/v1/products/:slug/review", app.requireAuthentication(app.createReviewHandler))
	router.HandlerFunc(http.MethodPut, "/v1/products/:id/review", app.requireAuthentication(app.updateReviewHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/products/:id/review", app.requireAuthentication(app.deleteReviewHandler))

	router.HandlerFunc(http.MethodPost, "/v1/products/:slug/rating", app.requireAuthentication(app.createRatingHandler))
	router.HandlerFunc(http.MethodPut, "/v1/products/:id/rating", app.requireAuthentication(app.updateRatingHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/products/:id/rating", app.requireAuthentication(app.deleteRatingHandler))

	router.HandlerFunc(http.MethodGet, "/v1/admin/users", app.requireRole("admin", app.getAllUsersHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/users/:id", app.requireRole("admin", app.getUserHandler))

	router.HandlerFunc(http.MethodGet, "/v1/my-orders", app.requireActivation(app.getOrdersOfAuthUserHandler))
	router.HandlerFunc(http.MethodGet, "/v1/my-orders/:id", app.requireActivation(app.getOrderByIDOfAuthUserHandler))
	router.HandlerFunc(http.MethodPost, "/v1/orders", app.requireActivation(app.createOrderHandler))

	router.HandlerFunc(http.MethodGet, "/v1/admin/orders", app.requireRole("admin", app.getAllOrdersHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/orders/:id", app.requireRole("admin", app.getOrderHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/user/:id/orders", app.requireRole("admin", app.getOrdersByUserHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/admin/orders/:id", app.requireRole("admin", app.editOrderHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/admin/orders/:id", app.requireRole("admin", app.deleteOrderHandler))

	router.HandlerFunc(http.MethodPost, "/v1/admin/roles", app.requireRole("admin", app.createRoleHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/roles", app.requireRole("admin", app.getAllRolesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/roles/:id", app.requireRole("admin", app.getRoleHandler))
	router.HandlerFunc(http.MethodPut, "/v1/admin/roles/:id", app.requireRole("admin", app.updateRoleHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/admin/roles/:id", app.requireRole("admin", app.deleteRoleHandler))
	router.HandlerFunc(http.MethodPut, "/v1/admin/users/:id/role", app.requireRole("admin", app.updateUserRoleHandler))

	router.HandlerFunc(http.MethodPost, "/v1/admin/categories", app.requireRole("admin", app.createCategoryHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/categories", app.requireRole("admin", app.getAllCategoriesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/categories/:id", app.requireRole("admin", app.getCategoryHandler))
	router.HandlerFunc(http.MethodPut, "/v1/admin/categories/:id", app.requireRole("admin", app.updateCategoryHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/admin/categories/:id", app.requireRole("admin", app.deleteCategoryHandler))

	router.HandlerFunc(http.MethodGet, "/v1/products", app.getAllProductsHandler)   // public
	router.HandlerFunc(http.MethodGet, "/v1/products/:slug", app.getProductHandler) // public
	router.HandlerFunc(http.MethodGet, "/v1/products/:slug/category", app.getProductsByCategoryHandler)
	router.HandlerFunc(http.MethodPost, "/v1/admin/products", app.requireRole("admin", app.createProductHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/admin/products/:id", app.requireRole("admin", app.updateProductHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/admin/products/:id", app.requireRole("admin", app.deleteProductHandler))

	return app.authenticate(router)
}
