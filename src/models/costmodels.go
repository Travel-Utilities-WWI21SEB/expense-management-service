package models

import "github.com/google/uuid"

// CostDTO Data transfer object for cost entries
type CostDTO struct {
	CostID         *uuid.UUID     `json:"costId"`
	CostCategoryID *uuid.UUID     `json:"costCategoryId"`
	Amount         string         `json:"amount"`
	CurrencyCode   string         `json:"currency"`
	Description    string         `json:"description"`
	CreationDate   string         `json:"createdAt"`
	DeductionDate  string         `json:"deductedAt"`
	EndDate        string         `json:"endDate"`
	Creditor       string         `json:"creditor"`
	Contributors   []*Contributor `json:"contributors"`
}

type CostDistributionDTO struct {
	CostCategoryName string `json:"costCategoryName"`
	Amount           string `json:"amount"`
}

type TripDistributionDTO struct {
	TripName string `json:"tripName"`
	Amount   string `json:"amount"`
}

type TripNameToIdDTO struct {
	TripName string    `json:"tripName"`
	Amount   string    `json:"amount"`
	TripId   uuid.UUID `json:"tripId"`
}

// CostOverviewDTO Data transfer object for cost overview
type CostOverviewDTO struct {
	TotalCosts                    string                 `json:"totalCosts"`
	AverageTripCosts              string                 `json:"averageTripCosts"`
	MostExpensiveTrip             *TripNameToIdDTO       `json:"mostExpensiveTrip"`
	LeastExpensiveTrip            *TripNameToIdDTO       `json:"leastExpensiveTrip"`
	AverageContributionPercentage string                 `json:"averageContributionPercentage"`
	TripDistribution              []*TripDistributionDTO `json:"tripDistribution"`
	CostDistribution              []*CostDistributionDTO `json:"costDistribution"`
}

type Contributor struct {
	Username string `json:"username"`
	Amount   string `json:"amount"`
}

// CostQueryParams Query parameters for cost entries
type CostQueryParams struct {
	TripId           *uuid.UUID
	CostCategoryId   *uuid.UUID
	CostCategoryName *string
	UserId           *uuid.UUID
	Username         *string
	MinAmount        *string // MinAmount steht für die minimale Kostenhöhe
	MaxAmount        *string // MinAmount steht für die maximale Kostenhöhe
	MinDeductionDate *string // MinDeductionDate steht für das früheste Datum, an dem die Kosten abgezogen wurden
	MaxDeductionDate *string // MaxDeductionDate steht für das späteste Datum, an dem die Kosten abgezogen wurden
	MinEndDate       *string // MinEndDate steht für das früheste Datum, an dem die Kosten enden
	MaxEndDate       *string // MaxEndDate steht für das späteste Datum, an dem die Kosten enden
	MinCreationDate  *string // MinCreationDate steht für das früheste Datum, an dem die Kosten erstellt wurden
	MaxCreationDate  *string // MaxCreationDate steht für das späteste Datum, an dem die Kosten erstellt wurden
	Page             int     // Page wird für die Paginierung verwendet
	PageSize         int     // PageSize steht für die Anzahl der Kosten, die pro Seite angezeigt werden
	SortBy           string  // SortBy steht für die Spalte, nach der die Kosten sortiert werden sollen
	SortOrder        string  // SortOrder steht für die Reihenfolge, in der die Kosten sortiert werden sollen
}
