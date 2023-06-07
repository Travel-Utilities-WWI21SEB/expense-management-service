package controllers

import (
	"context"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

// Exposed interface to the handler-package
type CostCtl interface {
	CreateCostEntry(ctx context.Context, createCostRequest models.CreateCostRequest) (*models.CostDetailsResponse, *models.ExpenseServiceError)
	PatchCostEntry(ctx context.Context) (*models.CostResponse, *models.ExpenseServiceError)
	PutCostEntry(ctx context.Context) (*models.CostResponse, *models.ExpenseServiceError)
	GetCostDetails(ctx context.Context) (*models.CostResponse, *models.ExpenseServiceError)
	GetTripCosts(ctx context.Context) (*models.CostResponse, *models.ExpenseServiceError)
	DeleteCostEntry(ctx context.Context) *models.ExpenseServiceError
}

// Cost Controller structure
type CostController struct {
	DatabaseMgr managers.DatabaseMgr
}

func (cc *CostController) CreateCostEntry(ctx context.Context, createCostRequest models.CreateCostRequest) (*models.CostDetailsResponse, *models.ExpenseServiceError) {
	// Workflow:
	// - Check if cost entry already has empty fields
	// - Generate costId
	// - Convert amount to decimal
	// - Check if amount is negative
	// - Check if currency code is valid
	// - Create cost entry
	// - Insert cost entry into database
	// - Insert cost entry and creator into cost_category_cost table
	// - Insert cost entry and debtor into user_cost_association table

	// Check if cost entry already has empty fields
	if utils.ContainsEmptyString(createCostRequest.Amount, createCostRequest.CurrencyCode) {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	// Generate costId
	costId, err := uuid.NewUUID()
	if err != nil {
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Convert amount to decimal
	amount, err := decimal.NewFromString(createCostRequest.Amount)
	if err != nil {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	// Check if amount is negative
	if amount.IsNegative() {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	// Check if currency code is valid
	if !utils.IsValidCurrencyCode(createCostRequest.CurrencyCode) {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	now := time.Now()

	if createCostRequest.DeductedAt == nil {
		createCostRequest.DeductedAt = &now
	}

	// Create cost entry
	costEntry := models.CostSchema{
		CostID:         &costId,
		Amount:         amount,
		Description:    createCostRequest.Description,
		CreationDate:   &now,
		DeductionDate:  createCostRequest.DeductedAt,
		EndDate:        createCostRequest.EndDate,
		CostCategoryID: createCostRequest.CostCategoryID,
	}

	// Insert cost entry into database
	insertString := "INSERT INTO cost (id, amount, description, created_at, deducted_at, end_date, id_cost_category) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	if _, err = cc.DatabaseMgr.ExecuteStatement(insertString, costEntry.CostID, costEntry.Amount, costEntry.Description, costEntry.CreationDate, costEntry.DeductionDate, costEntry.EndDate, costEntry.CostCategoryID); err != nil {
		// Check if cost category exists
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "foreign_key_violation" {
			return nil, expense_errors.EXPENSE_NOT_FOUND // Cost Category Not found
		}

		log.Printf("Error in controller.cost_controller.CreateCostEntry: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Get creator id
	creatorId := ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)
	// Check if creator is creditor
	isCreatorCreditor := createCostRequest.PaidBy == nil || *createCostRequest.PaidBy == *creatorId

	// Insert cost entry and creator into cost_category_cost table
	insertCreatorString := "INSERT INTO user_cost_association (id_user, id_cost, is_creditor) VALUES ($1, $2, $3)"
	if _, err = cc.DatabaseMgr.ExecuteStatement(insertCreatorString, creatorId, costEntry.CostID, isCreatorCreditor); err != nil {
		log.Printf("Error in controller.cost_controller.CreateCostEntry: %v", err)
		// Delete cost entry
		if _, err = cc.DatabaseMgr.ExecuteStatement("DELETE FROM cost WHERE id = $1", costEntry.CostID); err != nil {
			log.Printf("Error in controller.cost_controller.CreateCostEntry: %v", err)
		}
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Insert cost entry and *creditor* into cost_category_cost table, if creator is not creditor
	if !isCreatorCreditor {
		insertString = "INSERT INTO user_cost_association (id_user, id_cost, is_creditor) VALUES ($1, $2, $3)"
		if _, err = cc.DatabaseMgr.ExecuteStatement(insertString, createCostRequest.PaidBy, costEntry.CostID, true); err != nil {
			expenseErr := expense_errors.EXPENSE_INTERNAL_ERROR

			// If foreign key violation, user in paidBy field does not exist
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "foreign_key_violation" {
				expenseErr = expense_errors.EXPENSE_USER_NOT_FOUND // Paid by user not found
			} else {
				log.Printf("Error in controller.cost_controller.CreateCostEntry: %v", err)
			}

			// Delete cost entry
			deleteString := "DELETE FROM cost WHERE id = $1"
			if _, err = cc.DatabaseMgr.ExecuteStatement(deleteString, costEntry.CostID); err != nil {
				log.Printf("Error in controller.cost_controller.CreateCostEntry: %v", err)
			}
			return nil, expenseErr
		}
	}

	// Insert cost entry and *debtors* into cost_category_cost table
	for _, debtorId := range createCostRequest.PaidFor {
		insertString = "INSERT INTO user_cost_association (id_user, id_cost, is_creditor) VALUES ($1, $2, $3)"
		if _, err = cc.DatabaseMgr.ExecuteStatement(insertString, debtorId, costEntry.CostID, false); err != nil {
			expenseErr := expense_errors.EXPENSE_INTERNAL_ERROR

			// If foreign key violation, user in paidFor field does not exist
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "foreign_key_violation" {
				expenseErr = expense_errors.EXPENSE_USER_NOT_FOUND // Paid for user not found
			} else {
				log.Printf("Error in controller.cost_controller.CreateCostEntry: %v", err)
			}

			// Delete cost entry
			deleteString := "DELETE FROM cost WHERE id = $1"
			if _, err = cc.DatabaseMgr.ExecuteStatement(deleteString, costEntry.CostID); err != nil {
				log.Printf("Error in controller.cost_controller.CreateCostEntry: %v", err)
			}
			return nil, expenseErr
		}
	}

	// Return cost entry
	costResponse := &models.CostDetailsResponse{
		CostID:         costEntry.CostID,
		Amount:         costEntry.Amount.String(),
		Description:    costEntry.Description,
		CreationDate:   costEntry.CreationDate,
		DeductionDate:  costEntry.DeductionDate,
		EndDate:        costEntry.EndDate,
		CostCategoryID: costEntry.CostCategoryID,
		PaidBy:         creatorId,
		PaidFor:        createCostRequest.PaidFor,
	}

	return costResponse, nil
}

func (cc *CostController) PatchCostEntry(ctx context.Context) (*models.CostResponse, *models.ExpenseServiceError) {
	// TO-DO
	return nil, nil
}

func (cc *CostController) PutCostEntry(ctx context.Context) (*models.CostResponse, *models.ExpenseServiceError) {
	// TO-DO
	return nil, nil
}

func (cc *CostController) GetCostDetails(ctx context.Context) (*models.CostResponse, *models.ExpenseServiceError) {
	// TO-DO
	return nil, nil
}

func (cc *CostController) GetTripCosts(ctx context.Context) (*models.CostResponse, *models.ExpenseServiceError) {
	// TO-DO
	return nil, nil
}

func (cc *CostController) DeleteCostEntry(ctx context.Context) *models.ExpenseServiceError {
	// TO-DO
	return nil
}
