package handlers

import (
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controllers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"net/http"
)

func GetTransactionsHandler(transactionCtl controllers.TransactionCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get tripId from path
		tripId := uuid.MustParse(c.Param(models.ExpenseParamKeyTripId))

		response, serviceErr := transactionCtl.GetTransactionEntries(ctx, &tripId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func CreateTransactionHandler(transactionCtl controllers.TransactionCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var transactionRequest *models.TransactionDTO

		if err := c.ShouldBindJSON(&transactionRequest); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Check if amount is decimal
		if _, err := decimal.NewFromString(transactionRequest.Amount); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Get tripId from path
		tripId := uuid.MustParse(c.Param(models.ExpenseParamKeyTripId))

		// Get user id from context
		userId := ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)

		// Check if creditor and debtor are the same
		if userId.String() == transactionRequest.DebtorId.String() {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		response, serviceErr := transactionCtl.CreateTransactionEntry(ctx, &tripId, transactionRequest)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func DeleteTransactionHandler(transactionCtl controllers.TransactionCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get transactionId from path
		transactionId, err := uuid.Parse(c.Param(models.ExpenseParamKeyTransactionId))
		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		serviceErr := transactionCtl.DeleteTransactionEntry(ctx, &transactionId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.Status(http.StatusOK)
	}
}

func GetTransactionDetailsHandler(transactionCtl controllers.TransactionCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get transactionId from path
		transactionId, err := uuid.Parse(c.Param(models.ExpenseParamKeyTransactionId))
		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		response, serviceErr := transactionCtl.GetTransactionDetails(ctx, &transactionId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func AcceptTransaction(transactionCtl controllers.TransactionCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get transactionId from path
		transactionId, err := uuid.Parse(c.Param(models.ExpenseParamKeyTransactionId))
		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		response, serviceErr := transactionCtl.AcceptTransaction(ctx, &transactionId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func DeclineTransaction(transactionCtl controllers.TransactionCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get transactionId from path
		transactionId, err := uuid.Parse(c.Param(models.ExpenseParamKeyTransactionId))
		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		serviceErr := transactionCtl.DeleteTransactionEntry(ctx, &transactionId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
