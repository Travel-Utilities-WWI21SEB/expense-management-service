package expense_errors

import (
	models "github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
)

var (
	// EXPENSE_UPSTREAM_ERROR is used to indicate an error in 3rd party services
	EXPENSE_UPSTREAM_ERROR = &models.ExpenseServiceError{ErrorMessage: "UPSTREAM_ERROR", ErrorCode: "EM-001", Status: 500}
	// EXPENSE_INTERNAL_ERROR is used to indicate an internal, unclassified error
	EXPENSE_INTERNAL_ERROR = &models.ExpenseServiceError{ErrorMessage: "INTERNAL_ERROR", ErrorCode: "EM-002", Status: 500}
	// EXPENSE_NOT_FOUND is used to indicate that the requested resource was not found
	EXPENSE_NOT_FOUND = &models.ExpenseServiceError{ErrorMessage: "NOT_FOUND", ErrorCode: "EM-003", Status: 404}
	// EXPENSE_CONFLICT is used to indicate that the request could not be processed due to a conflict
	EXPENSE_CONFLICT = &models.ExpenseServiceError{ErrorMessage: "CONFLICT", ErrorCode: "EM-004", Status: 409}
	// EXPENSE_BAD_REQUEST is used to indicate that the request was malformed
	EXPENSE_BAD_REQUEST = &models.ExpenseServiceError{ErrorMessage: "BAD_REQUEST", ErrorCode: "EM-005", Status: 400}
	// EXPENSE_UNAUTHORIZED is used to indicate that the request was unauthorized
	EXPENSE_UNAUTHORIZED = &models.ExpenseServiceError{ErrorMessage: "UNAUTHORIZED", ErrorCode: "EM-006", Status: 401}
	// EXPENSE_CREDENTIALS_INVALID is used to indicate that the login credentials were invalid
	EXPENSE_CREDENTIALS_INVALID = &models.ExpenseServiceError{ErrorMessage: "CREDENTIALS_INVALID", ErrorCode: "EM-007", Status: 401}
	// EXPENSE_FORBIDDEN is used to indicate that the request was forbidden due to insufficient permissions
	EXPENSE_FORBIDDEN = &models.ExpenseServiceError{ErrorMessage: "FORBIDDEN", ErrorCode: "EM-008", Status: 403}
	// EXPENSE_USER_NOT_FOUND is used to indicate that the requested user was not found
	EXPENSE_USER_NOT_FOUND = &models.ExpenseServiceError{ErrorMessage: "USER_NOT_FOUND", ErrorCode: "EM-009", Status: 404}
	// EXPENSE_TRIP_NOT_FOUND is used to indicate that the requested trip was not found
	EXPENSE_TRIP_NOT_FOUND = &models.ExpenseServiceError{ErrorMessage: "TRIP_NOT_FOUND", ErrorCode: "EM-010", Status: 404}
	// EXPENSE_COST_NOT_FOUND is used to indicate that the requested cost was not found
	EXPENSE_COST_NOT_FOUND = &models.ExpenseServiceError{ErrorMessage: "COST_NOT_FOUND", ErrorCode: "EM-011", Status: 404}
	// EXPENSE_USER_EXISTS is used to indicate that the creation of a user failed because the user already exists
	EXPENSE_USER_EXISTS = &models.ExpenseServiceError{ErrorMessage: "USER_EXISTS", ErrorCode: "EM-012", Status: 409}
	// EXPENSE_USER_NOT_ACTIVATED is used to indicate that the user is not activated
	EXPENSE_USER_NOT_ACTIVATED = &models.ExpenseServiceError{ErrorMessage: "USER_NOT_ACTIVE", ErrorCode: "EM-013", Status: 403}
	// EXPENSE_MAIL_NOT_SENT is used to indicate that the mail could not be sent
	EXPENSE_MAIL_NOT_SENT = &models.ExpenseServiceError{ErrorMessage: "MAIL_NOT_SENT", ErrorCode: "EM-014", Status: 206}
	// EXPENSE_MAIL_ALREADY_VERIFIED is used to indicate that the mail was already verified
	EXPENSE_MAIL_ALREADY_VERIFIED = &models.ExpenseServiceError{ErrorMessage: "MAIL_ALREADY_VERIFIED", ErrorCode: "EM-015", Status: 409}
	// EXPENSE_ALREADY_ACCEPTED is used to indicate that the user was already accepted
	EXPENSE_ALREADY_ACCEPTED = &models.ExpenseServiceError{ErrorMessage: "ALREADY_ACCEPTED", ErrorCode: "EM-016", Status: 409}
	// EXPENSE_EMAIL_EXISTS is used to indicate that the email already exists
	EXPENSE_EMAIL_EXISTS = &models.ExpenseServiceError{ErrorMessage: "EMAIL_EXISTS", ErrorCode: "EM-017", Status: 409}
	// EXPENSE_USERNAME_EXISTS is used to indicate that the username already exists
	EXPENSE_USERNAME_EXISTS = &models.ExpenseServiceError{ErrorMessage: "USERNAME_EXISTS", ErrorCode: "EM-018", Status: 409}
	// EXPENSE_INVALID_ACTIVATION_TOKEN is used to indicate that the activation token is invalid
	EXPENSE_INVALID_ACTIVATION_TOKEN = &models.ExpenseServiceError{ErrorMessage: "INVALID_ACTIVATION_TOKEN", ErrorCode: "EM-019", Status: 400}
)
