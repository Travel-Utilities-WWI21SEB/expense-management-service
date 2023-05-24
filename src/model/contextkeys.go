package model

// ExpenseContextKey is the type for the context keys
type ExpenseContextKey string

const (
	// ExpenseContextKeyUserID is the key for the user id in the context
	ExpenseContextKeyUserID = ExpenseContextKey("userId")
)
