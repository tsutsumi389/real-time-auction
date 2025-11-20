package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get allowed origins from environment variable
		allowedOrigins := os.Getenv("CORS_ORIGINS")
		if allowedOrigins == "" {
			// Default for development
			allowedOrigins = "*"
		}

		// Get request origin
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		if allowedOrigins == "*" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			origins := strings.Split(allowedOrigins, ",")
			for _, o := range origins {
				if strings.TrimSpace(o) == origin {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		// Set CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
