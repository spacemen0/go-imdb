package helpers

import (
	"github.com/joho/godotenv"
	"log"
)

func ReadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
