package data

import (
	"errors"
	"time"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrDuplicateRecord = errors.New("duplicate record")
)

type CoreModel struct {
	ID        int64     `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Models struct {
	Users UserModel
}
