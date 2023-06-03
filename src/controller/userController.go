package controller

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expenseerror"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/manager"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/google/uuid"
)

// Exposed interface to the handler-package
type UserCtl interface {
	RegisterUser(ctx context.Context, registrationData model.RegistrationRequest) *model.ExpenseServiceError
	LoginUser(ctx context.Context, loginData model.LoginRequest) (*model.LoginResponse, *model.ExpenseServiceError)
	RefreshToken(ctx context.Context, userId *uuid.UUID) (*model.RefreshTokenResponse, *model.ExpenseServiceError)
	ActivateUser(ctx context.Context, token string) (*model.ActivationResponse, *model.ExpenseServiceError)
	UpdateUser(ctx context.Context) (*model.UserDetailsResponse, *model.ExpenseServiceError)
	DeleteUser(ctx context.Context) *model.ExpenseServiceError
	GetUserDetails(ctx context.Context) (*model.UserDetailsResponse, *model.ExpenseServiceError)
	SuggestUsers(ctx context.Context, query string) (*model.UserSuggestResponse, *model.ExpenseServiceError)
	CheckEmail(ctx context.Context, email string) *model.ExpenseServiceError
	CheckUsername(ctx context.Context, username string) *model.ExpenseServiceError
}

// User Controller structure
type UserController struct {
	MailMgr     manager.MailMgr
	DatabaseMgr manager.DatabaseMgr
}

const activationMailSubject = "Welcome to Costventures!"
const confirmationMailSubject = "Your mail has been verified!"

// RegisterUser creates a new user entry in the database
func (uc *UserController) RegisterUser(ctx context.Context, registrationData model.RegistrationRequest) *model.ExpenseServiceError {
	if utils.ContainsEmptyString(registrationData.Username, registrationData.Email, registrationData.Password) {
		return expenseerror.EXPENSE_BAD_REQUEST
	}

	// Check if user already exists
	queryString := "SELECT id FROM \"user\" WHERE email = $1 OR username = $2"
	row, err := uc.DatabaseMgr.ExecuteQuery(queryString, registrationData.Email, registrationData.Username)
	if err != nil {
		log.Printf("Error in userController.RegisterUser().ExecuteQuery(): %v", err.Error())
		return expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	if row.Next() {
		return expenseerror.EXPENSE_USER_EXISTS
	}

	// Create new user
	userId := uuid.New()
	hashedPassword, err := utils.HashPassword(registrationData.Password)
	if err != nil {
		log.Printf("Error in userController.RegisterUser().HashPassword(): %v", err.Error())
		return expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	user := &model.UserSchema{
		UserID:   &userId,
		UserName: registrationData.Username,
		Email:    registrationData.Email,
		Password: hashedPassword,
	}

	// Insert user into database
	queryString = "INSERT INTO \"user\" (id, username, email, password, activated) VALUES ($1, $2, $3, $4, $5)"
	if _, err := uc.DatabaseMgr.ExecuteStatement(queryString, user.UserID, user.UserName, user.Email, user.Password, false); err != nil {
		log.Printf("Error in userController.RegisterUser().ExecuteStatement(): %v", err.Error())
		return expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Generate activation token and send activation mail
	tokenId := uuid.New()
	var activationToken string
	now := time.Now()
	queryString = "SELECT COUNT(*) FROM activation_token WHERE token = $1"

	for {
		activationToken = utils.GenerateRandomString(6)
		row := uc.DatabaseMgr.ExecuteQueryRow(queryString, activationToken)

		var count int
		if err := row.Scan(&count); err != nil {
			log.Printf("Error in userController.RegisterUser().Scan(): %v", err.Error())
			return expenseerror.EXPENSE_UPSTREAM_ERROR
		}

		if count == 0 {
			// Token is unique
			break
		}
	}

	queryString = "INSERT INTO activation_token (id, token, created_at, id_user) VALUES ($1, $2, $3, $4)"
	if _, err := uc.DatabaseMgr.ExecuteStatement(queryString, tokenId, activationToken, now, user.UserID); err != nil {
		log.Printf("Error in userController.RegisterUser().ExecuteStatement(): %v", err.Error())
		return expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	activationMail := &model.ActivationMail{
		Username:        user.UserName,
		ActivationToken: activationToken,
		Subject:         activationMailSubject,
		Recipients:      []string{user.Email},
	}

	if err := uc.MailMgr.SendActivationMail(ctx, *activationMail); err != nil {
		log.Printf("Error in userController.RegisterUser().SendActivationMail(): %v", err.ErrorMessage)
		return err
	}

	return nil
}

// LoginUser checks if the user exists and if the password is correct
func (uc *UserController) LoginUser(ctx context.Context, loginData model.LoginRequest) (*model.LoginResponse, *model.ExpenseServiceError) {
	if utils.ContainsEmptyString(loginData.Email, loginData.Password) {
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	queryString := "SELECT id, password, activated FROM \"user\" WHERE email = $1"
	row := uc.DatabaseMgr.ExecuteQueryRow(queryString, loginData.Email)

	var userId uuid.UUID
	var hashedPassword string
	var activated bool

	if err := row.Scan(&userId, &hashedPassword, &activated); err != nil {
		log.Printf("Error in userController.LoginUser().Scan(): %v", err.Error())
		return nil, expenseerror.EXPENSE_USER_NOT_FOUND
	}

	if ok := utils.CheckPasswordHash(loginData.Password, hashedPassword); !ok {
		return nil, expenseerror.EXPENSE_CREDENTIALS_INVALID
	}

	if !activated {
		return nil, expenseerror.EXPENSE_USER_NOT_ACTIVATED
	}

	token, refreshToken, err := utils.GenerateJWT(&userId)
	if err != nil {
		log.Printf("Error in userController.LoginUser().GenerateJWT(): %v", err.Error())
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	return &model.LoginResponse{
		UserID:       &userId,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *UserController) RefreshToken(ctx context.Context, userId *uuid.UUID) (*model.RefreshTokenResponse, *model.ExpenseServiceError) {
	if utils.ContainsEmptyString(userId.String()) {
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	queryString := "SELECT COUNT(*) FROM \"user\" WHERE id = $1"
	var count int
	row := uc.DatabaseMgr.ExecuteQueryRow(queryString, userId)
	if err := row.Scan(&count); err != nil {
		log.Printf("Error in userController.RefreshToken().Scan(): %v", err.Error())
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	if count == 0 {
		return nil, expenseerror.EXPENSE_USER_NOT_FOUND
	}

	token, refreshToken, err := utils.GenerateJWT(userId)
	if err != nil {
		log.Printf("Error in userController.RefreshToken().GenerateJWT(): %v", err.Error())
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	return &model.RefreshTokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

// UpdateUser updates the user entry in the database
func (uc *UserController) UpdateUser(ctx context.Context) (*model.UserDetailsResponse, *model.ExpenseServiceError) {
	// TO-DO
	return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
}

// DeleteUser deletes the user entry in the database
func (uc *UserController) DeleteUser(ctx context.Context) *model.ExpenseServiceError {
	userId, ok := ctx.Value(model.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		return expenseerror.EXPENSE_BAD_REQUEST
	}

	queryString := "DELETE FROM \"user\" WHERE id = $1"
	result, err := uc.DatabaseMgr.ExecuteStatement(queryString, userId)
	if err != nil {
		log.Printf("Error in userController.DeleteUser().ExecuteStatement(): %v", err.Error())
		return expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error in userController.DeleteUser().RowsAffected(): %v", err.Error())
		return expenseerror.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected == 0 {
		return expenseerror.EXPENSE_USER_NOT_FOUND
	}

	return nil
}

// ActivateUser activates the user entry in the database
func (uc *UserController) ActivateUser(ctx context.Context, token string) (*model.ActivationResponse, *model.ExpenseServiceError) {
	queryString := "SELECT id, id_user, confirmed_at FROM activation_token WHERE token = $1"
	row := uc.DatabaseMgr.ExecuteQueryRow(queryString, token)

	var tokenId *uuid.UUID
	var userId *uuid.UUID
	var confirmedAt *time.Time
	if err := row.Scan(&tokenId, &userId, &confirmedAt); err != nil {
		log.Printf("Error in userController.ActivateUser().Scan(): %v", err.Error())
		return nil, expenseerror.EXPENSE_INVALID_ACTIVATION_TOKEN
	}

	if confirmedAt != nil {
		return nil, expenseerror.EXPENSE_MAIL_ALREADY_VERIFIED
	}

	// Select user from database
	queryString = "SELECT username, email FROM \"user\" WHERE id = $1"
	row = uc.DatabaseMgr.ExecuteQueryRow(queryString, userId)

	var username string
	var email string

	if err := row.Scan(&username, &email); err != nil {
		log.Printf("Error in userController.ActivateUser().Scan(): %v", err.Error())
		return nil, expenseerror.EXPENSE_USER_NOT_FOUND
	}

	// Activate user in database and save confirmation time
	queryString = "UPDATE \"user\" SET activated = $1 WHERE id = $2"
	_, err := uc.DatabaseMgr.ExecuteStatement(queryString, true, userId)
	if err != nil {
		log.Printf("Error in userController.ActivateUser().ExecuteStatement(): %v", err.Error())
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	now := time.Now()

	queryString = "UPDATE activation_token SET confirmed_at = $1 WHERE id = $2"
	_, err = uc.DatabaseMgr.ExecuteStatement(queryString, now, tokenId)
	if err != nil {
		log.Printf("Error in userController.ActivateUser().ExecuteStatement(): %v", err.Error())
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Generate JWT
	jwtToken, refreshToken, err := utils.GenerateJWT(userId)
	if err != nil {
		log.Printf("Error in userController.ActivateUser().GenerateJWT(): %v", err.Error())
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	// Send confirmation mail
	confirmationMail := &model.ConfirmationMail{
		Username:   username,
		Subject:    confirmationMailSubject,
		Recipients: []string{email},
	}

	if err := uc.MailMgr.SendConfirmationMail(ctx, *confirmationMail); err != nil {
		log.Printf("Error in userController.ActivateUser().SendConfirmationMail(): %v", err.ErrorMessage)
		return nil, err
	}

	return &model.ActivationResponse{
		Token:        jwtToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetUserDetails returns the user entry in the database
func (uc *UserController) GetUserDetails(ctx context.Context) (*model.UserDetailsResponse, *model.ExpenseServiceError) {
	userId, ok := ctx.Value(model.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	queryString := "SELECT username, email FROM \"user\" WHERE id = $1"
	row := uc.DatabaseMgr.ExecuteQueryRow(queryString, userId)

	var userDetailsResponse model.UserDetailsResponse
	if err := row.Scan(&userDetailsResponse.UserName, &userDetailsResponse.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, expenseerror.EXPENSE_USER_NOT_FOUND
		}

		log.Printf("Error in userController.GetUserDetails().Scan(): %v", err.Error())
		return nil, expenseerror.EXPENSE_INTERNAL_ERROR
	}

	return &userDetailsResponse, nil
}

// SuggestUsers returns the users which username contains the query string
func (uc *UserController) SuggestUsers(ctx context.Context, query string) (*model.UserSuggestResponse, *model.ExpenseServiceError) {
	if utils.ContainsEmptyString(query) {
		return nil, expenseerror.EXPENSE_BAD_REQUEST
	}

	queryString := "SELECT id, username FROM \"user\" WHERE username LIKE $1"
	rows, err := uc.DatabaseMgr.ExecuteQuery(queryString, fmt.Sprintf("%v%%", query))
	if err != nil {
		log.Printf("Error in userController.SuggestUsers().ExecuteQuery(): %v", err.Error())
		return nil, expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	var userSuggestResponse model.UserSuggestResponse

	for rows.Next() {
		var userId uuid.UUID
		var userName string

		if err := rows.Scan(&userId, &userName); err != nil {
			log.Printf("Error in userController.SuggestUsers().Scan(): %v", err.Error())
			return nil, expenseerror.EXPENSE_INTERNAL_ERROR
		}

		userSuggestResponse.UserSuggestions = append(userSuggestResponse.UserSuggestions, model.UserSuggestion{
			UserID:   &userId,
			Username: userName,
		})
	}

	return &userSuggestResponse, nil
}

func (uc *UserController) CheckEmail(ctx context.Context, email string) *model.ExpenseServiceError {
	if utils.ContainsEmptyString(email) {
		return expenseerror.EXPENSE_BAD_REQUEST
	}

	queryString := "SELECT COUNT(*) FROM \"user\" WHERE email = $1"
	var count int
	row := uc.DatabaseMgr.ExecuteQueryRow(queryString, email)
	if err := row.Scan(&count); err != nil {
		log.Printf("Error in userController.checkEmail().Scan(): %v", err.Error())
		return expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	if count != 0 {
		return expenseerror.EXPENSE_EMAIL_EXISTS
	}

	return nil
}

func (uc *UserController) CheckUsername(ctx context.Context, username string) *model.ExpenseServiceError {
	if utils.ContainsEmptyString(username) {
		return expenseerror.EXPENSE_BAD_REQUEST
	}

	queryString := "SELECT COUNT(*) FROM \"user\" WHERE username = $1"
	var count int
	row := uc.DatabaseMgr.ExecuteQueryRow(queryString, username)
	if err := row.Scan(&count); err != nil {
		log.Printf("Error in userController.checkEmail().Scan(): %v", err.Error())
		return expenseerror.EXPENSE_UPSTREAM_ERROR
	}

	if count != 0 {
		return expenseerror.EXPENSE_USERNAME_EXISTS
	}

	return nil
}
