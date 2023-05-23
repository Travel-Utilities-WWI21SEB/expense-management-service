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
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	UserID *uuid.UUID `json:"userId"`
}

type UserDetailsResponse struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type UserSuggestResponse struct {
	UserSuggestions []UserSuggestion `json:"userSuggestions"`
}

type UserSuggestion struct {
	UserID   *uuid.UUID `json:"userId"`
	Username string     `json:"username"`
}
