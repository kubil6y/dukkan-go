package data

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var AnonUser = &User{}

type User struct {
	Model
	FirstName   string `json:"first_name" gorm:"not null"`
	LastName    string `json:"last_name" gorm:"not null"`
	Email       string `json:"email" gorm:"uniqueIndex;not null"`
	Password    []byte `json:"password" gorm:"not null"` // TODO
	IsActivated bool   `json:"is_activated" gorm:"default:false;not null"`
	IsAdmin     bool   `json:"is_admin" gorm:"default:false;not null"`
	Role        Role   `json:"role" gorm:"not null"`
}

func (u *User) IsAnon() bool {
	return u == AnonUser
}

func (u *User) SetPassword(plain string) error {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = h
	return nil
}

func (u *User) ComparePassword(plain string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(u.Password, []byte(plain)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

type UserModel struct {
	DB *gorm.DB
}

func (m UserModel) Insert(u *User) error {
	return m.DB.Create(u).Error
}

func (m UserModel) GetByID(id int64) (*User, error) {
	var user User
	err := m.DB.Where("id=?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	var user User
	err := m.DB.Where("email=?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m UserModel) Update(u *User) error {
	return m.DB.Model(u).Updates(u).Error
}

func (m UserModel) Delete(u *User) error {
	return m.DB.Delete(u).Error
}
