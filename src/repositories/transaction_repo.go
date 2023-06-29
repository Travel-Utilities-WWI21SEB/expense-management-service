package repositories

import (
	"database/sql"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/google/uuid"
	"log"
)

type TransactionRepo interface {
	AddTx(tx *sql.Tx, transaction *models.TransactionSchema) *models.ExpenseServiceError
	DeleteTx(tx *sql.Tx, transactionId *uuid.UUID) *models.ExpenseServiceError

	GetTransactionsByTripIdAndUserId(tripId *uuid.UUID, userId *uuid.UUID) ([]*models.TransactionSchema, *models.ExpenseServiceError)
	GetTransactionById(id *uuid.UUID) (*models.TransactionSchema, *models.ExpenseServiceError)
}

type TransactionRepository struct {
	DatabaseMgr managers.DatabaseMgr
}

func (tr *TransactionRepository) AddTx(tx *sql.Tx, transaction *models.TransactionSchema) *models.ExpenseServiceError {
	query := "INSERT INTO transaction (id, id_creditor, id_debtor, id_trip, amount, created_at, currency_code, is_confirmed) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := tx.Exec(query, transaction.TransactionId, transaction.CreditorId, transaction.DebtorId, transaction.TripId, transaction.Amount, transaction.CreationDate, transaction.CurrencyCode, transaction.IsConfirmed)
	if err != nil {
		log.Printf("Error while executing query: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}
	return nil
}

func (tr *TransactionRepository) DeleteTx(tx *sql.Tx, transactionId *uuid.UUID) *models.ExpenseServiceError {
	query := "DELETE FROM transaction WHERE id = $1"
	_, err := tx.Exec(query, transactionId)
	if err != nil {
		log.Printf("Error while executing query: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}
	return nil
}

func (tr *TransactionRepository) GetTransactionById(id *uuid.UUID) (*models.TransactionSchema, *models.ExpenseServiceError) {
	query := "SELECT id, id_creditor, id_debtor, id_trip, amount, created_at, currency_code, is_confirmed FROM transaction WHERE id = $1"
	row := tr.DatabaseMgr.ExecuteQueryRow(query, id)

	var transaction models.TransactionSchema
	err := row.Scan(&transaction.TransactionId, &transaction.CreditorId, &transaction.DebtorId, &transaction.TripId, &transaction.Amount, &transaction.CreationDate, &transaction.CurrencyCode, &transaction.IsConfirmed)
	if err != nil {
		log.Printf("Error while scanning transaction: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return &transaction, nil
}

func (tr *TransactionRepository) GetTransactionsByTripIdAndUserId(tripId *uuid.UUID, userId *uuid.UUID) ([]*models.TransactionSchema, *models.ExpenseServiceError) {
	query := "SELECT id, id_creditor, id_debtor, id_trip, amount, created_at, is_confirmed FROM transaction WHERE id_trip = $1 AND (id_creditor = $2 OR id_debtor = $2)"
	rows, err := tr.DatabaseMgr.ExecuteQuery(query, tripId, userId)
	if err != nil {
		log.Printf("Error while executing query: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	transactions := make([]*models.TransactionSchema, 0)
	for rows.Next() {
		var transaction models.TransactionSchema
		err = rows.Scan(&transaction.TransactionId, &transaction.CreditorId, &transaction.DebtorId, &transaction.TripId, &transaction.Amount, &transaction.CreationDate, &transaction.IsConfirmed)
		if err != nil {
			log.Printf("Error while scanning transaction: %v", err)
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}
