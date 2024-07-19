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

// Enum for PrimaryProfession
type PrimaryProfession string

const (
	Accountant               PrimaryProfession = "accountant"
	Actor                    PrimaryProfession = "actor"
	Actress                  PrimaryProfession = "actress"
	AnimationDepartment      PrimaryProfession = "animation_department"
	ArchiveFootage           PrimaryProfession = "archive_footage"
	ArchiveSound             PrimaryProfession = "archive_sound"
	ArtDepartment            PrimaryProfession = "art_department"
	ArtDirector              PrimaryProfession = "art_director"
	Assistant                PrimaryProfession = "assistant"
	AssistantDirector        PrimaryProfession = "assistant_director"
	CameraDepartment         PrimaryProfession = "camera_department"
	CastingDepartment        PrimaryProfession = "casting_department"
	CastingDirector          PrimaryProfession = "casting_director"
	Choreographer            PrimaryProfession = "choreographer"
	Cinematographer          PrimaryProfession = "cinematographer"
	Composer                 PrimaryProfession = "composer"
	CostumeDepartment        PrimaryProfession = "costume_department"
	CostumeDesigner          PrimaryProfession = "costume_designer"
	Director                 PrimaryProfession = "director"
	Editor                   PrimaryProfession = "editor"
	EditorialDepartment      PrimaryProfession = "editorial_department"
	ElectricalDepartment     PrimaryProfession = "electrical_department"
	Executive                PrimaryProfession = "executive"
	Legal                    PrimaryProfession = "legal"
	LocationManagement       PrimaryProfession = "location_management"
	MakeUpDepartment         PrimaryProfession = "make_up_department"
	Manager                  PrimaryProfession = "manager"
	Miscellaneous            PrimaryProfession = "miscellaneous"
	MusicArtist              PrimaryProfession = "music_artist"
	MusicDepartment          PrimaryProfession = "music_department"
	Podcaster                PrimaryProfession = "podcaster"
	Producer                 PrimaryProfession = "producer"
	ProductionDepartment     PrimaryProfession = "production_department"
	ProductionDesigner       PrimaryProfession = "production_designer"
	ProductionManager        PrimaryProfession = "production_manager"
	Publicist                PrimaryProfession = "publicist"
	ScriptDepartment         PrimaryProfession = "script_department"
	SetDecorator             PrimaryProfession = "set_decorator"
	SoundDepartment          PrimaryProfession = "sound_department"
	Soundtrack               PrimaryProfession = "soundtrack"
	SpecialEffects           PrimaryProfession = "special_effects"
	Stunts                   PrimaryProfession = "stunts"
	TalentAgent              PrimaryProfession = "talent_agent"
	TransportationDepartment PrimaryProfession = "transportation_department"
	VisualEffects            PrimaryProfession = "visual_effects"
	Writer                   PrimaryProfession = "writer"
)

// Person represents a record in the IMDB dataset
type Person struct {
	ID                string  `gorm:"primaryKey;column:nconst"`
	PrimaryName       string  `gorm:"column:primaryName"`
	BirthYear         *int    `gorm:"column:birthYear"`
	DeathYear         *int    `gorm:"column:deathYear"`
	PrimaryProfession *string `gorm:"type:text[];column:primaryProfession"`
	KnownForTitles    *string `gorm:"type:text[];column:knownForTitles"`
}

func main() {
	// Database connection
	dsn := "host=localhost user=postgres password=Fracture123 dbname=IMDB port=5678 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&Person{})
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

		var birthYear, deathYear *int
		var professionArray, knownForTitlesArray *string

		// Convert the record to a Person struct
		if record[2] != "\\N" {
			birthYearVal, _ := strconv.Atoi(record[2])
			birthYear = &birthYearVal
		}
		if record[3] != "\\N" {
			deathYearVal, _ := strconv.Atoi(record[3])
			deathYear = &deathYearVal
		}
		if record[4] != "\\N" {
			professions := strings.Split(record[4], ",")
			professionStr := "{" + strings.Join(professions, ",") + "}"
			professionArray = &professionStr
		}
		if record[5] != "\\N" {
			knownForTitles := strings.Split(record[5], ",")
			knownForTitlesStr := "{" + strings.Join(knownForTitles, ",") + "}"
			knownForTitlesArray = &knownForTitlesStr
		}

		person := Person{
			ID:                record[0],
			PrimaryName:       record[1],
			BirthYear:         birthYear,
			DeathYear:         deathYear,
			PrimaryProfession: professionArray,
			KnownForTitles:    knownForTitlesArray,
		}

		db.Create(&person)
	}

	fmt.Println("Data insertion completed.")
}
