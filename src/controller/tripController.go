package controller

import (
	"context"
	"errors"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
	"github.com/google/uuid"
)

// Exposed interface to the handler-package
type TripCtl interface {
	createTripEntry(ctx context.Context) (*model.TripResponse, error)
	updateTripEntry(ctx context.Context, tripID *uuid.UUID) (*model.TripResponse, error)
	getTripDetails(ctx context.Context, tripID *uuid.UUID) (*model.TripResponse, error)
	getTripsByUser(ctx context.Context, userID *uuid.UUID) (*model.TripResponse, error)
	deleteTripEntry(ctx context.Context) error
}

// Cost Controller structure
type TripController struct {
}

func (tc *TripController) createTripEntry(ctx context.Context) (*model.TripResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (tc *TripController) updateTripEntry(ctx context.Context, tripID *uuid.UUID) (*model.TripResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (tc *TripController) getTripDetails(ctx context.Context, tripID *uuid.UUID) (*model.TripResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (tc *TripController) getTripsByUser(ctx context.Context, userID *uuid.UUID) (*model.TripResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (tc *TripController) deleteTripEntry(ctx context.Context, tripID *uuid.UUID) (*model.TripResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}
