package middlewares

import (
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
)

func TripValidationMiddleware(databaseMgr managers.DatabaseMgr) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("TripValidationMiddleware: %v", c.Request.URL)
		// Get tripId from path
		tripIdParam := c.Param("tripId")
		tripId, err := uuid.Parse(tripIdParam)
		if err != nil {
			log.Printf("Invalid tripId %s", tripIdParam)
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		userId, ok := c.Request.Context().Value(models.ExpenseContextKeyUserID).(*uuid.UUID)
		if !ok {
			log.Printf("Invalid userId %s", userId.String())
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_INTERNAL_ERROR)
			return
		}

		// Check if trip exists
		exists, err := databaseMgr.CheckIfExists("SELECT COUNT(*) FROM trip WHERE id = $1", &tripId)

		if err != nil {
			log.Printf("Error while checking if trip %s exists", tripId.String())
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_INTERNAL_ERROR)
			return
		}

		if !exists {
			log.Printf("Trip %s not found", tripId.String())
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_TRIP_NOT_FOUND)
			return
		}

		// Check if user is part of trip
		exists, err = databaseMgr.CheckIfExists("SELECT COUNT(*) FROM user_trip_association WHERE id_trip = $1 AND id_user = $2", &tripId, &userId)

		if err != nil {
			log.Printf("Error while checking if user %s is part of trip %s", userId.String(), tripId.String())
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_INTERNAL_ERROR)
			return
		}

		if !exists {
			log.Printf("User %s is not part of trip %s", userId.String(), tripId.String())
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_FORBIDDEN)
			return
		}
	}
}
