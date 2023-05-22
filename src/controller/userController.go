package controller

import (
	"context"
	"errors"
	"log"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/google/uuid"
)

// Exposed interface to the handler-package
type UserCtl interface {
	RegisterUser(ctx context.Context, registrationData model.RegistrationRequest) (*model.UserResponse, *model.ExpenseServiceError)
	LoginUser(ctx context.Context, loginData model.LoginRequest) (*model.UserResponse, *model.ExpenseServiceError)
	UpdateUser(ctx context.Context) (*model.UserResponse, *model.ExpenseServiceError)
	DeleteUser(ctx context.Context, userId *uuid.UUID) *model.ExpenseServiceError
	ActivateUser(ctx context.Context) *model.ExpenseServiceError
	GetUserDetails(ctx context.Context, userId *uuid.UUID) (*model.UserSchema, *model.ExpenseServiceError)
	SuggestUsers(ctx context.Context, query string) (*model.UserSuggestions, *model.ExpenseServiceError)
}

// User Controller structure
type UserController struct {
}

// RegisterUser creates a new user entry in the database
func (uc *UserController) RegisterUser(ctx context.Context, registrationData model.RegistrationRequest) (*model.UserResponse, *model.ExpenseServiceError) {
	if utils.ContainsEmptyString(registrationData.Username, registrationData.Email, registrationData.Password) {
		return nil, &model.ExpenseServiceError{Err: errors.New("username, email or password is empty"), Status: 400}
	}

	// Check if user already exists
	queryString := "SELECT id FROM \"user\" WHERE email = $1 OR username = $2"
	row, err := utils.ExecuteQuery(queryString, registrationData.Email, registrationData.Username)
	if err != nil {
		log.Printf("Error in userController.RegisterUser().ExecuteQuery(): %v", err.Error())
		return nil, &model.ExpenseServiceError{Err: err, Status: 500}
	}

	if row.Next() {
		return nil, &model.ExpenseServiceError{Err: errors.New("user already exists"), Status: 409}
	}

	// Create new user
	userId := uuid.New()
	hashedPassword, err := utils.HashPassword(registrationData.Password)
	if err != nil {
		log.Printf("Error in userController.RegisterUser().HashPassword(): %v", err.Error())
		return nil, &model.ExpenseServiceError{Err: err, Status: 500}
	}

	user := &model.UserSchema{
		UserID:   &userId,
		UserName: registrationData.Username,
		Email:    registrationData.Email,
		Password: hashedPassword,
	}

	// Insert user into database
	queryString = "INSERT INTO \"user\" (id, username, email, password) VALUES ($1, $2, $3, $4)"
	if err := utils.ExecuteStatement(queryString, user.UserID, user.UserName, user.Email, user.Password); err != nil {
		log.Printf("Error in userController.RegisterUser().ExecuteStatement(): %v", err.Error())
		return nil, &model.ExpenseServiceError{Err: err, Status: 500}
	}

	return &model.UserResponse{
		UserID: user.UserID,
	}, nil
}

// LoginUser checks if the user exists and if the password is correct
func (uc *UserController) LoginUser(ctx context.Context, loginData model.LoginRequest) (*model.UserResponse, *model.ExpenseServiceError) {
	if utils.ContainsEmptyString(loginData.Email, loginData.Password) {
		return nil, &model.ExpenseServiceError{Err: errors.New("email or password is empty"), Status: 400}
	}

	queryString := "SELECT id, password FROM \"user\" WHERE email = $1"
	row := utils.ExecuteQueryRow(queryString, loginData.Email)

	var userId uuid.UUID
	var hashedPassword string

	if err := row.Scan(&userId, &hashedPassword); err != nil {
		return nil, &model.ExpenseServiceError{Err: err, Status: 500}
	}

	if ok := utils.CheckPasswordHash(loginData.Password, hashedPassword); !ok {
		return nil, &model.ExpenseServiceError{Err: errors.New("invalid password"), Status: 401}
	}

	return &model.UserResponse{
		UserID: &userId,
	}, nil
}

// UpdateUser updates the user entry in the database
func (uc *UserController) UpdateUser(ctx context.Context) (*model.UserResponse, *model.ExpenseServiceError) {
	// TO-DO
	return nil, &model.ExpenseServiceError{Err: errors.New("not implemented"), Status: 501}
}

// DeleteUser deletes the user entry in the database
func (uc *UserController) DeleteUser(ctx context.Context, userId *uuid.UUID) *model.ExpenseServiceError {
	if utils.ContainsEmptyString(userId.String()) {
		return &model.ExpenseServiceError{Err: errors.New("userId is empty"), Status: 400}
	}

	queryString := "DELETE FROM \"user\" WHERE user_id = $1"
	if err := utils.ExecuteStatement(queryString, userId); err != nil {
		return &model.ExpenseServiceError{Err: err, Status: 500}
	}

	return nil
}

// ActivateUser activates the user entry in the database
func (uc *UserController) ActivateUser(ctx context.Context) *model.ExpenseServiceError {
	// TO-DO
	return &model.ExpenseServiceError{Err: errors.New("not implemented"), Status: 501}
}

// GetUserDetails returns the user entry in the database
func (uc *UserController) GetUserDetails(ctx context.Context, userId *uuid.UUID) (*model.UserSchema, *model.ExpenseServiceError) {
	if utils.ContainsEmptyString(userId.String()) {
		return nil, &model.ExpenseServiceError{Err: errors.New("userId is empty"), Status: 400}
	}

	queryString := "SELECT id, username, email FROM \"user\" WHERE user_id = $1"
	row, err := utils.ExecuteQuery(queryString, userId)
	if err != nil {
		return nil, &model.ExpenseServiceError{Err: err, Status: 500}
	}

	var user model.UserSchema
	if err := row.Scan(&user.UserID, &user.UserName, &user.Email); err != nil {
		return nil, &model.ExpenseServiceError{Err: err, Status: 500}
	}

	return &user, nil
}

// SuggestUsers returns the users which username contains the query string
func (uc *UserController) SuggestUsers(ctx context.Context, query string) (*model.UserSuggestions, *model.ExpenseServiceError) {
	if utils.ContainsEmptyString(query) {
		return nil, &model.ExpenseServiceError{Err: errors.New("query is empty"), Status: 400}
	}

	queryString := "SELECT id FROM \"user\" WHERE username LIKE $1"
	rows, err := utils.ExecuteQuery(queryString, "%"+query+"%")
	if err != nil {
		return nil, &model.ExpenseServiceError{Err: err, Status: 500}
	}

	var userIds []*uuid.UUID
	for rows.Next() {
		var userId uuid.UUID
		if err := rows.Scan(&userId); err != nil {
			return nil, &model.ExpenseServiceError{Err: err, Status: 500}
		}
		userIds = append(userIds, &userId)
	}

	return &model.UserSuggestions{
		UserIDs: userIds,
	}, nil
}
