package data

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kubil6y/dukkan-go/internal/validator"
	"gorm.io/gorm"
)

type Order struct {
	CoreModel
	UserID        int64       `json:"user_id" gorm:"not null"`
	User          *User       `json:"user,omitempty"`
	PaymentMethod string      `json:"payment_method" gorm:"not null"`
	IsPaid        bool        `json:"is_paid" gorm:"not null"`
	IsDelivered   bool        `json:"is_delivered" gorm:"not null"`
	PaidAt        time.Time   `json:"paid_at" gorm:"not null"`
	TotalPrice    float64     `json:"total_price" gorm:"not null"`
	DeliveredAt   time.Time   `json:"delivered_at" gorm:"not null"`
	OrderItems    []OrderItem `json:"order_items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

type OrderItem struct {
	CoreModel
	OrderID   int64    `json:"order_id" gorm:"not null"`
	Order     *Order   `json:"order,omitempty"`
	ProductID int64    `json:"product_id" gorm:"not null"`
	Product   *Product `json:"product,omitempty"`
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
	err := m.DB.Where("id=?", id).Preload("OrderItems.Product").First(&order).Error
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
	err := m.DB.Scopes(p.PaginatedResults).Preload("OrderItems.Product").Find(&orders).Error
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
	err := m.DB.Where("user_id=?", userID).Scopes(p.PaginatedResults).Preload("OrderItems.Product").Find(&orders).Error
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

type OrderItemDTO struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

type CreateOrderDTO struct {
	PaymentMethod string         `json:"payment_method"`
	OrderItems    []OrderItemDTO `json:"order_items"`
}

func (d *CreateOrderDTO) Validate(v *validator.Validator) {
	v.Check(d.PaymentMethod != "", "payment_method", "must be provided")
	v.Check(In([]string{"cash", "credit"}, strings.ToLower(strings.Trim(d.PaymentMethod, " "))), "payment_method", "must be cash or credit")
	v.Check(len(d.OrderItems) > 0, "order_items", "must be provided")
}

func (m OrderModel) CreateOrder(userID int64, dto CreateOrderDTO) (*Order, error) {
	var order Order
	order.UserID = userID
	order.PaymentMethod = dto.PaymentMethod

	tx := m.DB.Begin()
	err := tx.Create(&order).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var total float64

	for _, item := range dto.OrderItems {
		var product Product
		err := m.DB.Where("id=?", item.ProductID).First(&product).Error
		if err != nil {
			tx.Rollback()
			return nil, errors.New(fmt.Sprintf("product_id: %d does not exist\n", item.ProductID))
		}

		// check product's stock and save new stock count
		newProductCount := product.Count - item.Quantity
		if newProductCount >= 0 {
			product.Count = newProductCount
			tx.Save(&product)
		} else {
			tx.Rollback()
			return nil, ErrOutOfStock
		}

		total += (product.Price * float64(item.Quantity))

		orderItem := OrderItem{
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  item.Quantity,
		}

		order.OrderItems = append(order.OrderItems, orderItem)
	}

	order.TotalPrice = total

	err = tx.Save(&order).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &order, nil
}

func (m OrderModel) Save(order *Order) error {
	return m.DB.Save(order).Error
}
