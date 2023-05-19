package controller

import (
	"context"
	"errors"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
	"github.com/google/uuid"
)

// Exposed interface to the handler-package
type TripCtl interface {
	CreateTripEntry(ctx context.Context) (*model.TripResponse, error)
	UpdateTripEntry(ctx context.Context, tripID *uuid.UUID) (*model.TripResponse, error)
	GetTripDetails(ctx context.Context, tripID *uuid.UUID) (*model.TripResponse, error)
	DeleteTripEntry(ctx context.Context, tripID *uuid.UUID) error
}

// Trip Controller structure
type TripController struct {
}

func (tc *TripController) CreateTripEntry(ctx context.Context) (*model.TripResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
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
