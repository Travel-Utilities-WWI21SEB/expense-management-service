package controllers

import (
	"context"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/repositories"
	"time"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/google/uuid"
)

// TripCtl Exposed interface to the handler-package
type TripCtl interface {
	CreateTripEntry(ctx context.Context, tripData models.CreateTripRequest) (*models.TripResponse, *models.ExpenseServiceError)
	UpdateTripEntry(ctx context.Context, tripID *uuid.UUID, tripUpdateData models.UpdateTripRequest) (*models.TripResponse, *models.ExpenseServiceError)
	GetTripDetails(ctx context.Context, tripID *uuid.UUID) (*models.TripResponse, *models.ExpenseServiceError)
	DeleteTripEntry(ctx context.Context, tripID *uuid.UUID) *models.ExpenseServiceError
	GetTripEntries(ctx context.Context) ([]*models.TripResponse, *models.ExpenseServiceError)
	InviteUserToTrip(ctx context.Context, tripId *uuid.UUID, inviteUserRequest models.InviteUserRequest) (*models.TripResponse, *models.ExpenseServiceError)
	AcceptTripInvite(ctx context.Context, tripId *uuid.UUID) *models.ExpenseServiceError
}

// TripController Trip Controller structure
type TripController struct {
	DatabaseMgr      managers.DatabaseMgr
	TripRepo         repositories.TripRepo
	UserRepo         repositories.UserRepo
	CostRepo         repositories.CostRepo
	CostCategoryRepo repositories.CostCategoryRepo
}

func (tc *TripController) CreateTripEntry(ctx context.Context, tripRequest models.CreateTripRequest) (*models.TripResponse, *models.ExpenseServiceError) {
	// Create new trip
	tripID := uuid.New()
	tripStartDate, err := time.Parse(time.DateOnly, tripRequest.StartDate)
	if err != nil {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	tripEndDate, err := time.Parse(time.DateOnly, tripRequest.EndDate)
	if err != nil {
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	trip := &models.TripSchema{
		TripID:      &tripID,
		Name:        tripRequest.Name,
		Description: tripRequest.Description,
		Location:    tripRequest.Location,
		StartDate:   &tripStartDate,
		EndDate:     &tripEndDate,
	}

	// Insert trip into database
	if repoErr := tc.TripRepo.CreateTrip(trip); repoErr != nil {
		return nil, repoErr
	}

	// Insert user-trip association into database
	if repoErr := tc.TripRepo.AddUserToTrip(trip, ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID), true); repoErr != nil {
		// TODO: Delete trip from database
		return nil, repoErr
	}

	return tc.buildTripResponse(trip)
}

func (tc *TripController) GetTripEntries(ctx context.Context) ([]*models.TripResponse, *models.ExpenseServiceError) {
	// Get trips from database
	trips, repoErr := tc.TripRepo.GetTripsByUserId(ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID))
	if repoErr != nil {
		return nil, repoErr
	}

	// Iterate over rows and create trip response
	var tripResponses []*models.TripResponse
	for _, trip := range trips {

		// Append trip response to response array
		tripResponse, err := tc.buildTripResponse(trip)
		if err != nil {
			return nil, err
		}
		tripResponses = append(tripResponses, tripResponse)
	}

	return tripResponses, nil
}

func (tc *TripController) UpdateTripEntry(ctx context.Context, tripID *uuid.UUID, tripRequest models.UpdateTripRequest) (*models.TripResponse, *models.ExpenseServiceError) {
	// Check if user accepted trip invite
	if repoErr := tc.TripRepo.ValidateIfUserHasAccepted(tripID, ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)); repoErr != nil {
		return nil, repoErr
	}

	// Get trip from database
	trip, repoErr := tc.TripRepo.GetTripById(tripID)
	if repoErr != nil {
		return nil, repoErr
	}

	// Update trip data
	if tripRequest.Name != "" {
		trip.Name = tripRequest.Name
	}

	if tripRequest.Description != "" {
		trip.Description = tripRequest.Description
	}

	if tripRequest.Location != "" {
		trip.Location = tripRequest.Location
	}

	if tripRequest.StartDate != "" {
		*trip.StartDate, _ = time.Parse(time.DateOnly, tripRequest.StartDate)
	}

	if tripRequest.EndDate != "" {
		*trip.EndDate, _ = time.Parse(time.DateOnly, tripRequest.EndDate)
	}

	// Update trip in database
	repoErr = tc.TripRepo.UpdateTrip(trip)
	if repoErr != nil {
		return nil, repoErr
	}

	return tc.buildTripResponse(trip)
}

func (tc *TripController) GetTripDetails(_ context.Context, tripID *uuid.UUID) (*models.TripResponse, *models.ExpenseServiceError) {
	// Get trip from database
	trip, repoErr := tc.TripRepo.GetTripById(tripID)
	if repoErr != nil {
		return nil, repoErr
	}
	return tc.buildTripResponse(trip)
}

func (tc *TripController) DeleteTripEntry(ctx context.Context, tripID *uuid.UUID) *models.ExpenseServiceError {
	// Check if user accepted trip invite
	if repoErr := tc.TripRepo.ValidateIfUserHasAccepted(tripID, ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)); repoErr != nil {
		return repoErr
	}

	// Delete trip from database
	return tc.TripRepo.DeleteTrip(tripID)
}

func (tc *TripController) InviteUserToTrip(ctx context.Context, tripId *uuid.UUID, inviteUserRequest models.InviteUserRequest) (*models.TripResponse, *models.ExpenseServiceError) {
	// Check if user accepted trip invite
	if repoErr := tc.TripRepo.ValidateIfUserHasAccepted(tripId, ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)); repoErr != nil {
		return nil, repoErr
	}

	// Get invitedUser data from invite
	invitedUser := &models.UserSchema{
		Email:    inviteUserRequest.EMail,
		Username: inviteUserRequest.Username,
	}

	invitedUser, repoErr := tc.UserRepo.GetUserBySchema(invitedUser)
	if repoErr != nil {
		return nil, repoErr
	}

	// Get trip data from database
	trip, repoErr := tc.TripRepo.GetTripById(tripId)
	if repoErr != nil {
		return nil, repoErr
	}

	// Invite invitedUser to trip
	if repoErr := tc.TripRepo.AddUserToTrip(trip, invitedUser.UserID, false); repoErr != nil {
		return nil, repoErr
	}

	return tc.buildTripResponse(trip)
}

func (tc *TripController) AcceptTripInvite(ctx context.Context, tripId *uuid.UUID) *models.ExpenseServiceError {
	return tc.TripRepo.AcceptTripInvite(tripId, ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID))
}

func (tc *TripController) buildTripResponse(trip *models.TripSchema) (*models.TripResponse, *models.ExpenseServiceError) {
	// Get trip participants from database
	participants, repoErr := tc.TripRepo.GetTripParticipants(trip.TripID)
	if repoErr != nil {
		return nil, repoErr
	}

	// Build participant responses
	participationResponses := make([]models.TripParticipantResponse, len(participants))
	for i, participant := range participants {
		user, repoErr := tc.UserRepo.GetUserById(participant.UserID)
		if repoErr != nil {
			return nil, repoErr
		}

		participationResponses[i] = models.TripParticipantResponse{
			Username:          user.Username,
			HasAcceptedInvite: participant.HasAccepted,
			PresenceStartDate: participant.PresenceStartDate.Format(time.DateOnly),
			PresenceEndDate:   participant.PresenceEndDate.Format(time.DateOnly),
		}
	}

	// Get total cost of trip
	totalCostOfTrip, repoErr := tc.CostRepo.GetTotalCostByTripID(trip.TripID)
	if repoErr != nil {
		return nil, repoErr
	}

	// Get costcategories from database
	costCategories, repoErr := tc.CostCategoryRepo.GetCostCategoriesByTripID(trip.TripID)
	if repoErr != nil {
		return nil, repoErr
	}

	// Build cost category responses
	costCategoryResponses := make([]models.CostCategoryResponse, len(costCategories))
	for i, costCategory := range costCategories {
		// Get total cost for cost category
		totalCostOfCategory, repoErr := tc.CostRepo.GetTotalCostByCostCategoryID(costCategory.CostCategoryID)
		if repoErr != nil {
			return nil, repoErr
		}

		costCategoryResponses[i] = models.CostCategoryResponse{
			CostCategoryId: costCategory.CostCategoryID,
			Name:           costCategory.Name,
			Description:    costCategory.Description,
			Color:          costCategory.Color,
			Icon:           costCategory.Icon,
			TotalCost:      totalCostOfCategory.String(),
		}

	}

	// Build trip response
	return &models.TripResponse{
		TripID:         trip.TripID,
		Name:           trip.Name,
		Description:    trip.Description,
		Location:       trip.Location,
		StartDate:      trip.StartDate.Format(time.DateOnly),
		EndDate:        trip.EndDate.Format(time.DateOnly),
		CostCategories: costCategoryResponses,
		TotalCost:      totalCostOfTrip.String(),
		UserDebt:       "0.00",
		UserCredit:     "0.00",
		Participants:   participationResponses,
	}, nil
}
