package middleware

import (
	"context"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expenseerror"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if Authorization header is set
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_UNAUTHORIZED)
			return
		}

		// Check if Authorization header is valid
		tokenString, err := utils.ExtractToken(authHeader)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_UNAUTHORIZED)
			return
		}

		id, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_UNAUTHORIZED)
			return
		}

		// Add userId to context
		ctx := context.WithValue(c.Request.Context(), utils.ContextKeyUserID, id)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
