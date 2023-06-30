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
		// Get tripId from path
		tripId := uuid.MustParse(c.Param(models.ExpenseParamKeyTripId))

		// Get cost category entry from request body
		var costCategoryData models.CostCategoryPostRequest
		if err := c.ShouldBindJSON(&costCategoryData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Check if cost category entry already has empty fields
		if utils.ContainsEmptyString(costCategoryData.Name) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Create cost category entry
		ctx := c.Request.Context()
		response, serviceErr := costCategoryCtl.CreateCostCategory(ctx, &tripId, costCategoryData)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusCreated, response)
	}
}

func UpdateCostCategoryEntryHandler(costCategoryCtl controllers.CostCategoryCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get costCategoryId from path
		costCategoryId := uuid.MustParse(c.Param(models.ExpenseParamKeyCostCategoryId))

		// Get cost category entry from request body
		var costCategoryData models.CostCategoryPatchRequest
		if err := c.ShouldBindJSON(&costCategoryData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Update cost category entry
		ctx := c.Request.Context()
		response, serviceErr := costCategoryCtl.PatchCostCategory(ctx, &costCategoryId, costCategoryData)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetCostCategoryDetailsHandler(costCategoryCtl controllers.CostCategoryCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get costCategoryId from path
		costCategoryId := uuid.MustParse(c.Param(models.ExpenseParamKeyCostCategoryId))

		// Get cost category entry
		ctx := c.Request.Context()
		response, serviceErr := costCategoryCtl.GetCostCategoryDetails(ctx, &costCategoryId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func DeleteCostCategoryEntryHandler(costCategoryCtl controllers.CostCategoryCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get costCategoryId from path
		costCategoryId := uuid.MustParse(c.Param(models.ExpenseParamKeyCostCategoryId))

		// Delete cost category entry
		ctx := c.Request.Context()
		serviceErr := costCategoryCtl.DeleteCostCategory(ctx, &costCategoryId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func GetCostCategoryEntriesHandler(costCategoryCtl controllers.CostCategoryCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tripId from path
		tripId := uuid.MustParse(c.Param(models.ExpenseParamKeyTripId))

		// Get cost category entries
		ctx := c.Request.Context()
		response, serviceErr := costCategoryCtl.GetCostCategoryEntries(ctx, &tripId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
