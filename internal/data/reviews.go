package data

import "gorm.io/gorm"

type Review struct {
	CoreModel
	Text      string `json:"text" gorm:"not null"`
	UserID    int64  `json:"user_id" gorm:"not null"`
	ProductID int64  `json:"product_id" gorm:"not null"`
}

type ReviewModel struct {
	DB *gorm.DB
}
