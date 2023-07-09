package models

import "github.com/google/uuid"

type TransactionDTO struct {
	TransactionId *uuid.UUID   `json:"transactionId"`
	Creditor      *UserDto     `json:"creditor"`
	Debtor        *UserDto     `json:"debtor"`
	Trip          *SlimTripDTO `json:"trip"`
	// CreditorId    *uuid.UUID   `json:"creditorId,omitempty"` // Creditor is always context user
	DebtorId     *uuid.UUID `json:"debtorId,omitempty"`
	Amount       string     `json:"amount"`
	CreationDate string     `json:"createdAt"`
	IsConfirmed  bool       `json:"isConfirmed"`
}

type TransactionQueryParams struct {
	DebtorId         *uuid.UUID `json:"debtorId,omitempty"`
	DebtorUsername   string     `json:"debtorUsername,omitempty"`
	CreditorId       *uuid.UUID `json:"creditorId,omitempty"`
	CreditorUsername string     `json:"creditorUsername,omitempty"`
	SortBy           string     `json:"sortBy,omitempty"`
	SortOrder        string     `json:"order,omitempty"`
	IsConfirmed      *bool      `json:"isConfirmed,omitempty"`
}
