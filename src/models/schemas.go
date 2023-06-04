package models

import (
	"time"

	"github.com/google/uuid"
)

type CostSchema struct {
	CostID         *uuid.UUID `json:"costId" db:"id"`
	Amount         float32    `json:"amount" db:"amount"`
	CurrencyCode   string     `json:"currencyCode" db:"currency_code"`
	CreationDate   *time.Time `json:"createdAt" db:"created_at"`
	DeductionDate  *time.Time `json:"deductedAt" db:"deducted_at"`
	EndDate        *time.Time `json:"endDate,omitempty" db:"end_date"`
	CostCategoryID *uuid.UUID `json:"costCategoryId" db:"id_cost_category"`
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
	UserID   *uuid.UUID `json:"userId" db:"id"`
	UserName string     `json:"userName" db:"username"`
	Email    string     `json:"email" db:"email"`
	Password string     `json:"password" db:"password"`
}

type UserCostSchema struct {
	UserID     *uuid.UUID `json:"userId" db:"id_user"`
	CostID     *uuid.UUID `json:"costId" db:"id_cost"`
	IsCreditor bool       `json:"isCreditor" db:"is_creditor"`
}

type TripUserSchema struct {
	UserID    *uuid.UUID `json:"id_user" db:"id_user"`
	TripID    *uuid.UUID `json:"id_trip" db:"id_trip"`
	Accepted  bool       `json:"accepted" db:"accepted"`
	StartDate *time.Time `json:"startDate" db:"attendance_start_date"`
	EndDate   *time.Time `json:"endDate" db:"attendance_end_date"`
}

type Transaction struct {
	TransactionID *uuid.UUID `json:"transactionId" db:"id"`
	Amount        float32    `json:"amount" db:"amount"`
	CurrencyCode  string     `json:"currencyCode" db:"currency"`
	TransactionAt string     `json:"transactionAt" db:"transaction_at"`
	Sender        *uuid.UUID `json:"senderId" db:"sender_id"`
	Receiver      *uuid.UUID `json:"receiverId" db:"receiver_id"`
}
