package data

import "gorm.io/gorm"

type Rating struct {
	CoreModel
	Value     int64 `json:"value" gorm:"not null;check:value>=0 and value<=5"`
	UserID    int64 `json:"user_id" gorm:"not null"`
	ProductID int64 `json:"product_id" gorm:"not null"`
}

type RatingModel struct {
	DB *gorm.DB
}
