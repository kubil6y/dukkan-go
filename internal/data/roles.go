package data

import (
	"errors"

	"gorm.io/gorm"
)

type Role struct {
	CoreModel
	Name string `json:"name" gorm:"uniqueIndex;not null"`
}

type RoleModel struct {
	DB *gorm.DB
}

func (m RoleModel) Insert(r *Role) error {
	if err := m.DB.Create(&r).Error; err != nil {
		switch {
		case IsDuplicateRecord(err):
			return ErrDuplicateRecord
		default:
			return err
		}
	}
	return nil
}

func (m RoleModel) GetAll() ([]*Role, error) {
	var roles []*Role
	if err := m.DB.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (m RoleModel) GetByID(id int64) (*Role, error) {
	var role Role
	if err := m.DB.Where("id=?", id).First(&role).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &role, nil
}

func (m RoleModel) Update(role *Role) error {
	return m.DB.Updates(role).Error
}

func (m RoleModel) Delete(role *Role) error {
	return m.DB.Delete(role).Error
}
