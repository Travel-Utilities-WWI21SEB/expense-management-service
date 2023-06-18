package controllers

import (
	"context"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/repositories"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// CostCtl Exposed interface to the handler-package
type CostCtl interface {
	CreateCostEntry(ctx context.Context, createCostRequest models.CreateCostRequest) (*models.CostDetailsResponse, *models.ExpenseServiceError)
	GetCostDetails(ctx context.Context, costId *uuid.UUID) (*models.CostDetailsResponse, *models.ExpenseServiceError)
	// GetCostsByTrip(ctx context.Context, tripId *uuid.UUID) (*models.CostResponse, *models.ExpenseServiceError)
	// GetCostsByCostCategory(ctx context.Context, costCategoryId *uuid.UUID) (*models.CostResponse, *models.ExpenseServiceError)
	// GetCostsByContext(ctx context.Context) (*[]models.CostResponse, *models.ExpenseServiceError)
	PatchCostEntry(ctx context.Context) (*models.CostDetailsResponse, *models.ExpenseServiceError)
	PutCostEntry(ctx context.Context) (*models.CostDetailsResponse, *models.ExpenseServiceError)
	DeleteCostEntry(ctx context.Context) *models.ExpenseServiceError
}

// CostController Cost Controller structure
type CostController struct {
	DatabaseMgr managers.DatabaseMgr
	CostRepo    repositories.CostRepo
	UserRepo    repositories.UserRepo
	TripRepo    repositories.TripRepo
}

// CreateCostEntry Creates a cost entry and inserts it into the database
func (cc *CostController) CreateCostEntry(_ context.Context, createCostRequest models.CreateCostRequest) (*models.CostDetailsResponse, *models.ExpenseServiceError) {
	costId := uuid.New()
	now := time.Now()

	if createCostRequest.DeductedAt == nil {
		createCostRequest.DeductedAt = &now
	}

	// Create cost entry
	costEntry := &models.CostSchema{
		CostID:         &costId,
		Amount:         decimal.RequireFromString(createCostRequest.Amount),
		Description:    createCostRequest.Description,
		CreationDate:   &now,
		DeductionDate:  createCostRequest.DeductedAt,
		EndDate:        createCostRequest.EndDate,
		CostCategoryID: createCostRequest.CostCategoryID,
	}

	// Insert cost entry into database
	if repoErr := cc.CostRepo.CreateCost(costEntry); repoErr != nil {
		return nil, repoErr
	}

	contributors := make([]*models.Contributor, len(createCostRequest.Contributors))

	// Create cost contribution for contributors
	for i, contributor := range createCostRequest.Contributors {
		// Get user from database
		user, repoErr := cc.UserRepo.GetUserBySchema(&models.UserSchema{Username: contributor.Username})
		if repoErr != nil {
			return nil, repoErr
		}

		contribution := &models.CostContributionSchema{
			CostID:     &costId,
			UserID:     user.UserID,
			IsCreditor: contributor.Username == createCostRequest.Creditor,
			Amount:     decimal.RequireFromString(contributor.Amount),
		}

		// Insert cost contribution into database
		if repoErr := cc.CostRepo.AddCostContributor(contribution); repoErr != nil {
			return nil, repoErr
		}

		contributors[i] = &models.Contributor{Username: contributor.Username, Amount: contributor.Amount}
	}
	return cc.mapCostToResponse(costEntry, contributors), nil
}

func (cc *CostController) GetCostDetails(_ context.Context, costId *uuid.UUID) (*models.CostDetailsResponse, *models.ExpenseServiceError) {
	// Get cost entry from database
	cost, repoErr := cc.CostRepo.GetCostByID(costId)
	if repoErr != nil {
		return nil, repoErr
	}
	return cc.mapCostToResponse(cost, nil), nil
}

func (cc *CostController) GetCostsByTrip(ctx context.Context) (*[]models.CostDetailsResponse, *models.ExpenseServiceError) {
	// TO-DO
	return nil, nil
}

func (cc *CostController) PatchCostEntry(ctx context.Context) (*models.CostDetailsResponse, *models.ExpenseServiceError) {
	// TO-DO
	return nil, nil
}

func (cc *CostController) PutCostEntry(ctx context.Context) (*models.CostDetailsResponse, *models.ExpenseServiceError) {
	// TO-DO
	return nil, nil
}

func (cc *CostController) DeleteCostEntry(ctx context.Context) *models.ExpenseServiceError {
	// TO-DO
	return nil
}

// You can add optional parameters with: func (cc *CostController) GetCostDetails(ctx context.Context, costId *uuid.UUID, optionalParam string) (*models.CostDetailsResponse, *models.ExpenseServiceError) {
func (cc *CostController) mapCostToResponse(cost *models.CostSchema, contributors []*models.Contributor) *models.CostDetailsResponse {
	response := &models.CostDetailsResponse{
		CostID:         cost.CostID,
		Amount:         cost.Amount.String(),
		Description:    cost.Description,
		CreationDate:   cost.CreationDate,
		DeductionDate:  cost.DeductionDate,
		EndDate:        cost.EndDate,
		CostCategoryID: cost.CostCategoryID,
	}

	// If contributors are passed as parameter, use them
	if contributors != nil {
		response.Contributors = contributors
		return response
	}

	// Else get contributors from database
	contributions, _ := cc.CostRepo.GetCostContributors(cost.CostID)

	response.Contributors = make([]*models.Contributor, len(contributions))
	for i, contribution := range contributions {
		user, _ := cc.UserRepo.GetUserById(contribution.UserID)
		response.Contributors[i] = &models.Contributor{
			Username: user.Username,
			Amount:   contribution.Amount.String(),
		}

		if contribution.IsCreditor {
			response.Creditor = user.Username
		}
	}

	return response
}
