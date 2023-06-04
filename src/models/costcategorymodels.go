package models

import "github.com/google/uuid"

type CreateCostCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}

type CostCategoryResponse struct {
	CostCategoryId *uuid.UUID `json:"costCategoryId"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Icon           string     `json:"icon"`
	Color          string     `json:"color"`
}
