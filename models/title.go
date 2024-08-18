package models

import (
	"gorm.io/gorm"
)

// Title represents the title table in the database
type Title struct {
	ID             string    `gorm:"primaryKey;column:tconst;" json:"tConst"`
	TitleType      string    `json:"titleType"`
	PrimaryTitle   string    `json:"primaryTitle"`
	OriginalTitle  string    `json:"originalTitle"`
	IsAdult        bool      `json:"isAdult"`
	StartYear      *int      `json:"startYear"`
	EndYear        *int      `json:"endYear"`
	RuntimeMinutes *int      `json:"runtimeMinutes"`
	Genres         *string   `json:"genres"`
	Actors         []*Person `gorm:"many2many:filmography;constraint:OnDelete:CASCADE;" json:"actors"`
}

// CreateTitle creates a new title record in the database
func CreateTitle(db *gorm.DB, title *Title) error {
	return db.Create(title).Error
}

// GetTitle retrieves a title record by its ID
func GetTitle(db *gorm.DB, id string, verbose bool) (*Title, error) {
	var title Title
	if verbose {
		if err := db.Preload("Actors").First(&title, "tconst = ?", id).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Preload("Actors", func(db *gorm.DB) *gorm.DB {
			return db.Select("nconst")
		}).First(&title, "tconst = ?", id).Error; err != nil {
			return nil, err
		}
	}

	return &title, nil
}

// UpdateTitle updates an existing title record in the database
func UpdateTitle(db *gorm.DB, title *Title) error {
	var existingTitle Title
	if err := db.First(&existingTitle, "nconst = ?", title.ID).Error; err != nil {
		return err
	}
	if err := db.Model(&Title{}).Where("tconst = ?", title.ID).Updates(title).Error; err != nil {
		return err
	}
	if err := db.Model(title).Association("Actors").Replace(title.Actors); err != nil {
		return err
	}
	return nil
}

// DeleteTitle deletes a title record by its ID
func DeleteTitle(db *gorm.DB, id string) error {
	if err := db.Delete(&Title{}, "tconst = ?", id).Error; err != nil {
		return err
	}
	return nil
}
