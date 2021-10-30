package data

import (
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
