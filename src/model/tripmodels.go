package model

import (
	"github.com/google/uuid"
)

type TripsResponse struct {
	Trips []TripSchema `json:"trips"`
}

type TripResponse struct {
	TripID *uuid.UUID `json:"tripId"`
}
