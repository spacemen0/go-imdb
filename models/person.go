package models

import (
	"errors"

	"gorm.io/gorm"
)

// Person represents the person table in the database
type Person struct {
	ID                string  `gorm:"primaryKey;column:nconst"`
	PrimaryName       string  `gorm:"column:primaryName"`
	BirthYear         *int    `gorm:"column:birthYear"`
	DeathYear         *int    `gorm:"column:deathYear"`
	PrimaryProfession *string `gorm:"type:text[];column:primaryProfession"`
	KnownForTitles    *string `gorm:"type:text[];column:knownForTitles"`
}

// CreatePerson creates a new person record in the database
func CreatePerson(db *gorm.DB, person *Person) error {
	return db.Create(person).Error
}

// GetPerson retrieves a person record by its ID
func GetPerson(db *gorm.DB, id string) (*Person, error) {
	var person Person
	if err := db.First(&person, "nconst = ?", id).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

// UpdatePerson updates an existing person record in the database
func UpdatePerson(db *gorm.DB, person *Person) error {
	if db.Model(&Person{}).Where("nconst = ?", person.ID).Updates(person).RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}

// DeletePerson deletes a person record by its ID
func DeletePerson(db *gorm.DB, id string) error {
	if db.Delete(&Person{}, "nconst = ?", id).RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}
