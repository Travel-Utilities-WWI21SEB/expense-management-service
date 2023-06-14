package middlewares

import (
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ValidateUUID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Params
		params := c.Params

		// Loop through params
		for _, param := range params {
			// Check if param is uuid
			if _, err := uuid.Parse(param.Value); err != nil {
				utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
				return
			}
		}

		c.Next()
	}
}
