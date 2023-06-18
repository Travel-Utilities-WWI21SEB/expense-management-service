package handlers

import (
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controllers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

		// Distribute cost with all participants
		err := DistributeCosts(&costData)
		if err != nil {
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
		// ctx := c.Request.Context()

		/*response, err := costCtl.GetTripCosts(ctx, nil)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}*/

		c.JSON(http.StatusOK, nil)
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

		err := costCtl.DeleteCostEntry(ctx)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func ValidateAmount(amount string) decimal.Decimal {
	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return decimal.Zero
	}
	if amountDecimal.IsNegative() {
		return decimal.Zero
	}
	return amountDecimal.RoundDown(2)
}

func SumContributions(contributors []*models.Contributor) decimal.Decimal {
	sum := decimal.Zero

	// Sum of debtor contributions
	for _, contributor := range contributors {
		if contributor.Amount != "" {
			amount := ValidateAmount(contributor.Amount)
			sum = sum.Add(amount)
		}
	}
	return sum
}

func DistributeCosts(request *models.CreateCostRequest) *models.ExpenseServiceError {
	request.Amount = ValidateAmount(request.Amount).String()

	// Check if amount is 0
	if request.Amount == "0" {
		return expense_errors.EXPENSE_BAD_REQUEST
	}

	totalCost := ValidateAmount(request.Amount)
	sum := SumContributions(request.Contributors)

	// Check if sum of contributions is greater than total cost
	if sum.GreaterThan(totalCost) {
		return expense_errors.EXPENSE_BAD_REQUEST
	}

	// Get number of contributors with no amount
	numContributorsWithNoAmount := 0
	for _, contributor := range request.Contributors {
		if contributor.Amount == "" {
			numContributorsWithNoAmount++
		}
	}

	// Distribute remaining costs to contributors with no amount
	if numContributorsWithNoAmount > 0 {
		remainingCost := totalCost.Sub(sum)
		// Sum of distributed amounts
		distributedAmount := DistributeRemainingCosts(request.Contributors, remainingCost, numContributorsWithNoAmount, request.Creditor)

		// Check if the sum of contributions equals total cost after distribution
		sum = sum.Add(distributedAmount)
	}

	// Check if the sum of contributions equals total cost
	if !sum.Equal(totalCost) {
		log.Printf("Sum of contributions (%v) does not equal total cost (%v)", sum, totalCost)
		return expense_errors.EXPENSE_BAD_REQUEST
	}

	return nil
}

func DistributeRemainingCosts(contributors []*models.Contributor, remainingCost decimal.Decimal, numContributorsWithNoAmount int, creditor string) decimal.Decimal {
	amountPerContributor := remainingCost.Div(decimal.NewFromInt(int64(numContributorsWithNoAmount)))

	// Round amountPerContributor to 2 decimal places
	amountPerContributor = amountPerContributor.RoundDown(2)

	// Sum of distributed amounts
	distributedAmount := decimal.Zero

	// Check before distributing, if rounding difference is greater than 0.00
	roundingDifference := remainingCost.Sub(amountPerContributor.Mul(decimal.NewFromInt(int64(numContributorsWithNoAmount)))) // roundingDifference = remainingCost - (amountPerContributor * numContributorsWithNoAmount)
	log.Printf("Rounding difference: %v", roundingDifference)
	// Rounding difference can only be positive because of the way it is calculated (see above)

	// Distribute remaining cost to contributors with no amount
	for _, contributor := range contributors {
		if contributor.Amount == "" {
			realAmountPerContributor := amountPerContributor

			// Add rounding difference to creditor
			if contributor.Username == creditor {
				realAmountPerContributor = realAmountPerContributor.Add(roundingDifference)
			}
			log.Printf("Distributing %v to %v", realAmountPerContributor, contributor.Username)

			contributor.Amount = realAmountPerContributor.String()
			distributedAmount = distributedAmount.Add(realAmountPerContributor)
		}
	}

	return distributedAmount
}
