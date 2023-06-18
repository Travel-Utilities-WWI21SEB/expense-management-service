package models

import (
	"github.com/google/uuid"
)

type TripParticipantResponse struct {
	Username          string `json:"username"`
	HasAcceptedInvite bool   `json:"hasAcceptedInvite"`
	PresenceStartDate string `json:"presenceStartDate"`
	PresenceEndDate   string `json:"presenceEndDate"`
}

type CreateTripRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}

type UpdateTripRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}

type UpdateTripParticipantRequest struct {
	HasAcceptedInvite bool   `json:"hasAcceptedInvite"`
	PresenceStartDate string `json:"presenceStartDate"`
	PresenceEndDate   string `json:"presenceEndDate"`
}

type InviteUserRequest struct {
	Username string `json:"invitedUsername"`
	EMail    string `json:"invitedEmail"`
}

type TripDetailsResponse struct {
	TripID         *uuid.UUID                `json:"tripId"`
	Name           string                    `json:"name"`
	Description    string                    `json:"description"`
	Location       string                    `json:"location"`
	StartDate      string                    `json:"startDate"`
	EndDate        string                    `json:"endDate"`
	Participants   []TripParticipantResponse `json:"participants"`
	CostCategories []CostCategoryResponse    `json:"costCategories"`
	Costs          []CostDetailsResponse     `json:"costs"`
}

type TripResponse struct {
	TripID         *uuid.UUID                `json:"tripId"`
	Name           string                    `json:"name"`
	Description    string                    `json:"description"`
	Location       string                    `json:"location"`
	StartDate      string                    `json:"startDate"`
	EndDate        string                    `json:"endDate"`
	TotalCost      string                    `json:"totalCost"`
	UserDebt       string                    `json:"userDebt"`   // How much the user owes
	UserCredit     string                    `json:"userCredit"` // How much the user is owed
	CostCategories []CostCategoryResponse    `json:"costCategories"`
	Participants   []TripParticipantResponse `json:"participants"`
}
