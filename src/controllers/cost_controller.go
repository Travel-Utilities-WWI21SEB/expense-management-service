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
func (cc *CostController) CreateCostEntry(ctx context.Context, createCostRequest models.CreateCostRequest) (*models.CostDetailsResponse, *models.ExpenseServiceError) {
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

	// If no contributors, set creator as only contributor
	if createCostRequest.Contributors == nil {
		creator, repoErr := cc.UserRepo.GetUserByContext(ctx)
		if repoErr != nil {
			return nil, repoErr
		}
		// If no contributors, set creator as contributor
		createCostRequest.Contributors = []*models.ContributorsRequest{{Username: creator.Username, IsCreditor: true}}
	}

	// Iterate over contributors and insert them into database
	for _, contributor := range createCostRequest.Contributors {
		user, repoErr := cc.UserRepo.GetUserBySchema(&models.UserSchema{Username: contributor.Username})
		if repoErr != nil {
			return nil, repoErr
		}

		contribution := &models.CostContributionSchema{
			CostID:     &costId,
			UserID:     user.UserID,
			IsCreditor: contributor.IsCreditor,
		}

		if repoErr := cc.CostRepo.AddCostContributor(contribution); repoErr != nil {
			// Delete cost entry
			delError := cc.CostRepo.DeleteCostEntry(&costId)
			if delError != nil {
				return nil, delError
			}
			return nil, repoErr
		}
	}

	return cc.mapCostToResponse(costEntry), nil
}

func (cc *CostController) GetCostDetails(ctx context.Context, costId *uuid.UUID) (*models.CostDetailsResponse, *models.ExpenseServiceError) {
	// Get cost entry from database
	cost, repoErr := cc.CostRepo.GetCostByID(costId)
	if repoErr != nil {
		return nil, repoErr
	}
	return cc.mapCostToResponse(cost), nil
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

func (cc *CostController) mapCostToResponse(cost *models.CostSchema) *models.CostDetailsResponse {
	contributions, _ := cc.CostRepo.GetCostContributors(cost.CostID)
	response := &models.CostDetailsResponse{
		CostID:         cost.CostID,
		Amount:         cost.Amount.String(),
		Description:    cost.Description,
		CreationDate:   cost.CreationDate,
		DeductionDate:  cost.DeductionDate,
		EndDate:        cost.EndDate,
		CostCategoryID: cost.CostCategoryID,
	}
	response.Contributors = *new([]*models.ContributorsResponse)

	for _, contribution := range contributions {
		user, _ := cc.UserRepo.GetUserById(contribution.UserID)
		response.Contributors = append(response.Contributors, &models.ContributorsResponse{
			Username: user.Username,
			// Divide amount by number of contributors (prevents rounding errors for money)
			Amount:     cost.Amount.Div(decimal.NewFromInt(int64(len(contributions)))).String(),
			IsCreditor: contribution.IsCreditor,
		})
	}

	return response
}
