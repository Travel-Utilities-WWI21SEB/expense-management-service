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
	DatabaseMgr managers.DatabaseMgr
	TripRepo    repositories.TripRepo
	UserRepo    repositories.UserRepo
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
		TripID:    &tripID,
		Location:  tripRequest.Location,
		StartDate: &tripStartDate,
		EndDate:   &tripEndDate,
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

	return tc.buildTripResponse(&tripID)
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
		tripResponse, err := tc.buildTripResponse(trip.TripID)
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

	return tc.buildTripResponse(tripID)
}

func (tc *TripController) GetTripDetails(_ context.Context, tripID *uuid.UUID) (*models.TripResponse, *models.ExpenseServiceError) {
	return tc.buildTripResponse(tripID)
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

	return tc.buildTripResponse(tripId)
}

func (tc *TripController) AcceptTripInvite(ctx context.Context, tripId *uuid.UUID) *models.ExpenseServiceError {
	return tc.TripRepo.AcceptTripInvite(tripId, ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID))
}

func (tc *TripController) buildTripResponse(tripId *uuid.UUID) (*models.TripResponse, *models.ExpenseServiceError) {
	trip, repoErr := tc.TripRepo.GetTripById(tripId)
	if repoErr != nil {
		return nil, repoErr
	}

	participants, repoErr := tc.TripRepo.GetTripParticipants(trip.TripID)
	if repoErr != nil {
		return nil, repoErr
	}

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

	return &models.TripResponse{
		TripID:       trip.TripID,
		Location:     trip.Location,
		StartDate:    trip.StartDate.Format(time.DateOnly),
		EndDate:      trip.EndDate.Format(time.DateOnly),
		Participants: participationResponses,
	}, nil
}
