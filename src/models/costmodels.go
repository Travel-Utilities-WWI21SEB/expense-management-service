package models

import "github.com/google/uuid"

// CostDTO Data transfer object for cost entries
type CostDTO struct {
	CostID         *uuid.UUID     `json:"costId"`
	CostCategoryID *uuid.UUID     `json:"costCategoryId"`
	Amount         string         `json:"amount"`
	CurrencyCode   string         `json:"currency"`
	Description    string         `json:"description"`
	CreationDate   string         `json:"createdAt"`
	DeductionDate  string         `json:"deductedAt"`
	EndDate        string         `json:"endDate"`
	Creditor       string         `json:"creditor"`
	Contributors   []*Contributor `json:"contributors"`
}

type Contributor struct {
	Username string `json:"username"`
	Amount   string `json:"amount"`
}
