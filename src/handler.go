package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controller"
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
			utils.HandleErrorAndAbort(c, "invalid request body", http.StatusBadRequest, err)
			return
		}

		response, err := userCtl.RegisterUser(ctx, registrationData)
		if err != nil {
			utils.HandleErrorAndAbort(c, err.Err.Error(), err.Status, err.Err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func LoginUserHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var loginData model.LoginRequest
		if err := c.ShouldBindJSON(&loginData); err != nil {
			utils.HandleErrorAndAbort(c, "invalid request body", http.StatusBadRequest, err)
			return
		}

		response, err := userCtl.LoginUser(ctx, loginData)
		if err != nil {
			utils.HandleErrorAndAbort(c, err.Err.Error(), err.Status, err.Err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func UpdateUserHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		response, err := userCtl.UpdateUser(ctx)
		if err != nil {
			utils.HandleErrorAndAbort(c, err.Err.Error(), err.Status, err.Err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func DeleteUserHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		userIdParam := c.Param("userId")
		userId, err := uuid.Parse(userIdParam)
		if err != nil {
			utils.HandleErrorAndAbort(c, "invalid user id", http.StatusBadRequest, err)
			return
		}

		serviceError := userCtl.DeleteUser(ctx, &userId)
		if err != nil {
			utils.HandleErrorAndAbort(c, serviceError.Err.Error(), serviceError.Status, serviceError.Err)
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func ActivateUserHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		err := userCtl.ActivateUser(ctx)
		if err != nil {
			utils.HandleErrorAndAbort(c, err.Err.Error(), err.Status, err.Err)
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func GetUserDetailsHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var userId *uuid.UUID
		if err := c.ShouldBindUri(&userId); err != nil {
			utils.HandleErrorAndAbort(c, "invalid user id", http.StatusBadRequest, err)
			return
		}

		response, err := userCtl.GetUserDetails(ctx, userId)
		if err != nil {
			utils.HandleErrorAndAbort(c, err.Err.Error(), err.Status, err.Err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func SuggestUsersHandler(userCtl controller.UserCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		query := c.Param("q")
		response, err := userCtl.SuggestUsers(ctx, query)
		if err != nil {
			utils.HandleErrorAndAbort(c, err.Err.Error(), err.Status, err.Err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

/******************************************************************************************
 * TRIP ROUTES
 ******************************************************************************************/

func CreateTripEntryHandler(TripCtl controller.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		response, err := TripCtl.CreateTripEntry(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func UpdateTripEntryHandler(TripCtl controller.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		response, err := TripCtl.UpdateTripEntry(ctx, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetTripDetailsHandler(TripCtl controller.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		response, err := TripCtl.GetTripDetails(ctx, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func DeleteTripEntryHandler(TripCtl controller.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		err := TripCtl.DeleteTripEntry(ctx, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, nil)
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
