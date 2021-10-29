package data

import (
	"errors"

	"gorm.io/gorm"
)

// TODO
// handle prices decimal package?
// handle data racing problem when purchasing.
// handle relationships

var (
	CategoryElectronics = "electronics"
	CategoryComputers   = "computers"
	CategorySmartHome   = "smart home"
)

type Product struct {
	CoreModel
	Name        string   `json:"name" gorm:"not null"`
	Description string   `json:"description" gorm:"not null"`
	Brand       string   `json:"brand" gorm:"not null"`
	Image       string   `json:"image" gorm:"not null"`
	Price       float64  `json:"price" gorm:"not null"`
	Count       int64    `json:"count" gorm:"not null"`
	CategoryID  int64    `json:"category_id" gorm:"not null"`
	Category    Category `json:"category,omitempty"`
	Reviews     []Review `json:"reviews,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Ratings     []Rating `json:"ratings,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

// UserReviewed() is a function that checks,
// if user reviewed the product or not returns bool
func (p *Product) UserReviewed(user *User) bool {
	if len(p.Reviews) == 0 {
		return false
	}
	for _, review := range p.Reviews {
		if review.UserID == user.ID {
			return true
		}
	}
	return false
}

// UserRated() is a function that checks,
// if user reviewed the product or not returns bool
func (p *Product) UserRated(user *User) bool {
	if len(p.Ratings) == 0 {
		return false
	}
	for _, rating := range p.Ratings {
		if rating.UserID == user.ID {
			return true
		}
	}
	return false
}

type ProductModel struct {
	DB *gorm.DB
}

func (m ProductModel) Insert(p *Product) error {
	return m.DB.Create(p).Error
}

func (m ProductModel) GetAll(p *Paginate) ([]Product, Metadata, error) {
	var products []Product

	err := m.DB.Scopes(p.PaginatedResults).Find(&products).Error
	if err != nil {
		return nil, Metadata{}, nil
	}

	var total int64
	m.DB.Model(&Product{}).Count(&total)
	metadata := CalculateMetadata(p, int(total))
	return products, metadata, nil
}

func (m ProductModel) GetByID(id int64) (*Product, error) {
	var product Product

	// TODO

	err := m.DB.
		Where("id=?", id).
		Preload("Reviews").
		Preload("Reviews.User").
		Preload("Ratings").
		First(&product).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &product, nil
}

func (m ProductModel) Update(p *Product) error {
	return m.DB.Model(p).Updates(p).Error
}

func (m ProductModel) Delete(p *Product) error {
	return m.DB.Delete(p).Error
}
