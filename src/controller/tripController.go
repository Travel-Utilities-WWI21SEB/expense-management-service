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
	DeleteTripEntry(ctx context.Context, tripID *uuid.UUID) *model.ExpenseServiceError
	GetTripEntries(ctx context.Context) ([]*model.TripResponse, *model.ExpenseServiceError)
}

// Trip Controller structure
type TripController struct {
	DatabaseMgr manager.DatabaseMgr
}

func (tc *TripController) CreateTripEntry(ctx context.Context, tripData model.TripRequest) (*model.TripCreationResponse, *model.ExpenseServiceError) {
	if utils.ContainsEmptyString(tripData.Location, tripData.StartDate, tripData.EndDate) {
		log.Printf("Error in creating trip: %v", errors.New("empty string in request"))
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	// Create new trip
	tripID := uuid.New()
	tripStartDate, err := time.Parse(time.DateOnly, tripData.StartDate)
	if err != nil {
		log.Printf("Error in parsing trip start date: %v", err)
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	tripEndDate, err := time.Parse(time.DateOnly, tripData.EndDate)
	if err != nil {
		log.Printf("Error in parsing trip end date: %v", err)
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
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	// Insert user-trip association into database
	queryString = "INSERT INTO user_trip_association (id_user, id_trip, is_accepted, presence_start_date, presence_end_date) VALUES ($1, $2, $3, $4, $5)"
	if _, err := tc.DatabaseMgr.ExecuteStatement(queryString, tokenUserId, trip.TripID, true, trip.StartDate, trip.EndDate); err != nil {
		log.Printf("Error in tripController.CreateTripEntry.DatabaseMgr.ExecuteStatement(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Return trip response
	response := &model.TripCreationResponse{
		TripID: trip.TripID,
	}

	return response, nil
}

func (tc *TripController) GetTripEntries(ctx context.Context) ([]*model.TripResponse, *model.ExpenseServiceError) {
	// Get user id from context
	tokenUserId, ok := ctx.Value(model.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		log.Printf("Error in tripController.GetTripEntries.ctx.Value(): %v", ok)
		return nil, expenseerror.EXPENSE_INTERNAL_ERROR
	}

	// error if user is not logged in
	if tokenUserId == nil {
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	// Get trips from database
	queryString := "SELECT trip.id, trip.location, trip.start_date, trip.end_date FROM trip " +
		"INNER JOIN user_trip_association ON trip.id = user_trip_association.id_trip WHERE user_trip_association.id_user = $1"
	rows, err := tc.DatabaseMgr.ExecuteQuery(queryString, tokenUserId)
	if err != nil {
		log.Printf("Error in tripController.GetTripEntries.DatabaseMgr.ExecuteQuery(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Iterate over rows and create trip response
	var tripResponses []*model.TripResponse
	for rows.Next() {
		var tripResponse model.TripResponse
		if err := rows.Scan(&tripResponse.TripID, &tripResponse.Location, &tripResponse.StartDate, &tripResponse.EndDate); err != nil {
			log.Printf("Error in tripController.GetTripEntries.rows.Scan(): %v", err)
			return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
		}
		tripResponses = append(tripResponses, &tripResponse)
	}

	return tripResponses, nil
}

func (tc *TripController) UpdateTripEntry(ctx context.Context, tripID *uuid.UUID, tripUpdateDate model.TripUpdateRequest) (*model.TripResponse, *model.ExpenseServiceError) {
	if utils.ContainsEmptyString(tripID.String()) {
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	// Get authenticated user id from context
	tokenUserId, ok := ctx.Value(model.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		log.Printf("Error in tripController.UpdateTripEntry.ctx.Value(): %v", ok)
		return nil, expenseerror.EXPENSE_INTERNAL_ERROR
	}

	// Check if trip exists
	checkTripQueryString := "SELECT COUNT(*) FROM trip WHERE id = $1"
	row := tc.DatabaseMgr.ExecuteQueryRow(checkTripQueryString, tripID)
	var count int
	if err := row.Scan(&count); err != nil {
		log.Printf("Error in tripController.UpdateTripEntry.DatabaseMgr.ExecuteQueryRow(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	if count == 0 {
		return nil, expenseerror.EXPENSE_NOT_FOUND
	}

	// Check if user is associated with trip
	checkUserTripQueryString := "SELECT COUNT(*) FROM user_trip_association WHERE id_user = $1 AND id_trip = $2"
	row = tc.DatabaseMgr.ExecuteQueryRow(checkUserTripQueryString, tokenUserId, tripID)
	if err := row.Scan(&count); err != nil {
		log.Printf("Error in tripController.UpdateTripEntry.DatabaseMgr.ExecuteQueryRow(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	if count == 0 {
		return nil, expenseerror.EXPENSE_FORBIDDEN
	}

	// Get old trip data
	getTripQueryString := "SELECT location, start_date, end_date FROM trip WHERE id = $1"
	row = tc.DatabaseMgr.ExecuteQueryRow(getTripQueryString, tripID)
	var location string
	var startDate time.Time
	var endDate time.Time
	if err := row.Scan(&location, &startDate, &endDate); err != nil {
		log.Printf("Error in tripController.UpdateTripEntry.DatabaseMgr.ExecuteQueryRow(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Update trip data
	if tripUpdateDate.Location != nil {
		location = *tripUpdateDate.Location
	}

	if tripUpdateDate.StartDate != nil {
		startDate, _ = time.Parse(time.DateOnly, *tripUpdateDate.StartDate)
	}

	if tripUpdateDate.EndDate != nil {
		endDate, _ = time.Parse(time.DateOnly, *tripUpdateDate.EndDate)
	}

	// Update trip in database
	updateTripQueryString := "UPDATE trip SET location = $1, start_date = $2, end_date = $3 WHERE id = $4"
	if _, err := tc.DatabaseMgr.ExecuteStatement(updateTripQueryString, location, startDate, endDate, tripID); err != nil {
		log.Printf("Error in tripController.UpdateTripEntry.DatabaseMgr.ExecuteStatement(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Return trip response
	response := &model.TripResponse{
		TripID:    tripID,
		Location:  location,
		StartDate: startDate.String(),
		EndDate:   endDate.String(),
	}

	return response, nil
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
		log.Printf("Error in tripController.GetTripDetails.rows.Scan(): %v", err)
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
		log.Printf("Error in tripController.GetTripDetails.rows.Scan(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	log.Printf("tripResponse: %v", tripResponse)

	return &tripResponse, nil
}

func (tc *TripController) DeleteTripEntry(ctx context.Context, tripID *uuid.UUID) *model.ExpenseServiceError {
	if utils.ContainsEmptyString(tripID.String()) {
		return expenseerror.EXPENSE_INTERNAL_ERROR
	}

	// get user id from context
	tokenUserId, ok := ctx.Value(model.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		log.Printf("Error in tripController.DeleteTripEntry.ctx.Value(): %v", ok)
		return expenseerror.EXPENSE_INTERNAL_ERROR
	}

	// check if trip exists
	checkTripQueryString := "SELECT COUNT(*) FROM trip WHERE id = $1"
	row := tc.DatabaseMgr.ExecuteQueryRow(checkTripQueryString, tripID)
	var count int
	if err := row.Scan(&count); err != nil {
		log.Printf("Error in tripController.DeleteTripEntry.DatabaseMgr.ExecuteQueryRow(): %v", err)
		return expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	if count == 0 {
		return expenseerror.EXPENSE_NOT_FOUND
	}

	// check if user is part of trip
	checkUserTripQueryString := "SELECT COUNT(*) FROM user_trip_association WHERE id_trip = $1 AND id_user = $2"
	row = tc.DatabaseMgr.ExecuteQueryRow(checkUserTripQueryString, tripID, tokenUserId)
	if err := row.Scan(&count); err != nil {
		log.Printf("Error in tripController.DeleteTripEntry.DatabaseMgr.ExecuteQueryRow(): %v", err)
		return expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	if count == 0 {
		return expenseerror.EXPENSE_FORBIDDEN
	}

	// delete trip
	deleteTripQueryString := "DELETE FROM trip WHERE id = $1"
	if _, err := tc.DatabaseMgr.ExecuteStatement(deleteTripQueryString, tripID); err != nil {
		log.Printf("Error in tripController.DeleteTripEntry.DatabaseMgr.ExecuteStatement(): %v", err)
		return expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	return nil
}
