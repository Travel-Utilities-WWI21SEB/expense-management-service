package controller

import (
	"context"
	"errors"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/google/uuid"
)

// Exposed interface to the handler-package
type UserCtl interface {
	RegisterUser(ctx context.Context, registrationData model.RegistrationRequest) (*model.UserResponse, *model.ExpenseServiceError)
	LoginUser(ctx context.Context, loginData model.LoginRequest) (*model.UserResponse, *model.ExpenseServiceError)
	UpdateUser(ctx context.Context) (*model.UserResponse, error)
	DeleteUser(ctx context.Context) error
	ActivateUser(ctx context.Context) error
	GetUserDetails(ctx context.Context) (*model.UserSchema, error)
	SuggestUsers(ctx context.Context) (*model.UserSuggestions, error)
}

// User Controller structure
type UserController struct {
}

func (uc *UserController) RegisterUser(ctx context.Context, registrationData model.RegistrationRequest) (*model.UserResponse, *model.ExpenseServiceError) {
	if utils.ContainsEmptyString(registrationData.Username, registrationData.Email, registrationData.Password) {
		return nil, &model.ExpenseServiceError{Err: errors.New("username, email or password is empty"), Status: 400}
	}

	userId := uuid.New()
	hashedPassword, err := utils.HashPassword(registrationData.Password)
	if err != nil {
		return nil, &model.ExpenseServiceError{Err: err, Status: 500}
	}

	user := &model.UserSchema{
		UserID:   &userId,
		UserName: registrationData.Username,
		Email:    registrationData.Email,
		Password: hashedPassword,
	}

	queryString := "INSERT INTO users (user_id, username, email, password) VALUES ($1, $2, $3, $4)"
	if err := utils.ExecuteInsert(queryString, user.UserID, user.UserName, user.Email, user.Password); err != nil {
		return nil, &model.ExpenseServiceError{Err: err, Status: 500}
	}

	return &model.UserResponse{
		UserID: user.UserID,
	}, nil
}

func (uc *UserController) LoginUser(ctx context.Context, loginData model.LoginRequest) (*model.UserResponse, *model.ExpenseServiceError) {
	if utils.ContainsEmptyString(loginData.Email, loginData.Password) {
		return nil, &model.ExpenseServiceError{Err: errors.New("email or password is empty"), Status: 400}
	}

	queryString := "SELECT user_id, password FROM users WHERE email = $1"
	row, err := utils.ExecuteQuery(queryString, loginData.Email)
	if err != nil {
		return nil, &model.ExpenseServiceError{Err: err, Status: 500}
	}

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

func (uc *UserController) UpdateUser(ctx context.Context) (*model.UserResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (uc *UserController) DeleteUser(ctx context.Context) error {
	// TO-DO
	return errors.New("not implemented")
}

func (uc *UserController) ActivateUser(ctx context.Context) error {
	// TO-DO
	return errors.New("not implemented")
}

func (uc *UserController) GetUserDetails(ctx context.Context) (*model.UserSchema, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (uc *UserController) SuggestUsers(ctx context.Context) (*model.UserSuggestions, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}
