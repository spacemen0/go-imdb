package middlewares

import (
	"encoding/json"
	"net/http"
	"spacemen0.github.com/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// DataMiddleware processes JSON responses to omit empty fields in specific fields
func DataMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.DefaultQuery("verbose", "false") == "true" {
			c.Next()
			return
		}
		w := utils.NewResponseWriter(c.Writer)
		c.Writer = w

		c.Next()
		if !strings.Contains(c.Writer.Header().Get("Content-Type"), "application/json") {
			_, _ = w.ResponseWriter.Write(w.Body.Bytes())
			return
		}
		var data map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			_, _ = c.Writer.Write([]byte(`{"error": "Failed to process response"}`))
			return
		}
		if knownForTitles, ok := data["knownForTitles"].([]any); ok {
			processPreloadedData(knownForTitles)
			data["knownForTitles"] = knownForTitles
		}
		if actors, ok := data["actors"].([]any); ok {
			processPreloadedData(actors)
			data["actors"] = actors
		}
		newBody, err := json.Marshal(data)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			_, _ = c.Writer.Write([]byte(`{"error": "Failed to process response"}`))
			return
		}
		_, _ = w.ResponseWriter.WriteString(string(newBody))
		w.Body.Reset()
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
