package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Person{}, &Title{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestCreatePerson(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Create a new person
	person := &Person{
		ID:          "nm00001",
		PrimaryName: "John Doe",
		BirthYear:   nil,
		DeathYear:   nil,
	}

	err = CreatePerson(db, person)
	assert.NoError(t, err)

	// Verify the person was created
	var retrievedPerson Person
	err = db.First(&retrievedPerson, "nconst = ?", "nm00001").Error
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", retrievedPerson.PrimaryName)
}

func TestReadPerson(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Create a person to retrieve later
	person := &Person{
		ID:          "nm00002",
		PrimaryName: "Jane Doe",
		BirthYear:   nil,
		DeathYear:   nil,
	}

	err = CreatePerson(db, person)
	assert.NoError(t, err)

	// Retrieve the person
	retrievedPerson, err := GetPerson(db, "nm00002", false)
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", retrievedPerson.PrimaryName)
}

func TestUpdatePerson(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Create a person to update
	person := &Person{
		ID:          "nm00003",
		PrimaryName: "Mark Smith",
		BirthYear:   nil,
		DeathYear:   nil,
	}

	err = CreatePerson(db, person)
	assert.NoError(t, err)

	// Update the person's name
	person.PrimaryName = "Mark Johnson"
	err = UpdatePerson(db, person)
	assert.NoError(t, err)

	// Verify the update
	var updatedPerson Person
	err = db.First(&updatedPerson, "nconst = ?", "nm00003").Error
	assert.NoError(t, err)
	assert.Equal(t, "Mark Johnson", updatedPerson.PrimaryName)
}

func TestDeletePerson(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Create a person to delete
	person := &Person{
		ID:          "nm00004",
		PrimaryName: "Sarah Connor",
		BirthYear:   nil,
		DeathYear:   nil,
	}

	err = CreatePerson(db, person)
	assert.NoError(t, err)

	// Delete the person
	err = DeletePerson(db, "nm00004")
	assert.NoError(t, err)

	// Verify the person was deleted
	var deletedPerson Person
	err = db.First(&deletedPerson, "nconst = ?", "nm00004").Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestPersonWithTitles(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Create titles
	title1 := &Title{
		ID:            "tt00001",
		PrimaryTitle:  "Test Title 1",
		OriginalTitle: "Test Original Title 1",
		TitleType:     "movie",
	}

	title2 := &Title{
		ID:            "tt00002",
		PrimaryTitle:  "Test Title 2",
		OriginalTitle: "Test Original Title 2",
		TitleType:     "movie",
	}

	err = db.Create(title1).Error
	assert.NoError(t, err)
	err = db.Create(title2).Error
	assert.NoError(t, err)

	// Create a person with known titles
	person := &Person{
		ID:             "nm00005",
		PrimaryName:    "John Wick",
		KnownForTitles: []*Title{title1, title2},
	}

	err = CreatePerson(db, person)
	assert.NoError(t, err)

	// Retrieve the person with associated titles
	retrievedPerson, err := GetPerson(db, "nm00005", true)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(retrievedPerson.KnownForTitles))
	assert.Equal(t, "Test Title 1", retrievedPerson.KnownForTitles[0].PrimaryTitle)
	assert.Equal(t, "Test Title 2", retrievedPerson.KnownForTitles[1].PrimaryTitle)
}
