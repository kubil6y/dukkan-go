package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/gosimple/slug"
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

	furniture := data.Category{Name: "furniture"}
	electronics := data.Category{Name: "electronics"}
	beauty := data.Category{Name: "beauty"}
	deals := data.Category{Name: "deals"}

	furniture.Slug = slug.Make(furniture.Name)
	beauty.Slug = slug.Make(beauty.Name)
	electronics.Slug = slug.Make(electronics.Name)
	deals.Slug = slug.Make(deals.Name)

	categories := []data.Category{furniture, electronics, beauty, deals}

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
			Image:       "https://m.media-amazon.com/images/I/A1sKFc-P-6L._AC_UL320_.jpg",
			Price:       float64(rand.Intn(5000)),
			Count:       int64(rand.Intn(15)),
			CategoryID:  int64(rand.Intn(len(categories)) + 1),
		}
		product.Slug = data.Slugify(product.Name, 6)

		db.Create(&product)
	}
	fmt.Println("seeding products completed!")

}
