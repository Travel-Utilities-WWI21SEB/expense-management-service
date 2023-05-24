package manager

import (
	"context"
	"log"
	"os"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expenseerror"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/mailgun/mailgun-go/v4"
)

type MailMgr interface {
	SendActivationMail(ctx context.Context, mailData model.ActivationMail) *model.ExpenseServiceError
	SendConfirmationMail(ctx context.Context, mailData model.ConfirmationMail) *model.ExpenseServiceError
}

type MailManager struct {
	MailgunInstance *mailgun.MailgunImpl
}

const retryMailCount = 3
const emailSender = "Costventures Team <team@mail.costventures.works>"

func (mm *MailManager) SendActivationMail(ctx context.Context, mailData model.ActivationMail) *model.ExpenseServiceError {
	mailBody := utils.PrepareActivationMailBody(mailData.ActivationUrl, mailData.Username)

	// try sending mail 3 times
	for i := 0; i < retryMailCount; i++ {
		err := mm.sendMail(ctx, mailData.Recipients, mailData.Subject, mailBody)
		if err == nil {
			break
		}

		if i == retryMailCount-1 {
			log.Printf("Error in MailManager.SendActivationMail().SendMail(): %v", err.Error())
			return expenseerror.EXPENSE_MAIL_NOT_SENT
		}
	}

	return nil
}

func (mm *MailManager) SendConfirmationMail(ctx context.Context, mailData model.ConfirmationMail) *model.ExpenseServiceError {
	mailBody := utils.PrepareConfirmationMailBody(mailData.Username)

	// try sending mail 3 times
	for i := 0; i < retryMailCount; i++ {
		err := mm.sendMail(ctx, mailData.Recipients, mailData.Subject, mailBody)
		if err == nil {
			break
		}

		if i == retryMailCount-1 {
			log.Printf("Error in MailManager.SendConfirmationMail().SendMail(): %v", err.Error())
			return expenseerror.EXPENSE_MAIL_NOT_SENT
		}
	}

	return nil
}

func (mm *MailManager) sendMail(ctx context.Context, to []string, subject string, body string) error {
	message := mm.MailgunInstance.NewMessage(emailSender, subject, "", to...)
	message.AddHeader("Content-Type", "text/html")
	message.SetHtml(body)

	_, _, err := mm.MailgunInstance.Send(ctx, message)
	if err != nil {
		log.Println("Error in MailManager.SendMail().MailgunInstance.Send(): ", err.Error())
		return err
	}

	return nil
}

func InitializeMailgunClient() *mailgun.MailgunImpl {
	ApiKey := os.Getenv("MAILGUN_API_KEY")
	Domain := os.Getenv("MAILGUN_DOMAIN")

	log.Println("Initializing Mailgun client...")
	log.Println("Domain: ", Domain)

	mg := mailgun.NewMailgun(Domain, ApiKey)
	mg.SetAPIBase(mailgun.APIBaseEU)

	return mg
}
