package controller

import (
	"context"
	"errors"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
)

// Exposed interface to the handler-package
type UserCtl interface {
	registerUser(ctx context.Context) (*model.UserResponse, error)
	loginUser(ctx context.Context) (*model.UserResponse, error)
	updateUser(ctx context.Context) (*model.UserResponse, error)
	deleteUser(ctx context.Context) error
	activateUser(ctx context.Context) error
	getUserDetails(ctx context.Context) (*model.UserSchema, error)
	suggestUsers(ctx context.Context) (*model.UserSuggestions, error)
}

// Cost Controller structure
type UserController struct {
}

func (uc *UserController) registerUser(ctx context.Context) (*model.UserResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (uc *UserController) loginUser(ctx context.Context) (*model.UserResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (uc *UserController) getUserDetails(ctx context.Context) (*model.UserSchema, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (uc *UserController) getCostDetails(ctx context.Context) (*model.CostResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (cc *UserController) deleteCostEntry(ctx context.Context) (*model.CostResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}
