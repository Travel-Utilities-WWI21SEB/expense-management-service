package main

import (
	"net/http"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controller"
	"github.com/gin-gonic/gin"
)

// Controllers structure used to handle requests
type Controllers struct {
	UserController controller.UserCtl
	TripController controller.TripCtl
	CostController controller.CostCtl
}

func createRouter() *gin.Engine {
	router := gin.Default()
	apiv1 := router.Group("/api/v1")

	controllers := Controllers{
		UserController: &controller.UserController{},
		TripController: &controller.TripController{},
		CostController: &controller.CostController{},
	}

	router.Handle(http.MethodGet, "/lifecheck", LifeCheckHandler())

	// User Routes
	apiv1.Handle(http.MethodPost, "/user/register", RegisterUserHandler(controllers.UserController))
	apiv1.Handle(http.MethodPost, "/user/login", LoginUserHandler(controllers.UserController))
	apiv1.Handle(http.MethodPut, "/user/:userId", UpdateUserHandler(controllers.UserController))
	apiv1.Handle(http.MethodGet, "/user/:userId", GetUserDetailsHandler(controllers.UserController))
	apiv1.Handle(http.MethodDelete, "/user/:userId", DeleteUserHandler(controllers.UserController))
	apiv1.Handle(http.MethodPut, "/user/activate/:token", ActivateUserHandler(controllers.UserController))
	apiv1.Handle(http.MethodGet, "/user/suggest/:q", SuggestUsersHandler(controllers.UserController))

	// Trip Routes
	apiv1.Handle(http.MethodPost, "/trip/create", CreateTripEntryHandler(controllers.TripController))
	apiv1.Handle(http.MethodPut, "/trip/:tripId", UpdateTripEntryHandler(controllers.TripController))
	apiv1.Handle(http.MethodGet, "/trip/:tripId", GetTripDetailsHandler(controllers.TripController))
	apiv1.Handle(http.MethodDelete, "/trip/:tripId", DeleteTripEntryHandler(controllers.TripController))

	// Cost Routes
	apiv1.Handle(http.MethodPost, "/cost/create", CreateCostEntryHandler(controllers.CostController))
	apiv1.Handle(http.MethodPut, "/cost/:costId", UpdateCostEntryHandler(controllers.CostController))
	apiv1.Handle(http.MethodGet, "/cost/:costId", GetCostDetailsHandler(controllers.CostController))
	apiv1.Handle(http.MethodDelete, "/cost/:costId", DeleteCostEntryHandler(controllers.CostController))

	return router
}
