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
	"strconv"
	"time"
)

// CostCtl Exposed interface to the handler-package
type CostCtl interface {
	CreateCostEntry(ctx context.Context, tripId *uuid.UUID, createCostRequest models.CostDTO) (*models.CostDTO, *models.ExpenseServiceError)
	GetCostDetails(ctx context.Context, costId *uuid.UUID) (*models.CostDTO, *models.ExpenseServiceError)
	GetCostEntries(ctx context.Context, params *models.CostQueryParams) (*[]models.CostDTO, *models.ExpenseServiceError)
	PatchCostEntry(ctx context.Context, tripId *uuid.UUID, costId *uuid.UUID, request models.CostDTO) (*models.CostDTO, *models.ExpenseServiceError)
	DeleteCostEntry(ctx context.Context, costId *uuid.UUID) *models.ExpenseServiceError
}

// CostController Cost Controller structure
type CostController struct {
	DatabaseMgr      managers.DatabaseMgr
	CostRepo         repositories.CostRepo
	UserRepo         repositories.UserRepo
	TripRepo         repositories.TripRepo
	CostCategoryRepo repositories.CostCategoryRepo
}

// CreateCostEntry Creates a cost entry and inserts it into the database
func (cc *CostController) CreateCostEntry(_ context.Context, tripId *uuid.UUID, createCostRequest models.CostDTO) (*models.CostDTO, *models.ExpenseServiceError) {
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

	var deleteCostError *models.ExpenseServiceError

	// Delete trip from database if user is not added to trip
	defer func() {
		if deleteCostError != nil {
			cc.CostRepo.DeleteCostEntry(&costId)
		}
	}()

	contributors := make([]*models.Contributor, len(createCostRequest.Contributors))
	var creditorIsPartOfContributors bool
	// Create cost contribution for contributors
	for i, contributor := range createCostRequest.Contributors {
		// Get user from database
		user, repoErr := cc.UserRepo.GetUserBySchema(&models.UserSchema{Username: contributor.Username})
		if repoErr != nil {
			deleteCostError = repoErr
			return nil, repoErr
		}

		// Check if creditor is part of contributors
		if contributor.Username == createCostRequest.Creditor {
			creditorIsPartOfContributors = true
		}

		// Check if user is part of the trip
		repoErr = cc.TripRepo.ValidateIfUserHasAccepted(tripId, user.UserID)
		if repoErr != nil {
			deleteCostError = repoErr
			return nil, repoErr
		}

		contribution := &models.CostContributionSchema{
			CostID:     &costId,
			UserID:     user.UserID,
			IsCreditor: contributor.Username == createCostRequest.Creditor,
			Amount:     decimal.RequireFromString(contributor.Amount),
		}

		// Insert cost contribution into database
		if repoErr = cc.CostRepo.AddCostContributor(contribution); repoErr != nil {
			deleteCostError = repoErr
			return nil, repoErr
		}

		contributors[i] = &models.Contributor{Username: contributor.Username, Amount: contributor.Amount}
	}

	// Check if creditor is part of contributors
	if !creditorIsPartOfContributors {
		deleteCostError = expense_errors.EXPENSE_BAD_REQUEST
		return nil, deleteCostError
	}

	if deleteCostError != nil {
		return nil, deleteCostError
	}

	return cc.mapCostToResponse(costEntry), nil
}

func (cc *CostController) GetCostDetails(_ context.Context, costId *uuid.UUID) (*models.CostDTO, *models.ExpenseServiceError) {
	// Get cost entry from database
	cost, repoErr := cc.CostRepo.GetCostByID(costId)
	if repoErr != nil {
		return nil, repoErr
	}
	return cc.mapCostToResponse(cost), nil
}

func (cc *CostController) GetCostEntries(_ context.Context, params *models.CostQueryParams) (*[]models.CostDTO, *models.ExpenseServiceError) {
	// Mandatory parameters: tripId
	// Optional parameters: costCategoryId, username
	var costs []*models.CostSchema

	var args []interface{}
	query := `SELECT DISTINCT c.* FROM cost c INNER JOIN cost_category cc on c.id_cost_category = cc.id INNER JOIN user_cost_association uca on c.id = uca.id_cost WHERE id_trip = $1`
	args = append(args, params.TripId)

	if params.CostCategoryId != nil {
		query += ` AND c.id_cost_category = $` + strconv.Itoa(len(args)+1) // returns ' AND id_cost_category = $2'
		args = append(args, params.CostCategoryId)                         // returns ' AND id_cost_category = $2'
	}

	if params.CostCategoryName != nil {
		costCategory, repoErr := cc.CostCategoryRepo.GetCostCategoryByTripIdAndName(params.TripId, *params.CostCategoryName)
		if repoErr != nil {
			return nil, repoErr
		}

		if costCategory == nil {
			return nil, expense_errors.EXPENSE_BAD_REQUEST
		}

		query += ` AND c.id_cost_category = $` + strconv.Itoa(len(args)+1)
		args = append(args, costCategory.CostCategoryID)
	}

	if params.UserId != nil {
		query += ` AND uca.id_user = $` + strconv.Itoa(len(args)+1)
		args = append(args, params.UserId)
	}

	if params.Username != nil {
		user, repoErr := cc.UserRepo.GetUserBySchema(&models.UserSchema{Username: *params.Username})
		if repoErr != nil {
			return nil, repoErr
		}

		if user == nil {
			return nil, expense_errors.EXPENSE_BAD_REQUEST
		}

		query += ` AND uca.id_user = $` + strconv.Itoa(len(args)+1)
		args = append(args, user.UserID)
	}

	if params.MinAmount != nil {
		query += ` AND c.amount >= $` + strconv.Itoa(len(args)+1)
		args = append(args, params.MinAmount)
	}

	if params.MaxAmount != nil {
		query += ` AND c.amount <= $` + strconv.Itoa(len(args)+1)
		args = append(args, params.MaxAmount)
	}

	if params.MinDeductionDate != nil {
		query += ` AND c.deducted_at >= $` + strconv.Itoa(len(args)+1)
		args = append(args, params.MinDeductionDate)
	}

	if params.MaxDeductionDate != nil {
		query += ` AND c.deducted_at <= $` + strconv.Itoa(len(args)+1)
		args = append(args, params.MaxDeductionDate)
	}

	if params.MinCreationDate != nil {
		query += ` AND c.created_at >= $` + strconv.Itoa(len(args)+1)
		args = append(args, params.MinCreationDate)
	}

	if params.MaxCreationDate != nil {
		query += ` AND c.created_at <= $` + strconv.Itoa(len(args)+1)
		args = append(args, params.MaxCreationDate)
	}

	if params.MinEndDate != nil {
		query += ` AND c.end_date >= $` + strconv.Itoa(len(args)+1)
		args = append(args, params.MinEndDate)
	}

	if params.MaxEndDate != nil {
		query += ` AND c.end_date <= $` + strconv.Itoa(len(args)+1)
		args = append(args, params.MaxEndDate)
	}

	// Add order by clause
	if params.SortBy != "" {
		query += ` ORDER BY c.` + params.SortBy + ` ` + params.SortOrder
		/*
			Roses are red,
			Code should be neat,
			If not sanitized inputs,
			Then "DROP TABLE" you'll meet.
		*/
	}

	// Add limit and offset for pagination
	if params.PageSize > 0 && params.Page > 0 {
		query += ` LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
		args = append(args, params.PageSize, (params.Page-1)*params.PageSize)
	}

	log.Printf("Query: %v", query)
	rows, err := cc.DatabaseMgr.ExecuteQuery(query, args...)
	if err != nil {
		log.Printf("Error while executing query: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	for rows.Next() {
		var cost models.CostSchema
		err := rows.Scan(&cost.CostID, &cost.Amount, &cost.Description, &cost.CreationDate, &cost.DeductionDate, &cost.EndDate, &cost.CostCategoryID)
		if err != nil {
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}
		costs = append(costs, &cost)
	}

	// Map costs to response
	costsResponse := make([]models.CostDTO, len(costs))
	for i, cost := range costs {
		costsResponse[i] = *cc.mapCostToResponse(cost)
	}

	return &costsResponse, nil
}

func (cc *CostController) PatchCostEntry(_ context.Context, tripId *uuid.UUID, costId *uuid.UUID, request models.CostDTO) (*models.CostDTO, *models.ExpenseServiceError) {
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
		creditor, repoErr = cc.UserRepo.GetUserBySchema(&models.UserSchema{Username: request.Creditor})
		if repoErr != nil {
			return nil, repoErr
		}
		if repoErr = cc.TripRepo.ValidateIfUserHasAccepted(tripId, creditor.UserID); repoErr != nil {
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
	}

	// Check if creditor is in contributors
	if !checkIfCreditorIsContributor(request.Creditor, request.Contributors) {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	// Distribute cost among contributors
	if serviceErr := DistributeCosts(&request); serviceErr != nil {
		return nil, serviceErr
	}

	// Validate contributor input
	var creditorAvailable bool
	for _, contributor := range request.Contributors {
		if contributor.Username == request.Creditor {
			creditorAvailable = true
		}

		// Get user from database
		tempCreditor, repoErr := cc.UserRepo.GetUserBySchema(&models.UserSchema{Username: contributor.Username})
		if repoErr != nil {
			return nil, repoErr
		}

		if repoErr = cc.TripRepo.ValidateIfUserHasAccepted(tripId, tempCreditor.UserID); repoErr != nil {
			return nil, repoErr
		}
	}

	if !creditorAvailable {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	// Delete cost contributions
	if repoErr := cc.CostRepo.DeleteCostContributions(costId); repoErr != nil {
		return nil, repoErr
	}

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

	// Update cost entry in database
	if repoErr := cc.CostRepo.UpdateCost(cost); repoErr != nil {
		return nil, repoErr
	}

	return cc.mapCostToResponse(cost), nil
}

func (cc *CostController) DeleteCostEntry(_ context.Context, costId *uuid.UUID) *models.ExpenseServiceError {
	return cc.CostRepo.DeleteCostEntry(costId)
}

//**********************************************************************************************************************
// Helper functions
//**********************************************************************************************************************

// You can add optional parameters with: func (cc *CostController) GetCostDetails(ctx context.Context, costId *uuid.UUID, optionalParam string) (*models.CostDTO, *models.ExpenseServiceError) {
func (cc *CostController) mapCostToResponse(cost *models.CostSchema) *models.CostDTO {
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
