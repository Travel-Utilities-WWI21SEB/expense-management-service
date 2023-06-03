package handlers

import (
	"net/http"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controllers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCostCategoryEntryHandler(costCategoryCtl controllers.CostCategoryCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get tripId from path
		tripIdParam := c.Param("tripId")
		tripId, err := uuid.Parse(tripIdParam)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Get cost category entry from request body
		var costCategoryData models.CreateCostCategoryRequest
		if err := c.ShouldBindJSON(&costCategoryData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Create cost category entry
		response, serviceErr := costCategoryCtl.CreateCostCategory(ctx, &tripId, costCategoryData)
		if err != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
