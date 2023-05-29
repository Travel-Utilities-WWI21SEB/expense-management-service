package controller

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expenseerror"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/manager"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/lib/pq"
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
	InviteUserToTrip(ctx context.Context, tripId *uuid.UUID, inviteUserRequest model.InviteUserRequest) (*model.TripResponse, *model.ExpenseServiceError)
	AcceptTripInvite(ctx context.Context, tripId *uuid.UUID) (*model.TripResponse, *model.ExpenseServiceError)
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
		return nil, expenseerror.EXPENSE_TRIP_NOT_FOUND
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
	if tripUpdateDate.Location != "" {
		location = tripUpdateDate.Location
	}

	if tripUpdateDate.StartDate != "" {
		startDate, _ = time.Parse(time.DateOnly, tripUpdateDate.StartDate)
	}

	if tripUpdateDate.EndDate != "" {
		endDate, _ = time.Parse(time.DateOnly, tripUpdateDate.EndDate)
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
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	// Get authenticated user id from context
	tokenUserId, ok := ctx.Value(model.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		log.Printf("Error in tripController.GetTripDetails.ctx.Value(): %v", ok)
		return nil, expenseerror.EXPENSE_INTERNAL_ERROR
	}

	// Get trip details
	queryString := "SELECT id, location, start_date, end_date FROM trip WHERE id = $1"
	row := tc.DatabaseMgr.ExecuteQueryRow(queryString, tripID)

	var tripResponse model.TripResponse
	if err := row.Scan(&tripResponse.TripID, &tripResponse.Location, &tripResponse.StartDate, &tripResponse.EndDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, expenseerror.EXPENSE_TRIP_NOT_FOUND
		}
		log.Printf("Error in tripController.GetTripDetails.rows.Scan(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Check if user is part of trip (user-trip association)
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

	return &tripResponse, nil
}

func (tc *TripController) DeleteTripEntry(ctx context.Context, tripID *uuid.UUID) *model.ExpenseServiceError {
	if utils.ContainsEmptyString(tripID.String()) {
		return expenseerror.EXPENSE_BAD_REQUEST
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
		return expenseerror.EXPENSE_TRIP_NOT_FOUND
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

func (tc *TripController) InviteUserToTrip(ctx context.Context, tripId *uuid.UUID, inviteUserRequest model.InviteUserRequest) (*model.TripResponse, *model.ExpenseServiceError) {
	// Checks:
	// 1. Check if tripId is empty
	// 2. Get trip details (Checks if trip exists)
	// 3. Check if tokenUser is part of trip
	// 4. Check if user to invite is already invited to trip
	// Then insert into user_trip_association
	// Then return trip details

	if utils.ContainsEmptyString(tripId.String()) {
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	// Get trip details
	tripDetails := &model.TripResponse{}
	getTripDetailsQueryString := "SELECT id, location, start_date, end_date FROM trip WHERE id = $1"
	row := tc.DatabaseMgr.ExecuteQueryRow(getTripDetailsQueryString, tripId)
	if err := row.Scan(&tripDetails.TripID, &tripDetails.Location, &tripDetails.StartDate, &tripDetails.EndDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, expenseerror.EXPENSE_TRIP_NOT_FOUND
		}
		log.Printf("Error in tripController.InviteUserToTrip.DatabaseMgr.ExecuteQueryRow(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Get authenticated user id from context
	tokenUserId, ok := ctx.Value(model.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		log.Printf("Error in tripController.InviteUserToTrip.ctx.Value(): %v", ok)
		return nil, expenseerror.EXPENSE_INTERNAL_ERROR
	}

	// Get user id from inviteUserRequest
	getUserIdQueryString := "SELECT id FROM \"user\" WHERE username = $1 OR email = $2 LIMIT 1" // TODO: Check if username and email are for the same user
	row = tc.DatabaseMgr.ExecuteQueryRow(getUserIdQueryString, inviteUserRequest.Username, inviteUserRequest.Email)
	var userId uuid.UUID
	if err := row.Scan(&userId); err != nil {
		if err == sql.ErrNoRows {
			return nil, expenseerror.EXPENSE_USER_NOT_FOUND
		}
		log.Printf("Error in tripController.InviteUserToTrip.DatabaseMgr.ExecuteQueryRow(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Check if tokenUser is part of trip
	var count int
	checkUserTripQueryString := "SELECT COUNT(*) FROM user_trip_association WHERE id_trip = $1 AND id_user = $2"
	row = tc.DatabaseMgr.ExecuteQueryRow(checkUserTripQueryString, tripId, tokenUserId)
	if err := row.Scan(&count); err != nil {
		log.Printf("Error in tripController.InviteUserToTrip.DatabaseMgr.ExecuteQueryRow(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	if count == 0 {
		return nil, expenseerror.EXPENSE_FORBIDDEN
	}

	// Add user to trip
	addUserToTripQueryString := "INSERT INTO user_trip_association (id_trip, id_user, is_accepted) VALUES ($1, $2, $3)"
	if _, err := tc.DatabaseMgr.ExecuteStatement(addUserToTripQueryString, tripId, userId, false); err != nil {
		// if err is unique_violation, then user is already part of trip
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return nil, expenseerror.EXPENSE_CONFLICT
			}
		}
		log.Printf("Error in tripController.InviteUserToTrip.DatabaseMgr.ExecuteStatement(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}
	return tripDetails, nil
}

func (tc *TripController) AcceptTripInvite(ctx context.Context, tripId *uuid.UUID) (*model.TripResponse, *model.ExpenseServiceError) {
	if utils.ContainsEmptyString(tripId.String()) {
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	// Get trip details
	tripDetails := &model.TripResponse{}
	getTripDetailsQueryString := "SELECT id, location, start_date, end_date FROM trip WHERE id = $1"
	row := tc.DatabaseMgr.ExecuteQueryRow(getTripDetailsQueryString, tripId)
	if err := row.Scan(&tripDetails.TripID, &tripDetails.Location, &tripDetails.StartDate, &tripDetails.EndDate); err != nil {
		if sql.ErrNoRows == err {
			return nil, expenseerror.EXPENSE_TRIP_NOT_FOUND
		}
		log.Printf("Error in tripController.AcceptTripInvite.DatabaseMgr.ExecuteQueryRow(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Get authenticated user id from context
	tokenUserId, ok := ctx.Value(model.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		log.Printf("Error in tripController.AcceptTripInvite.ctx.Value(): %v", ok)
		return nil, expenseerror.EXPENSE_INTERNAL_ERROR
	}

	// Update user_trip_association
	updateUserTripQueryString := "UPDATE user_trip_association SET is_accepted = true WHERE id_trip = $1 AND id_user = $2"
	if query, err := tc.DatabaseMgr.ExecuteStatement(updateUserTripQueryString, tripId, tokenUserId); err != nil {
		// if affectedRows is 0, then user is already accepted or not invited to trip, error code is 409: Conflict
		if affectedRows, _ := query.RowsAffected(); affectedRows == 0 {
			return nil, expenseerror.EXPENSE_CONFLICT
		}
		log.Printf("Error in tripController.AcceptTripInvite.DatabaseMgr.ExecuteStatement(): %v", err)
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	return tripDetails, nil
}
