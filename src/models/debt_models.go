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

// DebtOverviewDTO Data transfer object for the debt overview in the UI
type DebtOverviewDTO struct {
	Debts            []*DebtDTO `json:"debts"`
	OpenDebtAmount   string     `json:"openDebtAmount"`
	OpenCreditAmount string     `json:"openCreditAmount"`
	TotalSpent       string     `json:"totalSpent"`
	TotalReceived    string     `json:"totalReceived"`
}
