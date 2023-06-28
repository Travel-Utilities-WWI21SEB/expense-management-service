package handlers

import (
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controllers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendContactMailHandler(mailCtl controllers.MailCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var mailData models.SendContactMailRequest
		if err := c.ShouldBindJSON(&mailData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		if utils.ContainsEmptyString(mailData.Email, mailData.Message, mailData.Name) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		err := mailCtl.SendMail(ctx, &mailData)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.AbortWithStatus(http.StatusOK)
	}
}
