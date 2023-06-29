package handlers

import (
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controllers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func GetDebtsHandler(DebtCtl controllers.DebtCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get tripId from request params
		tripId := uuid.MustParse(c.Param(models.ExpenseParamKeyTripId))

		// Get debts
		response, err := DebtCtl.GetDebtEntries(ctx, &tripId)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetDebtDetailsHandler(DebtCtl controllers.DebtCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
	}
}
