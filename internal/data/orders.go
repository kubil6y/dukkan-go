package data

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	CoreModel
	UserID        int64       `json:"user_id" gorm:"not null"`
	User          *User       `json:"user"`
	PaymentMethod string      `json:"payment_method" gorm:"not null"`
	IsPaid        bool        `json:"is_paid" gorm:"not null"`
	IsDelivered   bool        `json:"is_delivered" gorm:"not null"`
	PaidAt        time.Time   `json:"paid_at" gorm:"not null"`
	DeliveredAt   time.Time   `json:"delivered_at" gorm:"not null"`
	OrderItems    []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

type OrderItem struct {
	CoreModel
	OrderID   int64    `json:"order_id" gorm:"not null"`
	Order     *Order   `json:"order"`
	ProductID int64    `json:"product_id" gorm:"not null"`
	Product   *Product `json:"product"`
	Quantity  int64    `json:"quantity" gorm:"not null"`
}

type OrderModel struct {
	DB *gorm.DB
}

type OrderItemModel struct {
	DB *gorm.DB
}

func (m OrderModel) GetByID(id int64) (*Order, error) {
	var order Order
	err := m.DB.Where("id=?", id).First(&order).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &order, nil
}

func (m OrderModel) GetAllOrders(p *Paginate) ([]Order, Metadata, error) {
	var orders []Order
	err := m.DB.Scopes(p.PaginatedResults).Find(&orders).Error
	if err != nil {
		return nil, Metadata{}, err
	}

	var total int64
	m.DB.Model(&Order{}).Count(&total)
	metadata := CalculateMetadata(p, int(total))
	return orders, metadata, nil
}

func (m OrderModel) GetAllOrdersByUserID(p *Paginate, userID int64) ([]Order, Metadata, error) {
	var orders []Order
	err := m.DB.Where("user_id=?", userID).Scopes(p.PaginatedResults).Find(&orders).Error
	if err != nil {
		return nil, Metadata{}, err
	}

	var total int64
	m.DB.Model(&Order{}).Where("user_id=?", userID).Count(&total)
	metadata := CalculateMetadata(p, int(total))
	return orders, metadata, nil
}

func (m OrderModel) Insert(o *Order) error {
	return m.DB.Create(o).Error
}

func (m OrderModel) Update(o *Order) error {
	return m.DB.Updates(o).Error
}

func (m OrderModel) Delete(o *Order) error {
	return m.DB.Delete(o).Error
}
