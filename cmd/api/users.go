package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/kubil6y/dukkan-go/internal/data"
	"github.com/kubil6y/dukkan-go/internal/email"
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

	var user data.User
	if err := input.populate(&user); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// NOTE RoleID=3, default user
	// TODO
	user.RoleID = 3

	if err := app.models.Users.Insert(&user); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateRecord):
			v.AddError("email", "a user with the email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// create token
	token, err := app.models.Tokens.New(user.ID, 1*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// send email on the background

	app.background(func() {
		email.ActivationEmail(&user, token.Plaintext)
	})

	e := envelope{"user": user}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
