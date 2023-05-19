package controller

import (
	"context"
	"errors"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
)

// Exposed interface to the handler-package
type CostCtl interface {
	CreateCostEntry(ctx context.Context) (*model.CostResponse, error)
	UpdateCostEntry(ctx context.Context) (*model.CostResponse, error)
	GetCostDetails(ctx context.Context) (*model.CostResponse, error)
	DeleteCostEntry(ctx context.Context) error
}

// Cost Controller structure
type CostController struct {
}

func (cc *CostController) CreateCostEntry(ctx context.Context) (*model.CostResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (cc *CostController) UpdateCostEntry(ctx context.Context) (*model.CostResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (cc *CostController) GetCostDetails(ctx context.Context) (*model.CostResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (cc *CostController) DeleteCostEntry(ctx context.Context) error {
	// TO-DO
	return errors.New("not implemented")
}
