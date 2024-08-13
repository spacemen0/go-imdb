package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// OmitEmptyFieldsInPreloadedDataMiddleware processes JSON responses to omit empty fields in specific fields
func OmitEmptyFieldsInPreloadedDataMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture the response
		w := &responseWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		c.Next() // Process the request

		// Only process JSON responses
		if strings.Contains(c.Writer.Header().Get("Content-Type"), "application/json") {
			// Decode the JSON
			var data map[string]any
			if err := json.Unmarshal(w.body.Bytes(), &data); err != nil {
				c.Writer.WriteHeader(http.StatusInternalServerError)
				_, _ = c.Writer.Write([]byte(`{"error": "Failed to process response"}`))
				return
			}

			// Process specific fields
			if knownForTitles, ok := data["KnownForTitles"].([]any); ok {
				processPreloadedData(knownForTitles)
				data["KnownForTitles"] = knownForTitles
			}
			if actors, ok := data["Actors"].([]any); ok {
				processPreloadedData(actors)
				data["Actors"] = actors
			}
			// Encode the JSON back
			newBody, err := json.Marshal(data)
			if err != nil {
				c.Writer.WriteHeader(http.StatusInternalServerError)
				_, _ = c.Writer.Write([]byte(`{"error": "Failed to process response"}`))
				return
			}

			// Write the final response to the original ResponseWriter
			//c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", w.body.Len()))
			//c.Writer.Header().Set("Content-Type", "application/json")
			_, _ = w.ResponseWriter.WriteString(string(newBody))
			w.body.Reset()
		}
	}
}

// processPreloadedData processes a slice of preloaded data and removes empty fields
func processPreloadedData(items []interface{}) {
	for i, item := range items {
		if itemMap, ok := item.(map[string]interface{}); ok {
			removeEmptyFields(itemMap)
			if len(itemMap) == 0 {
				items[i] = nil
			} else {
				items[i] = itemMap
			}
		}
	}
}

// removeEmptyFields recursively removes empty fields from a map
func removeEmptyFields(data map[string]interface{}) {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			if v == "" {
				delete(data, key)
			}
		case []interface{}:
			if len(v) == 0 {
				delete(data, key)
			} else {
				processPreloadedData(v)
			}
		case map[string]interface{}:
			removeEmptyFields(v)
			if len(v) == 0 {
				delete(data, key)
			}
		case nil, bool:
			delete(data, key)
		}
	}
}

// responseWriter is a custom ResponseWriter that captures the response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}
