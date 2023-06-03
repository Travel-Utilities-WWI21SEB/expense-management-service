package controllers

import "github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"

// Exposed interface to the handler-package
type DebtCtl interface {
}

// Debt Controller structure
type DebtController struct {
	DatabaseMgr managers.DatabaseMgr
}
