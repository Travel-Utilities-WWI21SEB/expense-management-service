package main

import (
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
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
	MailController         controllers.MailCtl
	TransactionController  controllers.TransactionCtl
}

func createRouter(dbConnection *pgxpool.Pool) *gin.Engine {
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

	costRepo := &repositories.CostRepository{
		DatabaseMgr: databaseMgr,
	}

	debtRepo := &repositories.DebtRepository{
		DatabaseMgr: databaseMgr,
	}

	transactionRepo := &repositories.TransactionRepository{
		DatabaseMgr: databaseMgr,
	}

	controller := Controllers{
		UserController: &controllers.UserController{
			MailMgr:     mailMgr,
			DatabaseMgr: databaseMgr,
			UserRepo:    userRepo,
		},
		TripController: &controllers.TripController{
			DatabaseMgr:      databaseMgr,
			TripRepo:         tripRepo,
			UserRepo:         userRepo,
			CostRepo:         costRepo,
			CostCategoryRepo: costCategoryRepo,
			DebtRepo:         debtRepo,
		},
		CostCategoryController: &controllers.CostCategoryController{
			DatabaseMgr:      databaseMgr,
			CostCategoryRepo: costCategoryRepo,
			CostRepo:         costRepo,
		},
		CostController: &controllers.CostController{
			DatabaseMgr:      databaseMgr,
			CostRepo:         costRepo,
			UserRepo:         userRepo,
			TripRepo:         tripRepo,
			CostCategoryRepo: costCategoryRepo,
			DebtRepo:         debtRepo,
		},
		DebtController: &controllers.DebtController{
			DatabaseMgr: databaseMgr,
			DebtRepo:    debtRepo,
			UserRepo:    userRepo,
			TripRepo:    tripRepo,
		},
		TransactionController: &controllers.TransactionController{
			DatabaseMgr:     databaseMgr,
			TransactionRepo: transactionRepo,
			UserRepo:        userRepo,
			TripRepo:        tripRepo,
			DebtRepo:        debtRepo,
		},
		MailController: &controllers.MailController{
			MailMgr: mailMgr,
		},
	}

	router.Handle(http.MethodGet, "/lifecheck", handlers.LifeCheckHandler())
	apiv1.Handle(http.MethodPost, "/send-email", handlers.SendContactMailHandler(controller.MailController))

	// User Routes
	apiv1.Handle(http.MethodPost, "/users/register", handlers.RegisterUserHandler(controller.UserController))
	apiv1.Handle(http.MethodPost, "/users/login", handlers.LoginUserHandler(controller.UserController))
	apiv1.Handle(http.MethodPost, "/users/refresh", handlers.RefreshTokenHandler(controller.UserController))
	apiv1.Handle(http.MethodPost, "/users/resend-token", handlers.ResendTokenHandler(controller.UserController))
	apiv1.Handle(http.MethodPost, "/users/activate", handlers.ActivateUserHandler(controller.UserController))
	apiv1.Handle(http.MethodPost, "/users/check-email", handlers.CheckEmailHandler(controller.UserController))
	apiv1.Handle(http.MethodPost, "/users/check-username", handlers.CheckUsernameHandler(controller.UserController))
	apiv1.Handle(http.MethodPost, "/users/forgot-password", handlers.ForgotPasswordHandler(controller.UserController))
	apiv1.Handle(http.MethodPost, "/users/verify-reset-token", handlers.VerifyPasswordResetTokenHandler(controller.UserController))
	apiv1.Handle(http.MethodPost, "/users/reset-password", handlers.ResetPasswordHandler(controller.UserController))
	securedApiv1.Handle(http.MethodGet, "/users/suggest", handlers.SuggestUsersHandler(controller.UserController))
	securedApiv1.Handle(http.MethodGet, "/users", handlers.GetUserDetailsHandler(controller.UserController))
	securedApiv1.Handle(http.MethodPatch, "/users", handlers.UpdateUserHandler(controller.UserController))
	securedApiv1.Handle(http.MethodDelete, "/users", handlers.DeleteUserHandler(controller.UserController))

	// Trip Routes
	securedApiv1.Handle(http.MethodPost, "/trips", handlers.CreateTripEntryHandler(controller.TripController))
	securedApiv1.Handle(http.MethodGet, "/trips", handlers.GetTripEntriesHandler(controller.TripController))
	securedTripApiv1.Handle(http.MethodGet, "", handlers.GetTripDetailsHandler(controller.TripController))
	securedTripApiv1.Handle(http.MethodPatch, "", handlers.UpdateTripEntryHandler(controller.TripController))
	securedTripApiv1.Handle(http.MethodDelete, "", handlers.DeleteTripEntryHandler(controller.TripController))
	securedTripApiv1.Handle(http.MethodPost, "/invite", handlers.InviteUserToTripHandler(controller.TripController))
	securedTripApiv1.Handle(http.MethodPost, "/accept", handlers.AcceptTripInviteHandler(controller.TripController))
	securedTripApiv1.Handle(http.MethodPost, "/decline", handlers.DeclineTripInviteHandler(controller.TripController))

	// Cost Category Routes
	securedTripApiv1.Handle(http.MethodPost, "/cost-categories", handlers.CreateCostCategoryEntryHandler(controller.CostCategoryController))
	securedTripApiv1.Handle(http.MethodGet, "/cost-categories", handlers.GetCostCategoryEntriesHandler(controller.CostCategoryController))
	securedTripApiv1.Handle(http.MethodGet, "/cost-categories/:costCategoryId", handlers.GetCostCategoryDetailsHandler(controller.CostCategoryController))
	securedTripApiv1.Handle(http.MethodPatch, "/cost-categories/:costCategoryId", handlers.UpdateCostCategoryEntryHandler(controller.CostCategoryController))
	securedTripApiv1.Handle(http.MethodDelete, "/cost-categories/:costCategoryId", handlers.DeleteCostCategoryEntryHandler(controller.CostCategoryController))

	// Cost Routes
	securedApiv1.Handle(http.MethodGet, "/costs/overview", handlers.GetCostOverviewHandler(controller.CostController))
	securedTripApiv1.Handle(http.MethodPost, "/costs", handlers.CreateCostEntryHandler(controller.CostController))
	securedTripApiv1.Handle(http.MethodGet, "/costs", handlers.GetCostEntriesHandler(controller.CostController))
	securedTripApiv1.Handle(http.MethodGet, "/costs/:costId", handlers.GetCostDetailsHandler(controller.CostController))
	securedTripApiv1.Handle(http.MethodPatch, "/costs/:costId", handlers.UpdateCostEntryHandler(controller.CostController))
	securedTripApiv1.Handle(http.MethodDelete, "/costs/:costId", handlers.DeleteCostEntryHandler(controller.CostController))

	// Debts Routes
	securedTripApiv1.Handle(http.MethodGet, "/debts", handlers.GetDebtsHandler(controller.DebtController))
	securedTripApiv1.Handle(http.MethodGet, "/debts/:debtId", handlers.GetDebtDetailsHandler(controller.DebtController))

	// Transaction Routes
	securedTripApiv1.Handle(http.MethodPost, "/transactions", handlers.CreateTransactionHandler(controller.TransactionController))
	securedTripApiv1.Handle(http.MethodGet, "/transactions", handlers.GetTransactionsHandler(controller.TransactionController))
	securedTripApiv1.Handle(http.MethodGet, "/transactions/:transactionId", handlers.GetTransactionDetailsHandler(controller.TransactionController))
	securedTripApiv1.Handle(http.MethodDelete, "/transactions/:transactionId", handlers.DeleteTransactionHandler(controller.TransactionController))
	return router
}
