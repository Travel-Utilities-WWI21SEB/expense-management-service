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
		var costCategoryData models.CostCategoryPostRequest
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

		c.JSON(http.StatusCreated, response)
	}
}

func UpdateCostCategoryEntryHandler(costCategoryCtl controllers.CostCategoryCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get costCategoryId from path
		costCategoryIdParam := c.Param("costCategoryId")
		costCategoryId, err := uuid.Parse(costCategoryIdParam)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Get cost category entry from request body
		var costCategoryData models.CostCategoryPatchRequest
		if err := c.ShouldBindJSON(&costCategoryData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Update cost category entry
		response, serviceErr := costCategoryCtl.PatchCostCategory(ctx, &costCategoryId, costCategoryData)
		if err != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetCostCategoryDetailsHandler(costCategoryCtl controllers.CostCategoryCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get costCategoryId from path
		costCategoryIdParam := c.Param("costCategoryId")
		costCategoryId, err := uuid.Parse(costCategoryIdParam)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Get cost category entry
		response, serviceErr := costCategoryCtl.GetCostCategoryDetails(ctx, &costCategoryId)
		if err != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func DeleteCostCategoryEntryHandler(costCategoryCtl controllers.CostCategoryCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get costCategoryId from path
		costCategoryIdParam := c.Param("costCategoryId")
		costCategoryId, err := uuid.Parse(costCategoryIdParam)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Delete cost category entry
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
		ctx := c.Request.Context()

		// Get tripId from path
		tripIdParam := c.Param("tripId")
		tripId, err := uuid.Parse(tripIdParam)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Get cost category entries
		response, serviceErr := costCategoryCtl.GetCostCategoryEntries(ctx, &tripId)
		if err != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
