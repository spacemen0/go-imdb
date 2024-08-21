package helpers

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"spacemen0.github.com/models"
)

var db *gorm.DB

func InitDB() {
	if testing.Testing() {
		var err error
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			Log.Fatal("Failed to connect to in-memory database:", err)
		}
		err = db.AutoMigrate(&models.Title{}, &models.Person{})
		if err != nil {
			Log.Fatal("Failed to migrate database schema:", err)
		}
	} else {
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		sslmode := os.Getenv("SSL_MODE")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			Log.Fatal("Failed to connect to database:", err)
		}
		err = db.AutoMigrate(&models.Title{}, &models.Person{})
		if err != nil {
			Log.Fatal("Failed to migrate database schema:", err)
		}
		fullTextMigrations()
	}
}

func GetDB() *gorm.DB {
	return db
}

func fullTextMigrations() {
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_person_name ON people USING gin(to_tsvector('english', primary_name))").Error; err != nil {
		Log.Fatal("Failed to create index on people:", err)
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_title_text ON titles USING gin(to_tsvector('english', primary_title || ' ' || original_title))").Error; err != nil {
		Log.Fatal("Failed to create index on titles:", err)
	}
}
