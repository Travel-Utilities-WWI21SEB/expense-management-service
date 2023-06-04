package middlewares

import (
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
)

func TripValidationMiddleware(databaseMgr managers.DatabaseMgr) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if trip exists
		exists, err := databaseMgr.CheckIfExists("SELECT COUNT(*) FROM trip WHERE id_trip = $1", c.Param("tripId"))

		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_INTERNAL_ERROR)
			return
		}

		if !exists {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_TRIP_NOT_FOUND)
			return
		}

		// Check if user is part of trip
		exists, err = databaseMgr.CheckIfExists("SELECT COUNT(*) FROM user_trip_association WHERE id_trip = $1 AND id_user = $2", c.Param("tripId"), c.Request.Context().Value(models.ExpenseContextKeyUserID))

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
