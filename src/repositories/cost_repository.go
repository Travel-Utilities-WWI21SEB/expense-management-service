package repositories

import (
	"database/sql"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
)

type CostRepo interface {
	CreateCost(cost *models.CostSchema) *models.ExpenseServiceError
	GetCostByID(costId *uuid.UUID) (*models.CostSchema, *models.ExpenseServiceError)
	UpdateCost(cost *models.CostSchema) *models.ExpenseServiceError
	DeleteCostEntry(costId *uuid.UUID) *models.ExpenseServiceError

	GetCostsByTripID(tripId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError)
	GetCostsByTripIDAndContributorID(tripId *uuid.UUID, contributorId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError)
	GetCostsByContributorID(contributorId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError)
	GetCostsByCostCategoryID(costCategoryId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError)

	GetCostContributors(costId *uuid.UUID) ([]*models.CostContributionSchema, *models.ExpenseServiceError)
	AddCostContributor(contributor *models.CostContributionSchema) *models.ExpenseServiceError
	UpdateCostContributor(contributor *models.CostContributionSchema) *models.ExpenseServiceError
	RemoveCostContributor(costId *uuid.UUID, userId *uuid.UUID) *models.ExpenseServiceError
}

type CostRepository struct {
	DatabaseMgr managers.DatabaseMgr
}

//********************************************************************************************************************\\
// Cost																												  \\
//********************************************************************************************************************\\

// CreateCost Creates a new cost in the database
func (cr *CostRepository) CreateCost(cost *models.CostSchema) *models.ExpenseServiceError {
	_, err := cr.DatabaseMgr.ExecuteStatement("INSERT INTO cost (id, amount, description, created_at, deducted_at, end_date, id_cost_category) VALUES ($1, $2, $3, $4, $5, $6, $7)", cost.CostID, cost.Amount, cost.Description, cost.CreationDate, cost.DeductionDate, cost.EndDate, cost.CostCategoryID)
	if err != nil {

		// Check if cost category exists
		if pqErr := err.(*pq.Error); pqErr.Code.Name() == "foreign_key_violation" {
			return expense_errors.EXPENSE_NOT_FOUND // CostCategory not found
		}

		log.Printf("Error while inserting cost into database: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}
	return nil
}

// GetCostByID returns a cost by its id
func (cr *CostRepository) GetCostByID(costId *uuid.UUID) (*models.CostSchema, *models.ExpenseServiceError) {
	cost := &models.CostSchema{}

	row := cr.DatabaseMgr.ExecuteQueryRow("SELECT id, amount, description, created_at, deducted_at, end_date, id_cost_category FROM cost WHERE id = $1", costId)
	if err := row.Scan(&cost.CostID, &cost.Amount, &cost.Description, &cost.CreationDate, &cost.DeductionDate, &cost.EndDate, &cost.CostCategoryID); err != nil {
		if err == sql.ErrNoRows {
			return nil, expense_errors.EXPENSE_NOT_FOUND
		}

		log.Printf("Error while scanning row: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return cost, nil
}

// UpdateCost updates a cost in the database
func (cr *CostRepository) UpdateCost(cost *models.CostSchema) *models.ExpenseServiceError {
	result, err := cr.DatabaseMgr.ExecuteStatement("UPDATE cost SET amount = $1, description = $2, created_at = $3, deducted_at = $4, end_date = $5, id_cost_category = $6 WHERE id = $7", cost.Amount, cost.Description, cost.CreationDate, cost.DeductionDate, cost.EndDate, cost.CostCategoryID, cost.CostID)
	if err != nil {
		if pqErr := err.(*pq.Error); pqErr.Code.Name() == "foreign_key_violation" {
			return expense_errors.EXPENSE_NOT_FOUND // CostCategory not found
		}

		log.Printf("Error while updating cost in database: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return expense_errors.EXPENSE_NOT_FOUND
	}

	return nil
}

// DeleteCostEntry deletes a cost from the database
func (cr *CostRepository) DeleteCostEntry(costId *uuid.UUID) *models.ExpenseServiceError {
	result, err := cr.DatabaseMgr.ExecuteStatement("DELETE FROM cost WHERE id = $1", costId)
	if err != nil {
		log.Printf("Error while deleting cost from database: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return expense_errors.EXPENSE_NOT_FOUND
	}

	return nil
}

// GetCostsByTripID returns all costs associated with a trip through the cost_category database table
func (cr *CostRepository) GetCostsByTripID(tripId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError) {
	rows, err := cr.DatabaseMgr.ExecuteQuery("SELECT c.id, c.amount, c.description, c.created_at, c.deducted_at, c.end_date, c.id_cost_category FROM cost c INNER JOIN cost_category cc ON c.id_cost_category = cc.id WHERE cc.id_trip = $1", tripId)
	if err != nil {
		log.Printf("Error while querying database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return getCostsFromRows(rows)
}

// GetCostsByCostCategoryID returns all costs associated with a cost category
func (cr *CostRepository) GetCostsByCostCategoryID(costCategoryId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError) {
	rows, err := cr.DatabaseMgr.ExecuteQuery("SELECT id, amount, description, created_at, deducted_at, end_date, id_cost_category FROM cost WHERE id_cost_category = $1", costCategoryId)
	if err != nil {
		log.Printf("Error while querying database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return getCostsFromRows(rows)
}

// GetCostsByTripIDAndContributorID returns all costs associated with a trip and a contributor
func (cr *CostRepository) GetCostsByTripIDAndContributorID(tripId *uuid.UUID, contributorId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError) {
	rows, err := cr.DatabaseMgr.ExecuteQuery("SELECT c.id, c.amount, c.description, c.created_at, c.deducted_at, c.end_date, c.id_cost_category FROM cost c INNER JOIN user_cost_association uca ON c.id = uca.id_cost INNER JOIN cost_category cc ON c.id_cost_category = cc.id WHERE cc.id_trip = $1 AND uca.id_user = $2", tripId, contributorId)
	if err != nil {
		log.Printf("Error while querying database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return getCostsFromRows(rows)
}

// GetCostsByContributorID returns all costs associated with a contributor
func (cr *CostRepository) GetCostsByContributorID(contributorId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError) {
	rows, err := cr.DatabaseMgr.ExecuteQuery("SELECT c.id, c.amount, c.description, c.created_at, c.deducted_at, c.end_date, c.id_cost_category FROM cost c INNER JOIN user_cost_association uca ON c.id = uca.id_cost WHERE uca.id_user = $1", contributorId)
	if err != nil {
		log.Printf("Error while querying database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return getCostsFromRows(rows)
}

//********************************************************************************************************************\\
// Cost Contributor   																								  \\
//********************************************************************************************************************\\

// GetCostContributors returns all cost contributors associated with a cost
func (cr *CostRepository) GetCostContributors(costId *uuid.UUID) ([]*models.CostContributionSchema, *models.ExpenseServiceError) {
	rows, err := cr.DatabaseMgr.ExecuteQuery("SELECT id_user, id_cost, is_creditor FROM user_cost_association WHERE id_cost = $1", costId)
	if err != nil {
		log.Printf("Error while querying database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	contributors := make([]*models.CostContributionSchema, 0) // Empty slice
	for rows.Next() {
		var contributor models.CostContributionSchema
		err := rows.Scan(&contributor.UserID, &contributor.CostID, &contributor.IsCreditor)
		if err != nil {
			log.Printf("Error while scanning row: %v", err)
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}
		contributors = append(contributors, &contributor)
	}

	return contributors, nil
}

// AddCostContributor adds a cost contributor to the database table user_cost_association
func (cr *CostRepository) AddCostContributor(contributor *models.CostContributionSchema) *models.ExpenseServiceError {
	_, err := cr.DatabaseMgr.ExecuteStatement("INSERT INTO user_cost_association (id_user, id_cost, is_creditor) VALUES ($1, $2, $3)", contributor.UserID, contributor.CostID, contributor.IsCreditor)
	if err != nil {
		if pqErr := err.(*pq.Error); pqErr.Code.Name() == "foreign_key_violation" {
			return expense_errors.EXPENSE_NOT_FOUND // Cost or User not found
		}

		log.Printf("Error while inserting cost into database: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}
	return nil
}

// UpdateCostContributor updates a cost contributor in the database table user_cost_association
func (cr *CostRepository) UpdateCostContributor(contributor *models.CostContributionSchema) *models.ExpenseServiceError {
	result, err := cr.DatabaseMgr.ExecuteStatement("UPDATE user_cost_association SET is_creditor = $1 WHERE id_user = $2 AND id_cost = $3", contributor.IsCreditor, contributor.UserID, contributor.CostID)
	if err != nil {
		if pqErr := err.(*pq.Error); pqErr.Code.Name() == "foreign_key_violation" {
			log.Printf("Error while updating cost contributor: %v", err)
			return expense_errors.EXPENSE_NOT_FOUND // Cost or User not found
		}

		log.Printf("Error while inserting cost into database: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return expense_errors.EXPENSE_NOT_FOUND
	}

	return nil
}

// RemoveCostContributor removes a cost contributor from the database table user_cost_association
func (cr *CostRepository) RemoveCostContributor(costId *uuid.UUID, userId *uuid.UUID) *models.ExpenseServiceError {
	result, err := cr.DatabaseMgr.ExecuteStatement("DELETE FROM user_cost_association WHERE id_user = $1 AND id_cost = $2", userId, costId)
	if err != nil {
		log.Printf("Error while deleting cost contributor: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return expense_errors.EXPENSE_NOT_FOUND
	}

	return nil
}

//********************************************************************************************************************\\
// Helper Functions  																								  \\
//********************************************************************************************************************\\

func getCostsFromRows(rows *sql.Rows) ([]*models.CostSchema, *models.ExpenseServiceError) {
	costs := make([]*models.CostSchema, 0) // Empty slice
	for rows.Next() {
		var cost models.CostSchema
		if err := rows.Scan(&cost.CostID, &cost.Amount, &cost.Description, &cost.CreationDate, &cost.DeductionDate, &cost.EndDate, &cost.CostCategoryID); err != nil {
			log.Printf("Error while scanning row: %v", err)
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}
		costs = append(costs, &cost)
	}

	return costs, nil
}
