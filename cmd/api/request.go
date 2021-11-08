package main

import (
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gosimple/slug"
	"github.com/kubil6y/dukkan-go/internal/data"
	"github.com/kubil6y/dukkan-go/internal/validator"
)

// sanitize() trims spaces and transforms strings to lowercase
func sanitize(s string) string {
	return strings.ToLower(strings.Trim(s, " "))
}

type registerDTO struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Address         string `json:"address"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func (d *registerDTO) validate(v *validator.Validator) {
	v.Check(d.FirstName != "", "first_name", "must be provided")
	v.Check(d.LastName != "", "last_name", "must be provided")
	v.Check(d.Email != "", "email", "must be provided")
	v.Check(d.Address != "", "email", "must be provided")
	v.Check(d.Password != "", "password", "must be provided")
	v.Check(d.PasswordConfirm != "", "password_confirm", "must be provided")
	v.Check(d.Password == d.PasswordConfirm, "password", "passwords do not match")

	v.Check(govalidator.IsEmail(d.Email), "email", "must be a valid email address")
	v.Check(len(d.FirstName) >= 2, "first_name", "must be at least two characters")
	v.Check(len(d.LastName) >= 2, "last_name", "must be at least two characters")
	v.Check(len(d.Password) >= 6, "password", "must be at least six characters")
	v.Check(len(d.Address) > 12, "address", "must be longer than 12 characters")
}

func (d *registerDTO) populate(user *data.User) error {
	user.FirstName = sanitize(d.FirstName)
	user.LastName = sanitize(d.LastName)
	// NOTE email domains might be case sensitive
	user.Email = strings.ToLower(d.Email)
	user.Address = sanitize(d.Address)
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
	role.Name = sanitize(d.Name)
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
	role.Name = sanitize(d.Name)
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

type editProfileDTO struct {
	FirstName       *string `json:"first_name"`
	LastName        *string `json:"last_name"`
	Email           *string `json:"email"`
	Password        *string `json:"password"`
	PasswordConfirm *string `json:"password_confirm"`
	Address         *string `json:"address"`
}

func (d *editProfileDTO) validate(v *validator.Validator) {
	if d.FirstName != nil {
		v.Check(*d.FirstName != "", "first_name", "must be provided")
		v.Check(len(*d.FirstName) > 2, "first_name", "must be longer then two characters")
	}

	if d.LastName != nil {
		v.Check(*d.LastName != "", "last_name", "must be provided")
		v.Check(len(*d.LastName) > 2, "last_name", "must be longer then two characters")
	}

	if d.Email != nil {
		v.Check(*d.FirstName != "", "email", "must be provided")
		v.Check(govalidator.IsEmail(*d.Email), "email", "must be a valid email address")
	}

	if d.PasswordConfirm != nil {
		if d.Password == nil {
			v.Check(false, "password", "must be provided")
			return
		}
	}

	if d.Password != nil {
		if d.PasswordConfirm == nil {
			v.Check(false, "password_confirm", "must be provided")
			return
		}
		v.Check(*d.Password != "", "password", "must be provided")
		v.Check(*d.PasswordConfirm != "", "password_confirm", "must be provided")
		v.Check(len(*d.Password) >= 6, "password", "must be at least six characters")
		v.Check(*d.Password == *d.PasswordConfirm, "password", "passwords do not match")
	}

	if d.Address != nil {
		v.Check(*d.Address != "", "address", "must be provided")
		v.Check(len(*d.Address) > 12, "address", "must be longer than 12 characters")
	}
}

func (d *editProfileDTO) populate(user *data.User) {
	if d.FirstName != nil {
		user.FirstName = sanitize(*d.FirstName)
	}

	if d.LastName != nil {
		user.LastName = sanitize(*d.LastName)
	}
	if d.Email != nil {
		user.Email = sanitize(*d.Email)
	}

	if d.Password != nil {
		user.SetPassword(*d.Password)
	}

	if d.Address != nil {
		user.Address = sanitize(*d.Address)
	}
}

type createProductDTO struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Brand        string  `json:"brand"`
	CategoryName string  `json:"category_name"`
	Image        string  `json:"image"`
	Price        float64 `json:"price"`
	Count        int64   `json:"count"`
}

func (d *createProductDTO) validate(v *validator.Validator) {
	v.Check(d.Name != "", "name", "must be provided")
	v.Check(d.Description != "", "description", "must be provided")
	v.Check(d.Brand != "", "brand", "must be provided")
	v.Check(d.CategoryName != "", "category_name", "must be provided")
	v.Check(d.Image != "", "image", "must be provided")
	v.Check(d.Price != 0, "price", "must be provided")
	v.Check(d.Count != 0, "count", "must be provided")

	v.Check(govalidator.IsURL(d.Image), "image", "must be valid URL")
	v.Check(d.Price >= 0, "price", "must be valid value")
	v.Check(d.Count >= 0, "count", "must be valid value")
}

func (d *createProductDTO) populate(product *data.Product) {
	product.Name = sanitize(d.Name)
	product.Slug = data.Slugify(product.Name, 6)
	product.Description = sanitize(d.Description)
	product.Brand = sanitize(d.Brand)
	product.Image = sanitize(d.Image)
	product.Price = d.Price
	product.Count = d.Count
}

type updateProductDTO struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	Brand        *string  `json:"brand"`
	CategoryName *string  `json:"category_name"`
	Image        *string  `json:"image"`
	Price        *float64 `json:"price"`
	Count        *int64   `json:"count"`
}

func (d *updateProductDTO) validate(v *validator.Validator) {
	if d.Image != nil {
		v.Check(govalidator.IsURL(*d.Image), "image", "must be valid URL")
	}
	if d.Price != nil {
		v.Check(*d.Price >= 0, "price", "must be valid value")
	}
	if d.Count != nil {
		v.Check(*d.Count >= 0, "count", "must be valid value")
	}
}

func (d *updateProductDTO) populate(product *data.Product) {
	if d.Name != nil {
		product.Name = sanitize(*d.Name)
		product.Slug = data.Slugify(product.Name, 6)
	}
	if d.Description != nil {
		product.Description = *d.Description
	}
	if d.Brand != nil {
		product.Brand = *d.Brand
	}
	if d.Image != nil {
		product.Image = *d.Image
	}
	if d.Price != nil {
		product.Price = *d.Price
	}
	if d.Count != nil {
		product.Count = *d.Count
	}
}

type reviewDTO struct {
	Text string `json:"text"`
}

func (d *reviewDTO) validate(v *validator.Validator) {
	v.Check(d.Text != "", "text", "must be provided")
	v.Check(len(d.Text) > 3, "text", "must be longer than three characters")
}

type ratingDTO struct {
	Value int64 `json:"rating"`
}

func (d *ratingDTO) validate(v *validator.Validator) {
	v.Check(d.Value >= 1, "rating", "must be greater than or equal to one")
	v.Check(d.Value <= 5, "rating", "must be less than or equal to five")
}

type categoryDTO struct {
	Name string `json:"name"`
}

func (d categoryDTO) validate(v *validator.Validator) {
	v.Check(d.Name != "", "name", "must be provided")
	v.Check(len(d.Name) > 3, "name", "must be longer than three characters")
}

func (d categoryDTO) populate(category *data.Category) {
	category.Name = sanitize(d.Name)
	category.Slug = slug.Make(category.Name)
}

type editOrderDTO struct {
	PaymentMethod *string `json:"payment_method"`
	IsPaid        *bool   `json:"is_paid"`
	IsDelivered   *bool   `json:"is_delivered"`
}

func (d *editOrderDTO) validate(v *validator.Validator) {
	if d.PaymentMethod != nil {
		v.Check(data.In([]string{"cash", "credit"}, strings.ToLower(strings.Trim(*d.PaymentMethod, " "))), "payment_method", "must be cash or credit")
	}
}

func (d *editOrderDTO) populate(order *data.Order) {
	if d.PaymentMethod != nil {
		order.PaymentMethod = *d.PaymentMethod
	}

	if d.IsPaid != nil {
		if *d.IsPaid == true {
			order.IsPaid = true
			order.PaidAt = time.Now()
		}
	}

	if d.IsDelivered != nil {
		if *d.IsDelivered == true {
			order.IsDelivered = true
			order.DeliveredAt = time.Now()
		}
	}
}
