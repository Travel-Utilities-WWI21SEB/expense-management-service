package controllers

import (
	"context"
	"database/sql"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/repositories"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

type TransactionCtl interface {
	GetTransactionEntries(ctx context.Context, tripId *uuid.UUID) ([]*models.TransactionDTO, *models.ExpenseServiceError)
	GetTransactionDetails(ctx context.Context, transactionId *uuid.UUID) (*models.TransactionDTO, *models.ExpenseServiceError)
	CreateTransactionEntry(ctx context.Context, tripId *uuid.UUID, transactionRequest *models.TransactionDTO) (*models.TransactionDTO, *models.ExpenseServiceError)
	DeleteTransactionEntry(ctx context.Context, transactionId *uuid.UUID) *models.ExpenseServiceError
}

type TransactionController struct {
	DatabaseMgr     managers.DatabaseMgr
	TransactionRepo repositories.TransactionRepo
	UserRepo        repositories.UserRepo
	TripRepo        repositories.TripRepo
	DebtRepo        repositories.DebtRepo
}

func (tc *TransactionController) GetTransactionEntries(ctx context.Context, tripId *uuid.UUID) ([]*models.TransactionDTO, *models.ExpenseServiceError) {
	// Get user id from context
	userId := ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)

	transactions, repoErr := tc.TransactionRepo.GetTransactionsByTripIdAndUserId(tripId, userId)
	if repoErr != nil {
		return nil, repoErr
	}

	transactionDTOs := make([]*models.TransactionDTO, 0)
	for _, transaction := range transactions {
		transactionDto, repoErr := tc.mapTransactionToDto(transaction)
		if repoErr != nil {
			return nil, repoErr
		}

		transactionDTOs = append(transactionDTOs, transactionDto)
	}

	return transactionDTOs, nil
}

func (tc *TransactionController) GetTransactionDetails(ctx context.Context, transactionId *uuid.UUID) (*models.TransactionDTO, *models.ExpenseServiceError) {
	// Get user id from context
	userId := ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)

	transaction, repoErr := tc.TransactionRepo.GetTransactionById(transactionId)
	if repoErr != nil {
		return nil, repoErr
	}

	// Check if user is part of transaction
	if transaction.CreditorId.String() != userId.String() && transaction.DebtorId.String() != userId.String() {
		return nil, expense_errors.EXPENSE_FORBIDDEN
	}

	transactionDto, repoErr := tc.mapTransactionToDto(transaction)
	if repoErr != nil {
		return nil, repoErr
	}

	return transactionDto, nil
}

func (tc *TransactionController) CreateTransactionEntry(ctx context.Context, tripId *uuid.UUID, transactionRequest *models.TransactionDTO) (*models.TransactionDTO, *models.ExpenseServiceError) {
	// Begin transaction
	tx, err := tc.DatabaseMgr.BeginTx(ctx)
	if err != nil {
		log.Printf("Error while beginning transaction: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Make sure to rollback the transaction if it fails
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Printf("Error while rolling back transaction: %v", err)
		}
	}(tx)

	// Get user id from context
	userId := ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)
	log.Printf("User id: %v", userId)
	// Get creditor from request
	creditor, repoErr := tc.UserRepo.GetUserById(userId)
	if repoErr != nil {
		return nil, repoErr
	}

	// Check if user is part of trip
	if repoErr := tc.TripRepo.ValidateIfUserHasAccepted(tripId, userId); repoErr != nil {
		return nil, repoErr
	}

	// Get debtor from request
	debtor, repoErr := tc.UserRepo.GetUserById(transactionRequest.DebtorId)
	if repoErr != nil {
		return nil, repoErr
	}

	// Check if debtor is part of trip
	if repoErr = tc.TripRepo.ValidateIfUserHasAccepted(tripId, transactionRequest.DebtorId); repoErr != nil {
		return nil, repoErr
	}

	// Create transaction entry
	transactionId := uuid.New()
	now := time.Now()

	transaction := &models.TransactionSchema{
		TransactionId: &transactionId,
		CreditorId:    userId,
		DebtorId:      debtor.UserID,
		TripId:        tripId,
		Amount:        decimal.RequireFromString(transactionRequest.Amount),
		CreationDate:  &now,
		CurrencyCode:  "EUR",
		IsConfirmed:   false,
	}

	if repoErr := tc.TransactionRepo.AddTx(tx, transaction); repoErr != nil {
		return nil, repoErr
	}

	// Add debt to debtor
	if repoErr := tc.DebtRepo.CalculateDebt(tx, debtor.UserID, creditor.UserID, tripId, transaction.Amount); repoErr != nil {
		return nil, repoErr
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error while committing transaction: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	response, repoErr := tc.mapTransactionToDto(transaction)
	if repoErr != nil {
		return nil, repoErr
	}

	return response, nil
}

func (tc *TransactionController) mapTransactionToDto(transaction *models.TransactionSchema) (*models.TransactionDTO, *models.ExpenseServiceError) {
	response := &models.TransactionDTO{
		TransactionId: transaction.TransactionId,
		Creditor:      nil,
		Debtor:        nil,
		Trip:          nil,
		Amount:        transaction.Amount.String(),
		CreationDate:  transaction.CreationDate.String(),
		IsConfirmed:   transaction.IsConfirmed,
	}

	// Get creditor from database
	creditor, repoErr := tc.UserRepo.GetUserById(transaction.CreditorId)
	if repoErr != nil {
		return nil, repoErr
	}

	creditorDto := models.UserDto{
		UserID:   creditor.UserID,
		Username: creditor.Username,
		Email:    creditor.Email,
	}

	response.Creditor = &creditorDto

	// Get debtor from database
	debtor, repoErr := tc.UserRepo.GetUserById(transaction.DebtorId)
	if repoErr != nil {
		return nil, repoErr
	}

	debtorDto := models.UserDto{
		UserID:   debtor.UserID,
		Username: debtor.Username,
		Email:    debtor.Email,
	}

	response.Debtor = &debtorDto

	// Get trip from database
	trip, repoErr := tc.TripRepo.GetTripById(transaction.TripId)
	if repoErr != nil {
		return nil, repoErr
	}

	tripDto := models.SlimTripDTO{
		TripID:    trip.TripID,
		Name:      trip.Name,
		Location:  trip.Location,
		StartDate: trip.StartDate.String(),
		EndDate:   trip.EndDate.String(),
	}

	response.Trip = &tripDto

	return response, nil
}

func (tc *TransactionController) DeleteTransactionEntry(ctx context.Context, transactionId *uuid.UUID) *models.ExpenseServiceError {
	// Begin transaction
	tx, err := tc.DatabaseMgr.BeginTx(ctx)
	if err != nil {
		log.Printf("Error while beginning transaction: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Make sure to rollback the transaction if it fails
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Printf("Error while rolling back transaction: %v", err)
		}
	}(tx)

	// Get user id from context
	userId := ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)

	// Get transaction from database
	transaction, repoErr := tc.TransactionRepo.GetTransactionById(transactionId)
	if repoErr != nil {
		return repoErr
	}

	// Check if user is creditor or debtor
	if transaction.CreditorId.String() != userId.String() && transaction.DebtorId.String() != userId.String() {
		return expense_errors.EXPENSE_UNAUTHORIZED
	}

	// Delete transaction
	if repoErr := tc.TransactionRepo.DeleteTx(tx, transactionId); repoErr != nil {
		return repoErr
	}

	// Delete debt
	if repoErr := tc.DebtRepo.CalculateDebt(tx, transaction.CreditorId, transaction.DebtorId, transaction.TripId, transaction.Amount.Neg()); repoErr != nil {
		return repoErr
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error while committing transaction: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}
