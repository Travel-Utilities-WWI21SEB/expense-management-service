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
	Location  string `json:"location"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type UpdateTripRequest struct {
	Location  string `json:"location"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
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

type TripResponse struct {
	TripID       *uuid.UUID                `json:"tripId"`
	Location     string                    `json:"location"`
	StartDate    string                    `json:"startDate"`
	EndDate      string                    `json:"endDate"`
	Participants []TripParticipantResponse `json:"participants"`
}

type TripCreationResponse struct {
	TripID *uuid.UUID `json:"tripId"`
}
