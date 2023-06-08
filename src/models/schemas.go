package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CostSchema struct {
	CostID         *uuid.UUID      `json:"costId" db:"id"`
	Amount         decimal.Decimal `json:"amount" db:"amount"`
	Description    string          `json:"description" db:"description"`
	CreationDate   *time.Time      `json:"createdAt" db:"created_at"`
	DeductionDate  *time.Time      `json:"deductedAt" db:"deducted_at"`
	EndDate        *time.Time      `json:"endDate,omitempty" db:"end_date"`
	CostCategoryID *uuid.UUID      `json:"costCategoryId" db:"id_cost_category"`
}

type CostCategorySchema struct {
	CostCategoryID *uuid.UUID `json:"costCategoryId" db:"id"`
	Name           string     `json:"name" db:"name"`
	Description    string     `json:"description" db:"description"`
	Icon           string     `json:"icon" db:"icon"`
	Color          string     `json:"color" db:"color"`
	TripID         *uuid.UUID `json:"tripId" db:"id_trip"`
}

type TripSchema struct {
	TripID    *uuid.UUID `json:"tripId" db:"id"`
	Location  string     `json:"location" db:"location"`
	StartDate *time.Time `json:"startDate" db:"start_date"`
	EndDate   *time.Time `json:"endDate" db:"end_date"`
}

type UserSchema struct {
	UserID    *uuid.UUID `json:"userId" db:"id"`
	Username  string     `json:"userName" db:"username"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"password" db:"password"`
	Activated bool       `json:"activated" db:"activated"`
}

type ActivationTokenSchema struct {
	UserID      *uuid.UUID `json:"userId" db:"id_user"`
	Token       *string    `json:"token" db:"token"`
	CreatedAt   *time.Time `json:"createdAt" db:"created_at"`
	ConfirmedAt *time.Time `json:"confirmedAt" db:"confirmed_at"`
	ExpiresAt   *time.Time `json:"expiresAt" db:"expires_at"`
}

type UserCostSchema struct {
	UserID     *uuid.UUID `json:"userId" db:"id_user"`
	CostID     *uuid.UUID `json:"costId" db:"id_cost"`
	IsCreditor bool       `json:"isCreditor" db:"is_creditor"`
}

type UserTripSchema struct {
	UserID            *uuid.UUID `json:"id_user" db:"id_user"`
	TripID            *uuid.UUID `json:"id_trip" db:"id_trip"`
	HasAccepted       bool       `json:"accepted" db:"is_accepted"`
	PresenceStartDate *time.Time `json:"startDate" db:"presence_start_date"`
	PresenceEndDate   *time.Time `json:"endDate" db:"presence_end_date"`
}

type Transaction struct {
	TransactionID *uuid.UUID `json:"transactionId" db:"id"`
	Amount        float32    `json:"amount" db:"amount"`
	CurrencyCode  string     `json:"currencyCode" db:"currency"`
	TransactionAt string     `json:"transactionAt" db:"transaction_at"`
	Sender        *uuid.UUID `json:"senderId" db:"sender_id"`
	Receiver      *uuid.UUID `json:"receiverId" db:"receiver_id"`
}