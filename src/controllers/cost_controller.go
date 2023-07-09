package controllers

import (
	"context"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"log"
	"strconv"
	"time"
)

// CostCtl Exposed interface to the handler-package
type CostCtl interface {
	CreateCostEntry(ctx context.Context, tripId *uuid.UUID, createCostRequest models.CostDTO) (*models.CostDTO, *models.ExpenseServiceError)
	GetCostDetails(ctx context.Context, costId *uuid.UUID) (*models.CostDTO, *models.ExpenseServiceError)
	GetCostEntries(ctx context.Context, params *models.CostQueryParams) ([]*models.CostDTO, *models.ExpenseServiceError)
	PatchCostEntry(ctx context.Context, tripId *uuid.UUID, costId *uuid.UUID, request models.CostDTO) (*models.CostDTO, *models.ExpenseServiceError)
	DeleteCostEntry(ctx context.Context, tripId *uuid.UUID, costId *uuid.UUID) *models.ExpenseServiceError
	GetCostOverview(ctx context.Context) (*models.CostOverviewDTO, *models.ExpenseServiceError)
}

// CostController Cost Controller structure
type CostController struct {
	DatabaseMgr      managers.DatabaseMgr
	CostRepo         repositories.CostRepo
	UserRepo         repositories.UserRepo
	TripRepo         repositories.TripRepo
	CostCategoryRepo repositories.CostCategoryRepo
	DebtRepo         repositories.DebtRepo
}

func (cc *CostController) GetCostOverview(ctx context.Context) (*models.CostOverviewDTO, *models.ExpenseServiceError) {
	// Get user id from context
	userId, ok := ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	response, err := cc.CostRepo.GetCostOverview(ctx, userId)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateCostEntry Creates a cost entry and inserts it into the database
func (cc *CostController) CreateCostEntry(ctx context.Context, tripId *uuid.UUID, createCostRequest models.CostDTO) (*models.CostDTO, *models.ExpenseServiceError) {
	// Begin transaction
	tx, err := cc.DatabaseMgr.BeginTx(ctx)
	if err != nil {
		log.Printf("Error while beginning transaction: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Make sure to rollback the transaction if it fails
	defer func(tx pgx.Tx) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("Error while rolling back transaction: %v", err)
		}
	}(tx)

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
	if repoErr := cc.CostRepo.AddTx(ctx, tx, costEntry); repoErr != nil {
		return nil, repoErr
	}

	// Get creditor user from database
	creditorUser, repoErr := cc.UserRepo.GetUserBySchema(ctx, &models.UserSchema{Username: createCostRequest.Creditor})
	if repoErr != nil {
		return nil, repoErr
	}

	// Add creditor with zero amount to contributors if not already present
	var creditorFound bool
	for _, contributor := range createCostRequest.Debtors {
		if contributor.Username == creditorUser.Username {
			creditorFound = true
			break
		}
	}
	if !creditorFound {
		createCostRequest.Debtors = append(createCostRequest.Debtors, &models.Contributor{
			Username: creditorUser.Username,
			Amount:   "0.0",
		})
	}

	contributors := make([]*models.Contributor, len(createCostRequest.Debtors))
	for i, contributor := range createCostRequest.Debtors {
		contributorUser, repoErr := cc.UserRepo.GetUserBySchema(ctx, &models.UserSchema{Username: contributor.Username})
		if repoErr != nil {
			return nil, repoErr
		}

		// Check if user is part of the trip
		if repoErr = cc.TripRepo.ValidateIfUserHasAccepted(ctx, tripId, contributorUser.UserID); repoErr != nil {
			return nil, repoErr
		}

		contribution := &models.CostContributionSchema{
			CostID:     &costId,
			UserID:     contributorUser.UserID,
			IsCreditor: contributor.Username == creditorUser.Username,
			Amount:     decimal.RequireFromString(contributor.Amount),
		}

		// Insert cost contribution into database using transaction
		if repoErr = cc.CostRepo.AddCostContributorTx(ctx, tx, contribution); repoErr != nil {
			return nil, repoErr
		}

		// Calculate debt
		if contributorUser.Username != creditorUser.Username {
			if serviceErr := cc.DebtRepo.CalculateDebt(ctx, tx, creditorUser.UserID, contributorUser.UserID, tripId, contribution.Amount); serviceErr != nil {
				return nil, serviceErr
			}
		}

		contributors[i] = &models.Contributor{
			Username: contributor.Username,
			Amount:   contributor.Amount,
		}
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		log.Printf("Error while committing transaction: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return cc.mapCostToResponse(ctx, costEntry), nil
}

func (cc *CostController) GetCostDetails(ctx context.Context, costId *uuid.UUID) (*models.CostDTO, *models.ExpenseServiceError) {
	// Get cost entry from database
	cost, repoErr := cc.CostRepo.GetCostByID(ctx, costId)
	if repoErr != nil {
		return nil, repoErr
	}
	return cc.mapCostToResponse(ctx, cost), nil
}

func (cc *CostController) GetCostEntries(ctx context.Context, params *models.CostQueryParams) ([]*models.CostDTO, *models.ExpenseServiceError) {
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
		costCategory, repoErr := cc.CostCategoryRepo.GetCostCategoryByTripIdAndName(ctx, params.TripId, *params.CostCategoryName)
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
		user, repoErr := cc.UserRepo.GetUserBySchema(ctx, &models.UserSchema{Username: *params.Username})
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
		// Wunderschön Luca. :) - Kevin
	}

	// Add limit and offset for pagination
	if params.PageSize > 0 && params.Page > 0 {
		query += ` LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
		args = append(args, params.PageSize, (params.Page-1)*params.PageSize)
	}

	rows, err := cc.DatabaseMgr.ExecuteQuery(ctx, query, args...)
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
	costsResponse := make([]*models.CostDTO, len(costs))
	for i, cost := range costs {
		costsResponse[i] = cc.mapCostToResponse(ctx, cost)
	}

	return costsResponse, nil
}

func (cc *CostController) PatchCostEntry(ctx context.Context, tripId *uuid.UUID, costId *uuid.UUID, request models.CostDTO) (*models.CostDTO, *models.ExpenseServiceError) {
	// Begin transaction
	tx, err := cc.DatabaseMgr.BeginTx(ctx)
	if err != nil {
		log.Printf("Error while beginning transaction: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Make sure to rollback the transaction if it fails
	defer func(tx pgx.Tx) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("Error while rolling back transaction: %v", err)
		}
	}(tx)

	// Get cost entry from database
	cost, repoErr := cc.CostRepo.GetCostByID(ctx, costId)
	if repoErr != nil {
		return nil, repoErr
	}

	oldCreditorUser, repoErr := cc.CostRepo.GetCostCreditor(ctx, costId)
	if repoErr != nil {
		return nil, repoErr
	}

	amountChanged := request.Amount != "" && request.Amount != cost.Amount.String()
	creditorChanged := request.Creditor != "" && request.Creditor != oldCreditorUser.Username
	debtorsChanged := request.Debtors != nil && len(request.Debtors) > 0
	contributionsChanged := creditorChanged || debtorsChanged

	// Update cost entry if request contains new values
	if amountChanged {
		amount, err := ValidateAmount(request.Amount)
		if err != nil {
			return nil, err
		}
		cost.Amount = amount
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
		creditor, repoErr = cc.UserRepo.GetUserBySchema(ctx, &models.UserSchema{Username: request.Creditor})
		if repoErr != nil {
			return nil, repoErr
		}
		if repoErr = cc.TripRepo.ValidateIfUserHasAccepted(ctx, tripId, creditor.UserID); repoErr != nil {
			return nil, repoErr
		}
	} else {
		creditor = oldCreditorUser
	}
	request.Creditor = creditor.Username // Only a check, not needed for response

	// If only amount has changed, not the contributions, then get the contributions from the database
	if amountChanged && !contributionsChanged {
		contributions, repoErr := cc.CostRepo.GetCostContributors(ctx, costId)
		if repoErr != nil {
			return nil, repoErr
		}

		request.Debtors = make([]*models.Contributor, len(contributions))
		for i, contribution := range contributions {
			user, err := cc.UserRepo.GetUserById(ctx, contribution.UserID)
			if err != nil {
				return nil, err
			}

			request.Debtors[i] = &models.Contributor{
				Username: user.Username,
				Amount:   "", // Algorithm will calculate new amount
			}

			if contribution.Amount.IsZero() {
				request.Debtors[i].Amount = "0" // Stay zero if it was zero before
			}
		}
	}

	// Add creditor with zero amount to contributors if not already present
	var creditorFound bool
	for _, contributor := range request.Debtors {
		if contributor.Username == request.Creditor {
			creditorFound = true
			break
		}
	}
	if !creditorFound {
		request.Debtors = append(request.Debtors, &models.Contributor{
			Username: request.Creditor,
			Amount:   "0.0",
		})
	}

	// Distribute cost among contributors
	if serviceErr := DistributeCosts(&request); serviceErr != nil {
		return nil, serviceErr
	}

	// Delete old cost contributions and subtract debt from users
	contributors, repoErr := cc.CostRepo.GetCostContributors(ctx, costId)
	if repoErr != nil {
		return nil, repoErr
	}

	for _, contributor := range contributors {
		// Delete cost contribution from database
		if repoErr := cc.CostRepo.DeleteCostContributionTx(ctx, tx, contributor.UserID, costId); repoErr != nil {
			return nil, repoErr
		}

		// Subtract debt from user
		if repoErr := cc.DebtRepo.CalculateDebt(ctx, tx, creditor.UserID, contributor.UserID, tripId, contributor.Amount.Neg()); repoErr != nil {
			return nil, repoErr
		}
	}

	// Create cost contributions and add debt to users
	for _, contributor := range request.Debtors {
		// Get user from database
		user, repoErr := cc.UserRepo.GetUserBySchema(ctx, &models.UserSchema{Username: contributor.Username})
		if repoErr != nil {
			return nil, repoErr
		}

		contribution := &models.CostContributionSchema{
			CostID:     costId,
			UserID:     user.UserID,
			IsCreditor: contributor.Username == request.Creditor,
			Amount:     decimal.RequireFromString(contributor.Amount),
		}

		// Insert cost contribution into database
		if repoErr := cc.CostRepo.AddCostContributorTx(ctx, tx, contribution); repoErr != nil {
			return nil, repoErr
		}

		// Add debt to user
		if repoErr := cc.DebtRepo.CalculateDebt(ctx, tx, creditor.UserID, user.UserID, tripId, contribution.Amount); repoErr != nil {
			return nil, repoErr
		}
	}

	// Update cost entry in database
	if repoErr := cc.CostRepo.UpdateTx(ctx, tx, cost); repoErr != nil {
		return nil, repoErr
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		log.Printf("Error while committing transaction: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return cc.mapCostToResponse(ctx, cost), nil
}

func (cc *CostController) DeleteCostEntry(ctx context.Context, tripId *uuid.UUID, costId *uuid.UUID) *models.ExpenseServiceError {
	// Begin transaction
	tx, err := cc.DatabaseMgr.BeginTx(ctx)
	if err != nil {
		log.Printf("Error while beginning transaction: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Make sure to rollback the transaction if it fails
	defer func(tx pgx.Tx) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("Error while rolling back transaction: %v", err)
		}
	}(tx)

	// Get all cost contributions
	contributions, repoErr := cc.CostRepo.GetCostContributors(ctx, costId)
	if repoErr != nil {
		return repoErr
	}

	// Get creditor
	creditor, repoErr := cc.CostRepo.GetCostCreditor(ctx, costId)
	if repoErr != nil {
		return repoErr
	}

	// Remove debt from all contributors
	for _, contribution := range contributions {
		if contribution.IsCreditor {
			continue
		}
		// Subtract debt from user
		cc.DebtRepo.CalculateDebt(ctx, tx, creditor.UserID, contribution.UserID, tripId, contribution.Amount.Neg())
	}

	repoErr = cc.CostRepo.DeleteTx(ctx, tx, costId)
	if repoErr != nil {
		return repoErr
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		log.Printf("Error while committing transaction: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

// You can add optional parameters with: func (cc *CostController) GetCostDetails(ctx context.Context, costId *uuid.UUID, optionalParam string) (*models.CostDTO, *models.ExpenseServiceError) {
func (cc *CostController) mapCostToResponse(ctx context.Context, cost *models.CostSchema) *models.CostDTO {
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

	contributions, _ := cc.CostRepo.GetCostContributors(ctx, cost.CostID)

	response.Debtors = make([]*models.Contributor, len(contributions))
	for i, contribution := range contributions {
		user, _ := cc.UserRepo.GetUserById(ctx, contribution.UserID)
		response.Debtors[i] = &models.Contributor{
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

	sum, err := SumContributions(request.Debtors)
	if err != nil {
		return err
	}

	// Check if sum of contributions is greater than total cost
	if sum.GreaterThan(totalCost) {
		return expense_errors.EXPENSE_BAD_REQUEST
	}

	// Get number of contributors with no amount
	numContributorsWithNoAmount := 0
	for _, contributor := range request.Debtors {
		if contributor.Amount == "" {
			numContributorsWithNoAmount++
		}
	}

	// Distribute remaining costs to contributors with no amount
	if numContributorsWithNoAmount > 0 {
		remainingCost := totalCost.Sub(sum)
		// Sum of distributed amounts
		distributedAmount := DistributeRemainingCosts(request.Debtors, remainingCost, numContributorsWithNoAmount)

		// Add distributed amount to sum
		sum = sum.Add(distributedAmount)
	}

	// Check if the sum of contributions equals total cost
	if !sum.Equal(totalCost) {
		return expense_errors.EXPENSE_BAD_REQUEST
	}

	return nil
}

func DistributeRemainingCosts(contributors []*models.Contributor, remainingCost decimal.Decimal, numContributorsWithNoAmount int) decimal.Decimal {
	amountPerContributor := remainingCost.Div(decimal.NewFromInt(int64(numContributorsWithNoAmount)))

	// Round amountPerContributor to 2 decimal places
	amountPerContributor = amountPerContributor.RoundDown(2)

	// Sum of distributed amounts
	distributedAmount := decimal.Zero

	// Check before distributing, if rounding difference is greater than 0.00
	roundingDifference := remainingCost.Sub(amountPerContributor.Mul(decimal.NewFromInt(int64(numContributorsWithNoAmount)))) // roundingDifference = remainingCost - (amountPerContributor * numContributorsWithNoAmount)
	// Rounding difference can only be positive because of the way it is calculated (see above)

	// Distribute remaining cost to contributors with no amount
	for _, contributor := range contributors {
		if contributor.Amount == "" {
			realAmountPerContributor := amountPerContributor
			contributor.Amount = realAmountPerContributor.String()
			distributedAmount = distributedAmount.Add(realAmountPerContributor)
		}
	}

	// Write a while loop to distribute the rounding difference to the contributors with no amount
	var i int

	// "I'm a while loop
	// and I'm here to say
	// I'm gonna distribute
	// the rounding difference
	// in a very special way"
	for roundingDifference.GreaterThan(decimal.Zero) {
		// Add a cent to the contributor in the list at index i
		contributors[i].Amount = decimal.RequireFromString(contributors[i].Amount).Add(decimal.NewFromFloat(0.01)).String()

		// Add a cent to the distributed amount
		distributedAmount = distributedAmount.Add(decimal.NewFromFloat(0.01))

		// Subtract a cent from the rounding difference
		roundingDifference = roundingDifference.Sub(decimal.NewFromFloat(0.01))

		// Increment i
		i++

		if i >= len(contributors) {
			i = 0
		}
	}
	return distributedAmount
}
