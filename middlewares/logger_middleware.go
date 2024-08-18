package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"spacemen0.github.com/helpers"
	"time"

	"github.com/gin-gonic/gin"
	"spacemen0.github.com/utils" // Update this import path
)

// LoggerMiddleware logs request and response details to a file
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a buffer to capture the response body
		w := utils.NewResponseWriter(c.Writer)
		c.Writer = w

		// Log request details
		startTime := time.Now()

		// Read the request body
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		// Reset the request body so that it can be read again by handlers
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		// Proceed with the request
		c.Next()
		_, _ = w.ResponseWriter.Write(w.Body.Bytes())
		responseHeaders := ""
		for key, values := range w.Header() {
			responseHeaders += fmt.Sprintf("%s: %v\n", key, values)
		}
		// Log response details
		duration := time.Since(startTime)
		logEntry := fmt.Sprintf(
			"\nMethod: %s\nURL: %s\nHeaders: %v\nRequest Body: %s\nResponse Status: %d\nResponse Body: %s\nResponse Header: %sDuration: %v\n\n",
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Header,
			string(requestBody),
			c.Writer.Status(),
			string(w.Body.Bytes()),
			responseHeaders,
			duration,
		)
		helpers.Log.Println(logEntry)
	}
}
