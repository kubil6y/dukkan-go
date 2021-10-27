package data

import (
	"errors"

	"gorm.io/gorm"
)

type Address struct {
	CoreModel
	Title  string `json:"title" gorm:"not null"`
	Text   string `json:"text" gorm:"not null"`
	UserID int64  `json:"user_id"`
}

type AddressModel struct {
	DB *gorm.DB
}

func (m AddressModel) GetByID(id int64) (*Address, error) {
	var address Address
	if err := m.DB.Where("id=?", id).First(&address).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &address, nil
}

func (m AddressModel) Insert(a *Address) error {
	return m.DB.Create(&a).Error
}

func (m AddressModel) Update(a *Address) error {
	return m.DB.Model(a).Updates(a).Error
}

func (m AddressModel) Delete(a *Address) error {
	return m.DB.Delete(a).Error
}
