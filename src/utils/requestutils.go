package utils

import (
	"log"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/gin-gonic/gin"
)

func HandleErrorAndAbort(c *gin.Context, err models.ExpenseServiceError) {
	log.Printf("Error handling request: %v", err)
	c.AbortWithStatusJSON(err.Status, gin.H{"errorMessage": err.ErrorMessage, "errorCode": err.ErrorCode})
}
