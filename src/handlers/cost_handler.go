package handlers

import (
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controllers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

func CreateCostEntryHandler(costCtl controllers.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get cost entry from request body
		var costData models.CostDTO
		if err := c.ShouldBindJSON(&costData); err != nil {
			log.Printf("Error while binding JSON: %v", err)
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Check if cost entry already has empty fields
		if utils.ContainsEmptyString(costData.Amount, costData.CurrencyCode, costData.Creditor, costData.CostCategoryID.String()) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Check if dates are valid
		if !utils.IsValidDate(time.RFC3339, costData.DeductionDate, costData.EndDate) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Check if currency code is valid
		if !utils.IsValidCurrencyCode(costData.CurrencyCode) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Create cost entry
		response, serviceErr := costCtl.CreateCostEntry(ctx, costData)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusCreated, response)
	}
}

func UpdateCostEntryHandler(costCtl controllers.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get costId from request params
		costId, err := uuid.Parse(c.Param(models.ExpenseParamKeyCostId))
		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Get cost entry from request body
		var costData models.CostDTO
		if err := c.ShouldBindJSON(&costData); err != nil {
			log.Printf("Error while binding JSON: %v", err)
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Check if dates are valid
		if !utils.IsValidDate(time.RFC3339, costData.DeductionDate, costData.EndDate) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Check if currency code is valid
		if !utils.ContainsEmptyString(costData.CurrencyCode) && !utils.IsValidCurrencyCode(costData.CurrencyCode) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		response, serviceErr := costCtl.PatchCostEntry(ctx, &costId, costData)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetCostEntriesHandler(costCtl controllers.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get tripId from request params
		tripId := uuid.MustParse(c.Param(models.ExpenseParamKeyTripId))

		response, err := costCtl.GetCostEntriesByTrip(ctx, &tripId)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetCostDetailsHandler(costCtl controllers.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get costId from request params
		costId, err := uuid.Parse(c.Param(models.ExpenseParamKeyCostId))
		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		response, serviceErr := costCtl.GetCostDetails(ctx, &costId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func DeleteCostEntryHandler(costCtl controllers.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		// Get costId from request params
		costId, err := uuid.Parse(c.Param(models.ExpenseParamKeyCostId))
		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
		}

		serviceErr := costCtl.DeleteCostEntry(ctx, &costId)
		if err != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.AbortWithStatus(http.StatusNoContent)
	}
}
