package controller

import (
	"context"
	"errors"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
)

// Exposed interface to the handler-package
type CostCtl interface {
	createCostEntry(ctx context.Context) (*model.CostResponse, error)
	updateCostEntry(ctx context.Context) (*model.CostResponse, error)
	getTripCosts(ctx context.Context) (*model.CostResponse, error)
	getCostDetails(ctx context.Context) (*model.CostResponse, error)
	deleteCostEntry(ctx context.Context) error
}

// Cost Controller structure
type CostController struct {
}

func (cc *CostController) createCostEntry(ctx context.Context) (*model.CostResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (cc *CostController) updateCostEntry(ctx context.Context) (*model.CostResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (cc *CostController) getTripCosts(ctx context.Context) (*model.CostResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (cc *CostController) getCostDetails(ctx context.Context) (*model.CostResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}

func (cc *CostController) deleteCostEntry(ctx context.Context) (*model.CostResponse, error) {
	// TO-DO
	return nil, errors.New("not implemented")
}
