package model

import (
	"github.com/google/uuid"
)

type RegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	UserID *uuid.UUID `json:"userId"`
}

type UserSuggestions struct {
	UserIDs []*uuid.UUID `json:"userIds"`
}
