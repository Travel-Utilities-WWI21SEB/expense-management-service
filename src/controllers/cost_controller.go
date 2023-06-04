package controllers

import (
	"context"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
)

// Exposed interface to the handler-package
type CostCtl interface {
	CreateCostEntry(ctx context.Context) (*models.CostResponse, *models.ExpenseServiceError)
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

func (cc *CostController) CreateCostEntry(ctx context.Context) (*models.CostResponse, *models.ExpenseServiceError) {
	// TO-DO
	return nil, nil
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
