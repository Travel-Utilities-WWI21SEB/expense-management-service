package repositories

import (
	"context"
	"errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"
	"log"
)

type CostRepo interface {
	GetCostByID(ctx context.Context, costId *uuid.UUID) (*models.CostSchema, *models.ExpenseServiceError)

	AddTx(ctx context.Context, tx pgx.Tx, cost *models.CostSchema) *models.ExpenseServiceError
	UpdateTx(ctx context.Context, tx pgx.Tx, cost *models.CostSchema) *models.ExpenseServiceError
	DeleteTx(ctx context.Context, tx pgx.Tx, costId *uuid.UUID) *models.ExpenseServiceError

	GetCostsByTripID(ctx context.Context, tripId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError)
	GetCostsByTripIDAndContributorID(ctx context.Context, tripId *uuid.UUID, contributorId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError)
	GetCostsByContributorID(ctx context.Context, contributorId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError)
	GetCostsByCostCategoryID(ctx context.Context, costCategoryId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError)

	GetCostContributors(ctx context.Context, costId *uuid.UUID) ([]*models.CostContributionSchema, *models.ExpenseServiceError)
	AddCostContributor(ctx context.Context, contributor *models.CostContributionSchema) *models.ExpenseServiceError
	UpdateCostContributor(ctx context.Context, contributor *models.CostContributionSchema) *models.ExpenseServiceError
	GetCostCreditor(ctx context.Context, id *uuid.UUID) (*models.UserSchema, *models.ExpenseServiceError)

	AddCostContributorTx(ctx context.Context, tx pgx.Tx, contributor *models.CostContributionSchema) *models.ExpenseServiceError
	DeleteCostContributionTx(ctx context.Context, tx pgx.Tx, contributorId *uuid.UUID, costId *uuid.UUID) *models.ExpenseServiceError

	GetTotalCostByTripID(ctx context.Context, tripId *uuid.UUID) (*decimal.Decimal, *models.ExpenseServiceError)
	GetTotalCostByCostCategoryID(ctx context.Context, costCategoryId *uuid.UUID) (*decimal.Decimal, *models.ExpenseServiceError)
	DeleteCostContributions(ctx context.Context, costId *uuid.UUID) *models.ExpenseServiceError
	GetCostsByCostCategoryIDAndContributorID(ctx context.Context, costCategoryId *uuid.UUID, userId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError)
	GetCostOverview(ctx context.Context, userId *uuid.UUID) (*models.CostOverviewDTO, *models.ExpenseServiceError)
}

type CostRepository struct {
	DatabaseMgr managers.DatabaseMgr
}

//********************************************************************************************************************\\
// Cost																												  \\
//********************************************************************************************************************\\

// GetCostOverview returns an overview of all costs
func (cr *CostRepository) GetCostOverview(ctx context.Context, userId *uuid.UUID) (*models.CostOverviewDTO, *models.ExpenseServiceError) {
	response := &models.CostOverviewDTO{}
	var tripDistribution []*models.TripDistributionDTO
	var costDistribution []*models.CostDistributionDTO
	var mostExpensiveTrip *models.TripNameToIdDTO
	var leastExpensiveTrip *models.TripNameToIdDTO
	totalTripCosts := decimal.NewFromInt(0)
	allCosts := decimal.NewFromInt(0)

	// Get every trip the user is part of
	rows, err := cr.DatabaseMgr.ExecuteQuery(ctx, "SELECT id, name FROM trip WHERE id IN (SELECT id_trip FROM user_trip_association WHERE id_user = $1)", userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, expense_errors.EXPENSE_NOT_FOUND
		}
		log.Printf("Error while executing query: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}
	defer rows.Close()

	// Iterate over every trip
	for rows.Next() {
		// Get trip costs for every trip and cost categories respectively
		tripId := uuid.UUID{}
		var name string
		var tripCosts decimal.Decimal

		if err := rows.Scan(&tripId, &name); err != nil {
			log.Printf("Error while scanning row: %v", err)
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}

		// Get total costs for trip
		// The outer COALESCE is needed because the inner COALESCE returns NULL if there are no costs for the trip
		queryString := "SELECT COALESCE(SUM(COALESCE(amount, 0.0)), 0.0) FROM cost WHERE id_cost_category IN (SELECT id FROM cost_category WHERE id_trip = $1)"
		row := cr.DatabaseMgr.ExecuteQueryRow(ctx, queryString, tripId)

		var allCostsForTrip decimal.Decimal
		if err := row.Scan(&allCostsForTrip); err != nil {
			log.Printf("Error while scanning row: %v", err)
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}
		allCosts = allCosts.Add(allCostsForTrip)

		// Get total costs for trip that the user is a part of grouped by the cost category
		// The outer COALESCE is needed because the inner COALESCE returns NULL if there are no costs for the trip
		queryString = "SELECT COALESCE(SUM(COALESCE(amount, 0.0)), 0.0), id_cost_category FROM cost WHERE id_cost_category IN (SELECT id FROM cost_category WHERE id_trip = $1) AND id IN (SELECT id_cost FROM user_cost_association WHERE id_user = $2) GROUP BY id_cost_category"
		costRow, costErr := cr.DatabaseMgr.ExecuteQuery(ctx, queryString, tripId, userId)

		if costErr != nil {
			log.Printf("Error while executing query: %v", costErr)
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}

		for costRow.Next() {
			var costCategoryCosts decimal.Decimal
			var costCategoryID uuid.UUID

			if err := costRow.Scan(&costCategoryCosts, &costCategoryID); err != nil {
				log.Printf("Error while scanning row: %v", err)
				return nil, expense_errors.EXPENSE_INTERNAL_ERROR
			}

			queryString := "SELECT name FROM cost_category WHERE id = $1"
			nameRow := cr.DatabaseMgr.ExecuteQueryRow(ctx, queryString, costCategoryID)
			if err != nil {
				if err == pgx.ErrNoRows {
					return nil, expense_errors.EXPENSE_NOT_FOUND
				}
				log.Printf("Error while executing query: %v", err)
				return nil, expense_errors.EXPENSE_INTERNAL_ERROR
			}

			var costCategoryName string
			if err := nameRow.Scan(&costCategoryName); err != nil {
				log.Printf("Error while scanning row: %v", err)
				return nil, expense_errors.EXPENSE_INTERNAL_ERROR
			}

			tripCosts = tripCosts.Add(costCategoryCosts)

			// Add cost category to cost distribution
			costDistribution = append(costDistribution, &models.CostDistributionDTO{
				CostCategoryName: costCategoryName,
				Amount:           costCategoryCosts.String(),
			})
		}
		costRow.Close()

		// Add trip to trip distribution
		tripDistribution = append(tripDistribution, &models.TripDistributionDTO{
			TripName: name,
			Amount:   tripCosts.String(),
		})
		totalTripCosts = totalTripCosts.Add(tripCosts)

		// Check if trip is most expensive trip
		if mostExpensiveTrip == nil || tripCosts.GreaterThan(decimal.RequireFromString(mostExpensiveTrip.Amount)) {
			mostExpensiveTrip = &models.TripNameToIdDTO{
				TripID:   tripId,
				TripName: name,
				Amount:   tripCosts.String(),
			}
		}

		// Check if trip is least expensive trip
		if leastExpensiveTrip == nil || tripCosts.LessThan(decimal.RequireFromString(leastExpensiveTrip.Amount)) {
			leastExpensiveTrip = &models.TripNameToIdDTO{
				TripID:   tripId,
				TripName: name,
				Amount:   tripCosts.String(),
			}
		}
	}

	response.TripDistribution = tripDistribution
	response.CostDistribution = costDistribution
	response.MostExpensiveTrip = mostExpensiveTrip
	response.LeastExpensiveTrip = leastExpensiveTrip
	response.TotalCosts = totalTripCosts.String()

	if !totalTripCosts.IsZero() {
		response.AverageTripCosts = totalTripCosts.Div(decimal.NewFromInt(int64(len(tripDistribution)))).String()
		response.AverageContributionPercentage = totalTripCosts.Div(allCosts).Mul(decimal.NewFromInt(100)).String()
	} else {
		response.AverageTripCosts = "0"
		response.AverageContributionPercentage = "0"
	}

	return response, nil
}

// GetCostByID returns a cost by its id
func (cr *CostRepository) GetCostByID(ctx context.Context, costId *uuid.UUID) (*models.CostSchema, *models.ExpenseServiceError) {
	cost := &models.CostSchema{}

	row := cr.DatabaseMgr.ExecuteQueryRow(ctx, "SELECT id, amount, description, created_at, deducted_at, end_date, id_cost_category FROM cost WHERE id = $1", costId)
	if err := row.Scan(&cost.CostID, &cost.Amount, &cost.Description, &cost.CreationDate, &cost.DeductionDate, &cost.EndDate, &cost.CostCategoryID); err != nil {
		if err == pgx.ErrNoRows {
			return nil, expense_errors.EXPENSE_NOT_FOUND
		}

		log.Printf("Error while scanning row: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return cost, nil
}

func (*CostRepository) AddTx(ctx context.Context, tx pgx.Tx, cost *models.CostSchema) *models.ExpenseServiceError {
	query := "INSERT INTO cost (id, amount, description, created_at, deducted_at, end_date, id_cost_category) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := tx.Exec(ctx, query, cost.CostID, cost.Amount, cost.Description, cost.CreationDate, cost.DeductionDate, cost.EndDate, cost.CostCategoryID)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr); pgxErr.Code == "foreign_key_violation" {
			return expense_errors.EXPENSE_NOT_FOUND // CostCategory not found
		}

		log.Printf("Error while inserting cost into database: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}
	return nil
}

func (*CostRepository) UpdateTx(ctx context.Context, tx pgx.Tx, cost *models.CostSchema) *models.ExpenseServiceError {
	query := "UPDATE cost SET amount = $1, description = $2, deducted_at = $3, end_date = $4, id_cost_category = $5 WHERE id = $6"
	result, err := tx.Exec(ctx, query, cost.Amount, cost.Description, cost.DeductionDate, cost.EndDate, cost.CostCategoryID, cost.CostID)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr); pgxErr.Code == "foreign_key_violation" {
			return expense_errors.EXPENSE_NOT_FOUND // CostCategory not found
		}

		log.Printf("Error while updating cost in database: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return expense_errors.EXPENSE_NOT_FOUND
	}

	return nil
}

func (*CostRepository) DeleteTx(ctx context.Context, tx pgx.Tx, costId *uuid.UUID) *models.ExpenseServiceError {
	query := "DELETE FROM cost WHERE id = $1"
	result, err := tx.Exec(ctx, query, costId)
	if err != nil {
		log.Printf("Error while deleting cost from database: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return expense_errors.EXPENSE_NOT_FOUND
	}

	return nil
}

// GetCostsByTripID returns all costs associated with a trip through the cost_category database table
func (cr *CostRepository) GetCostsByTripID(ctx context.Context, tripId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError) {
	rows, err := cr.DatabaseMgr.ExecuteQuery(ctx, "SELECT c.id, c.amount, c.description, c.created_at, c.deducted_at, c.end_date, c.id_cost_category FROM cost c INNER JOIN cost_category cc ON c.id_cost_category = cc.id WHERE cc.id_trip = $1", tripId)
	if err != nil {
		log.Printf("Error while querying database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}
	defer rows.Close()

	return getCostsFromRows(rows)
}

// GetCostsByCostCategoryID returns all costs associated with a cost category
func (cr *CostRepository) GetCostsByCostCategoryID(ctx context.Context, costCategoryId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError) {
	rows, err := cr.DatabaseMgr.ExecuteQuery(ctx, "SELECT id, amount, description, created_at, deducted_at, end_date, id_cost_category FROM cost WHERE id_cost_category = $1", costCategoryId)
	if err != nil {
		log.Printf("Error while querying database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}
	defer rows.Close()

	return getCostsFromRows(rows)
}

// GetCostsByTripIDAndContributorID returns all costs associated with a trip and a contributor
func (cr *CostRepository) GetCostsByTripIDAndContributorID(ctx context.Context, tripId *uuid.UUID, contributorId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError) {
	rows, err := cr.DatabaseMgr.ExecuteQuery(ctx, "SELECT c.id, c.amount, c.description, c.created_at, c.deducted_at, c.end_date, c.id_cost_category FROM cost c INNER JOIN user_cost_association uca ON c.id = uca.id_cost INNER JOIN cost_category cc ON c.id_cost_category = cc.id WHERE cc.id_trip = $1 AND uca.id_user = $2", tripId, contributorId)
	if err != nil {
		log.Printf("Error while querying database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}
	defer rows.Close()

	return getCostsFromRows(rows)
}

// GetCostsByCostCategoryIDAndContributorID returns all costs associated with a cost category and a contributor
func (cr *CostRepository) GetCostsByCostCategoryIDAndContributorID(ctx context.Context, costCategoryId *uuid.UUID, contributorId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError) {
	rows, err := cr.DatabaseMgr.ExecuteQuery(ctx, "SELECT c.id, c.amount, c.description, c.created_at, c.deducted_at, c.end_date, c.id_cost_category FROM cost c INNER JOIN user_cost_association uca ON c.id = uca.id_cost WHERE c.id_cost_category = $1 AND uca.id_user = $2", costCategoryId, contributorId)
	if err != nil {
		log.Printf("Error while querying database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}
	defer rows.Close()

	return getCostsFromRows(rows)
}

// GetCostsByContributorID returns all costs associated with a contributor
func (cr *CostRepository) GetCostsByContributorID(ctx context.Context, contributorId *uuid.UUID) ([]*models.CostSchema, *models.ExpenseServiceError) {
	rows, err := cr.DatabaseMgr.ExecuteQuery(ctx, "SELECT c.id, c.amount, c.description, c.created_at, c.deducted_at, c.end_date, c.id_cost_category FROM cost c INNER JOIN user_cost_association uca ON c.id = uca.id_cost WHERE uca.id_user = $1", contributorId)
	if err != nil {
		log.Printf("Error while querying database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}
	defer rows.Close()

	return getCostsFromRows(rows)
}

//********************************************************************************************************************\\
// Cost Contributor   																								  \\
//********************************************************************************************************************\\

// GetCostContributors returns all cost contributors associated with a cost
func (cr *CostRepository) GetCostContributors(ctx context.Context, costId *uuid.UUID) ([]*models.CostContributionSchema, *models.ExpenseServiceError) {
	rows, err := cr.DatabaseMgr.ExecuteQuery(ctx, "SELECT id_user, id_cost, is_creditor, amount FROM user_cost_association WHERE id_cost = $1", costId)
	if err != nil {
		log.Printf("Error while querying database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}
	defer rows.Close()

	contributors := make([]*models.CostContributionSchema, 0)
	for rows.Next() {
		var contributor models.CostContributionSchema
		err := rows.Scan(&contributor.UserID, &contributor.CostID, &contributor.IsCreditor, &contributor.Amount)
		if err != nil {
			log.Printf("Error while scanning row: %v", err)
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}
		contributors = append(contributors, &contributor)
	}

	return contributors, nil
}

// AddCostContributor adds a cost contributor to the database table user_cost_association
func (cr *CostRepository) AddCostContributor(ctx context.Context, contributor *models.CostContributionSchema) *models.ExpenseServiceError {
	_, err := cr.DatabaseMgr.ExecuteStatement(ctx, "INSERT INTO user_cost_association (id_user, id_cost, is_creditor, amount) VALUES ($1, $2, $3, $4)", contributor.UserID, contributor.CostID, contributor.IsCreditor, contributor.Amount)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr); pgxErr.Code == "foreign_key_violation" {
			return expense_errors.EXPENSE_NOT_FOUND // Cost or User not found
		}

		log.Printf("Error while inserting cost into database: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}
	return nil
}

// UpdateCostContributor updates a cost contributor in the database table user_cost_association
func (cr *CostRepository) UpdateCostContributor(ctx context.Context, contributor *models.CostContributionSchema) *models.ExpenseServiceError {
	result, err := cr.DatabaseMgr.ExecuteStatement(ctx, "UPDATE user_cost_association SET is_creditor = $1 WHERE id_user = $2 AND id_cost = $3 AND amount = $4", contributor.IsCreditor, contributor.UserID, contributor.CostID, contributor.Amount)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr); pgxErr.Code == "foreign_key_violation" {
			log.Printf("Error while updating cost contributor: %v", err)
			return expense_errors.EXPENSE_NOT_FOUND // Cost or User not found
		}

		log.Printf("Error while inserting cost into database: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return expense_errors.EXPENSE_NOT_FOUND
	}

	return nil
}

// DeleteCostContributions deletes all cost contributions associated with a cost
func (cr *CostRepository) DeleteCostContributions(ctx context.Context, costId *uuid.UUID) *models.ExpenseServiceError {
	if _, err := cr.DatabaseMgr.ExecuteStatement(ctx, "DELETE FROM user_cost_association WHERE id_cost = $1", costId); err != nil {
		log.Printf("Error while deleting cost contributions: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

func (cr *CostRepository) GetCostCreditor(ctx context.Context, id *uuid.UUID) (*models.UserSchema, *models.ExpenseServiceError) {
	row := cr.DatabaseMgr.ExecuteQueryRow(ctx, "SELECT u.id, u.username, u.email FROM user_cost_association uca INNER JOIN \"user\" u ON uca.id_user = u.id WHERE uca.id_cost = $1 AND uca.is_creditor = true", id)

	var creditor models.UserSchema
	if err := row.Scan(&creditor.UserID, &creditor.Username, &creditor.Email); err != nil {
		log.Printf("Error while scanning row: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return &creditor, nil
}

func (*CostRepository) AddCostContributorTx(ctx context.Context, tx pgx.Tx, contributor *models.CostContributionSchema) *models.ExpenseServiceError {
	query := "INSERT INTO user_cost_association (id_user, id_cost, is_creditor, amount) VALUES ($1, $2, $3, $4)"
	_, err := tx.Exec(ctx, query, contributor.UserID, contributor.CostID, contributor.IsCreditor, contributor.Amount)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr); pgxErr.Code == "foreign_key_violation" {
			return expense_errors.EXPENSE_NOT_FOUND // Cost or User not found
		}

		log.Printf("Error while inserting cost into database: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}
	return nil
}

func (*CostRepository) DeleteCostContributionTx(ctx context.Context, tx pgx.Tx, contributorId *uuid.UUID, costId *uuid.UUID) *models.ExpenseServiceError {
	_, err := tx.Exec(ctx, "DELETE FROM user_cost_association WHERE id_user = $1 AND id_cost = $2", contributorId, costId)
	if err != nil {
		log.Printf("Error while deleting cost contribution: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

//********************************************************************************************************************\\
// Calculation  																									  \\
//********************************************************************************************************************\\

// GetTotalCostByTripID returns the total cost of a trip
func (cr *CostRepository) GetTotalCostByTripID(ctx context.Context, tripId *uuid.UUID) (*decimal.Decimal, *models.ExpenseServiceError) {
	row, err := cr.DatabaseMgr.ExecuteQuery(ctx, "SELECT COALESCE(SUM(c.amount),0) FROM cost c INNER JOIN cost_category cc ON c.id_cost_category = cc.id WHERE cc.id_trip = $1", tripId)
	if err != nil {
		log.Printf("Error while querying database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}
	defer row.Close()

	if !row.Next() {
		return nil, expense_errors.EXPENSE_NOT_FOUND // Trip not found
	}

	var totalCost decimal.Decimal
	err = row.Scan(&totalCost)
	if err != nil {
		log.Printf("Error while scanning row: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return &totalCost, nil
}

// GetTotalCostByCostCategoryID returns the total cost of a cost category
func (cr *CostRepository) GetTotalCostByCostCategoryID(ctx context.Context, costCategoryId *uuid.UUID) (*decimal.Decimal, *models.ExpenseServiceError) {
	var totalCost decimal.Decimal
	row := cr.DatabaseMgr.ExecuteQueryRow(ctx, "SELECT COALESCE(SUM(amount),0) FROM cost WHERE id_cost_category = $1", costCategoryId)
	err := row.Scan(&totalCost)
	if err != nil {
		log.Printf("Error while scanning row: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return &totalCost, nil
}

//********************************************************************************************************************\\
// Helper Functions  																								  \\
//********************************************************************************************************************\\

func getCostsFromRows(rows pgx.Rows) ([]*models.CostSchema, *models.ExpenseServiceError) {
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
