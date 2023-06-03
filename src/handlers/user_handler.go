package handlers

import (
	"net/http"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controllers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
)

func RegisterUserHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var registrationData models.RegistrationRequest
		if err := c.ShouldBindJSON(&registrationData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
			return
		}

		err := userCtl.RegisterUser(ctx, registrationData)
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

		// Generate new token and refresh token
		response, serviceErr := userCtl.RefreshToken(c.Request.Context(), id)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusCreated, response)
	}
}

func UpdateUserHandler(userCtl controllers.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		response, err := userCtl.UpdateUser(ctx)
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

		err := userCtl.CheckUsername(ctx, usernameData.Username)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.AbortWithStatus(http.StatusOK)
	}
}
