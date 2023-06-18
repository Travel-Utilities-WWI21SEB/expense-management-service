package models

import (
	"github.com/google/uuid"
	"time"
)

type Contributor struct {
	Username string `json:"username"`
	Amount   string `json:"amount"`
}

// CreateCostRequest is the request body for creating a cost
type CreateCostRequest struct {
	CostCategoryID *uuid.UUID     `json:"costCategoryId"`
	Amount         string         `json:"amount"`
	CurrencyCode   string         `json:"currency"`
	Description    string         `json:"description"`
	DeductedAt     *time.Time     `json:"deductedAt"`
	EndDate        *time.Time     `json:"endDate"`
	Creditor       string         `json:"creditor"`
	Contributors   []*Contributor `json:"contributors"`
}

// CostDetailsResponse is the response body for getting a cost
type CostDetailsResponse struct {
	CostID         *uuid.UUID     `json:"costId"`
	CreationDate   *time.Time     `json:"createdAt"`
	CostCategoryID *uuid.UUID     `json:"costCategoryId"`
	Amount         string         `json:"amount"`
	CurrencyCode   string         `json:"currency"`
	Description    string         `json:"description"`
	DeductionDate  *time.Time     `json:"deductedAt"`
	EndDate        *time.Time     `json:"endDate"`
	Creditor       string         `json:"creditor"`
	Contributors   []*Contributor `json:"contributors"`
}
