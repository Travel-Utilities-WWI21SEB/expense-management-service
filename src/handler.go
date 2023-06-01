package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controller"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expenseerror"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
)

func LifeCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := model.LifeCheckResponse{
			Alive:   true,
			Version: "1.0.0",
		}

		c.JSON(http.StatusOK, response)
	}
}

/******************************************************************************************
 * TO-DO: Implement the following handlers:
 * - UserHandler
 * - TripHandler
 * - CostHandler
 ******************************************************************************************/

/******************************************************************************************
 * USER ROUTES
 ******************************************************************************************/

func RegisterUserHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var registrationData model.RegistrationRequest
		if err := c.ShouldBindJSON(&registrationData); err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		response, err := userCtl.RegisterUser(ctx, registrationData)
		if err != nil {
			// Return partial response if user was created but mail was not sent
			if err == expenseerror.EXPENSE_MAIL_NOT_SENT {
				c.JSON(http.StatusPartialContent, response)
				return
			}

			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.JSON(http.StatusCreated, response)
	}
}

func LoginUserHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var loginData model.LoginRequest
		if err := c.ShouldBindJSON(&loginData); err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
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

func RefreshTokenHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		var refreshTokenData *model.RefreshTokenRequest
		if err := c.ShouldBindJSON(&refreshTokenData); err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		id, err := utils.ValidateToken(refreshTokenData.RefreshToken)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_UNAUTHORIZED)
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

func UpdateUserHandler(userCtl controller.UserCtl) gin.HandlerFunc {
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

func DeleteUserHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO: Needs to be re-implemented after trip and cost routes are implemented
		ctx := c.Request.Context()

		userIdParam := c.Param(model.ExpenseParamKeyUserId)
		userId, err := uuid.Parse(userIdParam)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		serviceErr := userCtl.DeleteUser(ctx, &userId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func ActivateUserHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		tokenString := c.Query(model.ExpenseQueryKeyToken)
		token, err := uuid.Parse(tokenString)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		serviceErr := userCtl.ActivateUser(ctx, &token)
		if serviceErr != nil {
			// Return partial response if user was created but mail was not sent
			if serviceErr == expenseerror.EXPENSE_MAIL_NOT_SENT {
				c.JSON(http.StatusAccepted, gin.H{"message": "User successfully activated but activation mail was not sent"})
				return
			}

			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User successfully activated"})
	}
}

func GetUserDetailsHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userIdParam := c.Param(model.ExpenseParamKeyUserId)
		userId, err := uuid.Parse(userIdParam)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		response, serviceErr := userCtl.GetUserDetails(ctx, &userId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func SuggestUsersHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		query := c.Query(model.ExpenseQueryKeyQueryString)
		response, err := userCtl.SuggestUsers(ctx, query)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func CheckEmailHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Extract email from body
		var emailData *model.CheckEmailRequest
		if err := c.ShouldBindJSON(&emailData); err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		err := userCtl.CheckEmail(ctx, emailData.Email)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email is available"})
	}
}

func CheckUsernameHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Extract username from body
		var usernameData *model.CheckUsernameRequest
		if err := c.ShouldBindJSON(&usernameData); err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		err := userCtl.CheckUsername(ctx, usernameData.Username)
		if err != nil {
			utils.HandleErrorAndAbort(c, *err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Username is available"})
	}
}

/******************************************************************************************
 * TRIP ROUTES
 ******************************************************************************************/

func CreateTripEntryHandler(tripCtl controller.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var tripData model.TripRequest
		if err := c.ShouldBindJSON(&tripData); err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		response, serviceErr := tripCtl.CreateTripEntry(ctx, tripData)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetTripEntriesHandler(tripCtl controller.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		response, serviceErr := tripCtl.GetTripEntries(ctx)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetTripDetailsHandler(TripCtl controller.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get the trip id from the context
		tripIdParam := c.Param(model.ExpenseParamTripId)
		// Convert the trip id (string) to uuid
		tripId, err := uuid.Parse(tripIdParam)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		// Call the service to get the trip details
		response, serviceErr := TripCtl.GetTripDetails(ctx, &tripId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func UpdateTripEntryHandler(TripCtl controller.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		tripIdParam := c.Param(model.ExpenseParamTripId)
		tripId, err := uuid.Parse(tripIdParam)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		var tripUpdateRequest model.TripUpdateRequest
		if err := c.ShouldBindJSON(&tripUpdateRequest); err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		response, serviceErr := TripCtl.UpdateTripEntry(ctx, &tripId, tripUpdateRequest)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func DeleteTripEntryHandler(TripCtl controller.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		tripIdParam := c.Param(model.ExpenseParamTripId)
		tripId, err := uuid.Parse(tripIdParam)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		serviceErr := TripCtl.DeleteTripEntry(ctx, &tripId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func InviteUserToTripHandler(TripCtl controller.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		tripIdParam := c.Param(model.ExpenseParamTripId)
		tripId, err := uuid.Parse(tripIdParam)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		var inviteUserRequest model.InviteUserRequest
		if err := c.ShouldBindJSON(&inviteUserRequest); err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		response, serviceErr := TripCtl.InviteUserToTrip(ctx, &tripId, inviteUserRequest)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func AcceptTripInviteHandler(TripCtl controller.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		tripIdParam := c.Param(model.ExpenseParamTripId)
		tripId, err := uuid.Parse(tripIdParam)
		if err != nil {
			utils.HandleErrorAndAbort(c, *expenseerror.EXPENSE_BAD_REQUEST)
			return
		}

		serviceErr := TripCtl.AcceptTripInvite(ctx, &tripId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

/******************************************************************************************
 * COST ROUTES
 ******************************************************************************************/

func CreateCostEntryHandler(costCtl controller.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		response, err := costCtl.CreateCostEntry(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func UpdateCostEntryHandler(costCtl controller.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		response, err := costCtl.UpdateCostEntry(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetCostDetailsHandler(costCtl controller.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		response, err := costCtl.GetCostDetails(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func DeleteCostEntryHandler(costCtl controller.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		err := costCtl.DeleteCostEntry(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}
