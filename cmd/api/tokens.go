package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kubil6y/dukkan-go/internal/data"
	"github.com/kubil6y/dukkan-go/internal/email"
	"github.com/kubil6y/dukkan-go/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input createAuthenticationTokenDTO

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	matches, err := user.ComparePassword(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !matches {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Deleting auth tokens, keeping only last 5 on the background with a goroutine.
	app.background(func() {
		fmt.Println("trying to delete tokens") // TODO
		if err := app.models.Tokens.KeepLastFiveAuthTokens(user.ID); err != nil {
			app.logger.Error("Error deleting tokens")
		}
	})

	e := envelope{"authentication_token": map[string]interface{}{
		"token":  token.Plaintext,
		"expiry": token.Expiry,
	}}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) activateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var input activateAccountDTO

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.Code)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.IsActivated = true

	if err := app.models.Tokens.ActivateUserAndDeleteToken(input.Code, user); err != nil {
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

func (app *application) generateActivationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input generateActivationTokenDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := app.models.Tokens.New(user.ID, 1*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		email.ActivationEmail(user, token.Plaintext)
	})

	e := envelope{
		"code": token.Plaintext,
	}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
