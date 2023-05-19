package controller

import (
	"context"
	"errors"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
)

// Exposed interface to the handler-package
type UserCtl interface {
	RegisterUser(ctx context.Context) (*model.UserResponse, error)
	LoginUser(ctx context.Context) (*model.UserResponse, error)
	UpdateUser(ctx context.Context) (*model.UserResponse, error)
	DeleteUser(ctx context.Context) error
	ActivateUser(ctx context.Context) error
	GetUserDetails(ctx context.Context) (*model.UserSchema, error)
	SuggestUsers(ctx context.Context) (*model.UserSuggestions, error)
}

// User Controller structure
type UserController struct {
}

func (uc *UserController) RegisterUser(ctx context.Context) (*model.UserResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (uc *UserController) LoginUser(ctx context.Context) (*model.UserResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
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
