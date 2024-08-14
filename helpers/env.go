package helpers

import (
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		Log.Fatal("Error loading .env file")
	}
}
