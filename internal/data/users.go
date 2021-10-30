package data

import (
	"crypto/sha256"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var AnonUser = &User{}

type User struct {
	CoreModel
	FirstName   string   `json:"first_name" gorm:"not null"`
	LastName    string   `json:"last_name" gorm:"not null"`
	Email       string   `json:"email" gorm:"uniqueIndex;not null"`
	Password    []byte   `json:"-" gorm:"not null"`
	Address     string   `json:"address" gorm:"not null"`
	IsActivated bool     `json:"is_activated" gorm:"default:false;not null"`
	RoleID      int64    `json:"-" gorm:"not null"`
	Role        *Role    `json:"role"`
	Tokens      []Token  `json:"tokens,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Reviews     []Review `json:"reviews,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
	Ratings     []Rating `json:"ratings,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
	Orders      []Order  `json:"orders,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
}

func (u *User) IsAnon() bool {
	return u == AnonUser
}

func (u *User) FullName() string {
	return strings.Title(u.FirstName + " " + u.LastName)
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

func (u *User) DidOrderProduct(product *Product) (bool, error) {
	if u.Orders == nil {
		return false, errors.New("could not load orders for the user.")
	}

	// check if user has bought the product
	for _, order := range u.Orders {
		for _, orderItem := range order.OrderItems {
			if orderItem.ProductID == product.ID {
				return true, nil
			}
		}
	}
	return false, nil
}

type UserModel struct {
	DB *gorm.DB
}

func (m UserModel) Insert(u *User) error {
	if err := m.DB.Create(u).Error; err != nil {
		switch {
		case IsDuplicateRecord(err):
			return ErrDuplicateRecord
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetAll(p *Paginate) ([]User, Metadata, error) {
	var users []User

	err := m.DB.Scopes(p.PaginatedResults).Preload("Role").Find(&users).Error
	if err != nil {
		return nil, Metadata{}, nil
	}

	var total int64
	m.DB.Model(&User{}).Count(&total)
	metadata := CalculateMetadata(p, int(total))
	return users, metadata, nil
}

func (m UserModel) GetByID(id int64) (*User, error) {
	var user User
	err := m.DB.Preload("Role").Where("id=?", id).First(&user).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	var user User
	err := m.DB.Where("email=?", email).First(&user).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) Update(u *User) error {
	return m.DB.Model(u).Updates(u).Error
}

func (m UserModel) Delete(u *User) error {
	return m.DB.Delete(u).Error
}

func (m UserModel) GetForToken(scope string, tokenPlaintext string) (*User, error) {
	sizedTokenHash := sha256.Sum256([]byte(tokenPlaintext))
	tokenHash := sizedTokenHash[:]

	var token Token
	err := m.DB.Where("hash=? and scope=? and expiry > ?", tokenHash, scope, time.Now()).First(&token).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	var user User
	err = m.DB.
		Where("id=?", token.UserID).
		Preload("Role").
		First(&user).Error

	return &user, nil
}

func (m UserModel) GetUserWithOrders(id int64) (*User, error) {
	var user User
	err := m.DB.Preload("Orders.OrderItems").Where("id = ?", id).First(&user).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
