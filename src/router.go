package main

import (
	"database/sql"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/repositories"
	"net/http"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controllers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/handlers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/middlewares"
	"github.com/gin-gonic/gin"
)

// Controllers structure used to handle requests
type Controllers struct {
	UserController         controllers.UserCtl
	TripController         controllers.TripCtl
	CostController         controllers.CostCtl
	CostCategoryController controllers.CostCategoryCtl
	DebtController         controllers.DebtCtl
}

func createRouter(dbConnection *sql.DB) *gin.Engine {
	router := gin.New()

	// Attach logger middleware
	router.Use(gin.Logger())

	// Attach recovery middleware
	router.Use(gin.Recovery())

	// Configure CORS
	router.Use(middlewares.CorsMiddleware())

	apiv1 := router.Group("/api/v1")
	apiv1.Use(middlewares.ValidateUUID())

	securedApiv1 := router.Group("/api/v1")
	securedApiv1.Use(middlewares.ValidateUUID(), middlewares.JwtAuthMiddleware())

	databaseMgr := &managers.DatabaseManager{
		Connection: dbConnection,
	}

	securedTripApiv1 := securedApiv1.Group("/trips/:tripId")
	securedTripApiv1.Use(middlewares.TripValidationMiddleware(databaseMgr))

	// Initialize Mailgun client
	mgInstance := managers.InitializeMailgunClient()
	if mgInstance == nil {
		panic("Could not initialize Mailgun instance")
	}

	mailMgr := &managers.MailManager{
		MailgunInstance: mgInstance,
	}

	userRepo := &repositories.UserRepository{
		DatabaseMgr: databaseMgr,
	}

	tripRepo := &repositories.TripRepository{
		DatabaseMgr: databaseMgr,
	}

	costCategoryRepo := &repositories.CostCategoryRepository{
		DatabaseMgr: databaseMgr,
	}

	controllers := Controllers{
		UserController: &controllers.UserController{
			MailMgr:     mailMgr,
			DatabaseMgr: databaseMgr,
			UserRepo:    userRepo,
		},
		TripController: &controllers.TripController{
			DatabaseMgr: databaseMgr,
			TripRepo:    tripRepo,
			UserRepo:    userRepo,
		},
		CostCategoryController: &controllers.CostCategoryController{
			DatabaseMgr:      databaseMgr,
			CostCategoryRepo: costCategoryRepo,
		},
		CostController: &controllers.CostController{
			DatabaseMgr: databaseMgr,
		},
		DebtController: &controllers.DebtController{
			DatabaseMgr: databaseMgr,
		},
	}

	router.Handle(http.MethodGet, "/lifecheck", handlers.LifeCheckHandler())

	// User Routes
	apiv1.Handle(http.MethodPost, "/users/register", handlers.RegisterUserHandler(controllers.UserController))
	apiv1.Handle(http.MethodPost, "/users/login", handlers.LoginUserHandler(controllers.UserController))
	apiv1.Handle(http.MethodPost, "/users/refresh", handlers.RefreshTokenHandler(controllers.UserController))
	apiv1.Handle(http.MethodPost, "/users/resend-token", handlers.ResendTokenHandler(controllers.UserController))
	apiv1.Handle(http.MethodPost, "/users/activate", handlers.ActivateUserHandler(controllers.UserController))
	apiv1.Handle(http.MethodPost, "/users/check-email", handlers.CheckEmailHandler(controllers.UserController))
	apiv1.Handle(http.MethodPost, "/users/check-username", handlers.CheckUsernameHandler(controllers.UserController))
	securedApiv1.Handle(http.MethodGet, "/users/suggest", handlers.SuggestUsersHandler(controllers.UserController))
	securedApiv1.Handle(http.MethodGet, "/users", handlers.GetUserDetailsHandler(controllers.UserController))
	securedApiv1.Handle(http.MethodPatch, "/users", handlers.UpdateUserHandler(controllers.UserController))
	securedApiv1.Handle(http.MethodDelete, "/users", handlers.DeleteUserHandler(controllers.UserController))

	// Trip Routes
	securedApiv1.Handle(http.MethodPost, "/trips", handlers.CreateTripEntryHandler(controllers.TripController))
	securedApiv1.Handle(http.MethodGet, "/trips", handlers.GetTripEntriesHandler(controllers.TripController))
	securedTripApiv1.Handle(http.MethodGet, "", handlers.GetTripDetailsHandler(controllers.TripController))
	securedTripApiv1.Handle(http.MethodPatch, "", handlers.UpdateTripEntryHandler(controllers.TripController))
	securedTripApiv1.Handle(http.MethodDelete, "", handlers.DeleteTripEntryHandler(controllers.TripController))
	securedTripApiv1.Handle(http.MethodPost, "/invite", handlers.InviteUserToTripHandler(controllers.TripController))
	securedTripApiv1.Handle(http.MethodPost, "/accept", handlers.AcceptTripInviteHandler(controllers.TripController))

	// Cost Category Routes
	securedTripApiv1.Handle(http.MethodPost, "/cost-categories", handlers.CreateCostCategoryEntryHandler(controllers.CostCategoryController))
	securedTripApiv1.Handle(http.MethodGet, "/cost-categories", handlers.GetCostCategoryEntriesHandler(controllers.CostCategoryController))
	securedTripApiv1.Handle(http.MethodGet, "/cost-categories/:costCategoryId", handlers.GetCostCategoryDetailsHandler(controllers.CostCategoryController))
	securedTripApiv1.Handle(http.MethodPatch, "/cost-categories/:costCategoryId", handlers.UpdateCostCategoryEntryHandler(controllers.CostCategoryController))
	securedTripApiv1.Handle(http.MethodDelete, "/cost-categories/:costCategoryId", handlers.DeleteCostCategoryEntryHandler(controllers.CostCategoryController))

	// Cost Routes
	securedTripApiv1.Handle(http.MethodPost, "/costs", handlers.CreateCostEntryHandler(controllers.CostController))
	securedTripApiv1.Handle(http.MethodGet, "/costs", handlers.GetCostEntriesHandler(controllers.CostController))
	securedTripApiv1.Handle(http.MethodPatch, "/costs/:costId", handlers.UpdateCostEntryHandler(controllers.CostController))
	securedTripApiv1.Handle(http.MethodGet, "/costs/:costId", handlers.GetCostDetailsHandler(controllers.CostController))
	securedTripApiv1.Handle(http.MethodDelete, "/costs/:costId", handlers.DeleteCostEntryHandler(controllers.CostController))

	// Debts Routes
	securedTripApiv1.Handle(http.MethodPost, "/debts", handlers.CreateDebtHandler(controllers.DebtController))
	securedTripApiv1.Handle(http.MethodGet, "/debts", handlers.GetDebtsHandler(controllers.DebtController))
	securedTripApiv1.Handle(http.MethodGet, "/debts/:debtId", handlers.GetDebtDetailsHandler(controllers.DebtController))
	securedTripApiv1.Handle(http.MethodPatch, "/debts", handlers.UpdateDebtHandler(controllers.DebtController))

	return router
}
