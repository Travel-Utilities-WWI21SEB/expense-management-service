package models

import "github.com/google/uuid"

// DebtDTO Data transfer object for debt entries
type DebtDTO struct {
	DebtID       *uuid.UUID   `json:"debtId"`
	Creditor     *UserDto     `json:"creditor"`
	Debtor       *UserDto     `json:"debtor"`
	Trip         *SlimTripDTO `json:"trip"`
	Amount       string       `json:"amount"`
	CurrencyCode string       `json:"currency"`
	CreationDate string       `json:"createdAt"`
	UpdateDate   string       `json:"updatedAt"`
}
