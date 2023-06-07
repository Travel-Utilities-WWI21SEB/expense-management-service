package controllers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/lib/pq"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/google/uuid"
)

// Exposed interface to the handler-package
type TripCtl interface {
	CreateTripEntry(ctx context.Context, tripData models.CreateTripRequest) (*models.TripCreationResponse, *models.ExpenseServiceError)
	UpdateTripEntry(ctx context.Context, tripID *uuid.UUID, tripUpdateData models.UpdateTripRequest) (*models.TripResponse, *models.ExpenseServiceError)
	GetTripDetails(ctx context.Context, tripID *uuid.UUID) (*models.TripResponse, *models.ExpenseServiceError)
	DeleteTripEntry(ctx context.Context, tripID *uuid.UUID) *models.ExpenseServiceError
	GetTripEntries(ctx context.Context) ([]*models.TripResponse, *models.ExpenseServiceError)
	InviteUserToTrip(ctx context.Context, tripId *uuid.UUID, inviteUserRequest models.InviteUserRequest) (*models.TripResponse, *models.ExpenseServiceError)
	AcceptTripInvite(ctx context.Context, tripId *uuid.UUID) *models.ExpenseServiceError
}

// Trip Controller structure
type TripController struct {
	DatabaseMgr managers.DatabaseMgr
}

func (tc *TripController) CreateTripEntry(ctx context.Context, tripData models.CreateTripRequest) (*models.TripCreationResponse, *models.ExpenseServiceError) {
	if utils.ContainsEmptyString(tripData.Location, tripData.StartDate, tripData.EndDate) {
		log.Printf("Error in creating trip: %v", errors.New("empty string in request"))
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	// Create new trip
	tripID := uuid.New()
	tripStartDate, err := time.Parse(time.DateOnly, tripData.StartDate)
	if err != nil {
		log.Printf("Error in parsing trip start date: %v", err)
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	tripEndDate, err := time.Parse(time.DateOnly, tripData.EndDate)
	if err != nil {
		log.Printf("Error in parsing trip end date: %v", err)
		return nil, expense_errors.EXPENSE_BAD_REQUEST
	}

	trip := &models.TripSchema{
		TripID:    &tripID,
		Location:  tripData.Location,
		StartDate: &tripStartDate,
		EndDate:   &tripEndDate,
	}

	// Insert trip into database
	queryString := "INSERT INTO trip (id, location, start_date, end_date) VALUES ($1, $2, $3, $4)"
	if _, err := tc.DatabaseMgr.ExecuteStatement(queryString, trip.TripID, trip.Location, trip.StartDate, trip.EndDate); err != nil {
		log.Printf("Error in tripController.CreateTripEntry.DatabaseMgr.ExecuteStatement(): %v", err)
		return nil, expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	// Get user id from context
	tokenUserId, ok := ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		log.Printf("Error in tripController.CreateTripEntry.ctx.Value(): %v", ok)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Insert user-trip association into database
	queryString = "INSERT INTO user_trip_association (id_user, id_trip, is_accepted, presence_start_date, presence_end_date) VALUES ($1, $2, $3, $4, $5)"
	if _, err := tc.DatabaseMgr.ExecuteStatement(queryString, tokenUserId, trip.TripID, true, trip.StartDate, trip.EndDate); err != nil {
		log.Printf("Error in tripController.CreateTripEntry.DatabaseMgr.ExecuteStatement(): %v", err)
		return nil, expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	// Return trip response
	response := &models.TripCreationResponse{
		TripID: trip.TripID,
	}

	return response, nil
}

func (tc *TripController) GetTripEntries(ctx context.Context) ([]*models.TripResponse, *models.ExpenseServiceError) {
	// Get user id from context
	tokenUserId, ok := ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		log.Printf("Error in tripController.GetTripEntries.ctx.Value(): %v", ok)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Get trips from database
	queryString := "SELECT trip.id, trip.location, trip.start_date, trip.end_date FROM trip " +
		"INNER JOIN user_trip_association ON trip.id = user_trip_association.id_trip WHERE user_trip_association.id_user = $1"
	rows, err := tc.DatabaseMgr.ExecuteQuery(queryString, tokenUserId)
	if err != nil {
		log.Printf("Error in tripController.GetTripEntries.DatabaseMgr.ExecuteQuery(): %v", err)
		return nil, expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	// Iterate over rows and create trip response
	var tripResponses []*models.TripResponse
	for rows.Next() {
		var tripResponse models.TripResponse
		if err := rows.Scan(&tripResponse.TripID, &tripResponse.Location, &tripResponse.StartDate, &tripResponse.EndDate); err != nil {
			log.Printf("Error in tripController.GetTripEntries.rows.Scan(): %v", err)
			return nil, expense_errors.EXPENSE_UPSTREAM_ERROR
		}
		tripResponses = append(tripResponses, &tripResponse)
	}

	return tripResponses, nil
}

func (tc *TripController) UpdateTripEntry(ctx context.Context, tripID *uuid.UUID, tripUpdateDate models.UpdateTripRequest) (*models.TripResponse, *models.ExpenseServiceError) {
	// Get old trip data
	getTripQueryString := "SELECT location, start_date, end_date FROM trip WHERE id = $1"
	row := tc.DatabaseMgr.ExecuteQueryRow(getTripQueryString, tripID)
	var location string
	var startDate time.Time
	var endDate time.Time
	if err := row.Scan(&location, &startDate, &endDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, expense_errors.EXPENSE_TRIP_NOT_FOUND
		}
		log.Printf("Error in tripController.UpdateTripEntry.DatabaseMgr.ExecuteQueryRow(): %v", err)
		return nil, expense_errors.EXPENSE_UPSTREAM_ERROR
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
		return nil, expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	// Return trip response
	response := &models.TripResponse{
		TripID:    tripID,
		Location:  location,
		StartDate: startDate.String(),
		EndDate:   endDate.String(),
	}

	return response, nil
}

func (tc *TripController) GetTripDetails(ctx context.Context, tripID *uuid.UUID) (*models.TripResponse, *models.ExpenseServiceError) {
	// Get trip details
	queryString := "SELECT id, location, start_date, end_date FROM trip WHERE id = $1"
	row := tc.DatabaseMgr.ExecuteQueryRow(queryString, tripID)

	var tripResponse models.TripResponse
	if err := row.Scan(&tripResponse.TripID, &tripResponse.Location, &tripResponse.StartDate, &tripResponse.EndDate); err != nil {
		log.Printf("Error in tripController.GetTripDetails.rows.Scan(): %v", err)
		return nil, expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	return &tripResponse, nil
}

func (tc *TripController) DeleteTripEntry(ctx context.Context, tripID *uuid.UUID) *models.ExpenseServiceError {
	// delete trip
	deleteString := "DELETE FROM trip WHERE id = $1"
	if _, err := tc.DatabaseMgr.ExecuteStatement(deleteString, tripID); err != nil {
		log.Printf("Error in tripController.DeleteTripEntry.DatabaseMgr.ExecuteStatement(): %v", err)
		return expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	return nil
}

func (tc *TripController) InviteUserToTrip(ctx context.Context, tripId *uuid.UUID, inviteUserRequest models.InviteUserRequest) (*models.TripResponse, *models.ExpenseServiceError) {
	// Get trip details
	tripDetails := &models.TripResponse{}
	getTripDetailsQueryString := "SELECT id, location, start_date, end_date FROM trip WHERE id = $1"
	row := tc.DatabaseMgr.ExecuteQueryRow(getTripDetailsQueryString, tripId)
	if err := row.Scan(&tripDetails.TripID, &tripDetails.Location, &tripDetails.StartDate, &tripDetails.EndDate); err != nil {
		log.Printf("Error in tripController.InviteUserToTrip.DatabaseMgr.ExecuteQueryRow(): %v", err)
		return nil, expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	// Add user to trip
	addUserToTripQueryString := "INSERT INTO user_trip_association (id_trip, id_user, is_accepted) VALUES ($1, $2, $3)"
	if _, err := tc.DatabaseMgr.ExecuteStatement(addUserToTripQueryString, tripId, inviteUserRequest.UserID, false); err != nil {
		// if err is unique_violation, then user is already part of trip
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, expense_errors.EXPENSE_CONFLICT
		}

		// if err is foreign_key_violation, then user does not exist
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return nil, expense_errors.EXPENSE_USER_NOT_FOUND
		}

		log.Printf("Error in tripController.InviteUserToTrip.DatabaseMgr.ExecuteStatement(): %v", err)
		return nil, expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	return tripDetails, nil
}

func (tc *TripController) AcceptTripInvite(ctx context.Context, tripId *uuid.UUID) *models.ExpenseServiceError {
	if utils.ContainsEmptyString(tripId.String()) {
		return expense_errors.EXPENSE_BAD_REQUEST
	}

	// Get trip details
	tripDetails := &models.TripResponse{}
	getTripDetailsQueryString := "SELECT id, location, start_date, end_date FROM trip WHERE id = $1"
	row := tc.DatabaseMgr.ExecuteQueryRow(getTripDetailsQueryString, tripId)
	if err := row.Scan(&tripDetails.TripID, &tripDetails.Location, &tripDetails.StartDate, &tripDetails.EndDate); err != nil {
		if sql.ErrNoRows == err {
			return expense_errors.EXPENSE_TRIP_NOT_FOUND
		}
		log.Printf("Error in tripController.AcceptTripInvite.DatabaseMgr.ExecuteQueryRow(): %v", err)
		return expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	// Get authenticated user id from context
	tokenUserId, ok := ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		log.Printf("Error in tripController.AcceptTripInvite.ctx.Value(): %v", ok)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Update user_trip_association
	updateUserTripQueryString := "UPDATE user_trip_association SET is_accepted = true WHERE id_trip = $1 AND id_user = $2"
	result, err := tc.DatabaseMgr.ExecuteStatement(updateUserTripQueryString, tripId, tokenUserId)
	if err != nil {
		log.Printf("Error in tripController.AcceptTripInvite.DatabaseMgr.ExecuteStatement(): %v", err)
		return expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	//if affectedRows is 0, then user is already accepted or not invited to trip, error code is 409: Conflict
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error in tripController.AcceptTripInvite.RowsAffected(): %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected == 0 {
		return expense_errors.EXPENSE_ALREADY_ACCEPTED
	}

	return nil
}
