package repositories

import (
	"database/sql"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/google/uuid"
)

type DebtRepo interface {
	GetDebtById(debtId *uuid.UUID) (*models.DebtSchema, *models.ExpenseServiceError)
	GetDebts() ([]*models.DebtSchema, *models.ExpenseServiceError)
	AddTx(tx *sql.Tx, debt *models.DebtSchema) *models.ExpenseServiceError
	UpdateTx(tx *sql.Tx, debt *models.DebtSchema) *models.ExpenseServiceError
	DeleteTx(tx *sql.Tx, debtId *uuid.UUID) *models.ExpenseServiceError
}

type DebtRepository struct {
	DatabaseMgr *managers.DatabaseManager
}

func (dr *DebtRepository) GetDebtById(debtId *uuid.UUID) (*models.DebtSchema, *models.ExpenseServiceError) {
	query := "SELECT * FROM debt WHERE id = $1"
	row := dr.DatabaseMgr.ExecuteQueryRow(query, debtId)
	debt := &models.DebtSchema{}

	err := row.Scan(&debt.DebtID, &debt.CreditorId, &debt.DebtorId, &debt.TripId, &debt.Amount, &debt.CurrencyCode, &debt.CreationDate, &debt.UpdateDate)
	if err != nil {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	return debt, nil
}

func (dr *DebtRepository) GetDebts() ([]*models.DebtSchema, *models.ExpenseServiceError) {
	query := "SELECT * FROM debt"
	rows, err := dr.DatabaseMgr.ExecuteQuery(query)
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

func (dr *DebtRepository) AddTx(tx *sql.Tx, debt *models.DebtSchema) *models.ExpenseServiceError {
	query := "INSERT INTO debt (id, id_creditor, id_debtor, id_trip, amount, currency_code, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := tx.Exec(query, debt.DebtID, debt.CreditorId, debt.DebtorId, debt.TripId, debt.Amount, debt.CurrencyCode, debt.CreationDate, debt.UpdateDate)
	if err != nil {
		return expense_errors.EXPENSE_BAD_REQUEST
	}

	return nil
}

func (dr *DebtRepository) UpdateTx(tx *sql.Tx, debt *models.DebtSchema) *models.ExpenseServiceError {
	query := "UPDATE debt SET id_creditor = $1, id_debtor = $2, id_trip = $3, amount = $4, currency_code = $5, updated_at = $6 WHERE id = $7"
	_, err := tx.Exec(query, debt.CreditorId, debt.DebtorId, debt.TripId, debt.Amount, debt.CurrencyCode, debt.UpdateDate, debt.DebtID)
	if err != nil {
		return expense_errors.EXPENSE_BAD_REQUEST
	}

	return nil
}

func (dr *DebtRepository) DeleteTx(tx *sql.Tx, debtId *uuid.UUID) *models.ExpenseServiceError {
	query := "DELETE FROM debt WHERE id = $1"
	_, err := tx.Exec(query, debtId)
	if err != nil {
		return expense_errors.EXPENSE_BAD_REQUEST
	}

	return nil
}
