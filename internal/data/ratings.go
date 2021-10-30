package data

import (
	"errors"

	"gorm.io/gorm"
)

type Rating struct {
	CoreModel
	Value     int64 `json:"value" gorm:"not null;check:value>=0 and value<=5"`
	UserID    int64 `json:"user_id" gorm:"not null"`
	ProductID int64 `json:"product_id" gorm:"not null"`
}

type RatingModel struct {
	DB *gorm.DB
}

func (m RatingModel) Insert(r *Rating) error {
	return m.DB.Create(r).Error
}

func (m RatingModel) GetByID(id int64) (*Rating, error) {
	var rating Rating

	if err := m.DB.Where("id=?", id).First(&rating).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &rating, nil
}

func (m RatingModel) Update(r *Rating) error {
	return m.DB.Model(r).Updates(r).Error
}

func (m RatingModel) Delete(r *Rating) error {
	return m.DB.Delete(r).Error
}
