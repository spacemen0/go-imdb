package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Title struct {
	ID             string  `gorm:"primaryKey;column:tconst"`
	TitleType      string  `gorm:"column:titleType"`
	PrimaryTitle   string  `gorm:"column:primaryTitle"`
	OriginalTitle  string  `gorm:"column:originalTitle"`
	IsAdult        bool    `gorm:"column:isAdult"`
	StartYear      *int    `gorm:"column:startYear"`
	EndYear        *int    `gorm:"column:endYear"`
	RuntimeMinutes *int    `gorm:"column:runTimeMinutes"`
	Genres         *string `gorm:"type:text[];column:genres"`
}

func main() {
	dsn := "host=localhost user=postgres password=Fracture123 dbname=IMDB port=5678 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&Title{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	// Open the TSV file
	file, err := os.Open("data.tsv")
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	// Create a new CSV reader with tab as the delimiter
	reader := csv.NewReader(file)
	reader.Comma = '\t'

	// Read the header
	header, err := reader.Read()
	if err != nil {
		log.Fatal("Error reading header:", err)
	}

	fmt.Println("Header:", header)

	// Read and process the records
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Error reading record:", err)
		}

		var startYear, endYear, runtimeMinutes *int
		var genres *string

		// Convert fields to appropriate types
		if record[5] != "\\N" {
			startYearVal, _ := strconv.Atoi(record[5])
			startYear = &startYearVal
		}
		if record[6] != "\\N" {
			endYearVal, _ := strconv.Atoi(record[6])
			endYear = &endYearVal
		}
		if record[7] != "\\N" {
			runtimeMinutesVal, _ := strconv.Atoi(record[7])
			runtimeMinutes = &runtimeMinutesVal
		}
		if record[8] != "\\N" {
			genreList := strings.Split(record[8], ",")
			genreStr := "{" + strings.Join(genreList, ",") + "}"
			genres = &genreStr
		}

		isAdult, _ := strconv.ParseBool(record[4])

		title := Title{
			ID:             record[0],
			TitleType:      record[1],
			PrimaryTitle:   record[2],
			OriginalTitle:  record[3],
			IsAdult:        isAdult,
			StartYear:      startYear,
			EndYear:        endYear,
			RuntimeMinutes: runtimeMinutes,
			Genres:         genres,
		}

		db.Create(&title)
	}

	fmt.Println("Data insertion completed.")
}
