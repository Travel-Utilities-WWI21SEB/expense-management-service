package middlewares

import (
	"github.com/google/uuid"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
)

func TripValidationMiddleware(databaseMgr managers.DatabaseMgr) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tripId from path
		tripIdParam := c.Param("tripId")
		tripId, err := uuid.Parse(tripIdParam)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Get userId from context
		userId, ok := c.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)
		if !ok {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_INTERNAL_ERROR)
			return
		}

		// Check if trip exists
		exists, err := databaseMgr.CheckIfExists("SELECT COUNT(*) FROM trip WHERE id = $1", &tripId)

		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_INTERNAL_ERROR)
			return
		}

		if !exists {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_TRIP_NOT_FOUND)
			return
		}

		// Check if user is part of trip
		exists, err = databaseMgr.CheckIfExists("SELECT COUNT(*) FROM user_trip_association WHERE id_trip = $1 AND id_user = $2", &tripId, &userId)

		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_INTERNAL_ERROR)
			return
		}

		if !exists {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_UNAUTHORIZED)
			return
		}
	}
}
