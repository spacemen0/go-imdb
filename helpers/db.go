package helpers

import (
	"fmt"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"spacemen0.github.com/models"
)

var DB *gorm.DB

func InitDB() {
	if testing.Testing() {
		var err error
		DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			Log.Fatal("Failed to connect to in-memory database:", err)
		}
		err = DB.AutoMigrate(&models.Title{}, &models.Person{})
		if err != nil {
			Log.Fatal("Failed to migrate database schema:", err)
		}
	} else {
		dbConfig := AppConfig.Database

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.Port, dbConfig.SSLMode)
		var err error
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			Log.Fatal("Failed to connect to database:", err)
		}
		err = DB.AutoMigrate(&models.Title{}, &models.Person{})
		if err != nil {
			Log.Fatal("Failed to migrate database schema:", err)
		}
		fullTextMigrations()
	}
}

func fullTextMigrations() {
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_person_name ON people USING gin(to_tsvector('english', primary_name))").Error; err != nil {
		Log.Fatal("Failed to create index on people:", err)
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_title_text ON titles USING gin(to_tsvector('english', primary_title || ' ' || original_title))").Error; err != nil {
		Log.Fatal("Failed to create index on titles:", err)
	}
}
