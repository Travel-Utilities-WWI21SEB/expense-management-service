package controllers

import (
	"context"
	"fmt"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/repositories"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"strings"
	"time"
)

// UserCtl Exposed interface to the handler-package
type UserCtl interface {
	RegisterUser(ctx context.Context, registrationData models.RegistrationRequest, form *multipart.Form) *models.ExpenseServiceError
	LoginUser(ctx context.Context, loginData models.LoginRequest) (*models.LoginResponse, *models.ExpenseServiceError)
	RefreshToken(ctx context.Context, userId *uuid.UUID) (*models.RefreshTokenResponse, *models.ExpenseServiceError)
	ForgotPassword(ctx context.Context, email string) *models.ExpenseServiceError
	VerifyPasswordResetToken(ctx context.Context, email, token string) *models.ExpenseServiceError
	ResetPassword(ctx context.Context, email, password, token string) *models.ExpenseServiceError
	ResendToken(ctx context.Context, email string) *models.ExpenseServiceError
	ActivateUser(ctx context.Context, token string) (*models.ActivationResponse, *models.ExpenseServiceError)
	UpdateUser(ctx context.Context, request *models.UpdateUserRequest) (*models.UserDetailsResponse, *models.ExpenseServiceError)

	DeleteUser(ctx context.Context) *models.ExpenseServiceError
	GetUserDetails(ctx context.Context) (*models.UserDetailsResponse, *models.ExpenseServiceError)
	SuggestUsers(ctx context.Context, query string) (*[]models.UserSuggestion, *models.ExpenseServiceError)
	CheckEmail(ctx context.Context, email string) *models.ExpenseServiceError
	CheckUsername(ctx context.Context, username string) *models.ExpenseServiceError
}

// UserController User Controller structure
type UserController struct {
	MailMgr     managers.MailMgr
	DatabaseMgr managers.DatabaseMgr
	ImageMgr    managers.ImageMgr
	UserRepo    repositories.UserRepo
}

const activationMailSubject = "Welcome to Costventures!"
const confirmationMailSubject = "Your mail has been verified!"

const passwordResetMailSubject = "Reset your password!"
const passwordResetConfirmationMailSubject = "Your password has been reset!"

const (
	activationToken    = "activationToken"
	resetPasswordToken = "forgotPasswordToken"
)

// RegisterUser creates a new user entry in the database
func (uc *UserController) RegisterUser(ctx context.Context, registrationData models.RegistrationRequest, form *multipart.Form) *models.ExpenseServiceError {
	// Create user object
	userId := uuid.New()
	hashedPassword, err := utils.HashPassword(registrationData.Password)
	if err != nil {
		log.Printf("Error in userController.RegisterUser().HashPassword(): %v", err.Error())
		return expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	birthDate, _ := time.Parse(time.DateOnly, registrationData.Birthday)
	creationDate, _ := time.Parse(time.DateOnly, time.Now().String())

	user := &models.UserSchema{
		UserID:    &userId,
		Username:  registrationData.Username,
		FirstName: registrationData.FirstName,
		LastName:  registrationData.LastName,
		Location:  registrationData.Location,
		Email:     registrationData.Email,
		Password:  hashedPassword,
		Activated: false,
		Birthday:  &birthDate,
		CreatedAt: &creationDate,
	}

	// Get file and header
	file := (*form).File["profilePicture"][0]

	// Set profile picture path to avatar_<userId>.<fileExtension>
	fileExtension := file.Filename[strings.LastIndex(file.Filename, ".")+1:]
	user.ProfilePicture = fmt.Sprintf("avatar_%s.%s", user.UserID, fileExtension)

	// Upload profile picture
	if file != nil {
		if err := uc.ImageMgr.UploadImage(file, user.UserID); err != nil {
			return err
		}
	} else {
		if err := uc.ImageMgr.UploadDefaultProfilePicture(user.UserID); err != nil {
			return err
		}
	}

	// Insert user into database
	if repoErr := uc.UserRepo.CreateUser(ctx, user); repoErr != nil {
		return repoErr
	}

	// Insert token into database
	token, repoErr := uc.UserRepo.CreateTokenByUserIdAndType(ctx, user.UserID, activationToken)
	if repoErr != nil {
		return repoErr
	}

	activationMail := &models.ActivationMail{
		Username:        user.Username,
		ActivationToken: token.Token,
		Subject:         activationMailSubject,
		Recipients:      []string{user.Email},
	}

	return uc.MailMgr.SendActivationMail(ctx, *activationMail)
}

// LoginUser checks if the user exists and if the password is correct
func (uc *UserController) LoginUser(ctx context.Context, loginData models.LoginRequest) (*models.LoginResponse, *models.ExpenseServiceError) {
	user := &models.UserSchema{
		Email: loginData.Email,
	}

	// Get user from database
	user, repoErr := uc.UserRepo.GetUserBySchema(ctx, user)
	if repoErr != nil {
		return nil, repoErr
	}

	// Check if password is correct
	if ok := utils.CheckPasswordHash(loginData.Password, user.Password); !ok {
		return nil, expense_errors.EXPENSE_CREDENTIALS_INVALID
	}

	// Check if user is activated
	if repoErr := uc.UserRepo.ValidateIfUserIsActivated(ctx, user.UserID); repoErr != nil {
		return nil, repoErr
	}

	// Generate JWT token
	token, refreshToken, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		log.Printf("Error in userController.LoginUser().GenerateJWT(): %v", err.Error())
		return nil, expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	return &models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *UserController) RefreshToken(ctx context.Context, userId *uuid.UUID) (*models.RefreshTokenResponse, *models.ExpenseServiceError) {
	// Check if user exists
	if repoErr := uc.UserRepo.ValidateIfUserExists(ctx, userId); repoErr != nil {
		return nil, repoErr
	}

	// Generate JWT token
	token, refreshToken, err := utils.GenerateJWT(userId)
	if err != nil {
		log.Printf("Error in userController.RefreshToken().GenerateJWT(): %v", err.Error())
		return nil, expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	return &models.RefreshTokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *UserController) ForgotPassword(ctx context.Context, email string) *models.ExpenseServiceError {
	user := &models.UserSchema{
		Email: email,
	}

	// Get user from database
	user, repoErr := uc.UserRepo.GetUserBySchema(ctx, user)
	if repoErr != nil {
		return repoErr
	}

	// Delete old token
	if _, repoErr := uc.UserRepo.DeleteTokenByUserIdAndType(ctx, user.UserID, resetPasswordToken); repoErr != nil {
		return repoErr
	}

	// Insert token into database
	token, repoErr := uc.UserRepo.CreateTokenByUserIdAndType(ctx, user.UserID, resetPasswordToken)
	if repoErr != nil {
		return repoErr
	}

	passwordResetMail := &models.PasswordResetMail{
		Username:   user.Username,
		ResetToken: token.Token,
		Subject:    passwordResetMailSubject,
		Recipients: []string{user.Email},
	}

	return uc.MailMgr.SendPasswordResetMail(ctx, passwordResetMail)
}

func (uc *UserController) VerifyPasswordResetToken(ctx context.Context, email, token string) *models.ExpenseServiceError {
	// Get token from database
	tokenSchema, repoErr := uc.UserRepo.GetTokenByTokenAndType(ctx, token, resetPasswordToken)
	if repoErr != nil {
		return repoErr
	}

	// Get user email from database
	user, repoErr := uc.UserRepo.GetUserById(ctx, tokenSchema.UserID)
	if repoErr != nil {
		return repoErr
	}

	// Validate token
	if user.Email != email {
		return expense_errors.EXPENSE_FORBIDDEN
	}

	return nil
}

func (uc *UserController) ResetPassword(ctx context.Context, email, password, token string) *models.ExpenseServiceError {
	// Get token from database
	tokenSchema, repoErr := uc.UserRepo.GetTokenByTokenAndType(ctx, token, resetPasswordToken)
	if repoErr != nil {
		return repoErr
	}

	// Get user from database
	user, repoErr := uc.UserRepo.GetUserById(ctx, tokenSchema.UserID)
	if repoErr != nil {
		return repoErr
	}

	// Validate token
	if user.Email != email {
		return expense_errors.EXPENSE_FORBIDDEN
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("Error in userController.ResetPassword().HashPassword(): %v", err.Error())
		return expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	// Update password
	if repoErr := uc.UserRepo.UpdatePassword(ctx, user.UserID, hashedPassword); repoErr != nil {
		return repoErr
	}

	passwordResetConfirmationMail := &models.ResetPasswordConfirmationMail{
		Username: user.Username,
		Subject:  passwordResetConfirmationMailSubject,
		Recipients: []string{
			user.Email,
		},
	}

	return uc.MailMgr.SendResetPasswordConfirmationMail(ctx, passwordResetConfirmationMail)
}

func (uc *UserController) ResendToken(ctx context.Context, email string) *models.ExpenseServiceError {
	user := &models.UserSchema{
		Email: email,
	}

	// Get user from database
	user, repoErr := uc.UserRepo.GetUserBySchema(ctx, user)
	if repoErr != nil {
		return repoErr
	}

	// Check if user is activated
	if user.Activated {
		return expense_errors.EXPENSE_MAIL_ALREADY_VERIFIED
	}

	// Delete old token
	if _, repoErr := uc.UserRepo.DeleteTokenByUserIdAndType(ctx, user.UserID, activationToken); repoErr != nil {
		return repoErr
	}

	// Insert token into database
	token, repoErr := uc.UserRepo.CreateTokenByUserIdAndType(ctx, user.UserID, activationToken)
	if repoErr != nil {
		return repoErr
	}

	activationMail := &models.ActivationMail{
		Username:        email,
		ActivationToken: token.Token,
		Subject:         activationMailSubject,
		Recipients:      []string{email},
	}

	if err := uc.MailMgr.SendActivationMail(ctx, *activationMail); err != nil {
		log.Printf("Error in userController.ResendToken().SendActivationMail(): %v", err.ErrorMessage)
		return err
	}

	return nil
}

// UpdateUser updates the user entry in the database
func (uc *UserController) UpdateUser(ctx context.Context, request *models.UpdateUserRequest) (*models.UserDetailsResponse, *models.ExpenseServiceError) {
	// Get user from database
	user, repoErr := uc.UserRepo.GetUserByContext(ctx)
	if repoErr != nil {
		return nil, repoErr
	}

	// Check if user is activated
	if repoErr := uc.UserRepo.ValidateIfUserIsActivated(ctx, user.UserID); repoErr != nil {
		return nil, repoErr
	}

	// Patch user
	if request.Username != "" {
		user.Username = request.Username
	}

	if request.Email != "" {
		user.Email = request.Email
	}

	if request.Password != "" {
		hashedPassword, err := utils.HashPassword(request.Password)
		if err != nil {
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}
		user.Password = hashedPassword
	}

	if repoErr = uc.UserRepo.UpdateUser(ctx, user); repoErr != nil {
		return nil, repoErr
	}

	return buildUserResponse(user), nil
}

func (uc *UserController) DeleteUser(ctx context.Context) *models.ExpenseServiceError {
	userId, ok := ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return uc.UserRepo.DeleteUser(ctx, userId)
}

func (uc *UserController) ActivateUser(ctx context.Context, tokenString string) (*models.ActivationResponse, *models.ExpenseServiceError) {
	// Get UserID from token
	token, repoErr := uc.UserRepo.GetTokenByTokenAndType(ctx, tokenString, activationToken)
	if repoErr != nil {
		return nil, repoErr
	}

	// Activate user
	if repoErr := uc.UserRepo.ActivateUser(ctx, token.UserID); repoErr != nil {
		return nil, repoErr
	}

	// Confirm token
	if repoErr := uc.UserRepo.ConfirmTokenByType(ctx, token.UserID, activationToken); repoErr != nil {
		return nil, repoErr
	}

	// Generate JWT
	jwtToken, refreshToken, err := utils.GenerateJWT(token.UserID)
	if err != nil {
		log.Printf("Error in userController.ActivateUser().GenerateJWT(): %v", err.Error())
		return nil, expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	// Get user from database
	user, repoErr := uc.UserRepo.GetUserById(ctx, token.UserID)
	if repoErr != nil {
		return nil, repoErr
	}

	// Send confirmation mail
	confirmationMail := &models.ConfirmationMail{
		Username:   user.Username,
		Subject:    confirmationMailSubject,
		Recipients: []string{user.Email},
	}

	if err := uc.MailMgr.SendConfirmationMail(ctx, *confirmationMail); err != nil {
		log.Printf("Error in userController.ActivateUser().SendConfirmationMail(): %v", err.ErrorMessage)
		return nil, err
	}

	return &models.ActivationResponse{
		Token:        jwtToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *UserController) GetUserDetails(ctx context.Context) (*models.UserDetailsResponse, *models.ExpenseServiceError) {
	// Get user from database
	user, repoErr := uc.UserRepo.GetUserByContext(ctx)
	if repoErr != nil {
		return nil, repoErr
	}

	// Build response
	return buildUserResponse(user), nil
}

func (uc *UserController) SuggestUsers(ctx context.Context, query string) (*[]models.UserSuggestion, *models.ExpenseServiceError) {
	// Find users like query
	users, repoErr := uc.UserRepo.FindUsersLikeUsername(ctx, query)
	if repoErr != nil {
		return nil, repoErr
	}

	// Build response
	userSuggestionsResponse := make([]models.UserSuggestion, len(users))
	for i, user := range users {
		userSuggestionsResponse[i] = models.UserSuggestion{
			UserID:   user.UserID,
			Username: user.Username,
		}
	}
	return &userSuggestionsResponse, nil
}

func (uc *UserController) CheckEmail(ctx context.Context, email string) *models.ExpenseServiceError {
	return uc.UserRepo.ValidateEmailExistence(ctx, email)
}

func (uc *UserController) CheckUsername(ctx context.Context, username string) *models.ExpenseServiceError {
	return uc.UserRepo.ValidateUsernameExistence(ctx, username)
}

// UTITLITY FUNCTIONS

func buildUserResponse(user *models.UserSchema) *models.UserDetailsResponse {
	return &models.UserDetailsResponse{
		ID:             user.UserID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Location:       user.Location,
		UserName:       user.Username,
		Email:          user.Email,
		Birthday:       user.Birthday.Format(time.DateOnly),
		ProfilePicture: "https://raw.githubusercontent.com/Travel-Utilities-WWI21SEB/expense-management-ui/main/static/justinIcon.jpg",
		CreatedAt:      user.CreatedAt.Format(time.DateOnly),
		OpenDebts:      2,
		TripsJoined:    2,
	}
}
