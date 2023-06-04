package handlers

import (
	"net/http"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controllers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTripEntryHandler(tripCtl controllers.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var tripData models.TripRequest

		if err := c.ShouldBindJSON(&tripData); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
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

func GetTripEntriesHandler(tripCtl controllers.TripCtl) gin.HandlerFunc {
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

func GetTripDetailsHandler(TripCtl controllers.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get the tripId from the path
		tripId := uuid.MustParse(c.Param(models.ExpenseParamKeyTripId))

		// Call the service to get the trip details
		response, serviceErr := TripCtl.GetTripDetails(ctx, &tripId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func UpdateTripEntryHandler(TripCtl controllers.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get the tripId from the path
		tripId := uuid.MustParse(c.Param(models.ExpenseParamKeyTripId))

		var tripUpdateRequest models.TripUpdateRequest
		if err := c.ShouldBindJSON(&tripUpdateRequest); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
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

func DeleteTripEntryHandler(TripCtl controllers.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get the tripId from the path
		tripId := uuid.MustParse(c.Param(models.ExpenseParamKeyTripId))

		serviceErr := TripCtl.DeleteTripEntry(ctx, &tripId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.AbortWithStatus(http.StatusNoContent)
	}
}

func InviteUserToTripHandler(TripCtl controllers.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get the tripId from the path
		tripId := uuid.MustParse(c.Param(models.ExpenseParamKeyTripId))

		var inviteUserRequest models.InviteUserRequest
		if err := c.ShouldBindJSON(&inviteUserRequest); err != nil {
			utils.HandleErrorAndAbort(c, *expense_errors.EXPENSE_BAD_REQUEST)
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

func AcceptTripInviteHandler(TripCtl controllers.TripCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get the tripId from the path
		tripId := uuid.MustParse(c.Param(models.ExpenseParamKeyTripId))

		serviceErr := TripCtl.AcceptTripInvite(ctx, &tripId)
		if serviceErr != nil {
			utils.HandleErrorAndAbort(c, *serviceErr)
			return
		}

		c.AbortWithStatus(http.StatusNoContent)
	}
}
