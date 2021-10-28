package data

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrDuplicateRecord = errors.New("duplicate record")
)

type CoreModel struct {
	ID        int64     `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

type Models struct {
	Users    UserModel
	Tokens   TokenModel
	Roles    RoleModel
	Products ProductModel
	Reviews  ReviewModel
}

func NewModels(db *gorm.DB) Models {
	return Models{
		Users:    UserModel{DB: db},
		Tokens:   TokenModel{DB: db},
		Roles:    RoleModel{DB: db},
		Products: ProductModel{DB: db},
		Reviews:  ReviewModel{DB: db},
	}
}
