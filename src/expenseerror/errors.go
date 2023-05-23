package expenseerror

import (
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
)

var (
	EXPENSE_UPSTREAM_ERROR      = &model.ExpenseServiceError{ErrorMessage: "UPSTREAM_ERROR", ErrorCode: "EM-001", Status: 500}
	EXPENSE_INTERNAL_ERROR      = &model.ExpenseServiceError{ErrorMessage: "INTERNAL_ERROR", ErrorCode: "EM-002", Status: 500}
	EXPENSE_NOT_FOUND           = &model.ExpenseServiceError{ErrorMessage: "NOT_FOUND", ErrorCode: "EM-003", Status: 404}
	EXPENSE_CONFLICT            = &model.ExpenseServiceError{ErrorMessage: "CONFLICT", ErrorCode: "EM-004", Status: 409}
	EXPENSE_BAD_REQUEST         = &model.ExpenseServiceError{ErrorMessage: "BAD_REQUEST", ErrorCode: "EM-005", Status: 400}
	EXPENSE_UNAUTHORIZED        = &model.ExpenseServiceError{ErrorMessage: "UNAUTHORIZED", ErrorCode: "EM-006", Status: 401}
	EXPENSE_CREDENTIALS_INVALID = &model.ExpenseServiceError{ErrorMessage: "CREDENTIALS_INVALID", ErrorCode: "EM-007", Status: 401}
	EXPENSE_FORBIDDEN           = &model.ExpenseServiceError{ErrorMessage: "FORBIDDEN", ErrorCode: "EM-008", Status: 403}
	EXPENSE_USER_NOT_FOUND      = &model.ExpenseServiceError{ErrorMessage: "USER_NOT_FOUND", ErrorCode: "EM-009", Status: 404}
	EXPENSE_TRIP_NOT_FOUND      = &model.ExpenseServiceError{ErrorMessage: "TRIP_NOT_FOUND", ErrorCode: "EM-010", Status: 404}
	EXPENSE_COST_NOT_FOUND      = &model.ExpenseServiceError{ErrorMessage: "COST_NOT_FOUND", ErrorCode: "EM-011", Status: 404}
	EXPENSE_USER_EXISTS         = &model.ExpenseServiceError{ErrorMessage: "USER_EXISTS", ErrorCode: "EM-012", Status: 409}
	EXPENSE_USER_NOT_ACTIVE     = &model.ExpenseServiceError{ErrorMessage: "USER_NOT_ACTIVE", ErrorCode: "EM-013", Status: 403}
)
