package main

import (
	"database/sql"
	"net/http"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controller"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/manager"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/middleware"
	"github.com/gin-gonic/gin"
)

// Controllers structure used to handle requests
type Controllers struct {
	UserController controller.UserCtl
	TripController controller.TripCtl
	CostController controller.CostCtl
}

func createRouter(dbConnection *sql.DB) *gin.Engine {
	router := gin.Default()
	apiv1 := router.Group("/api/v1")
	securedApiv1 := router.Group("/api/v1")

	securedApiv1.Use(middleware.JwtAuthMiddleware())

	databaseMgr := &manager.DatabaseManager{
		Connection: dbConnection,
	}

	// Initialize Mailgun client
	mgInstance := manager.InitializeMailgunClient()
	if mgInstance == nil {
		panic("Could not initialize Mailgun instance")
	}

	mailMgr := &manager.MailManager{
		MailgunInstance: mgInstance,
	}

	controllers := Controllers{
		UserController: &controller.UserController{
			MailMgr:     mailMgr,
			DatabaseMgr: databaseMgr,
		},
		TripController: &controller.TripController{
			DatabaseMgr: databaseMgr,
		},
		CostController: &controller.CostController{},
	}

	router.Handle(http.MethodGet, "/lifecheck", LifeCheckHandler())

	// User Routes
	apiv1.Handle(http.MethodPost, "/user/register", RegisterUserHandler(controllers.UserController))
	apiv1.Handle(http.MethodPost, "/user/login", LoginUserHandler(controllers.UserController))
	apiv1.Handle(http.MethodPost, "/refresh", RefreshTokenHandler(controllers.UserController))
	apiv1.Handle(http.MethodPost, "/user/activate", ActivateUserHandler(controllers.UserController))
	apiv1.Handle(http.MethodPost, "user/check-email", CheckEmailHandler(controllers.UserController))
	apiv1.Handle(http.MethodPost, "user/check-username", CheckUsernameHandler(controllers.UserController))
	securedApiv1.Handle(http.MethodGet, "/user/suggest", SuggestUsersHandler(controllers.UserController))
	securedApiv1.Handle(http.MethodGet, "/user/:userId", GetUserDetailsHandler(controllers.UserController))
	securedApiv1.Handle(http.MethodPatch, "/user/:userId", UpdateUserHandler(controllers.UserController))
	securedApiv1.Handle(http.MethodDelete, "/user/:userId", DeleteUserHandler(controllers.UserController))

	// Trip Routes
	securedApiv1.Handle(http.MethodPost, "/trips", CreateTripEntryHandler(controllers.TripController))
	securedApiv1.Handle(http.MethodGet, "/trips", GetTripEntriesHandler(controllers.TripController))
	securedApiv1.Handle(http.MethodGet, "/trips/:tripId", GetTripDetailsHandler(controllers.TripController))
	securedApiv1.Handle(http.MethodPatch, "/trips/:tripId", UpdateTripEntryHandler(controllers.TripController))
	securedApiv1.Handle(http.MethodDelete, "/trips/:tripId", DeleteTripEntryHandler(controllers.TripController))
	securedApiv1.Handle(http.MethodPost, "/trips/:tripId/invite", InviteUserToTripHandler(controllers.TripController))
	securedApiv1.Handle(http.MethodPost, "/trips/:tripId/accept", AcceptTripInviteHandler(controllers.TripController))

	// Cost Routes
	securedApiv1.Handle(http.MethodPost, "/cost/create", CreateCostEntryHandler(controllers.CostController))
	securedApiv1.Handle(http.MethodPatch, "/cost/:costId", UpdateCostEntryHandler(controllers.CostController))
	securedApiv1.Handle(http.MethodGet, "/cost/:costId", GetCostDetailsHandler(controllers.CostController))
	securedApiv1.Handle(http.MethodDelete, "/cost/:costId", DeleteCostEntryHandler(controllers.CostController))

	return router
}
