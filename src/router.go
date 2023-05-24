package main

import (
	"net/http"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controller"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/middleware"
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
	securedApiv1 := router.Group("/api/v1")

	securedApiv1.Use(middleware.JwtAuthMiddleware())

	controllers := Controllers{
		UserController: &controller.UserController{},
		TripController: &controller.TripController{},
		CostController: &controller.CostController{},
	}

	router.Handle(http.MethodGet, "/lifecheck", LifeCheckHandler())

	// User Routes
	apiv1.Handle(http.MethodPost, "/user/register", RegisterUserHandler(controllers.UserController))
	apiv1.Handle(http.MethodPost, "/user/login", LoginUserHandler(controllers.UserController))
	apiv1.Handle(http.MethodPost, "/user/activate", ActivateUserHandler(controllers.UserController))
	securedApiv1.Handle(http.MethodGet, "/user/suggest", SuggestUsersHandler(controllers.UserController))
	securedApiv1.Handle(http.MethodGet, "/user/:userId", GetUserDetailsHandler(controllers.UserController))
	securedApiv1.Handle(http.MethodPatch, "/user/:userId", UpdateUserHandler(controllers.UserController))
	securedApiv1.Handle(http.MethodDelete, "/user/:userId", DeleteUserHandler(controllers.UserController))

	// Trip Routes
	securedApiv1.Handle(http.MethodPost, "/trip/create", CreateTripEntryHandler(controllers.TripController))
	securedApiv1.Handle(http.MethodPatch, "/trip/:tripId", UpdateTripEntryHandler(controllers.TripController))
	securedApiv1.Handle(http.MethodGet, "/trip/:tripId", GetTripDetailsHandler(controllers.TripController))
	securedApiv1.Handle(http.MethodDelete, "/trip/:tripId", DeleteTripEntryHandler(controllers.TripController))

	// Cost Routes
	securedApiv1.Handle(http.MethodPost, "/cost/create", CreateCostEntryHandler(controllers.CostController))
	securedApiv1.Handle(http.MethodPatch, "/cost/:costId", UpdateCostEntryHandler(controllers.CostController))
	securedApiv1.Handle(http.MethodGet, "/cost/:costId", GetCostDetailsHandler(controllers.CostController))
	securedApiv1.Handle(http.MethodDelete, "/cost/:costId", DeleteCostEntryHandler(controllers.CostController))

	return router
}
