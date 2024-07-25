package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"spacemen0.github.com/helpers"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Person struct {
	ID                string `gorm:"primaryKey;column:nconst"`
	PrimaryName       string
	BirthYear         *int
	DeathYear         *int
	PrimaryProfession *string  `gorm:"type:text[]"`
	KnownForTitles    []*Title `gorm:"many2many:filmography;"`
}

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

func main() {
	helpers.InitDB()
	db := helpers.GetDB()
	err := readTitles("title.tsv", db)
	if err != nil {
		log.Fatalf("failed to read titles: %v", err)
	}
	err = readPeople("person.tsv", db)
	if err != nil {
		log.Fatalf("failed to read people: %v", err)
	}
}

func readPeople(filename string, db *gorm.DB) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Print(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	scanner.Scan() // Skip header line

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		birthYear, deathYear := parseYear(fields[2]), parseYear(fields[3])
		var profession *string
		if fields[4] != "\\N" {
			professionList := strings.Split(fields[4], ",")
			professionStr := "{" + strings.Join(professionList, ",") + "}"
			profession = &professionStr
		}

		person := Person{
			ID:                fields[0],
			PrimaryName:       fields[1],
			BirthYear:         birthYear,
			DeathYear:         deathYear,
			PrimaryProfession: profession,
			KnownForTitles:    []*Title{},
		}

		if err := db.Create(&person).Error; err != nil {
			log.Printf("failed to create person: %v", err)
		}
		if fields[4] != "\\N" {
			knownForTitles := strings.Split(fields[5], ",")
			for _, titleID := range knownForTitles {
				var title Title
				if titleID != "\\N" {
					if err := db.First(&title, "tconst = ?", titleID).Error; err == nil {
						err := db.Model(&person).Association("KnownForTitles").Append(&title)
						if err != nil {
							return err
						}
					}
				}

			}
		}

	}

	return scanner.Err()
}

func readTitles(filename string, db *gorm.DB) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	scanner.Scan() // Skip header line

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		startYear, endYear := parseYear(fields[5]), parseYear(fields[6])
		runtimeMinutes := parseInt(fields[7])
		var genres *string
		if fields[8] != "\\N" {
			genreList := strings.Split(fields[8], ",")
			genreStr := "{" + strings.Join(genreList, ",") + "}"
			genres = &genreStr
		}
		isAdult := fields[4] == "1"

		title := Title{
			ID:             fields[0],
			TitleType:      fields[1],
			PrimaryTitle:   fields[2],
			OriginalTitle:  fields[3],
			IsAdult:        isAdult,
			StartYear:      startYear,
			EndYear:        endYear,
			RuntimeMinutes: runtimeMinutes,
			Genres:         genres,
			Actors:         []*Person{},
		}

		if err := db.Create(&title).Error; err != nil {
			log.Printf("failed to create title: %v", err)
		}
	}

	return scanner.Err()
}

func parseYear(yearStr string) *int {
	if yearStr == "\\N" {
		return nil
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Printf("failed to parse year: %v", err)
		return nil
	}
	return &year
}

func parseInt(intStr string) *int {
	if intStr == "\\N" {
		return nil
	}
	value, err := strconv.Atoi(intStr)
	if err != nil {
		log.Printf("failed to parse int: %v", err)
		return nil
	}
	return &value
}

func parseString(str string) *string {
	if str == "\\N" {
		return nil
	}
	return &str
}
