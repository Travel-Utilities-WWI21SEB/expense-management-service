package model

import (
	"github.com/google/uuid"
)

type CostResponse struct {
	TripID *uuid.UUID `json:"tripId"`
	CostID *uuid.UUID `json:"costId"`
}
