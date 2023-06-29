package controllers

import "github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"

// DebtCtl Exposed interface to the handler-package
type DebtCtl interface {
}

// DebtController Debt Controller structure
type DebtController struct {
	DatabaseMgr managers.DatabaseMgr
}
