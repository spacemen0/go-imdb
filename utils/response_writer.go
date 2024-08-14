package utils

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

// ResponseWriter is a custom response writer to capture the response Body
type ResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

// NewResponseWriter creates a new ResponseWriter instance
func NewResponseWriter(w gin.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		Body:           &bytes.Buffer{},
	}
}

// Write captures the response Body
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	return rw.Body.Write(b)
}
