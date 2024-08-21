package helpers

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		Port     int    `yaml:"port"`
		SSLMode  string `yaml:"ssl_mode"`
	} `yaml:"database"`
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
}

var AppConfig Config

func LoadConfig() {
	// Open the YAML configuration file
	file, err := os.Open("config.yaml")
	if err != nil {
		Log.Fatalf("Error opening config.yaml file: %v", err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			Log.Printf("Error closing file: %v", err)
		}
	}(file)

	// Decode the YAML file into the AppConfig struct
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		Log.Fatalf("Error decoding YAML file: %v", err)
	}

	// (Optional) Print to verify the configuration
	Log.Printf("Loaded Config: %+v\n", AppConfig)
}
