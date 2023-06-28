package repositories

import (
	"database/sql"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
)

type TripRepo interface {
	GetTripById(tripId *uuid.UUID) (*models.TripSchema, *models.ExpenseServiceError)
	GetTripsByUserId(userId *uuid.UUID) ([]*models.TripSchema, *models.ExpenseServiceError)
	CreateTrip(trip *models.TripSchema) *models.ExpenseServiceError
	UpdateTrip(trip *models.TripSchema) *models.ExpenseServiceError
	DeleteTrip(tripId *uuid.UUID) *models.ExpenseServiceError

	AddUserToTrip(trip *models.TripSchema, invitedUserId *uuid.UUID, isCreator bool) *models.ExpenseServiceError
	AcceptTripInvite(tripId *uuid.UUID, userId *uuid.UUID) *models.ExpenseServiceError
	DeclineTripInvite(tripId *uuid.UUID, userId *uuid.UUID) *models.ExpenseServiceError

	ValidateIfTripExists(tripId *uuid.UUID) *models.ExpenseServiceError
	ValidateIfUserHasAccepted(tripId *uuid.UUID, userId *uuid.UUID) *models.ExpenseServiceError

	GetTripParticipant(tripId *uuid.UUID, userId *uuid.UUID) (*models.UserTripSchema, *models.ExpenseServiceError)
	GetTripParticipants(tripId *uuid.UUID) ([]*models.UserTripSchema, *models.ExpenseServiceError)
	GetAcceptedTripParticipants(tripId *uuid.UUID) ([]*models.UserTripSchema, *models.ExpenseServiceError)
	UpdateTripParticipant(userTrip *models.UserTripSchema) *models.ExpenseServiceError
	UpdateTripParticipantTx(tx *sql.Tx, userTrip *models.UserTripSchema) *models.ExpenseServiceError
}

type TripRepository struct {
	DatabaseMgr managers.DatabaseMgr
}

func (tr *TripRepository) GetTripById(tripId *uuid.UUID) (*models.TripSchema, *models.ExpenseServiceError) {
	row := tr.DatabaseMgr.ExecuteQueryRow("SELECT id, name, description, location, start_date, end_date FROM trip WHERE id = $1", tripId)
	return rowToTripSchema(row)
}

func (tr *TripRepository) GetTripsByUserId(userId *uuid.UUID) ([]*models.TripSchema, *models.ExpenseServiceError) {
	rows, err := tr.DatabaseMgr.ExecuteQuery("SELECT t.id, t.name, t.description, t.location, t.start_date, t.end_date FROM trip t JOIN user_trip_association uta on t.id = uta.id_trip WHERE uta.id_user = $1", userId)
	if err != nil {
		log.Printf("Error while querying trips: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}
	return rowsToTripSchema(rows)
}

func (tr *TripRepository) CreateTrip(trip *models.TripSchema) *models.ExpenseServiceError {
	result, err := tr.DatabaseMgr.ExecuteStatement("INSERT INTO trip (id, name, description, location, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6)", trip.TripID, trip.Name, trip.Description, trip.Location, trip.StartDate, trip.EndDate)
	if err != nil {
		log.Printf("Error while inserting trip: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		log.Printf("Error while inserting trip: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

func (tr *TripRepository) UpdateTrip(trip *models.TripSchema) *models.ExpenseServiceError {
	// Workaround for postgres not supporting ON CONFLICT DO UPDATE
	// https://stackoverflow.com/questions/17267417/how-to-upsert-merge-insert-on-duplicate-update-in-postgresql

	// Update trip
	result, err := tr.DatabaseMgr.ExecuteStatement("UPDATE trip SET location = $1, start_date = $2, end_date = $3 WHERE id = $4", trip.Location, trip.StartDate, trip.EndDate, trip.TripID)
	if err != nil {
		log.Printf("Error while updating trip: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		log.Printf("Error while updating trip: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

func (tr *TripRepository) DeleteTrip(tripId *uuid.UUID) *models.ExpenseServiceError {
	result, err := tr.DatabaseMgr.ExecuteStatement("DELETE FROM trip WHERE id = $1", tripId)
	if err != nil {
		log.Printf("Error while deleting trip: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		log.Printf("Error while deleting trip: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

func (tr *TripRepository) AddUserToTrip(trip *models.TripSchema, invitedUserId *uuid.UUID, isCreator bool) *models.ExpenseServiceError {
	// Insert user into user_trip_association
	_, err := tr.DatabaseMgr.ExecuteStatement("INSERT INTO user_trip_association (id_user, id_trip, is_accepted, presence_start_date, presence_end_date) VALUES ($1, $2, $3, $4, $5)", invitedUserId, trip.TripID, isCreator, trip.StartDate, trip.EndDate)
	if err != nil {
		// If user is already invited, return conflict
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return expense_errors.EXPENSE_CONFLICT
		}

		log.Printf("Error while inserting user_trip_association: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

func (tr *TripRepository) AcceptTripInvite(tripId *uuid.UUID, userId *uuid.UUID) *models.ExpenseServiceError {
	result, err := tr.DatabaseMgr.ExecuteStatement("UPDATE user_trip_association SET is_accepted = $1 WHERE id_user = $2 AND id_trip = $3", true, userId, tripId)
	if err != nil {
		log.Printf("Error while updating user_trip_association: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// If no rows were affected, User already accepted the invite
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return expense_errors.EXPENSE_CONFLICT
	}

	return nil
}

func (tr *TripRepository) DeclineTripInvite(tripId *uuid.UUID, userId *uuid.UUID) *models.ExpenseServiceError {
	result, err := tr.DatabaseMgr.ExecuteStatement("DELETE FROM user_trip_association WHERE id_user = $1 AND id_trip = $2", userId, tripId)
	if err != nil {
		log.Printf("Error while deleting user_trip_association: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// If no rows were affected, User already declined the invite
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return expense_errors.EXPENSE_CONFLICT
	}

	return nil
}

func (tr *TripRepository) ValidateIfTripExists(tripId *uuid.UUID) *models.ExpenseServiceError {
	rows, err := tr.DatabaseMgr.ExecuteQuery("SELECT id FROM trip WHERE id = $1", tripId)
	if err != nil {
		log.Printf("Error while querying trip: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rows.Next() {
		return nil
	}

	return expense_errors.EXPENSE_TRIP_NOT_FOUND
}

func (tr *TripRepository) ValidateIfUserHasAccepted(tripId *uuid.UUID, userId *uuid.UUID) *models.ExpenseServiceError {
	rows, err := tr.DatabaseMgr.ExecuteQuery("SELECT id_user FROM user_trip_association WHERE id_user = $1 AND id_trip = $2 AND is_accepted = $3", userId, tripId, true)
	if err != nil {
		log.Printf("Error while querying user_trip_association: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rows.Next() {
		return nil
	}

	return expense_errors.EXPENSE_FORBIDDEN
}

func (tr *TripRepository) CheckIfUserIsInvited(tripId *uuid.UUID, userId *uuid.UUID) (bool, *models.ExpenseServiceError) {
	rows, err := tr.DatabaseMgr.ExecuteQuery("SELECT id_user FROM user_trip_association WHERE id_user = $1 AND id_trip = $2 AND is_accepted = $3", userId, tripId, false)
	if err != nil {
		log.Printf("Error while querying user_trip_association: %v", err)
		return false, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rows.Next() {
		return true, nil
	}

	return false, nil
}

func (tr *TripRepository) GetTripParticipant(tripId *uuid.UUID, userId *uuid.UUID) (*models.UserTripSchema, *models.ExpenseServiceError) {
	rows := tr.DatabaseMgr.ExecuteQueryRow("SELECT id_user, id_trip, is_accepted, presence_start_date, presence_end_date FROM user_trip_association WHERE id_user = $1 AND id_trip = $2", userId, tripId)

	var participant models.UserTripSchema
	if err := rows.Scan(&participant.UserID, &participant.TripID, &participant.HasAccepted, &participant.PresenceStartDate, &participant.PresenceEndDate); err != nil {
		log.Printf("Error while scanning user_trip_association: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return &participant, nil
}

func (tr *TripRepository) GetTripParticipants(tripId *uuid.UUID) ([]*models.UserTripSchema, *models.ExpenseServiceError) {
	rows, err := tr.DatabaseMgr.ExecuteQuery("SELECT id_user, id_trip, is_accepted, presence_start_date, presence_end_date FROM user_trip_association WHERE id_trip = $1", tripId)
	if err != nil {
		log.Printf("Error while querying user_trip_association: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	var participants []*models.UserTripSchema
	for rows.Next() {
		var participant models.UserTripSchema
		err := rows.Scan(&participant.UserID, &participant.TripID, &participant.HasAccepted, &participant.PresenceStartDate, &participant.PresenceEndDate)
		if err != nil {
			log.Printf("Error while scanning user_trip_association: %v", err)
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}
		participants = append(participants, &participant)
	}

	return participants, nil
}

func (tr *TripRepository) GetAcceptedTripParticipants(tripId *uuid.UUID) ([]*models.UserTripSchema, *models.ExpenseServiceError) {
	query := "SELECT id_user, id_trip, is_accepted, presence_start_date, presence_end_date FROM user_trip_association WHERE id_trip = $1 AND is_accepted = $2"
	rows, err := tr.DatabaseMgr.ExecuteQuery(query, tripId, true)
	if err != nil {
		log.Printf("Error while querying user_trip_association: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	participants := make([]*models.UserTripSchema, 0)
	for rows.Next() {
		var participant models.UserTripSchema
		err := rows.Scan(&participant.UserID, &participant.TripID, &participant.HasAccepted, &participant.PresenceStartDate, &participant.PresenceEndDate)
		if err != nil {
			log.Printf("Error while scanning user_trip_association: %v", err)
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}
		participants = append(participants, &participant)
	}

	return participants, nil
}

func (tr *TripRepository) UpdateTripParticipant(userTrip *models.UserTripSchema) *models.ExpenseServiceError {
	// Update user_trip_association
	result, err := tr.DatabaseMgr.ExecuteStatement("UPDATE user_trip_association SET is_accepted = $1, presence_start_date = $2, presence_end_date = $3 WHERE id_user = $4 AND id_trip = $5", userTrip.HasAccepted, userTrip.PresenceStartDate, userTrip.PresenceEndDate, userTrip.UserID, userTrip.TripID)
	if err != nil {
		log.Printf("Error while updating user_trip_association: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		log.Printf("Error while updating user_trip_association: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

func (tr *TripRepository) UpdateTripParticipantTx(tx *sql.Tx, userTrip *models.UserTripSchema) *models.ExpenseServiceError {
	query := "UPDATE user_trip_association SET is_accepted = $1, presence_start_date = $2, presence_end_date = $3 WHERE id_user = $4 AND id_trip = $5"
	_, err := tx.Exec(query, userTrip.HasAccepted, userTrip.PresenceStartDate, userTrip.PresenceEndDate, userTrip.UserID, userTrip.TripID)
	if err != nil {
		log.Printf("Error while updating user_trip_association: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

// ************************************************************
// ********************* Helper Functions *********************
// ************************************************************

// rowToTripSchema converts a row to a TripSchema
func rowToTripSchema(row *sql.Row) (*models.TripSchema, *models.ExpenseServiceError) {
	trip := models.TripSchema{}
	if err := row.Scan(&trip.TripID, &trip.Name, &trip.Description, &trip.Location, &trip.StartDate, &trip.EndDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, expense_errors.EXPENSE_TRIP_NOT_FOUND
		}

		log.Printf("Error while scanning trip: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return &trip, nil
}

// rowsToTripSchema converts a set of rows to a slice of TripSchema
func rowsToTripSchema(rows *sql.Rows) ([]*models.TripSchema, *models.ExpenseServiceError) {
	trips := make([]*models.TripSchema, 0) // It is important to initialize the slice with 0 length so that it is serialized to [] instead of null
	for rows.Next() {
		var trip models.TripSchema
		err := rows.Scan(&trip.TripID, &trip.Name, &trip.Description, &trip.Location, &trip.StartDate, &trip.EndDate)
		if err != nil {
			log.Printf("Error while scanning trip: %v", err)
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}
		trips = append(trips, &trip)
	}

	return trips, nil
}
