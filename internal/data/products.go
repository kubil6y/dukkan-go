package data

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	CategoryElectronics = "electronics"
	CategoryComputers   = "computers"
	CategorySmartHome   = "smart home"
)

type ProductWrapper struct {
	Product
	RatingAverage float64 `json:"rating_average"`
	RatingCount   int     `json:"rating_count"`
	ReviewCount   int     `json:"review_count"`
}

func NewProductWrapper(product *Product) *ProductWrapper {
	return &ProductWrapper{
		Product:       *product,
		RatingAverage: product.CalculateRating(),
		RatingCount:   len(product.Ratings),
		ReviewCount:   len(product.Reviews),
	}
}

type Product struct {
	CoreModel
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"uniqueIndex;not null"`
	Description string    `json:"description" gorm:"not null"`
	Brand       string    `json:"brand" gorm:"not null"`
	Image       string    `json:"image" gorm:"not null"`
	Price       float64   `json:"price" gorm:"not null"`
	Count       int64     `json:"count" gorm:"not null"`
	CategoryID  int64     `json:"category_id" gorm:"not null"`
	Category    *Category `json:"category,omitempty"`
	Reviews     []Review  `json:"reviews,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Ratings     []Rating  `json:"ratings,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

func (p *Product) CalculateRating() float64 {
	if len(p.Ratings) == 0 {
		return 0
	}

	var rating int64
	for _, v := range p.Ratings {
		rating += v.Value
	}
	return float64(rating) / float64(len(p.Ratings))
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

func (m ProductModel) GetAll(p *Paginate, searchTerm string) ([]Product, Metadata, error) {
	var products []Product

	err := m.DB.Scopes(p.PaginatedResults).Preload("Category").Where("name ILIKE ?", fmt.Sprintf("%%%s%%", searchTerm)).Find(&products).Error
	if err != nil {
		return nil, Metadata{}, err
	}

	var total int64
	m.DB.Model(&Product{}).Where("name ILIKE ?", fmt.Sprintf("%%%s%%", searchTerm)).Count(&total)
	metadata := CalculateMetadata(p, int(total))
	return products, metadata, nil
}

func (m ProductModel) GetBySlug(slug string) (*Product, error) {
	var product Product
	err := m.DB.
		Preload("Reviews").
		Preload("Ratings").
		Where("slug=?", slug).
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

func (m ProductModel) GetByCategory(p *Paginate, categoryID int64) ([]Product, Metadata, error) {
	var products []Product
	err := m.DB.Scopes(p.PaginatedResults).Where("category_id=?", categoryID).Find(&products).Error
	if err != nil {
		return nil, Metadata{}, err
	}

	var total int64
	m.DB.Model(&Product{}).Where("category_id=?", categoryID).Count(&total)
	metadata := CalculateMetadata(p, int(total))
	return products, metadata, nil
}

func (m ProductModel) GetByID(id int64) (*Product, error) {
	var product Product

	err := m.DB.
		Where("id=?", id).First(&product).Error
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
