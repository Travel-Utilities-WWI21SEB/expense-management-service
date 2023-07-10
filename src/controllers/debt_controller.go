package controllers

import (
	"context"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/repositories"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"log"
)

// DebtCtl Exposed interface to the handler-package
type DebtCtl interface {
	GetDebtOverview(ctx context.Context, userId *uuid.UUID) (*models.DebtOverviewDTO, *models.ExpenseServiceError)
	GetDebtEntries(ctx context.Context, tripId *uuid.UUID) ([]*models.DebtDTO, *models.ExpenseServiceError)
	GetDebtDetails(ctx context.Context, debtId *uuid.UUID) (*models.DebtDTO, *models.ExpenseServiceError)
}

// DebtController Debt Controller structure
type DebtController struct {
	DatabaseMgr     managers.DatabaseMgr
	DebtRepo        repositories.DebtRepo
	UserRepo        repositories.UserRepo
	TransactionRepo repositories.TransactionRepo
	TripRepo        repositories.TripRepo
}

func (dc *DebtController) GetDebtOverview(ctx context.Context, userId *uuid.UUID) (*models.DebtOverviewDTO, *models.ExpenseServiceError) {
	debtEntries, err := dc.DebtRepo.GetDebtEntries(ctx, userId)
	if err != nil {
		log.Printf("Error while getting debt entries: %v", err)
		return nil, err
	}

	// Get all transactions from database
	transactions, err := dc.TransactionRepo.GetAllTransactions(ctx, userId)
	if err != nil {
		log.Printf("Error while getting transactions: %v", err)
		return nil, err
	}

	// Sum up all transactions, to get the total amount of money spent and received
	var totalAmountSpent decimal.Decimal
	var totalAmountReceived decimal.Decimal
	for _, transaction := range transactions {
		if transaction.DebtorId == userId {
			totalAmountReceived = totalAmountReceived.Add(transaction.Amount)
		} else {
			totalAmountSpent = totalAmountSpent.Add(transaction.Amount)
		}
	}

	// Sum up all debt entries, to get the total amount of money owed and owing
	var totalAmountOwed decimal.Decimal
	var totalAmountOwing decimal.Decimal
	for _, debtEntry := range debtEntries {
		parsedAmount, err := decimal.NewFromString(debtEntry.Amount)
		if err != nil {
			log.Printf("Error while parsing amount: %v", err)
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}

		if parsedAmount.GreaterThanOrEqual(decimal.Zero) {
			totalAmountOwed = totalAmountOwed.Add(parsedAmount)
		} else {
			totalAmountOwing = totalAmountOwing.Add(parsedAmount)
		}
	}

	return &models.DebtOverviewDTO{
		Debts:            debtEntries,
		OpenDebtAmount:   totalAmountOwing.String(),
		OpenCreditAmount: totalAmountOwed.String(),
		TotalSpent:       totalAmountSpent.String(),
		TotalReceived:    totalAmountReceived.String(),
	}, nil
}

// GetDebtEntries Get all debt entries for a trip
func (dc *DebtController) GetDebtEntries(ctx context.Context, tripId *uuid.UUID) ([]*models.DebtDTO, *models.ExpenseServiceError) {
	// Get all debt entries from database
	debtEntries, err := dc.DebtRepo.GetDebtEntriesByTripId(ctx, tripId)
	if err != nil {
		return nil, err
	}

	// Get Trip from database
	trip, err := dc.TripRepo.GetTripById(ctx, tripId)
	if err != nil {
		return nil, err
	}

	slimTripDto := &models.SlimTripDTO{
		TripID:    trip.TripID,
		Name:      trip.Name,
		Location:  trip.Location,
		StartDate: trip.StartDate.String(),
		EndDate:   trip.EndDate.String(),
	}

	// Convert debt entries to DTOs
	var debtDTOs []*models.DebtDTO
	for _, debtEntry := range debtEntries {
		debtDto := &models.DebtDTO{
			DebtID:       debtEntry.DebtID,
			Trip:         slimTripDto,
			Amount:       debtEntry.Amount.String(),
			CurrencyCode: debtEntry.CurrencyCode,
			CreationDate: debtEntry.CreationDate.String(),
			UpdateDate:   debtEntry.UpdateDate.String(),
		}

		creditor, err := dc.UserRepo.GetUserById(ctx, debtEntry.CreditorId)
		if err != nil {
			return nil, err
		}
		creditorDto := &models.UserDto{
			UserID:   creditor.UserID,
			Username: creditor.Username,
			Email:    creditor.Email,
		}
		debtDto.Creditor = creditorDto

		debtor, err := dc.UserRepo.GetUserById(ctx, debtEntry.DebtorId)
		if err != nil {
			return nil, err
		}
		debtorDto := &models.UserDto{
			UserID:   debtor.UserID,
			Username: debtor.Username,
			Email:    debtor.Email,
		}
		debtDto.Debtor = debtorDto

		debtDTOs = append(debtDTOs, debtDto)
	}

	return debtDTOs, nil
}

// GetDebtDetails Get details of a debt entry
func (dc *DebtController) GetDebtDetails(ctx context.Context, debtId *uuid.UUID) (*models.DebtDTO, *models.ExpenseServiceError) {
	// Get debt entry from database
	debtEntry, err := dc.DebtRepo.GetDebtById(ctx, debtId)
	if err != nil {
		return nil, err
	}

	// Get Trip from database
	trip, err := dc.TripRepo.GetTripById(ctx, debtEntry.TripId)
	if err != nil {
		return nil, err
	}

	slimTripDto := &models.SlimTripDTO{
		TripID:    trip.TripID,
		Name:      trip.Name,
		Location:  trip.Location,
		StartDate: trip.StartDate.String(),
		EndDate:   trip.EndDate.String(),
	}

	// Convert debt entry to DTO
	debtDto := &models.DebtDTO{
		DebtID:       debtEntry.DebtID,
		Trip:         slimTripDto,
		Amount:       debtEntry.Amount.String(),
		CurrencyCode: debtEntry.CurrencyCode,
		CreationDate: debtEntry.CreationDate.String(),
		UpdateDate:   debtEntry.UpdateDate.String(),
	}

	creditor, err := dc.UserRepo.GetUserById(ctx, debtEntry.CreditorId)
	if err != nil {
		return nil, err
	}
	creditorDto := &models.UserDto{
		UserID:   creditor.UserID,
		Username: creditor.Username,
		Email:    creditor.Email,
	}
	debtDto.Creditor = creditorDto

	debtor, err := dc.UserRepo.GetUserById(ctx, debtEntry.DebtorId)
	if err != nil {
		return nil, err
	}
	debtorDto := &models.UserDto{
		UserID:   debtor.UserID,
		Username: debtor.Username,
		Email:    debtor.Email,
	}
	debtDto.Debtor = debtorDto

	return debtDto, nil
}
