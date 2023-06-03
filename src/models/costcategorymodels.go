package models

type CreateCostCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}

type CostCategoryResponse struct {
	CostCategoryId int    `json:"costCategoryId"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Icon           string `json:"icon"`
	Color          string `json:"color"`
}
