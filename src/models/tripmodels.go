package models

import (
	"github.com/google/uuid"
)

type CreateTripRequest struct {
	Location  string `json:"location"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type UpdateTripRequest struct {
	Location  string `json:"location"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type InviteUserRequest struct {
	UserID *uuid.UUID `json:"invitedUserId"`
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
