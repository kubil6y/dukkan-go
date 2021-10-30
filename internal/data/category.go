package data

import (
	"errors"

	"gorm.io/gorm"
)

type Category struct {
	CoreModel
	Name     string    `json:"name" gorm:"uniqueIndex;not null"`
	Slug     string    `json:"slug" gorm:"uniqueIndex;not null"`
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`
}

type CategoryModel struct {
	DB *gorm.DB
}

func (m CategoryModel) Insert(c *Category) error {
	if err := m.DB.Create(c).Error; err != nil {
		switch {
		case IsDuplicateRecord(err):
			return ErrDuplicateRecord
		default:
			return err
		}
	}
	return nil
}

func (m CategoryModel) GetAll() ([]Category, error) {
	var categories []Category

	err := m.DB.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (m CategoryModel) GetByID(id int64) (*Category, error) {
	var category Category
	if err := m.DB.Where("id=?", id).First(&category).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &category, nil
}

// TODO
func (m CategoryModel) GetBySlug(slug string) (*Category, error) {
	var category Category
	if err := m.DB.Where("slug=?", slug).First(&category).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &category, nil
}

func (m CategoryModel) GetByName(name string) (*Category, error) {
	var category Category
	if err := m.DB.Where("name=?", name).First(&category).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &category, nil
}

func (m CategoryModel) Update(c *Category) error {
	return m.DB.Updates(c).Error
}

func (m CategoryModel) Delete(c *Category) error {
	return m.DB.Delete(c).Error
}
