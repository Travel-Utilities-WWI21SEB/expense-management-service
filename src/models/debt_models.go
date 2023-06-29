package models

import "github.com/google/uuid"

// DebtDTO Data transfer object for debt entries
type DebtDTO struct {
	DebtID       *uuid.UUID `json:"debtId"`
	CreditorId   *uuid.UUID `json:"creditorId"`
	DebtorId     *uuid.UUID `json:"debtorId"`
	TripId       *uuid.UUID `json:"tripId"`
	Amount       string     `json:"amount"`
	CurrencyCode string     `json:"currency"`
	CreationDate string     `json:"createdAt"`
	UpdateDate   string     `json:"updatedAt"`
}
