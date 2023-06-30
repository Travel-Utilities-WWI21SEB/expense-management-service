package repositories

import (
	"context"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log"
)

type TransactionRepo interface {
	AddTx(ctx context.Context, tx pgx.Tx, transaction *models.TransactionSchema) *models.ExpenseServiceError
	DeleteTx(ctx context.Context, tx pgx.Tx, transactionId *uuid.UUID) *models.ExpenseServiceError

	GetTransactionsByTripIdAndUserId(ctx context.Context, tripId *uuid.UUID, userId *uuid.UUID) ([]*models.TransactionSchema, *models.ExpenseServiceError)
	GetTransactionById(ctx context.Context, id *uuid.UUID) (*models.TransactionSchema, *models.ExpenseServiceError)
}

type TransactionRepository struct {
	DatabaseMgr managers.DatabaseMgr
}

func (*TransactionRepository) AddTx(ctx context.Context, tx pgx.Tx, transaction *models.TransactionSchema) *models.ExpenseServiceError {
	query := "INSERT INTO transaction (id, id_creditor, id_debtor, id_trip, amount, created_at, currency_code, is_confirmed) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := tx.Exec(ctx, query, transaction.TransactionId, transaction.CreditorId, transaction.DebtorId, transaction.TripId, transaction.Amount, transaction.CreationDate, transaction.CurrencyCode, transaction.IsConfirmed)
	if err != nil {
		log.Printf("Error while executing query: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}
	return nil
}

func (*TransactionRepository) DeleteTx(ctx context.Context, tx pgx.Tx, transactionId *uuid.UUID) *models.ExpenseServiceError {
	query := "DELETE FROM transaction WHERE id = $1"
	_, err := tx.Exec(ctx, query, transactionId)
	if err != nil {
		log.Printf("Error while executing query: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}
	return nil
}

func (tr *TransactionRepository) GetTransactionById(ctx context.Context, id *uuid.UUID) (*models.TransactionSchema, *models.ExpenseServiceError) {
	query := "SELECT id, id_creditor, id_debtor, id_trip, amount, created_at, currency_code, is_confirmed FROM transaction WHERE id = $1"
	row := tr.DatabaseMgr.ExecuteQueryRow(ctx, query, id)

	var transaction models.TransactionSchema
	err := row.Scan(&transaction.TransactionId, &transaction.CreditorId, &transaction.DebtorId, &transaction.TripId, &transaction.Amount, &transaction.CreationDate, &transaction.CurrencyCode, &transaction.IsConfirmed)
	if err != nil {
		log.Printf("Error while scanning transaction: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return &transaction, nil
}

func (tr *TransactionRepository) GetTransactionsByTripIdAndUserId(ctx context.Context, tripId *uuid.UUID, userId *uuid.UUID) ([]*models.TransactionSchema, *models.ExpenseServiceError) {
	query := "SELECT id, id_creditor, id_debtor, id_trip, amount, created_at, is_confirmed FROM transaction WHERE id_trip = $1 AND (id_creditor = $2 OR id_debtor = $2)"
	rows, err := tr.DatabaseMgr.ExecuteQuery(ctx, query, tripId, userId)
	if err != nil {
		log.Printf("Error while executing query: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}
	defer rows.Close()

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
