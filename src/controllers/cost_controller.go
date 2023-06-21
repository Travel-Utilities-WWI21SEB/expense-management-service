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
	"time"
)

// CostCtl Exposed interface to the handler-package
type CostCtl interface {
	CreateCostEntry(ctx context.Context, createCostRequest models.CostDTO) (*models.CostDTO, *models.ExpenseServiceError)
	GetCostDetails(ctx context.Context, costId *uuid.UUID) (*models.CostDTO, *models.ExpenseServiceError)
	// GetCostEntriesByTrip(ctx context.Context, tripId *uuid.UUID) (*models.CostDTO, *models.ExpenseServiceError)
	// GetCostEntriesByCostCategory(ctx context.Context, costCategoryId *uuid.UUID) (*models.CostDTO, *models.ExpenseServiceError)
	// GetCostEntriesByContext(ctx context.Context) (*[]models.CostDTO, *models.ExpenseServiceError)
	PatchCostEntry(ctx context.Context, costId *uuid.UUID, request models.CostDTO) (*models.CostDTO, *models.ExpenseServiceError)
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
func (cc *CostController) CreateCostEntry(_ context.Context, createCostRequest models.CostDTO) (*models.CostDTO, *models.ExpenseServiceError) {
	costId := uuid.New()
	now := time.Now()

	deductionDate := now

	if createCostRequest.DeductionDate != "" {
		deductionDate, _ = time.Parse(time.RFC3339, createCostRequest.DeductionDate)
	}

	// Distribute cost among contributors
	if serviceErr := DistributeCosts(&createCostRequest); serviceErr != nil {
		return nil, serviceErr
	}

	// Create cost entry
	costEntry := &models.CostSchema{
		CostID:         &costId,
		Amount:         decimal.RequireFromString(createCostRequest.Amount),
		Description:    createCostRequest.Description,
		CreationDate:   &now,
		DeductionDate:  &deductionDate,
		CostCategoryID: createCostRequest.CostCategoryID,
	}

	if createCostRequest.EndDate != "" {
		endDate, _ := time.Parse(time.RFC3339, createCostRequest.EndDate)
		costEntry.EndDate = &endDate
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

func (cc *CostController) GetCostDetails(_ context.Context, costId *uuid.UUID) (*models.CostDTO, *models.ExpenseServiceError) {
	// Get cost entry from database
	cost, repoErr := cc.CostRepo.GetCostByID(costId)
	if repoErr != nil {
		return nil, repoErr
	}
	return cc.mapCostToResponse(cost, nil), nil
}

func (cc *CostController) GetCostEntriesByTrip(ctx context.Context, tripId *uuid.UUID) (*[]models.CostDTO, *models.ExpenseServiceError) {
	// TO-DO
	return nil, nil
}

func (cc *CostController) PatchCostEntry(ctx context.Context, costId *uuid.UUID, request models.CostDTO) (*models.CostDTO, *models.ExpenseServiceError) {
	// Get cost entry from database
	cost, repoErr := cc.CostRepo.GetCostByID(costId)
	if repoErr != nil {
		return nil, repoErr
	}

	amountChanged := request.Amount != ""
	creditorChanged := request.Creditor != ""
	contributorsChanged := request.Contributors != nil && len(request.Contributors) > 0

	// Update cost entry if request contains new values
	if amountChanged {
		amount, err := ValidateAmount(request.Amount)
		if err != nil {
			return nil, err
		}
		cost.Amount = amount
	} else {
		request.Amount = cost.Amount.String()
	}

	if request.Description != "" {
		cost.Description = request.Description
	}

	if request.DeductionDate != "" {
		deductionDate, _ := time.Parse(time.RFC3339, request.DeductionDate)
		cost.DeductionDate = &deductionDate
	}

	if request.EndDate != "" {
		endDate, _ := time.Parse(time.RFC3339, request.EndDate)
		cost.EndDate = &endDate
	}

	if request.CostCategoryID != nil {
		cost.CostCategoryID = request.CostCategoryID
	}

	var creditor *models.UserSchema
	if creditorChanged {
		if creditor, repoErr = cc.UserRepo.GetUserBySchema(&models.UserSchema{Username: request.Creditor}); repoErr != nil {
			return nil, repoErr
		}
	} else {
		if creditor, repoErr = cc.CostRepo.GetCostCreditor(costId); repoErr != nil {
			return nil, repoErr
		}
	}
	request.Creditor = creditor.Username // Only a check, not needed for response

	// If contributors have changed, update cost contributions and distribute cost among contributors
	// Rule: You can only replace all contributors at once
	if amountChanged || contributorsChanged {
		// Copy contributors from request if they have not changed
		if !contributorsChanged {
			contributions, repoErr := cc.CostRepo.GetCostContributors(costId)
			if repoErr != nil {
				return nil, repoErr
			}

			request.Contributors = make([]*models.Contributor, len(contributions))

			for i, contribution := range contributions {
				user, err := cc.UserRepo.GetUserById(contribution.UserID)
				if err != nil {
					return nil, err
				}
				request.Contributors[i] = &models.Contributor{
					Username: user.Username,
					Amount:   "", // Algorithm will calculate new amount
				}
			}
		}

		// Check if creditor is in contributors
		if !checkIfCreditorIsContributor(request.Creditor, request.Contributors) {
			return nil, expense_errors.EXPENSE_BAD_REQUEST
		}

		// Distribute cost among contributors
		if serviceErr := DistributeCosts(&request); serviceErr != nil {
			return nil, serviceErr
		}

		// Delete cost contributions
		if repoErr := cc.CostRepo.DeleteCostContributions(costId); repoErr != nil {
			return nil, repoErr
		}

		var creditorAvailable bool

		// Create cost contributions
		for _, contributor := range request.Contributors {
			// Get user from database
			user, repoErr := cc.UserRepo.GetUserBySchema(&models.UserSchema{Username: contributor.Username})
			if repoErr != nil {
				return nil, repoErr
			}

			if creditor.Username == user.Username {
				creditorAvailable = true
			}

			contribution := &models.CostContributionSchema{
				CostID:     costId,
				UserID:     user.UserID,
				IsCreditor: contributor.Username == request.Creditor,
				Amount:     decimal.RequireFromString(contributor.Amount),
			}

			// Insert cost contribution into database
			if repoErr := cc.CostRepo.AddCostContributor(contribution); repoErr != nil {
				return nil, repoErr
			}
		}

		if !creditorAvailable {
			return nil, expense_errors.EXPENSE_BAD_REQUEST
		}
	}

	// Update cost entry in database
	if repoErr := cc.CostRepo.UpdateCost(cost); repoErr != nil {
		return nil, repoErr
	}

	return cc.mapCostToResponse(cost, request.Contributors), nil
}

func (cc *CostController) DeleteCostEntry(ctx context.Context) *models.ExpenseServiceError {
	// TO-DO
	return nil
}

// You can add optional parameters with: func (cc *CostController) GetCostDetails(ctx context.Context, costId *uuid.UUID, optionalParam string) (*models.CostDTO, *models.ExpenseServiceError) {
func (cc *CostController) mapCostToResponse(cost *models.CostSchema, contributors []*models.Contributor) *models.CostDTO {
	response := &models.CostDTO{
		CostID:         cost.CostID,
		Amount:         cost.Amount.String(),
		Description:    cost.Description,
		CreationDate:   cost.CreationDate.String(),
		DeductionDate:  cost.DeductionDate.String(),
		CostCategoryID: cost.CostCategoryID,
	}

	if cost.EndDate != nil {
		response.EndDate = cost.EndDate.String()
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

// ValidateAmount validates the amount of a cost entry, converts it to decimal and rounds it down to 2 decimal places
func ValidateAmount(amount string) (decimal.Decimal, *models.ExpenseServiceError) {
	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return decimal.Zero, expense_errors.EXPENSE_BAD_REQUEST
	}
	if amountDecimal.IsNegative() {
		return decimal.Zero, expense_errors.EXPENSE_BAD_REQUEST
	}
	if amountDecimal.IsZero() {
		return decimal.Zero, expense_errors.EXPENSE_BAD_REQUEST
	}
	return amountDecimal.RoundDown(2), nil
}

// SumContributions sums up the contributions of all debtors
func SumContributions(contributors []*models.Contributor) (decimal.Decimal, *models.ExpenseServiceError) {
	sum := decimal.Zero

	// Sum of debtor contributions
	for _, contributor := range contributors {
		if contributor.Amount != "" {
			amount, err := ValidateAmount(contributor.Amount)
			if err != nil {
				return decimal.Zero, err
			}
			sum = sum.Add(amount)
		}
	}
	return sum, nil
}

func DistributeCosts(request *models.CostDTO) *models.ExpenseServiceError {
	// Validate total cost amount
	totalCost, err := ValidateAmount(request.Amount)
	if err != nil {
		return err
	}
	request.Amount = totalCost.String()

	sum, err := SumContributions(request.Contributors)
	if err != nil {
		return err
	}

	// Check if sum of contributions is greater than total cost
	if sum.GreaterThan(totalCost) {
		return expense_errors.EXPENSE_BAD_REQUEST
	}

	// Get number of contributors with no amount
	numContributorsWithNoAmount := 0
	for _, contributor := range request.Contributors {
		if contributor.Amount == "" {
			numContributorsWithNoAmount++
		}
	}

	// Distribute remaining costs to contributors with no amount
	if numContributorsWithNoAmount > 0 {
		remainingCost := totalCost.Sub(sum)
		// Sum of distributed amounts
		distributedAmount := DistributeRemainingCosts(request.Contributors, remainingCost, numContributorsWithNoAmount, request.Creditor)

		// Add distributed amount to sum
		sum = sum.Add(distributedAmount)
	}

	// Check if the sum of contributions equals total cost
	if !sum.Equal(totalCost) {
		log.Printf("Sum of contributions (%v) does not equal total cost (%v)", sum, totalCost)
		return expense_errors.EXPENSE_BAD_REQUEST
	}

	return nil
}

func DistributeRemainingCosts(contributors []*models.Contributor, remainingCost decimal.Decimal, numContributorsWithNoAmount int, creditor string) decimal.Decimal {
	amountPerContributor := remainingCost.Div(decimal.NewFromInt(int64(numContributorsWithNoAmount)))

	// Round amountPerContributor to 2 decimal places
	amountPerContributor = amountPerContributor.RoundDown(2)

	// Sum of distributed amounts
	distributedAmount := decimal.Zero

	// Check before distributing, if rounding difference is greater than 0.00
	roundingDifference := remainingCost.Sub(amountPerContributor.Mul(decimal.NewFromInt(int64(numContributorsWithNoAmount)))) // roundingDifference = remainingCost - (amountPerContributor * numContributorsWithNoAmount)
	log.Printf("Rounding difference: %v", roundingDifference)
	// Rounding difference can only be positive because of the way it is calculated (see above)

	// Distribute remaining cost to contributors with no amount
	for _, contributor := range contributors {
		if contributor.Amount == "" {
			realAmountPerContributor := amountPerContributor

			// Add rounding difference to creditor
			if contributor.Username == creditor {
				realAmountPerContributor = realAmountPerContributor.Add(roundingDifference)
			}
			log.Printf("Distributing %v to %v", realAmountPerContributor, contributor.Username)

			contributor.Amount = realAmountPerContributor.String()
			distributedAmount = distributedAmount.Add(realAmountPerContributor)
		}
	}

	return distributedAmount
}

func checkIfCreditorIsContributor(creditor string, contributors []*models.Contributor) bool {
	for _, contributor := range contributors {
		if contributor.Username == creditor {
			return true
		}
	}
	return false
}
