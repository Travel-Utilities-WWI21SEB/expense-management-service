package controllers

import (
	"context"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
)

// MailCtl exposed interface to the handler-package
type MailCtl interface {
	SendMail(ctx context.Context, mailData *models.SendContactMailRequest) *models.ExpenseServiceError
}

// MailController Mail Controller struct
type MailController struct {
	MailMgr *managers.MailManager
}

func (mc *MailController) SendMail(ctx context.Context, mailData *models.SendContactMailRequest) *models.ExpenseServiceError {
	return mc.MailMgr.SendContactMail(ctx, mailData)
}
