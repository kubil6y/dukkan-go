package main

import (
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
