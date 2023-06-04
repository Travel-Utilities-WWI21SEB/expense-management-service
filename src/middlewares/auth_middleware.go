package middlewares

import (
	"context"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
	"log"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("JwtAuthMiddleware: %v", c.Request.URL)
		// Check if Authorization header is set
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("JwtAuthMiddleware: Authorization header not set")
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_UNAUTHORIZED)
			return
		}

		// Check if Authorization header is valid
		tokenString, err := utils.ExtractToken(authHeader)
		if err != nil {
			log.Printf("JwtAuthMiddleware: Authorization header not valid")
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_UNAUTHORIZED)
			return
		}

		id, err := utils.ValidateToken(tokenString)
		if err != nil {
			log.Printf("JwtAuthMiddleware: Token not valid")
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_UNAUTHORIZED)
			return
		}

		// Add userId to context
		ctx := context.WithValue(c.Request.Context(), models.ExpenseContextKeyUserID, id)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
