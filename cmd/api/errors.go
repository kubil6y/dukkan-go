package main

import (
	"fmt"
	"net/http"
)

// logError() logs errors
func (app *application) logError(r *http.Request, err error) {
	app.logger.Errorw(err.Error(),
		"request_method", r.Method,
		"request_url", r.URL.String(),
	)
}

// Dynamic error response generator
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	out := app.outERR(message)
	if err := app.writeJSON(w, status, out, nil); err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// 400 - StatusBadRequest
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// 500 - StatusInternalServerError
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// 404 - StatusNotFound
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// 405 - StatusMethodNotAllowed
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// 422 - StatusUnprocessableEntity
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors interface{}) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// 401 - StatusUnauthorized
func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// 401 - StatusUnauthorized
func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// 401 - StatusUnauthorized
func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// 403 - StatusForbidden
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

// 403 - StatusForbidden
func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

// 403 - StatusForbidden
func (app *application) alreadyReviewedResponse(w http.ResponseWriter, r *http.Request) {
	message := "you have already reviewed this product"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

// 403 - StatusForbidden
func (app *application) alreadyRatedResponse(w http.ResponseWriter, r *http.Request) {
	message := "you have already rated this product"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

// 403 - StatusForbidden
func (app *application) notPurchasedResponse(w http.ResponseWriter, r *http.Request) {
	message := "you have not purchased this product"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

// 401 - StatusBadRequest
func (app *application) outOfStockResponse(w http.ResponseWriter, r *http.Request) {
	message := "product is out of stock"
	app.errorResponse(w, r, http.StatusBadRequest, message)
}

// 429 - StatusTooManyRequests
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}
