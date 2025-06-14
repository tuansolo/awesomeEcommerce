package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs the incoming HTTP request and response
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Read the request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Restore the request body for further processing
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create a custom response writer to capture the response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request details
		log.Printf("[%s] %s %s %d %s %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			latency,
			c.Errors.String(),
		)

		// In development, we might want to log request and response bodies
		// In production, this should be conditional or disabled for privacy and performance
		if gin.Mode() == gin.DebugMode {
			// Log request body if it's not too large
			if len(requestBody) > 0 && len(requestBody) < 10000 {
				log.Printf("Request Body: %s", string(requestBody))
			}

			// Log response body if it's not too large
			if blw.body.Len() > 0 && blw.body.Len() < 10000 {
				log.Printf("Response Body: %s", blw.body.String())
			}
		}
	}
}

// bodyLogWriter is a custom response writer that captures the response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body
func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
