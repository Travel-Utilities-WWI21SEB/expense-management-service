package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func HandleErrorAndAbort(c *gin.Context, message string, status int, err error) {
	log.Printf("Error handling request: %v", err)
	c.AbortWithStatusJSON(status, gin.H{"error": message})
}
