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

	v.Check(govalidator.IsEmail(d.Email), "email", "must be a valid email")
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
	v.Check(govalidator.IsEmail(d.Email), "email", "must be a valid email")
}
