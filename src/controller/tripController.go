package controller

import (
	"context"
	"errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expenseerror"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/manager"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"time"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
	"github.com/google/uuid"
)

// Exposed interface to the handler-package
type TripCtl interface {
	CreateTripEntry(ctx context.Context, tripData model.TripRequest) (*model.TripResponse, *model.ExpenseServiceError)
	UpdateTripEntry(ctx context.Context, tripID *uuid.UUID) (*model.TripResponse, error)
	GetTripDetails(ctx context.Context, tripID *uuid.UUID) (*model.TripResponse, error)
	DeleteTripEntry(ctx context.Context, tripID *uuid.UUID) error
}

// Trip Controller structure
type TripController struct {
	DatabaseMgr manager.DatabaseMgr
}

func (tc *TripController) CreateTripEntry(ctx context.Context, tripData model.TripRequest) (*model.TripResponse, *model.ExpenseServiceError) {
	if utils.ContainsEmptyString(tripData.Location, tripData.StartDate, tripData.EndDate) {
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	// Create new trip
	tripID := uuid.New()
	tripStartDate, err := time.Parse(time.DateOnly, tripData.StartDate)
	if err != nil {
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	tripEndDate, err := time.Parse(time.DateOnly, tripData.EndDate)
	if err != nil {
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	trip := &model.TripSchema{
		TripID:    &tripID,
		Location:  tripData.Location,
		StartDate: &tripStartDate,
		EndDate:   &tripEndDate,
	}

	// Insert trip into database
	queryString := "INSERT INTO trip (id, location, start_date, end_date) VALUES ($1, $2, $3, $4)"
	if _, err := tc.DatabaseMgr.ExecuteStatement(queryString, trip.TripID, trip.Location, trip.StartDate, trip.EndDate); err != nil {
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Get user id from context
	tokenUserId, ok := ctx.Value(model.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		return nil, expenseerror.EXPENSE_INTERNAL_ERROR
	}

	// error if user is not logged in
	if tokenUserId == nil {
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	// Insert user-trip association into database
	queryString = "INSERT INTO user_trip_association (id_user, id_trip, is_accepted, presence_start_date, presence_end_date) VALUES ($1, $2, $3, $4, $5)"
	if _, err := tc.DatabaseMgr.ExecuteStatement(queryString, tokenUserId, trip.TripID, true, trip.StartDate, trip.EndDate); err != nil {
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Return trip response
	response := &model.TripResponse{
		TripID: trip.TripID,
	}

	return response, nil
}

func (tc *TripController) UpdateTripEntry(ctx context.Context, tripID *uuid.UUID) (*model.TripResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (tc *TripController) GetTripDetails(ctx context.Context, tripID *uuid.UUID) (*model.TripResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (tc *TripController) DeleteTripEntry(ctx context.Context, tripID *uuid.UUID) error {
	// TO-DO
	return errors.New("not implemented")
}
