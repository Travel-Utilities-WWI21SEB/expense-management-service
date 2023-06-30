package repositories

import (
	"context"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

type DebtRepo interface {
	GetDebtById(ctx context.Context, debtId *uuid.UUID) (*models.DebtSchema, *models.ExpenseServiceError)
	GetDebts(ctx context.Context) ([]*models.DebtSchema, *models.ExpenseServiceError)
	AddTx(ctx context.Context, tx pgx.Tx, debt *models.DebtSchema) *models.ExpenseServiceError
	UpdateTx(ctx context.Context, tx pgx.Tx, debt *models.DebtSchema) *models.ExpenseServiceError
	DeleteTx(ctx context.Context, tx pgx.Tx, debtId *uuid.UUID) *models.ExpenseServiceError

	GetDebtByCreditorId(ctx context.Context, creditorId *uuid.UUID) (*models.DebtSchema, *models.ExpenseServiceError)
	GetDebtByCreditorIdAndDebtorIdAndTripId(ctx context.Context, creditorId *uuid.UUID, debtorId *uuid.UUID, tripId *uuid.UUID) (*models.DebtSchema, *models.ExpenseServiceError)
	GetDebtEntriesByTripId(ctx context.Context, tripId *uuid.UUID) ([]*models.DebtSchema, *models.ExpenseServiceError)

	CalculateDebt(ctx context.Context, tx pgx.Tx, creditorId *uuid.UUID, debtorId *uuid.UUID, tripId *uuid.UUID, amountToAdd decimal.Decimal) *models.ExpenseServiceError
}

type DebtRepository struct {
	DatabaseMgr *managers.DatabaseManager
}

func (dr *DebtRepository) GetDebtById(ctx context.Context, debtId *uuid.UUID) (*models.DebtSchema, *models.ExpenseServiceError) {
	query := "SELECT id, id_creditor, id_debtor, id_trip, amount, currency_code, created_at, updated_at FROM debt WHERE id = $1"
	row := dr.DatabaseMgr.ExecuteQueryRow(ctx, query, debtId)
	debt := &models.DebtSchema{}

	err := row.Scan(&debt.DebtID, &debt.CreditorId, &debt.DebtorId, &debt.TripId, &debt.Amount, &debt.CurrencyCode, &debt.CreationDate, &debt.UpdateDate)
	if err != nil {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	return debt, nil
}

func (dr *DebtRepository) GetDebts(ctx context.Context) ([]*models.DebtSchema, *models.ExpenseServiceError) {
	query := "SELECT * FROM debt"
	rows, err := dr.DatabaseMgr.ExecuteQuery(ctx, query)
	if err != nil {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	var debts []*models.DebtSchema
	for rows.Next() {
		debt := &models.DebtSchema{}
		err := rows.Scan(&debt.DebtID, &debt.CreditorId, &debt.DebtorId, &debt.TripId, &debt.Amount, &debt.CurrencyCode, &debt.CreationDate, &debt.UpdateDate)
		if err != nil {
			return nil, expense_errors.EXPENSE_BAD_REQUEST
		}

		debts = append(debts, debt)
	}

	return debts, nil
}

func (dr *DebtRepository) AddTx(ctx context.Context, tx pgx.Tx, debt *models.DebtSchema) *models.ExpenseServiceError {
	query := "INSERT INTO debt (id, id_creditor, id_debtor, id_trip, amount, currency_code, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := tx.Exec(ctx, query, debt.DebtID, debt.CreditorId, debt.DebtorId, debt.TripId, debt.Amount, debt.CurrencyCode, debt.CreationDate, debt.UpdateDate)
	if err != nil {
		return expense_errors.EXPENSE_BAD_REQUEST
	}

	return nil
}

func (dr *DebtRepository) UpdateTx(ctx context.Context, tx pgx.Tx, debt *models.DebtSchema) *models.ExpenseServiceError {
	query := "UPDATE debt SET id_creditor = $1, id_debtor = $2, id_trip = $3, amount = $4, currency_code = $5, updated_at = $6 WHERE id = $7"
	log.Printf("Debt: %v \t time: %v", debt.DebtID.String(), debt.UpdateDate)
	_, err := tx.Exec(ctx, query, debt.CreditorId, debt.DebtorId, debt.TripId, debt.Amount, debt.CurrencyCode, debt.UpdateDate, debt.DebtID)
	if err != nil {
		return expense_errors.EXPENSE_BAD_REQUEST
	}

	return nil
}

func (dr *DebtRepository) DeleteTx(ctx context.Context, tx pgx.Tx, debtId *uuid.UUID) *models.ExpenseServiceError {
	query := "DELETE FROM debt WHERE id = $1"
	_, err := tx.Exec(ctx, query, debtId)
	if err != nil {
		return expense_errors.EXPENSE_BAD_REQUEST
	}

	return nil
}

func (dr *DebtRepository) GetDebtByCreditorId(ctx context.Context, creditorId *uuid.UUID) (*models.DebtSchema, *models.ExpenseServiceError) {
	query := "SELECT * FROM debt WHERE id_creditor = $1"
	row := dr.DatabaseMgr.ExecuteQueryRow(ctx, query, creditorId)
	debt := &models.DebtSchema{}

	err := row.Scan(&debt.DebtID, &debt.CreditorId, &debt.DebtorId, &debt.TripId, &debt.Amount, &debt.CurrencyCode, &debt.CreationDate, &debt.UpdateDate)
	if err != nil {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	return debt, nil
}

func (dr *DebtRepository) GetDebtByCreditorIdAndDebtorIdAndTripId(ctx context.Context, creditorId *uuid.UUID, debtorId *uuid.UUID, tripId *uuid.UUID) (*models.DebtSchema, *models.ExpenseServiceError) {
	query := "SELECT id, id_creditor, id_debtor, id_trip, amount, currency_code, created_at, updated_at FROM debt WHERE id_creditor = $1 AND id_debtor = $2 AND id_trip = $3"
	row := dr.DatabaseMgr.ExecuteQueryRow(ctx, query, creditorId, debtorId, tripId)
	debt := &models.DebtSchema{}

	err := row.Scan(&debt.DebtID, &debt.CreditorId, &debt.DebtorId, &debt.TripId, &debt.Amount, &debt.CurrencyCode, &debt.CreationDate, &debt.UpdateDate)
	if err != nil {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	return debt, nil
}

func (dr *DebtRepository) GetDebtEntriesByTripId(ctx context.Context, tripId *uuid.UUID) ([]*models.DebtSchema, *models.ExpenseServiceError) {
	query := "SELECT id, id_creditor, id_debtor, id_trip, amount, currency_code, created_at, updated_at FROM debt WHERE id_trip = $1"
	rows, err := dr.DatabaseMgr.ExecuteQuery(ctx, query, tripId)
	if err != nil {
		log.Printf("Error while getting debt entries by trip id: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	var debts []*models.DebtSchema
	for rows.Next() {
		debt := &models.DebtSchema{}
		err := rows.Scan(&debt.DebtID, &debt.CreditorId, &debt.DebtorId, &debt.TripId, &debt.Amount, &debt.CurrencyCode, &debt.CreationDate, &debt.UpdateDate)
		if err != nil {
			return nil, expense_errors.EXPENSE_BAD_REQUEST
		}

		debts = append(debts, debt)
	}

	return debts, nil
}

func (dr *DebtRepository) CalculateDebt(ctx context.Context, tx pgx.Tx, creditorId *uuid.UUID, debtorId *uuid.UUID, tripId *uuid.UUID, amountToAdd decimal.Decimal) *models.ExpenseServiceError {
	// Check if creditor and debtor are the same
	if creditorId.String() == debtorId.String() {
		return nil
	}

	now := time.Now()
	debt, repoErr := dr.GetDebtByCreditorIdAndDebtorIdAndTripId(ctx, creditorId, debtorId, tripId)
	if repoErr != nil {
		return repoErr
	}

	debt.Amount = debt.Amount.Add(amountToAdd)
	debt.UpdateDate = &now
	repoErr = dr.UpdateTx(ctx, tx, debt)
	if repoErr != nil {
		return repoErr
	}

	otherDebt, repoErr := dr.GetDebtByCreditorIdAndDebtorIdAndTripId(ctx, debtorId, creditorId, tripId)
	if repoErr != nil {
		return repoErr
	}

	// Update existing debt
	otherDebt.Amount = otherDebt.Amount.Sub(amountToAdd)
	otherDebt.UpdateDate = &now
	repoErr = dr.UpdateTx(ctx, tx, otherDebt)
	if repoErr != nil {
		return repoErr
	}

	return nil
}
