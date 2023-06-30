package controllers

import (
	"context"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/repositories"
	"github.com/google/uuid"
)

// CostCategoryCtl Exposed interface to the handler-package
type CostCategoryCtl interface {
	CreateCostCategory(ctx context.Context, tripId *uuid.UUID, costCategoryRequest models.CostCategoryPostRequest) (*models.CostCategoryResponse, *models.ExpenseServiceError)
	PatchCostCategory(ctx context.Context, costCategoryId *uuid.UUID, costCategoryRequest models.CostCategoryPatchRequest) (*models.CostCategoryResponse, *models.ExpenseServiceError)
	GetCostCategoryDetails(ctx context.Context, costCategoryId *uuid.UUID) (*models.CostCategoryResponse, *models.ExpenseServiceError)
	DeleteCostCategory(ctx context.Context, costCategoryId *uuid.UUID) *models.ExpenseServiceError
	GetCostCategoryEntries(ctx context.Context, tripId *uuid.UUID) ([]*models.CostCategoryResponse, *models.ExpenseServiceError)
}

// CostCategoryController Cost Category Controller structure
type CostCategoryController struct {
	DatabaseMgr      managers.DatabaseMgr
	CostCategoryRepo repositories.CostCategoryRepo
	CostRepo         repositories.CostRepo
}

func (ccc *CostCategoryController) CreateCostCategory(ctx context.Context, tripId *uuid.UUID, createCostCategoryRequest models.CostCategoryPostRequest) (*models.CostCategoryResponse, *models.ExpenseServiceError) {
	// Generate cost category id
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

	if err := ccc.CostCategoryRepo.CreateCostCategory(ctx, costCategory); err != nil {
		return nil, err
	}

	return ccc.responseBuilder(ctx, costCategory), nil
}

func (ccc *CostCategoryController) PatchCostCategory(ctx context.Context, costCategoryId *uuid.UUID, costCategoryPatchRequest models.CostCategoryPatchRequest) (*models.CostCategoryResponse, *models.ExpenseServiceError) {
	// Get cost category from database
	costCategory, err := ccc.CostCategoryRepo.GetCostCategoryByID(ctx, costCategoryId)
	if err != nil {
		return nil, err
	}

	// Update cost category object
	if costCategoryPatchRequest.Name != "" {
		costCategory.Name = costCategoryPatchRequest.Name
	}

	if costCategoryPatchRequest.Description != "" {
		costCategory.Description = costCategoryPatchRequest.Description
	}

	if costCategoryPatchRequest.Icon != "" {
		costCategory.Icon = costCategoryPatchRequest.Icon
	}

	if costCategoryPatchRequest.Color != "" {
		costCategory.Color = costCategoryPatchRequest.Color
	}

	// Update cost category in database
	if err := ccc.CostCategoryRepo.UpdateCostCategory(ctx, costCategory); err != nil {
		return nil, err
	}

	// Return updated cost category
	return ccc.responseBuilder(ctx, costCategory), nil
}

func (ccc *CostCategoryController) GetCostCategoryDetails(ctx context.Context, costCategoryId *uuid.UUID) (*models.CostCategoryResponse, *models.ExpenseServiceError) {
	costCategory, repoErr := ccc.CostCategoryRepo.GetCostCategoryByID(ctx, costCategoryId)
	if repoErr != nil {
		return nil, repoErr
	}
	return ccc.responseBuilder(ctx, costCategory), nil
}

func (ccc *CostCategoryController) DeleteCostCategory(ctx context.Context, costCategoryId *uuid.UUID) *models.ExpenseServiceError {
	return ccc.CostCategoryRepo.DeleteCostCategory(ctx, costCategoryId)
}

func (ccc *CostCategoryController) GetCostCategoryEntries(ctx context.Context, tripId *uuid.UUID) ([]*models.CostCategoryResponse, *models.ExpenseServiceError) {
	costCategories, repoErr := ccc.CostCategoryRepo.GetCostCategoriesByTripID(ctx, tripId)
	if repoErr != nil {
		return nil, repoErr
	}

	costCategoriesReponse := make([]*models.CostCategoryResponse, 0, len(costCategories))
	for _, costCategory := range costCategories {
		costCategoriesReponse = append(costCategoriesReponse, ccc.responseBuilder(ctx, &costCategory))
	}

	return costCategoriesReponse, nil
}

func (ccc *CostCategoryController) responseBuilder(ctx context.Context, costCategories *models.CostCategorySchema) *models.CostCategoryResponse {
	// Get total cost of cost category
	totalCost, err := ccc.CostRepo.GetTotalCostByCostCategoryID(ctx, costCategories.CostCategoryID)
	if err != nil {
		return nil
	}

	return &models.CostCategoryResponse{
		CostCategoryId: costCategories.CostCategoryID,
		Name:           costCategories.Name,
		Description:    costCategories.Description,
		Icon:           costCategories.Icon,
		Color:          costCategories.Color,
		TotalCost:      totalCost.String(),
	}
}
