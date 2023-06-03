package controllers

import (
	"context"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/google/uuid"
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

	// Check if trip exists
	trip, err := ccc.DatabaseMgr.GetTripById(ctx, tripId)

}
