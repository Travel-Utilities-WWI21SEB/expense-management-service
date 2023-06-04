package models

import "github.com/google/uuid"

type CostCategoryPostRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}

type CostCategoryPatchRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Color       string `json:"color,omitempty"`
}

type CostCategoryResponse struct {
	CostCategoryId *uuid.UUID `json:"costCategoryId"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Icon           string     `json:"icon"`
	Color          string     `json:"color"`
}
