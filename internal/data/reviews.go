package data

import (
	"errors"

	"gorm.io/gorm"
)

type Review struct {
	CoreModel
	Text      string `json:"text" gorm:"not null"`
	UserID    int64  `json:"-" gorm:"not null"`
	User      *User  `json:"user,omitempty"`
	ProductID int64  `json:"-" gorm:"not null"`
}

type ReviewModel struct {
	DB *gorm.DB
}

func (m ReviewModel) Insert(r *Review) error {
	return m.DB.Create(r).Error
}

func (m ReviewModel) GetByID(id int64) (*Review, error) {
	var review Review

	if err := m.DB.Where("id=?", id).First(&review).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &review, nil
}

func (m ReviewModel) Update(r *Review) error {
	return m.DB.Model(r).Updates(r).Error
}

func (m ReviewModel) Delete(r *Review) error {
	return m.DB.Delete(r).Error
}
