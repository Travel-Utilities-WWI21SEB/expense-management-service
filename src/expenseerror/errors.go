package expenseerror

import (
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
)

var (
	// EXPENSE_UPSTREAM_ERROR is used to indicate an error in 3rd party services
	EXPENSE_UPSTREAM_ERROR = &model.ExpenseServiceError{ErrorMessage: "UPSTREAM_ERROR", ErrorCode: "EM-001", Status: 500}
	// EXPENSE_INTERNAL_ERROR is used to indicate an internal, unclassified error
	EXPENSE_INTERNAL_ERROR = &model.ExpenseServiceError{ErrorMessage: "INTERNAL_ERROR", ErrorCode: "EM-002", Status: 500}
	// EXPENSE_NOT_FOUND is used to indicate that the requested resource was not found
	EXPENSE_NOT_FOUND = &model.ExpenseServiceError{ErrorMessage: "NOT_FOUND", ErrorCode: "EM-003", Status: 404}
	// EXPENSE_CONFLICT is used to indicate that the request could not be processed due to a conflict
	EXPENSE_CONFLICT = &model.ExpenseServiceError{ErrorMessage: "CONFLICT", ErrorCode: "EM-004", Status: 409}
	// EXPENSE_BAD_REQUEST is used to indicate that the request was malformed
	EXPENSE_BAD_REQUEST = &model.ExpenseServiceError{ErrorMessage: "BAD_REQUEST", ErrorCode: "EM-005", Status: 400}
	// EXPENSE_UNAUTHORIZED is used to indicate that the request was unauthorized
	EXPENSE_UNAUTHORIZED = &model.ExpenseServiceError{ErrorMessage: "UNAUTHORIZED", ErrorCode: "EM-006", Status: 401}
	// EXPENSE_CREDENTIALS_INVALID is used to indicate that the login credentials were invalid
	EXPENSE_CREDENTIALS_INVALID = &model.ExpenseServiceError{ErrorMessage: "CREDENTIALS_INVALID", ErrorCode: "EM-007", Status: 401}
	// EXPENSE_FORBIDDEN is used to indicate that the request was forbidden due to insufficient permissions
	EXPENSE_FORBIDDEN = &model.ExpenseServiceError{ErrorMessage: "FORBIDDEN", ErrorCode: "EM-008", Status: 403}
	// EXPENSE_USER_NOT_FOUND is used to indicate that the requested user was not found
	EXPENSE_USER_NOT_FOUND = &model.ExpenseServiceError{ErrorMessage: "USER_NOT_FOUND", ErrorCode: "EM-009", Status: 404}
	// EXPENSE_TRIP_NOT_FOUND is used to indicate that the requested trip was not found
	EXPENSE_TRIP_NOT_FOUND = &model.ExpenseServiceError{ErrorMessage: "TRIP_NOT_FOUND", ErrorCode: "EM-010", Status: 404}
	// EXPENSE_COST_NOT_FOUND is used to indicate that the requested cost was not found
	EXPENSE_COST_NOT_FOUND = &model.ExpenseServiceError{ErrorMessage: "COST_NOT_FOUND", ErrorCode: "EM-011", Status: 404}
	// EXPENSE_USER_EXISTS is used to indicate that the creation of a user failed because the user already exists
	EXPENSE_USER_EXISTS = &model.ExpenseServiceError{ErrorMessage: "USER_EXISTS", ErrorCode: "EM-012", Status: 409}
	// EXPENSE_USER_NOT_ACTIVATED is used to indicate that the user is not activated
	EXPENSE_USER_NOT_ACTIVATED = &model.ExpenseServiceError{ErrorMessage: "USER_NOT_ACTIVE", ErrorCode: "EM-013", Status: 403}
	// EXPENSE_MAIL_NOT_SENT is used to indicate that the mail could not be sent
	EXPENSE_MAIL_NOT_SENT = &model.ExpenseServiceError{ErrorMessage: "MAIL_NOT_SENT", ErrorCode: "EM-014", Status: 500}
	// EXPENSE_MAIL_ALREADY_VERIFIED is used to indicate that the mail was already verified
	EXPENSE_MAIL_ALREADY_VERIFIED = &model.ExpenseServiceError{ErrorMessage: "MAIL_ALREADY_VERIFIED", ErrorCode: "EM-015", Status: 409}
)
