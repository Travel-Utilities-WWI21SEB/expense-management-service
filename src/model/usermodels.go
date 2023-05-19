package model

import (
	"github.com/google/uuid"
)

type UserResponse struct {
	UserID *uuid.UUID `json:"userId"`
}

type UserSuggestions struct {
	UserIDs []*uuid.UUID `json:"userIds"`
}
