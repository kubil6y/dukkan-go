package data

type Role struct {
	CoreModel
	Name   string `json:"name" gorm:"not null"`
	UserID int64  `json:"user_id" gorm:"not null"`
}
