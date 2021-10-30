package main

import (
	"errors"
	"net/http"

	"github.com/kubil6y/dukkan-go/internal/data"
	"github.com/kubil6y/dukkan-go/internal/validator"
)

func (app *application) createReviewHandler(w http.ResponseWriter, r *http.Request) {
	slug := app.parseSlugParam(r)

	var input reviewDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// product check...
	product, err := app.models.Products.GetBySlug(slug)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user := app.getUserContext(r)
	user, err = app.models.Users.GetUserWithOrders(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	didOrder, err := user.DidOrderProduct(product)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !didOrder {
		app.notPurchasedResponse(w, r)
		return
	}

	if product.UserReviewed(user) {
		app.alreadyReviewedResponse(w, r)
		return
	}

	review := data.Review{}
	review.Text = input.Text
	review.UserID = user.ID
	review.ProductID = product.ID

	if err := app.models.Reviews.Insert(&review); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{"review": review}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input reviewDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	review, err := app.models.Reviews.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user := app.getUserContext(r)
	if review.UserID != user.ID {
		// user is not the owner of this review
		app.notPermittedResponse(w, r)
		return
	}

	review.Text = input.Text
	if err := app.models.Reviews.Update(review); err != nil {
		app.serverErrorResponse(w, r, err)
	}

	e := envelope{"review": review}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	review, err := app.models.Reviews.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user := app.getUserContext(r)
	if review.UserID != user.ID {
		// user is not the owner of this review
		app.notPermittedResponse(w, r)
		return
	}

	if err := app.models.Reviews.Delete(review); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{"message": "success"}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
