package handlers

import (
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controllers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net/http"
)

func RegisterUserHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get the multipart form
		log.Printf("RegisterUserHandler: Getting multipart form")
		if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		form := c.Request.MultipartForm

		var registrationData models.RegistrationRequest
		// Read the form data into the registrationData struct
		log.Printf("RegisterUserHandler: Reading form data into struct")
		if err := c.ShouldBindWith(&registrationData, binding.Form); err != nil {
			log.Printf("RegisterUserHandler: Error while binding form data to struct")
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		} else {
			log.Printf("RegisterUserHandler: Successfully bound form data to struct")
		}

		log.Printf("RegisterUserHandler: Checking if registration data is empty")

		if utils.ContainsEmptyString(registrationData.Username, registrationData.Password, registrationData.Email,
			registrationData.FirstName, registrationData.LastName, registrationData.Birthday, registrationData.Location) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		err := userCtl.RegisterUser(ctx, registrationData, form)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.AbortWithStatus(http.StatusCreated)
	}
}

func LoginUserHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var loginData models.LoginRequest
		if err := c.ShouldBindJSON(&loginData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		if utils.ContainsEmptyString(loginData.Password, loginData.Email) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
		}

		response, err := userCtl.LoginUser(ctx, loginData)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func RefreshTokenHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		var refreshTokenData *models.RefreshTokenRequest
		if err := c.ShouldBindJSON(&refreshTokenData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		id, err := utils.ValidateToken(refreshTokenData.RefreshToken)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_UNAUTHORIZED)
			return
		}

		if utils.ContainsEmptyString(id.String()) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		// Generate new token and refresh token
		response, serviceErr := userCtl.RefreshToken(c.Request.Context(), id)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func ForgotPasswordHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var forgotPasswordData *models.ForgotPasswordRequest
		if err := c.ShouldBindJSON(&forgotPasswordData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		if utils.ContainsEmptyString(forgotPasswordData.Email) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		err := userCtl.ForgotPassword(ctx, forgotPasswordData.Email)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.AbortWithStatus(http.StatusOK)
	}
}

func VerifyPasswordResetTokenHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var verifyResetTokenData *models.VerifyPasswordResetTokenRequest
		if err := c.ShouldBindJSON(&verifyResetTokenData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		if utils.ContainsEmptyString(verifyResetTokenData.Token) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		err := userCtl.VerifyPasswordResetToken(ctx, verifyResetTokenData.Email, verifyResetTokenData.Token)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.AbortWithStatus(http.StatusOK)
	}
}

func ResetPasswordHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var resetPasswordData *models.ResetPasswordRequest
		if err := c.ShouldBindJSON(&resetPasswordData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		if utils.ContainsEmptyString(resetPasswordData.Email, resetPasswordData.Token, resetPasswordData.Password) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		err := userCtl.ResetPassword(ctx, resetPasswordData.Email, resetPasswordData.Password, resetPasswordData.Token)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.AbortWithStatus(http.StatusOK)
	}
}

func ResendTokenHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var resendTokenData *models.ResendTokenRequest
		if err := c.ShouldBindJSON(&resendTokenData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		if utils.ContainsEmptyString(resendTokenData.Email) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		err := userCtl.ResendToken(ctx, resendTokenData.Email)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.AbortWithStatus(http.StatusOK)
	}
}

func UpdateUserHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var updateData *models.UpdateUserRequest
		if err := c.ShouldBindJSON(&updateData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		response, err := userCtl.UpdateUser(ctx, updateData)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func DeleteUserHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO: Needs to be re-implemented after trip and cost routes are implemented
		ctx := c.Request.Context()

		serviceErr := userCtl.DeleteUser(ctx)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.AbortWithStatus(http.StatusNoContent)
	}
}

func ActivateUserHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		token := c.Query(models.ExpenseQueryKeyToken)

		if utils.ContainsEmptyString(token) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		response, serviceErr := userCtl.ActivateUser(ctx, token)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetUserDetailsHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		response, serviceErr := userCtl.GetUserDetails(ctx)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func SuggestUsersHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		query := c.Query(models.ExpenseQueryKeyQueryString)
		if utils.ContainsEmptyString(query) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		response, err := userCtl.SuggestUsers(ctx, query)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func CheckEmailHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Extract email from body
		var emailData *models.CheckEmailRequest
		if err := c.ShouldBindJSON(&emailData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		if utils.ContainsEmptyString(emailData.Email) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		err := userCtl.CheckEmail(ctx, emailData.Email)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.AbortWithStatus(http.StatusOK)
	}
}

func CheckUsernameHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Extract username from body
		var usernameData *models.CheckUsernameRequest
		if err := c.ShouldBindJSON(&usernameData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		if utils.ContainsEmptyString(usernameData.Username) {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		err := userCtl.CheckUsername(ctx, usernameData.Username)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.AbortWithStatus(http.StatusOK)
	}
}
