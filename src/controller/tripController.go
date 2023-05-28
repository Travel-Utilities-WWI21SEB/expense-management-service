package controller

import (
	"context"
	"errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expenseerror"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/manager"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"log"
	"time"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
	"github.com/google/uuid"
)

// Exposed interface to the handler-package
type TripCtl interface {
	CreateTripEntry(ctx context.Context, tripData model.TripRequest) (*model.TripCreationResponse, *model.ExpenseServiceError)
	UpdateTripEntry(ctx context.Context, tripID *uuid.UUID, tripUpdateData model.TripUpdateRequest) (*model.TripResponse, *model.ExpenseServiceError)
	GetTripDetails(ctx context.Context, tripID *uuid.UUID) (*model.TripResponse, *model.ExpenseServiceError)
	DeleteTripEntry(ctx context.Context, tripID *uuid.UUID) error
}

// Trip Controller structure
type TripController struct {
	DatabaseMgr manager.DatabaseMgr
}

func (tc *TripController) CreateTripEntry(ctx context.Context, tripData model.TripRequest) (*model.TripCreationResponse, *model.ExpenseServiceError) {
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
		log.Printf("Error in tripController.CreateTripEntry.DatabaseMgr.ExecuteStatement(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Get user id from context
	tokenUserId, ok := ctx.Value(model.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		log.Printf("Error in tripController.CreateTripEntry.ctx.Value(): %v", ok)
		return nil, expenseerror.EXPENSE_INTERNAL_ERROR
	}

	// error if user is not logged in
	if tokenUserId == nil {
		log.Printf("Error in tripController.CreateTripEntry.ctx.Value(): %v", errors.New("user not logged in"))
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	// Insert user-trip association into database
	queryString = "INSERT INTO user_trip_association (id_user, id_trip, is_accepted, presence_start_date, presence_end_date) VALUES ($1, $2, $3, $4, $5)"
	if _, err := tc.DatabaseMgr.ExecuteStatement(queryString, tokenUserId, trip.TripID, true, trip.StartDate, trip.EndDate); err != nil {
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Return trip response
	response := &model.TripCreationResponse{
		TripID: trip.TripID,
	}

	return response, nil
}

func (tc *TripController) UpdateTripEntry(ctx context.Context, tripID *uuid.UUID, tripUpdateDate model.TripUpdateRequest) (*model.TripResponse, *model.ExpenseServiceError) {
	// TO-DO
	return nil, expenseerror.EXPENSE_BAD_REQUEST
}

func (tc *TripController) GetTripDetails(ctx context.Context, tripID *uuid.UUID) (*model.TripResponse, *model.ExpenseServiceError) {
	if utils.ContainsEmptyString(tripID.String()) {
		return nil, expenseerror.EXPENSE_INTERNAL_ERROR
	}

	// Get authenticated user id from context
	tokenUserId, ok := ctx.Value(model.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		log.Printf("Error in tripController.GetTripDetails.ctx.Value(): %v", ok)
		return nil, expenseerror.EXPENSE_INTERNAL_ERROR
	}
	log.Printf("tokenUserId: %v", tokenUserId)

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++//
	// Check if trip exists                                                                          //
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++//
	checkTripQueryString := "SELECT COUNT(*) FROM trip WHERE id = $1"
	row := tc.DatabaseMgr.ExecuteQueryRow(checkTripQueryString, tripID)

	var tripCount int
	if err := row.Scan(&tripCount); err != nil {
		log.Printf("Error in tripController.GetTripDetails.rows.Scan(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}
	if tripCount == 0 {
		log.Printf("Error in tripController.GetTripDetails.tripCount: %v", tripCount)
		return nil, expenseerror.EXPENSE_NOT_FOUND
	}

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++//
	// Check if user is part of trip (user-trip association)                                         //
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++//
	checkQueryString := "SELECT COUNT(*) FROM user_trip_association WHERE id_trip = $1 AND id_user = $2"
	row = tc.DatabaseMgr.ExecuteQueryRow(checkQueryString, tripID, tokenUserId)

	var associationCount int
	if err := row.Scan(&associationCount); err != nil {
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	if associationCount == 0 {
		return nil, expenseerror.EXPENSE_FORBIDDEN
	}

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++//
	// Get trip details                                                                              //
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++//
	queryString := "SELECT id, location, start_date, end_date FROM trip WHERE id = $1"
	row = tc.DatabaseMgr.ExecuteQueryRow(queryString, tripID)

	var tripResponse model.TripResponse
	if err := row.Scan(&tripResponse.TripID, &tripResponse.Location, &tripResponse.StartDate, &tripResponse.EndDate); err != nil {
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	log.Printf("tripResponse: %v", tripResponse)

	return &tripResponse, nil
}

func (tc *TripController) DeleteTripEntry(ctx context.Context, tripID *uuid.UUID) error {
	// TO-DO
	return errors.New("not implemented")
}
