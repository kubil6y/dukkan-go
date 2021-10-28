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
	Category    string   `json:"category" gorm:"not null"`
	Image       string   `json:"image" gorm:"not null"`
	Price       float64  `json:"price" gorm:"not null"`
	Count       int64    `json:"count" gorm:"not null"`
	Reviews     []Review `json:"reviews,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Ratings     []Rating `json:"ratings,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

type ProductModel struct {
	DB *gorm.DB
}

func (m ProductModel) Insert(p *Product) error {
	if err := m.DB.Create(p).Error; err != nil {
		// no duplicate check because there are no unique constraints...
		return err
	}
	return nil
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

	err := m.DB.Where("id=?", id).First(&product).Error
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
