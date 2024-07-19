package helpers

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"spacemen0.github.com/models"
)

var db *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=Fracture123 dbname=IMDB port=5678 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Title{}, &models.Person{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
}

func GetDB() *gorm.DB {
	return db
}
