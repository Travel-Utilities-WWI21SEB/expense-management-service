package models

type ExpenseServiceError struct {
	ErrorMessage string
	ErrorCode    string
	Status       int
}
