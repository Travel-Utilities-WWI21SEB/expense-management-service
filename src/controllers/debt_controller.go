package controllers

import (
	"context"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/repositories"
	"github.com/google/uuid"
)

// DebtCtl Exposed interface to the handler-package
type DebtCtl interface {
	GetDebtEntries(ctx context.Context, tripId *uuid.UUID) ([]*models.DebtDTO, *models.ExpenseServiceError)
	GetDebtDetails(ctx context.Context, debtId *uuid.UUID) (*models.DebtDTO, *models.ExpenseServiceError)
}

// DebtController Debt Controller structure
type DebtController struct {
	DatabaseMgr managers.DatabaseMgr
	DebtRepo    repositories.DebtRepo
	UserRepo    repositories.UserRepo
	TripRepo    repositories.TripRepo
}

// GetDebtEntries Get all debt entries for a trip
func (dc *DebtController) GetDebtEntries(ctx context.Context, tripId *uuid.UUID) ([]*models.DebtDTO, *models.ExpenseServiceError) {
	// Get all debt entries from database
	debtEntries, err := dc.DebtRepo.GetDebtEntriesByTripId(tripId)
	if err != nil {
		return nil, err
	}

	// Get Trip from database
	trip, err := dc.TripRepo.GetTripById(tripId)
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

		creditor, err := dc.UserRepo.GetUserById(debtEntry.CreditorId)
		if err != nil {
			return nil, err
		}
		creditorDto := &models.UserDto{
			UserID:   creditor.UserID,
			Username: creditor.Username,
			Email:    creditor.Email,
		}
		debtDto.Creditor = creditorDto

		debtor, err := dc.UserRepo.GetUserById(debtEntry.DebtorId)
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
	debtEntry, err := dc.DebtRepo.GetDebtById(debtId)
	if err != nil {
		return nil, err
	}

	// Get Trip from database
	trip, err := dc.TripRepo.GetTripById(debtEntry.TripId)
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

	creditor, err := dc.UserRepo.GetUserById(debtEntry.CreditorId)
	if err != nil {
		return nil, err
	}
	creditorDto := &models.UserDto{
		UserID:   creditor.UserID,
		Username: creditor.Username,
		Email:    creditor.Email,
	}
	debtDto.Creditor = creditorDto

	debtor, err := dc.UserRepo.GetUserById(debtEntry.DebtorId)
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
