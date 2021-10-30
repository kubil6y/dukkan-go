package main

import (
	"errors"
	"net/http"

	"github.com/kubil6y/dukkan-go/internal/data"
	"github.com/kubil6y/dukkan-go/internal/validator"
)

func (app *application) createRoleHandler(w http.ResponseWriter, r *http.Request) {
	var input createRoleDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var role data.Role
	input.populate(&role)
	if err := app.models.Roles.Insert(&role); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateRecord):
			v.AddError("name", "already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	e := envelope{"message": "success"}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusCreated, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getAllRolesHandler(w http.ResponseWriter, r *http.Request) {
	v := validator.New()
	p := data.NewPaginate(r, v, 10, 1)

	if data.ValidatePaginate(p, v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	roles, metadata, err := app.models.Roles.GetAll(p)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{
		"roles":    roles,
		"metadata": metadata,
	}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) getRoleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	role, err := app.models.Roles.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	e := envelope{"role": role}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input updateRoleDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	role, err := app.models.Roles.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	input.populate(role)
	if err := app.models.Roles.Update(role); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{"role": role}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteRoleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	role, err := app.models.Roles.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.models.Roles.Delete(role); err != nil {
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

func (app *application) updateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := app.parseIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input updateUserRoleDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByID(userID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	v := validator.New()
	if input.validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	role, err := app.models.Roles.GetByID(input.RoleID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Role = role

	if err := app.models.Users.Update(user); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	out := envelope{"user": user}
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
