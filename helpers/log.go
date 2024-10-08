package helpers

import (
	"log"
	"os"
)

var Log *log.Logger

func InitLogger(path string) {
	// Open the log file
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	// Create a new logger that writes to the file
	Log = log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)
}
