package models

import (
	"github.com/google/uuid"
)

type TripDTO struct {
	TripID         *uuid.UUID             `json:"tripId"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Location       string                 `json:"location"`
	StartDate      string                 `json:"startDate"`
	EndDate        string                 `json:"endDate"`
	TotalCost      string                 `json:"totalCost,omitempty"`
	UserDebt       string                 `json:"userDebt,omitempty"`   // How much the user owes
	UserCredit     string                 `json:"userCredit,omitempty"` // How much the user is owed
	CostCategories []CostCategoryResponse `json:"costCategories"`
	Participants   []TripParticipationDTO `json:"participants"`
}

type TripParticipationDTO struct {
	UserID            *uuid.UUID `json:"userId"`
	Username          string     `json:"username"`
	HasAcceptedInvite bool       `json:"hasAcceptedInvite"`
	PresenceStartDate string     `json:"presenceStartDate"`
	PresenceEndDate   string     `json:"presenceEndDate"`
}

type SlimTripDTO struct {
	TripID      *uuid.UUID `json:"tripId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Location    string     `json:"location"`
	StartDate   string     `json:"startDate"`
	EndDate     string     `json:"endDate"`
}
