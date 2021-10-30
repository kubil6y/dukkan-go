package main

import (
	"errors"
	"net/http"

	"github.com/kubil6y/dukkan-go/internal/data"
	"github.com/kubil6y/dukkan-go/internal/validator"
)

func (app *application) getOrdersOfAuthUserHandler(w http.ResponseWriter, r *http.Request) {
	v := validator.New()
	p := data.NewPaginate(r, v, 10, 1)

	if data.ValidatePaginate(p, v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user := app.getUserContext(r)

	orders, metadata, err := app.models.Orders.GetAllOrdersByUserID(p, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{
		"orders":   orders,
		"metadata": metadata,
	}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getOrderByIDOfAuthUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.getUserContext(r)

	order, err := app.models.Orders.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case order.UserID != user.ID:
			// NOTE not found is a better response for security reasons,
			// but for now lets keep it this way.
			app.notPermittedResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	e := envelope{"order": order}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getAllOrdersHandler(w http.ResponseWriter, r *http.Request) {
	v := validator.New()
	p := data.NewPaginate(r, v, 10, 1)

	orders, metadata, err := app.models.Orders.GetAllOrders(p)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{
		"orders":   orders,
		"metadata": metadata,
	}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	order, err := app.models.Orders.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	e := envelope{"order": order}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getOrdersByUserHandler(w http.ResponseWriter, r *http.Request) {
	// user_id
	id, err := app.parseIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	p := data.NewPaginate(r, v, 10, 1)

	orders, metadata, err := app.models.Orders.GetAllOrdersByUserID(p, id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{
		"orders":   orders,
		"metadata": metadata,
	}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) editOrderHandler(w http.ResponseWriter, r *http.Request) {}

func (app *application) deleteOrderHandler(w http.ResponseWriter, r *http.Request) {}
