package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/kubil6y/dukkan-go/internal/data"
	"gorm.io/gorm"
)

func (app *application) seed(db *gorm.DB) {
	seedRolesAndUsers(db)
	seedCategoriesAndProducts(db)
}

func seedRolesAndUsers(db *gorm.DB) {
	adminRole := &data.Role{Name: "admin"}
	userRole := &data.Role{Name: "user"}
	roles := []*data.Role{adminRole, userRole}

	fmt.Println("seeding roles...")
	for _, role := range roles {
		db.Create(role)
	}
	fmt.Println("seeded roles completed!")

	fmt.Println("seeding admin...")
	admin := data.User{
		FirstName:   "admin",
		LastName:    "admin",
		Email:       "lieqb2@gmail.com",
		Address:     "kadikoy/istanbul",
		IsActivated: true,
		RoleID:      1,
	}
	admin.SetPassword("random")
	db.Create(&admin)
	fmt.Println("seeding admin completed")

	fmt.Println("seeding users...")
	for i := 0; i < 10; i++ {
		user := data.User{
			FirstName:   faker.FirstName(),
			LastName:    faker.LastName(),
			Email:       faker.Email(),
			Address:     "somewhere over there",
			IsActivated: i%2 == 0,
			RoleID:      2,
		}
		user.SetPassword("random")
		db.Create(&user)
	}
	fmt.Println("seeded users completed!")
}

func seedCategoriesAndProducts(db *gorm.DB) {
	rand.Seed(time.Now().UnixNano())

	computers := data.Category{Name: "computers"}
	electronics := data.Category{Name: "electronics"}
	clothing := data.Category{Name: "clothing"}
	categories := []data.Category{computers, electronics, clothing}

	fmt.Println("seeding categories...")
	for _, category := range categories {
		db.Create(&category)
	}
	fmt.Println("seeding categories completed!")

	fmt.Println("seeding products...")
	for i := 0; i < 30; i++ {
		product := data.Product{
			Name:        faker.Username(),
			Description: faker.Username(),
			Brand:       faker.Username(),
			Image:       "https://" + faker.Username(),
			Price:       float64(rand.Intn(5000)),
			Count:       int64(rand.Intn(15)),
			CategoryID:  int64(rand.Intn(3) + 1),
		}

		db.Create(&product)
	}
	fmt.Println("seeding products completed!")

}
