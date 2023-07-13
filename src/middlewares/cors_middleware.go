package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("CorsMiddleware: Setting CORS headers")

		allowedOrigins := []string{
			"http://localhost:4173",          // SvelteKit Vite Prod Preview
			"http://localhost:4174",          // SvelteKit Docker Compose
			"http://localhost:5173",          // SvelteKit Vite Dev Preview
			"https://costventures.works.net", // SvelteKit Vite Prod Build
		}
		origin := c.GetHeader("Origin")

		// Check if the origin is allowed
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, PATCH, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
