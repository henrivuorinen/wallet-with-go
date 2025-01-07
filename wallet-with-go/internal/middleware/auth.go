package middleware

import (
	"crypto/subtle"
	_ "crypto/subtle"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "net/http"
	"os"
	_ "os"
)

// authmiddleware provides basic auth for engine api
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//read credentials
		expectedApiKey := os.Getenv("GAME_ENGINE_API_KEY")
		if expectedApiKey == "" {
			//fallback
			expectedApiKey = "default-secure-api-key"
		}

		//get api key from header
		apiKey := c.GetHeader("X-Game-Engine-Api-Key")
		log.Println("Received API Key:", apiKey)

		//use constant-time comparison
		if apiKey == "" || subtle.ConstantTimeCompare(
			[]byte(apiKey),
			[]byte(expectedApiKey),
		) != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// CORSMiddleware adds CORS headers to the response
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080") // Your frontend URL
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Game-Engine-Api-Key")

		// Handle preflight (OPTIONS) requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
