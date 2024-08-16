package models

import (
	"gorm.io/gorm"
)

// Person represents the person table in the database
type Person struct {
	ID                string   `gorm:"primaryKey;column:nconst;" json:"nConst"`
	PrimaryName       string   `json:"primaryName"`
	BirthYear         *int     `json:"birthYear"`
	DeathYear         *int     `json:"deathYear"`
	PrimaryProfession *string  `gorm:"type:text[];" json:"primaryProfession"`
	KnownForTitles    []*Title `gorm:"many2many:filmography;constraint:OnDelete:CASCADE;" json:"knownForTitles"`
}

// CreatePerson creates a new person record in the database
func CreatePerson(db *gorm.DB, person *Person) error {
	return db.Create(person).Error
}

// GetPerson retrieves a person record by its ID
func GetPerson(db *gorm.DB, id string, verbose bool) (*Person, error) {
	var person Person
	if verbose {
		// If verbose is true, preload all associated data
		if err := db.Preload("KnownForTitles").Preload("KnownForTitles.Actors").First(&person, "nconst = ?", id).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Preload("KnownForTitles", func(db *gorm.DB) *gorm.DB {
			return db.Select("tconst")
		}).First(&person, "nconst = ?", id).Error; err != nil {
			return nil, err
		}
	}
	return &person, nil
}

// UpdatePerson updates an existing person record in the database
func UpdatePerson(db *gorm.DB, person *Person) error {
	if err := db.Model(&Person{}).Where("nconst = ?", person.ID).Updates(person).Error; err != nil {
		return err
	}
	return nil
}

// DeletePerson deletes a person record by its ID
func DeletePerson(db *gorm.DB, id string) error {
	if err := db.Delete(&Person{}, "nconst = ?", id).Error; err != nil {
		return err
	}
	return nil
}
