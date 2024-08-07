package models

import (
	"errors"

	"gorm.io/gorm"
)

// Title represents the title table in the database
type Title struct {
	ID             string `gorm:"primaryKey;column:tconst"`
	TitleType      string
	PrimaryTitle   string
	OriginalTitle  string
	IsAdult        bool
	StartYear      *int
	EndYear        *int
	RuntimeMinutes *int
	Genres         *string
	Actors         []*Person `gorm:"many2many:filmography;"`
}

// CreateTitle creates a new title record in the database
func CreateTitle(db *gorm.DB, title *Title) error {
	return db.Create(title).Error
}

// GetTitle retrieves a title record by its ID
func GetTitle(db *gorm.DB, id string) (*Title, error) {
	var title Title
	if err := db.Preload("Actors", func(db *gorm.DB) *gorm.DB {
		return db.Select("nconst")
	}).First(&title, "tconst = ?", id).Error; err != nil {
		return nil, err
	}
	return &title, nil
}

// UpdateTitle updates an existing title record in the database
func UpdateTitle(db *gorm.DB, title *Title) error {
	if db.Model(&Title{}).Where("tconst = ?", title.ID).Updates(title).RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}

// DeleteTitle deletes a title record by its ID
func DeleteTitle(db *gorm.DB, id string) error {
	if db.Delete(&Title{}, "tconst = ?", id).RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}
