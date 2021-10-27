package main

import (
	"github.com/kubil6y/dukkan-go/internal/data"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDatabase(cfg config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.db.dsn))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&data.User{},
		&data.Role{},
		&data.Token{},
		&data.Address{},
	)
}
