package main

import (
	"log"

	"github.com/bxcodec/faker"
	"github.com/kubil6y/dukkan-go/internal/data"
	"gorm.io/gorm"
)

type config struct {
	port   string
	env    string
	domain string
	db     struct {
		dsn string
	}
}

func main() {
	var cfg config
	setupFlags(&cfg)

	db, err := connectDatabase(cfg)
	if err != nil {
		log.Fatal("database connection failed")
	}

	seedRolesAndUsers(db)
}

func seedRolesAndUsers(db *gorm.DB) {
	adminRole := &data.Role{Name: "admin"}
	userRole := &data.Role{Name: "user"}
	roles := []*data.Role{adminRole, userRole}

	for _, role := range roles {
		db.Create(role)
	}

	for i := 0; i < 30; i++ {
		if i == 0 {
			user := data.User{
				FirstName:   "admin",
				LastName:    "admin",
				Email:       "lieqb2@gmail.com",
				Address:     faker.Address(),
				IsActivated: true,
				RoleID:      1,
				//Role:        *adminRole,
			}
			user.SetPassword("random")
			continue
		}
		// TODO password
		user := data.User{
			FirstName:   faker.FirstName(),
			LastName:    faker.LastName(),
			Email:       faker.Email(),
			Address:     faker.Address(),
			IsActivated: i%2 == 0,
			RoleID:      2,
			//Role:        *userRole,
		}
		user.SetPassword("random")
	}
}
