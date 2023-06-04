package controllers

import (
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
)

type CostCategoryCtl interface {
	CreateCostCategory(tripId *uuid.UUID, costCategoryRequest models.CostCategoryPostRequest) (*models.CostCategoryResponse, *models.ExpenseServiceError)
	PatchCostCategory(costCategoryId *uuid.UUID, costCategoryRequest models.CostCategoryPatchRequest) (*models.CostCategoryResponse, *models.ExpenseServiceError)
	GetCostCategoryDetails(costCategoryId *uuid.UUID) (*models.CostCategoryResponse, *models.ExpenseServiceError)
	DeleteCostCategory(costCategoryId *uuid.UUID) *models.ExpenseServiceError
	GetCostCategoryEntries(tripId *uuid.UUID) ([]*models.CostCategoryResponse, *models.ExpenseServiceError)
}

type CostCategoryController struct {
	DatabaseMgr managers.DatabaseMgr
}

func (ccc *CostCategoryController) CreateCostCategory(tripId *uuid.UUID, createCostCategoryRequest models.CostCategoryPostRequest) (*models.CostCategoryResponse, *models.ExpenseServiceError) {
	if utils.ContainsEmptyString(createCostCategoryRequest.Name) {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

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

	// Insert cost category into database
	insertString := "INSERT INTO cost_category (id, name, description, icon, color, id_trip) VALUES ($1, $2, $3, $4, $5, $6)"
	if _, err := ccc.DatabaseMgr.ExecuteQuery(insertString, costCategory.CostCategoryID, costCategory.Name, costCategory.Description, costCategory.Icon, costCategory.Color, costCategory.TripID); err != nil {
		// If the constraint was violated, the cost category already exists
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // 23505 = unique_violation
			log.Printf("CostCategoryController: Cost category already exists")
			return nil, expense_errors.EXPENSE_CONFLICT
		}
		log.Printf("CostCategoryController: Error while inserting cost category into database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Return created cost category
	response := &models.CostCategoryResponse{
		CostCategoryId: costCategory.CostCategoryID,
		Name:           costCategory.Name,
		Description:    costCategory.Description,
		Icon:           costCategory.Icon,
		Color:          costCategory.Color,
	}

	return response, nil
}

func (ccc *CostCategoryController) PatchCostCategory(costCategoryId *uuid.UUID, costCategoryPatchRequest models.CostCategoryPatchRequest) (*models.CostCategoryResponse, *models.ExpenseServiceError) {
	// Get cost category from database
	var costCategory models.CostCategorySchema

	queryString := "SELECT id, name, description, icon, color, id_trip FROM cost_category WHERE id = $1"
	row := ccc.DatabaseMgr.ExecuteQueryRow(queryString, costCategoryId)
	if err := row.Scan(&costCategory.CostCategoryID, &costCategory.Name, &costCategory.Description, &costCategory.Icon, &costCategory.Color, &costCategory.TripID); err != nil {
		log.Printf("CostCategoryController: Error while getting cost category from database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Update cost category
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
	updateString := "UPDATE cost_category SET name = $1, description = $2, icon = $3, color = $4 WHERE id = $5"
	if _, err := ccc.DatabaseMgr.ExecuteQuery(updateString, costCategory.Name, costCategory.Description, costCategory.Icon, costCategory.Color, costCategory.CostCategoryID); err != nil {
		log.Printf("CostCategoryController: Error while updating cost category in database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Return updated cost category
	response := &models.CostCategoryResponse{
		CostCategoryId: costCategory.CostCategoryID,
		Name:           costCategory.Name,
		Description:    costCategory.Description,
		Icon:           costCategory.Icon,
		Color:          costCategory.Color,
	}

	return response, nil
}

func (ccc *CostCategoryController) GetCostCategoryDetails(costCategoryId *uuid.UUID) (*models.CostCategoryResponse, *models.ExpenseServiceError) {
	// Get cost category from database
	var costCategoryResponse models.CostCategoryResponse

	queryString := "SELECT id, name, description, icon, color, id_trip FROM cost_category WHERE id = $1"
	row := ccc.DatabaseMgr.ExecuteQueryRow(queryString, costCategoryId)
	if err := row.Scan(&costCategoryResponse.CostCategoryId, &costCategoryResponse.Name, &costCategoryResponse.Description, &costCategoryResponse.Icon, &costCategoryResponse.Color); err != nil {
		log.Printf("CostCategoryController: Error while getting cost category from database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return &costCategoryResponse, nil
}

func (ccc *CostCategoryController) DeleteCostCategory(costCategoryId *uuid.UUID) *models.ExpenseServiceError {
	// Delete cost category from database
	deleteString := "DELETE FROM cost_category WHERE id = $1"
	if _, err := ccc.DatabaseMgr.ExecuteQuery(deleteString, costCategoryId); err != nil {
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

func (ccc *CostCategoryController) GetCostCategoryEntries(tripId *uuid.UUID) ([]*models.CostCategoryResponse, *models.ExpenseServiceError) {
	// Get cost categories from database
	var costCategoriesReponse []*models.CostCategoryResponse

	queryString := "SELECT id, name, description, icon, color FROM cost_category WHERE id_trip = $1"
	rows, err := ccc.DatabaseMgr.ExecuteQuery(queryString, &tripId)
	if err != nil {
		log.Printf("CostCategoryController: Error while getting cost categories from database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	for rows.Next() {
		var costCategoryResponse models.CostCategoryResponse
		if err := rows.Scan(&costCategoryResponse.CostCategoryId, &costCategoryResponse.Name, &costCategoryResponse.Description, &costCategoryResponse.Icon, &costCategoryResponse.Color); err != nil {
			log.Printf("CostCategoryController: Error while scanning cost category from database: %v", err)
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}
		costCategoriesReponse = append(costCategoriesReponse, &costCategoryResponse)
	}

	return costCategoriesReponse, nil
}
