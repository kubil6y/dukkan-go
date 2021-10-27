package main

import (
	"github.com/asaskevich/govalidator"
	"github.com/kubil6y/dukkan-go/internal/data"
	"github.com/kubil6y/dukkan-go/internal/validator"
)

type registerDTO struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func (d *registerDTO) validate(v *validator.Validator) {
	v.Check(d.FirstName != "", "first_name", "must be provided")
	v.Check(d.LastName != "", "last_name", "must be provided")
	v.Check(d.Email != "", "email", "must be provided")
	v.Check(d.Password != "", "password", "must be provided")
	v.Check(d.PasswordConfirm != "", "password_confirm", "must be provided")
	v.Check(d.Password == d.PasswordConfirm, "password", "passwords do not match")

	v.Check(govalidator.IsEmail(d.Email), "email", "must be a valid email address")
	v.Check(len(d.FirstName) >= 2, "first_name", "must be at least two characters")
	v.Check(len(d.LastName) >= 2, "last_name", "must be at least two characters")
	v.Check(len(d.Password) >= 6, "password", "must be at least six characters")
}

func (d *registerDTO) populate(user *data.User) error {
	user.FirstName = d.FirstName
	user.LastName = d.LastName
	user.Email = d.Email
	user.IsActivated = false
	user.IsAdmin = false
	user.RoleID = 2
	// hashing
	err := user.SetPassword(d.Password)
	return err
}

type createAuthenticationTokenDTO struct {
	Email    string
	Password string
}

func (d *createAuthenticationTokenDTO) validate(v *validator.Validator) {
	v.Check(d.Email != "", "email", "must be provided")
	v.Check(d.Password != "", "password", "must be provided")
	v.Check(govalidator.IsEmail(d.Email), "email", "must be a valid email address")
}

type createRoleDTO struct {
	Name string `json:"name"`
}

func (d *createRoleDTO) validate(v *validator.Validator) {
	v.Check(d.Name != "", "name", "must be provided")
	v.Check(len(d.Name) > 2, "name", "must be longer than two characters")
}

func (d *createRoleDTO) populate(role *data.Role) {
	role.Name = d.Name
}

// seems like duplication but in the future,
// role fields might change, and there might be
// fields that we dont allow them to change.
type updateRoleDTO struct {
	Name string `json:"name"`
}

func (d *updateRoleDTO) validate(v *validator.Validator) {
	v.Check(d.Name != "", "name", "must be provided")
	v.Check(len(d.Name) > 2, "name", "must be longer than two characters")
}

func (d *updateRoleDTO) populate(role *data.Role) {
	role.Name = d.Name
}

type updateUserRoleDTO struct {
	RoleID int64 `json:"role_id"`
}

func (d *updateUserRoleDTO) validate(v *validator.Validator) {
	v.Check(d.RoleID != 0, "role_id", "must be provided")
	v.Check(d.RoleID > 0, "role_id", "invalid role id value")
}

func validateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type activateAccountDTO struct {
	Code string `json:"code"`
}

func (d *activateAccountDTO) validate(v *validator.Validator) {
	v.Check(d.Code != "", "code", "must be provided")
	v.Check(len(d.Code) == 26, "code", "must be 26 bytes long")
}

type generateActivationTokenDTO struct {
	Email string `json:"email"`
}

func (d *generateActivationTokenDTO) validate(v *validator.Validator) {
	v.Check(d.Email != "", "email", "must be provided")
	v.Check(govalidator.IsEmail(d.Email), "email", "must be a valid email address")
}
