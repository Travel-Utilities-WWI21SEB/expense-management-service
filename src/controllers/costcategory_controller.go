package controllers

import (
	"context"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Exposed interface to the handler-package
type CostCategoryCtl interface {
	CreateCostCategory(ctx context.Context, tripId *uuid.UUID, costCategoryRequest models.CreateCostCategoryRequest) (*models.CostCategoryResponse, *models.ExpenseServiceError)
}

// Cost Category Controller structure
type CostCategoryController struct {
	DatabaseMgr managers.DatabaseMgr
}

func (ccc *CostCategoryController) CreateCostCategory(ctx context.Context, tripId *uuid.UUID, createCostCategoryRequest models.CreateCostCategoryRequest) (*models.CostCategoryResponse, *models.ExpenseServiceError) {
	if utils.ContainsEmptyString(createCostCategoryRequest.Name) {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	// Create cost category
	costCategoryId := uuid.New()

	// Create cost category object based on CostCategorySchema
	costCategory := &models.CostCategorySchema{
		CostCategoryID: &costCategoryId,
		Name:           createCostCategoryRequest.Name,
		Description:    createCostCategoryRequest.Description,
		Icon:           createCostCategoryRequest.Icon,
		Color:          createCostCategoryRequest.Color,
		TripID:         tripId,
	}

	// Insert cost category into database
	insertString := "INSERT INTO cost_category (id, name, description, icon, color, trip_id) VALUES ($1, $2, $3, $4, $5, $6)"
	if _, err := ccc.DatabaseMgr.ExecuteQuery(insertString, costCategory.CostCategoryID, costCategory.Name, costCategory.Description, costCategory.Icon, costCategory.Color, costCategory.TripID); err != nil {
		// If the constraint was violated, the cost category already exists
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // 23505 = unique_violation
			return nil, expense_errors.EXPENSE_CONFLICT
		}

		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Return created cost category
	response := models.CostCategoryResponse{
		CostCategoryId: &costCategoryId,
		Name:           createCostCategoryRequest.Name,
		Description:    createCostCategoryRequest.Description,
		Icon:           createCostCategoryRequest.Icon,
		Color:          createCostCategoryRequest.Color,
	}

	return &response, nil
}
