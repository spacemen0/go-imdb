package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB initializes a test database using SQLite in-memory
func SetupTestDB() (*gorm.DB, error) {
	// Open an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(&Title{}, &Person{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// TestTitleCRUD tests the CRUD operations for Title
func TestTitleCRUD(t *testing.T) {
	db, err := SetupTestDB()
	assert.NoError(t, err)

	// Test Create
	title := &Title{
		ID:             "tt0000001",
		TitleType:      "movie",
		PrimaryTitle:   "Test Movie",
		OriginalTitle:  "Test Movie Original",
		IsAdult:        false,
		StartYear:      nil,
		EndYear:        nil,
		RuntimeMinutes: nil,
		Genres:         nil,
	}
	err = db.Create(title).Error
	assert.NoError(t, err)

	// Test Read
	var readTitle Title
	err = db.First(&readTitle, "tconst = ?", title.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, title.PrimaryTitle, readTitle.PrimaryTitle)

	// Test Update
	newTitle := "Updated Movie Title"
	readTitle.PrimaryTitle = newTitle
	err = db.Save(&readTitle).Error
	assert.NoError(t, err)

	// Verify Update
	var updatedTitle Title
	err = db.First(&updatedTitle, "tconst = ?", title.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, newTitle, updatedTitle.PrimaryTitle)

	// Test Delete
	err = db.Delete(&updatedTitle).Error
	assert.NoError(t, err)

	// Verify Delete
	var deletedTitle Title
	err = db.First(&deletedTitle, "tconst = ?", title.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
