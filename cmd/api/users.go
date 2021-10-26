package main

import (
	"fmt"
	"net/http"

	"github.com/kubil6y/dukkan-go/internal/validator"
)

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	var input registerDTO

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	out := app.outOK(input)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login handler")
}
