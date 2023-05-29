package model

import (
	"github.com/google/uuid"
)

type TripRequest struct {
	Location  string `json:"location"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type TripUpdateRequest struct {
	Location  string `json:"location"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type InviteUserRequest struct {
	Username string `json:"username"`
}

type TripsResponse struct {
	Trips []TripSchema `json:"trips"`
}

type TripResponse struct {
	TripID    *uuid.UUID `json:"tripId"`
	Location  string     `json:"location"`
	StartDate string     `json:"startDate"`
	EndDate   string     `json:"endDate"`
}

type TripCreationResponse struct {
	TripID *uuid.UUID `json:"tripId"`
}
