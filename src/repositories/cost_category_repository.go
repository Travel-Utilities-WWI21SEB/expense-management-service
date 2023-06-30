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
	"log"
)

type CostCategoryRepo interface {
	CreateCostCategory(ctx context.Context, costCategory *models.CostCategorySchema) *models.ExpenseServiceError
	GetCostCategoryByID(ctx context.Context, uuid *uuid.UUID) (*models.CostCategorySchema, *models.ExpenseServiceError)
	GetCostCategoriesByTripID(ctx context.Context, uuid *uuid.UUID) ([]models.CostCategorySchema, *models.ExpenseServiceError)
	UpdateCostCategory(ctx context.Context, costCategory *models.CostCategorySchema) *models.ExpenseServiceError
	DeleteCostCategory(ctx context.Context, uuid *uuid.UUID) *models.ExpenseServiceError

	GetCostCategoryByTripIdAndName(ctx context.Context, tripId *uuid.UUID, name string) (*models.CostCategorySchema, *models.ExpenseServiceError)
}

type CostCategoryRepository struct {
	DatabaseMgr managers.DatabaseMgr
}

func (ccr *CostCategoryRepository) CreateCostCategory(ctx context.Context, costCategory *models.CostCategorySchema) *models.ExpenseServiceError {
	_, err := ccr.DatabaseMgr.ExecuteStatement(ctx, "INSERT INTO cost_category (id, name, description, icon, color, id_trip) VALUES ($1, $2, $3, $4, $5, $6)", costCategory.CostCategoryID, costCategory.Name, costCategory.Description, costCategory.Icon, costCategory.Color, costCategory.TripID)
	if err != nil {
		// Check if cost category already exists
		var pqxErr *pgconn.PgError
		if errors.As(err, &pqxErr); pqxErr.Code == "unique_violation" {
			return expense_errors.EXPENSE_CONFLICT
		}
		log.Printf("Error while inserting cost category into database: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}
	return nil
}

func (ccr *CostCategoryRepository) GetCostCategoryByID(ctx context.Context, uuid *uuid.UUID) (*models.CostCategorySchema, *models.ExpenseServiceError) {
	schema := &models.CostCategorySchema{}

	row := ccr.DatabaseMgr.ExecuteQueryRow(ctx, "SELECT * FROM cost_category WHERE id = $1", uuid)
	if err := row.Scan(&schema.CostCategoryID, &schema.Name, &schema.Description, &schema.Icon, &schema.Color, &schema.TripID); err != nil {
		// Check if no cost category was found
		if err == pgx.ErrNoRows {
			return nil, expense_errors.EXPENSE_NOT_FOUND
		}

		log.Printf("Error while scanning cost category from database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return schema, nil
}

func (ccr *CostCategoryRepository) GetCostCategoriesByTripID(ctx context.Context, tripId *uuid.UUID) ([]models.CostCategorySchema, *models.ExpenseServiceError) {
	schemas := make([]models.CostCategorySchema, 0)

	rows, err := ccr.DatabaseMgr.ExecuteQuery(ctx, "SELECT * FROM cost_category WHERE id_trip = $1", tripId)
	if err != nil {
		log.Printf("Error while scanning cost categories from database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}
	defer rows.Close()

	for rows.Next() {
		schema := models.CostCategorySchema{}
		if err := rows.Scan(&schema.CostCategoryID, &schema.Name, &schema.Description, &schema.Icon, &schema.Color, &schema.TripID); err != nil {
			log.Printf("Error while scanning cost categories from database: %v", err)
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}
		schemas = append(schemas, schema)
	}

	return schemas, nil
}

func (ccr *CostCategoryRepository) UpdateCostCategory(ctx context.Context, costCategory *models.CostCategorySchema) *models.ExpenseServiceError {
	result, err := ccr.DatabaseMgr.ExecuteStatement(ctx, "UPDATE cost_category SET name = $1, description = $2, icon = $3, color = $4 WHERE id = $5", costCategory.Name, costCategory.Description, costCategory.Icon, costCategory.Color, costCategory.CostCategoryID)
	if err != nil {
		// Check if cost category already exists
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr); pgxErr.Code == "unique_violation" {
			return expense_errors.EXPENSE_CONFLICT
		}
		log.Printf("Error while updating cost category in database: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return expense_errors.EXPENSE_NOT_FOUND
	}

	return nil
}

func (ccr *CostCategoryRepository) DeleteCostCategory(ctx context.Context, uuid *uuid.UUID) *models.ExpenseServiceError {
	result, err := ccr.DatabaseMgr.ExecuteStatement(ctx, "DELETE FROM cost_category WHERE id = $1", uuid)
	if err != nil {
		log.Printf("Error while deleting cost category from database: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return expense_errors.EXPENSE_NOT_FOUND
	}

	return nil
}

func (ccr *CostCategoryRepository) GetCostCategoryByTripIdAndName(ctx context.Context, tripId *uuid.UUID, name string) (*models.CostCategorySchema, *models.ExpenseServiceError) {
	schema := &models.CostCategorySchema{}

	row := ccr.DatabaseMgr.ExecuteQueryRow(ctx, "SELECT * FROM cost_category WHERE id_trip = $1 AND name = $2", tripId, name)
	if err := row.Scan(&schema.CostCategoryID, &schema.Name, &schema.Description, &schema.Icon, &schema.Color, &schema.TripID); err != nil {
		// Check if no cost category was found
		if err == pgx.ErrNoRows {
			return nil, expense_errors.EXPENSE_NOT_FOUND
		}

		log.Printf("Error while scanning cost category from database: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return schema, nil
}
