package models

import (
	"github.com/google/uuid"
)

// CreateCostRequest is the request body for creating a cost
type CreateCostRequest struct {
	CostCategoryID *uuid.UUID `json:"costCategoryId"`
	Cost           float64    `json:"cost"`
}

// CreateCostResponse is the response body for creating a cost
type CreateCostResponse struct {
	CostID *uuid.UUID `json:"costId"`
}

// CostAmount shows the amount of a cost per person
type CostAmount struct {
	UserID       *uuid.UUID `json:"userId"`
	Amount       float64    `json:"amount"`
	CurrencyCode string     `json:"currency"`
}

// CostDetailsResponse is the response body for getting a cost
type CostDetailsResponse struct {
	CostCategoryID *uuid.UUID   `json:"costCategoryId"`
	Amount         float64      `json:"amount"`
	CurrencyCode   string       `json:"currency"`
	Icon           string       `json:"icon"`
	Color          string       `json:"color"`
	PaidBy         *uuid.UUID   `json:"paidBy"`
	PaidFor        []CostAmount `json:"paidFor"`
}

// TripCostsResponse is the response body for getting all costs of a trip
type TripCostsResponse struct {
	Costs []CostDetailsResponse `json:"costs"`
}

type CostResponse struct {
	TripID *uuid.UUID `json:"tripId"`
	CostID *uuid.UUID `json:"costId"`
}
