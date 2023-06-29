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
	response := &models.CostOverviewDTO{
		TotalCosts:       "1321,23",
		AverageTripCosts: "123,23",
		MostExpensiveTrip: &models.TripNameToIdDTO{
			TripName: "Trip to Berlin",
			Amount:   "18321,23",
			TripId:   uuid.New(),
		},
		LeastExpensiveTrip: &models.TripNameToIdDTO{
			TripName: "Trip to Palo Alto",
			Amount:   "923,23",
			TripId:   uuid.New(),
		},
		AverageContributionPercentage: "23,23",
		TripDistribution: []*models.TripDistributionDTO{
			{
				TripName: "Trip to Berlin",
				Amount:   "18321,23",
			},
			{
				TripName: "Trip to Palo Alto",
				Amount:   "923,23",
			},
			{
				TripName: "Trip to San Francisco",
				Amount:   "1234,23",
			},
			{
				TripName: "Trip to New York",
				Amount:   "12312,23",
			},
		},
		CostDistribution: []*models.CostDistributionDTO{
			{
				CostCategoryName: "Food",
				Amount:           "2123,23",
			},
			{
				CostCategoryName: "Accommodation",
				Amount:           "142,23",
			},
			{
				CostCategoryName: "Transportation",
				Amount:           "2312,23",
			},
		},
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
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
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
	if repoErr := cc.CostRepo.AddTx(tx, costEntry); repoErr != nil {
		return nil, repoErr
	}

	// Get creditor user from database
	creditorUser, repoErr := cc.UserRepo.GetUserBySchema(&models.UserSchema{Username: createCostRequest.Creditor})
	if repoErr != nil {
		return nil, repoErr
	}

	// Create cost contribution for contributors
	contributors := make([]*models.Contributor, len(createCostRequest.Contributors))
	var creditorIsPartOfContributors bool
	for i, contributor := range createCostRequest.Contributors {
		contributorUser, repoErr := cc.UserRepo.GetUserBySchema(&models.UserSchema{Username: contributor.Username})
		if repoErr != nil {
			return nil, repoErr
		}

		// Check if creditor is part of contributors
		if contributor.Username == createCostRequest.Creditor {
			creditorIsPartOfContributors = true
		}

		// Check if user is part of the trip
		if repoErr = cc.TripRepo.ValidateIfUserHasAccepted(tripId, contributorUser.UserID); repoErr != nil {
			return nil, repoErr
		}

		contribution := &models.CostContributionSchema{
			CostID:     &costId,
			UserID:     contributorUser.UserID,
			IsCreditor: contributor.Username == creditorUser.Username,
			Amount:     decimal.RequireFromString(contributor.Amount),
		}

		// Insert cost contribution into database using transaction
		if repoErr = cc.CostRepo.AddCostContributorTx(tx, contribution); repoErr != nil {
			return nil, repoErr
		}

		// Calculate debt
		if contributorUser.Username != creditorUser.Username {
			log.Printf("Calculating debt for user %v and creditor %v", contributorUser.Username, creditorUser.Username)
			if serviceErr := cc.calculateDebt(tx, creditorUser, contributorUser, tripId, contribution.Amount); serviceErr != nil {
				return nil, serviceErr
			}
		}

		contributors[i] = &models.Contributor{
			Username: contributor.Username,
			Amount:   contributor.Amount,
		}
	}

	// Check if creditor is part of contributors
	if !creditorIsPartOfContributors {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error while committing transaction: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
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

func (cc *CostController) GetCostEntries(_ context.Context, params *models.CostQueryParams) ([]*models.CostDTO, *models.ExpenseServiceError) {
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
		log.Printf("Cost: %v", cost.DeductionDate)
		costs = append(costs, &cost)
	}

	// Map costs to response
	costsResponse := make([]*models.CostDTO, len(costs))
	for i, cost := range costs {
		costsResponse[i] = cc.mapCostToResponse(cost)
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
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Printf("Error while rolling back transaction: %v", err)
		}
	}(tx)

	// Get cost entry from database
	cost, repoErr := cc.CostRepo.GetCostByID(costId)
	if repoErr != nil {
		return nil, repoErr
	}

	amountChanged := request.Amount != ""
	creditorChanged := request.Creditor != ""
	contributorsChanged := (request.Contributors != nil && len(request.Contributors) > 0) || creditorChanged

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

	// If only amount has changed, not the contributors, then get the contributors from the database
	if amountChanged && !contributorsChanged {
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

	/*// Delete cost contributions
	if repoErr := cc.CostRepo.DeleteCostContributions(costId); repoErr != nil {
		return nil, repoErr
	}*/

	// Delete old cost contributions and subtract debt from users
	contributors, repoErr := cc.CostRepo.GetCostContributors(costId)
	if repoErr != nil {
		return nil, repoErr
	}
	for _, contributor := range contributors {
		// Get user from database
		user, repoErr := cc.UserRepo.GetUserById(contributor.UserID)
		if repoErr != nil {
			return nil, repoErr
		}

		// Delete cost contribution from database
		if repoErr := cc.CostRepo.DeleteCostContributionTx(tx, user.UserID); repoErr != nil {
			return nil, repoErr
		}

		// Subtract debt from user
		if repoErr := cc.calculateDebt(tx, creditor, user, tripId, contributor.Amount.Neg()); repoErr != nil {
			return nil, repoErr
		}
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
		if repoErr := cc.CostRepo.AddCostContributorTx(tx, contribution); repoErr != nil {
			return nil, repoErr
		}

		// Add debt to user
		if repoErr := cc.calculateDebt(tx, creditor, user, tripId, contribution.Amount); repoErr != nil {
			return nil, repoErr
		}
	}

	// Update cost entry in database
	if repoErr := cc.CostRepo.UpdateTx(tx, cost); repoErr != nil {
		return nil, repoErr
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error while committing transaction: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return cc.mapCostToResponse(cost), nil
}

func (cc *CostController) DeleteCostEntry(ctx context.Context, tripId *uuid.UUID, costId *uuid.UUID) *models.ExpenseServiceError {
	// Begin transaction
	tx, err := cc.DatabaseMgr.BeginTx(ctx)
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

	// Get all cost contributions
	contributions, repoErr := cc.CostRepo.GetCostContributors(costId)
	if repoErr != nil {
		return repoErr
	}

	// Get creditor
	creditor, repoErr := cc.CostRepo.GetCostCreditor(costId)
	if repoErr != nil {
		return repoErr
	}

	// Remove debt from all contributors
	for _, contribution := range contributions {
		if contribution.IsCreditor {
			continue
		}

		// Get user from database
		user, repoErr := cc.UserRepo.GetUserById(contribution.UserID)
		if repoErr != nil {
			return repoErr
		}

		// Subtract debt from user
		cc.calculateDebt(tx, creditor, user, tripId, contribution.Amount.Neg())
	}

	repoErr = cc.CostRepo.DeleteTx(tx, costId)
	if repoErr != nil {
		return repoErr
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error while committing transaction: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

//**********************************************************************************************************************
// Helper functions
//**********************************************************************************************************************

// mapCostToResponse maps a cost entry to a cost response
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

func (cc *CostController) calculateDebt(tx *sql.Tx, creditor *models.UserSchema, debtor *models.UserSchema, tripId *uuid.UUID, amountToAdd decimal.Decimal) *models.ExpenseServiceError {
	if creditor.Username == debtor.Username {
		return nil
	}
	// Check if debt already exists
	log.Printf("Debug: Creditor: %v, Debtor: %v, TripId: %v, Amount: %v", creditor.UserID, debtor.UserID, tripId, amountToAdd)
	debt, repoErr := cc.DebtRepo.GetDebtByCreditorIdAndDebtorIdAndTripId(creditor.UserID, debtor.UserID, tripId)
	if repoErr != nil {
		return repoErr
	}

	// Update existing debt
	debt.Amount = debt.Amount.Add(amountToAdd)
	repoErr = cc.DebtRepo.UpdateTx(tx, debt)
	if repoErr != nil {
		return repoErr
	}

	otherDebt, repoErr := cc.DebtRepo.GetDebtByCreditorIdAndDebtorIdAndTripId(debtor.UserID, creditor.UserID, tripId)
	if repoErr != nil {
		return repoErr
	}

	// Update existing debt
	otherDebt.Amount = otherDebt.Amount.Sub(amountToAdd)
	repoErr = cc.DebtRepo.UpdateTx(tx, otherDebt)
	if repoErr != nil {
		return repoErr
	}

	return nil
}
