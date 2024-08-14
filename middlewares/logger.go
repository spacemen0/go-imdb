package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"spacemen0.github.com/utils" // Update this import path
)

// LoggerToFile logs request and response details to a file
func LoggerToFile(filePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a buffer to capture the response body
		buffer := new(bytes.Buffer)

		// Use the custom ResponseWriter to capture response body
		w := utils.NewResponseWriter(c.Writer)
		w.Body = buffer
		c.Writer = w

		// Log request details
		startTime := time.Now()

		// Read the request body
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		// Proceed with the request
		c.Next()

		// Log response details
		duration := time.Since(startTime)
		logEntry := fmt.Sprintf(
			"Time: %s\nMethod: %s\nURL: %s\nHeaders: %v\nRequest Body: %s\nResponse Status: %d\nResponse Body: %s\nDuration: %v\n\n",
			startTime.Format(time.RFC3339),
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Header,
			string(requestBody),
			c.Writer.Status(),
			buffer.String(),
			duration,
		)

		// Write log to file
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening log file:", err)
			return
		}
		defer func(f *os.File) {
			if err := f.Close(); err != nil {
				fmt.Println("Error closing log file:", err)
			}
		}(f)

		if _, err := f.WriteString(logEntry); err != nil {
			fmt.Println("Error writing to log file:", err)
		}
	}
}
