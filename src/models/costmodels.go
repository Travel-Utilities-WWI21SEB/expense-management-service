package models

import (
	"github.com/google/uuid"
	"time"
)

type PaidFor struct {
	UserID *uuid.UUID `json:"userId"`
	Amount string     `json:"amount"`
}

// CreateCostRequest is the request body for creating a cost
type CreateCostRequest struct {
	CostCategoryID *uuid.UUID `json:"costCategoryId"`
	Amount         string     `json:"amount"`
	CurrencyCode   string     `json:"currency"`
	Description    string     `json:"description"`
	DeductedAt     *time.Time `json:"deductedAt"`
	EndDate        *time.Time `json:"endDate"`
	PaidBy         *uuid.UUID `json:"paidBy"`
	PaidFor        []PaidFor  `json:"paidFor"`
}

// CostAmount shows the amount of a cost per person
type CostAmount struct {
	UserID       *uuid.UUID `json:"userId"`
	Amount       string     `json:"amount"`
	CurrencyCode string     `json:"currency"`
}

// CostDetailsResponse is the response body for getting a cost
type CostDetailsResponse struct {
	CostID         *uuid.UUID `json:"costId"`
	CostCategoryID *uuid.UUID `json:"costCategoryId"`
	Amount         string     `json:"amount"`
	CurrencyCode   string     `json:"currency"`
	Description    string     `json:"description"`
	CreationDate   *time.Time `json:"createdAt"`
	DeductionDate  *time.Time `json:"deductedAt"`
	EndDate        *time.Time `json:"endDate"`
	PaidBy         *uuid.UUID `json:"paidBy"`
	PaidFor        []PaidFor  `json:"paidFor"`
}

// TripCostsResponse is the response body for getting all costs of a trip
type TripCostsResponse struct {
	Costs []CostDetailsResponse `json:"costs"`
}

type CostResponse struct {
	TripID *uuid.UUID `json:"tripId"`
	CostID *uuid.UUID `json:"costId"`
}
