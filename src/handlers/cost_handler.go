package handlers

import (
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controllers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CreateCostEntryHandler(costCtl controllers.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get cost entry from request body
		var costData models.CreateCostRequest
		if err := c.ShouldBindJSON(&costData); err != nil {
			log.Printf("Error while binding JSON: %v", err)
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Check if cost entry already has empty fields
		if utils.ContainsEmptyString(costData.Amount, costData.CurrencyCode) {
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
		// TO-DO
		ctx := c.Request.Context()

		response, err := costCtl.PatchCostEntry(ctx)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetCostEntriesHandler(costCtl controllers.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		response, err := costCtl.GetTripCosts(ctx)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetCostDetailsHandler(costCtl controllers.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		response, err := costCtl.GetCostDetails(ctx)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func DeleteCostEntryHandler(costCtl controllers.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		err := costCtl.DeleteCostEntry(ctx)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}
