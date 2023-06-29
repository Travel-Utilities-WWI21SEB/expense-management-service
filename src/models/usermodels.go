package models

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

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type VerifyPasswordResetTokenRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type UserDto struct {
	UserID   *uuid.UUID `json:"userId,omitempty"`
	Username string     `json:"username,omitempty"`
	Email    string     `json:"email,omitempty"`
	Password string     `json:"password,omitempty"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type ResendTokenRequest struct {
	Email string `json:"email"`
}

type ActivationResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
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

type CheckEmailRequest struct {
	Email string `json:"email"`
}

type CheckEmailResponse struct {
	EmailExists bool `json:"emailExists"`
}

type CheckUsernameRequest struct {
	Username string `json:"username"`
}

type CheckUsernameResponse struct {
	UsernameExists bool `json:"usernameExists"`
}
