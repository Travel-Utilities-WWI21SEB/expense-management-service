package utils

import (
	"log"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
	"github.com/gin-gonic/gin"
)

func HandleErrorAndAbort(c *gin.Context, err model.ExpenseServiceError) {
	log.Printf("Error handling request: %v", err)
	c.AbortWithStatusJSON(err.Status, gin.H{"errorMessage": err.ErrorMessage, "errorCode": err.ErrorCode})
}

func IsDateFormatCorrect(date ...string) bool {
	for _, d := range date {
		// correct format: YYYY-MM-DD
		// check if string has correct length
		if len(d) != 10 {
			return false
		}

		// check if string has correct format
		if d[4] != '-' || d[7] != '-' {
			return false
		}

		// check if string contains only numbers
		if !IsNumeric(d[0:4]) || !IsNumeric(d[5:7]) || !IsNumeric(d[8:10]) {
			return false
		}
	}
	return true
}

func IsNumeric(s string) bool {
	// check if string is numeric
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
