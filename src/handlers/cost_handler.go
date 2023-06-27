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
	"strconv"
	"strings"
	"time"
)

func CreateCostEntryHandler(costCtl controllers.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get tripId from request params
		tripId, err := uuid.Parse(c.Param(models.ExpenseParamKeyTripId))
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
		response, serviceErr := costCtl.CreateCostEntry(ctx, &tripId, costData)
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

		// Get tripId from request params
		tripId := uuid.MustParse(c.Param(models.ExpenseParamKeyTripId))

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

		response, serviceErr := costCtl.PatchCostEntry(ctx, &tripId, &costId, costData)
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

		// Create query params
		queryParams := models.CostQueryParams{
			TripId:           &tripId,
			CostCategoryId:   nil,
			CostCategoryName: nil,
			UserId:           nil,
			Username:         nil,
			MinAmount:        nil,
			MaxAmount:        nil,
			MinDeductionDate: nil,
			MaxDeductionDate: nil,
			MinEndDate:       nil,
			MaxEndDate:       nil,
			MinCreationDate:  nil,
			MaxCreationDate:  nil,
			Page:             0,
			PageSize:         0,
			SortBy:           "created_at",
			SortOrder:        "DESC",
		}

		// Get all query params from request
		CostCategoryIdStr := c.Query("costCategoryId")
		CostCategoryNameStr := c.Query("costCategoryName")
		UserIdStr := c.Query("userId")
		UsernameStr := c.Query("username")
		MinAmountStr := c.Query("minAmount")
		MaxAmountStr := c.Query("maxAmount")
		MinDeductionDateStr := c.Query("minDeductionDate")
		MaxDeductionDateStr := c.Query("maxDeductionDate")
		MinEndDateStr := c.Query("minEndDate")
		MaxEndDateStr := c.Query("maxEndDate")
		MinCreationDateStr := c.Query("minCreationDate")
		MaxCreationDateStr := c.Query("maxCreationDate")
		PageStr := c.Query("page")
		PageSizeStr := c.Query("pageSize")
		SortByStr := c.Query("sortBy")
		SortOrderStr := c.Query("sortOrder")

		if CostCategoryIdStr != "" {
			id, err := uuid.Parse(CostCategoryIdStr)
			if err != nil {
				utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
				return
			}
			queryParams.CostCategoryId = &id
		}

		if CostCategoryNameStr != "" {
			queryParams.CostCategoryName = &CostCategoryNameStr
		}

		if UserIdStr != "" {
			userId, err := uuid.Parse(UserIdStr)
			if err != nil {
				utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
				return
			}
			queryParams.UserId = &userId
		}

		if UsernameStr != "" {
			queryParams.Username = &UsernameStr
		}

		if MinAmountStr != "" {
			_, err := strconv.ParseFloat(MinAmountStr, 64)
			if err != nil {
				utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
				return
			}
			queryParams.MinAmount = &MinAmountStr
		}

		if MaxAmountStr != "" {
			_, err := strconv.ParseFloat(MaxAmountStr, 64)
			if err != nil {
				utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
				return
			}
			queryParams.MaxAmount = &MaxAmountStr
		}

		if MinDeductionDateStr != "" {
			_, err := time.Parse(time.RFC3339, MinDeductionDateStr)
			if err != nil {
				utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
				return
			}
			queryParams.MinDeductionDate = &MinDeductionDateStr
		}

		if MaxDeductionDateStr != "" {
			_, err := time.Parse(time.RFC3339, MaxDeductionDateStr)
			if err != nil {
				log.Printf("Error while parsing date: %v", err)
				utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
				return
			}
			queryParams.MaxDeductionDate = &MaxDeductionDateStr
		}

		if MinEndDateStr != "" {
			_, err := time.Parse(time.RFC3339, MinEndDateStr)
			if err != nil {
				utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
				return
			}
			queryParams.MinEndDate = &MinEndDateStr
		}

		if MaxEndDateStr != "" {
			_, err := time.Parse(time.RFC3339, MaxEndDateStr)
			if err != nil {
				utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
				return
			}
			queryParams.MaxEndDate = &MaxEndDateStr
		}

		if MinCreationDateStr != "" {
			_, err := time.Parse(time.RFC3339, MinCreationDateStr)
			if err != nil {
				utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
				return
			}
			queryParams.MinCreationDate = &MinCreationDateStr
		}

		if MaxCreationDateStr != "" {
			_, err := time.Parse(time.RFC3339, MaxCreationDateStr)
			if err != nil {
				utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
				return
			}
			queryParams.MaxCreationDate = &MaxCreationDateStr
		}

		if PageStr != "" {
			page, err := strconv.Atoi(PageStr)
			if err != nil {
				utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
				return
			}
			queryParams.Page = page
		}

		if PageSizeStr != "" {
			pageSize, err := strconv.Atoi(PageSizeStr)
			if err != nil {
				utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
				return
			}
			queryParams.PageSize = pageSize
		}

		if SortByStr != "" {
			queryParams.SortBy = SortByStr
		}

		if strings.ToUpper(SortOrderStr) == "ASC" {
			queryParams.SortOrder = "ASC"
		}

		// Pass query params to the controller
		response, err := costCtl.GetCostEntries(ctx, &queryParams)
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
