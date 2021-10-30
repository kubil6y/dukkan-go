package main

import (
	"errors"
	"net/http"

	"github.com/kubil6y/dukkan-go/internal/data"
	"github.com/kubil6y/dukkan-go/internal/validator"
)

func (app *application) createRatingHandler(w http.ResponseWriter, r *http.Request) {
	slug := app.parseSlugParam(r)

	var input ratingDTO
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

	if product.UserRated(user) {
		app.alreadyRatedResponse(w, r)
		return
	}

	rating := data.Rating{}
	rating.Value = input.Value
	rating.UserID = user.ID
	rating.ProductID = product.ID

	if err := app.models.Ratings.Insert(&rating); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{"rating": rating}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) updateRatingHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input ratingDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	rating, err := app.models.Ratings.GetByID(id)
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
	if rating.UserID != user.ID {
		// user is not the owner of this rating
		app.notPermittedResponse(w, r)
		return
	}

	rating.Value = input.Value
	if err := app.models.Ratings.Update(rating); err != nil {
		app.serverErrorResponse(w, r, err)
	}

	e := envelope{"rating": rating}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) deleteRatingHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	rating, err := app.models.Ratings.GetByID(id)
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
	if rating.UserID != user.ID {
		// user is not the owner of this rating
		app.notPermittedResponse(w, r)
		return
	}

	if err := app.models.Ratings.Delete(rating); err != nil {
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
