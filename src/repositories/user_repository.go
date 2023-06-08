package repositories

import (
	"context"
	"database/sql"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
	"time"
)

type UserRepo interface {
	GetUserById(userId *uuid.UUID) (*models.UserSchema, *models.ExpenseServiceError)
	GetUserBySchema(user *models.UserSchema) (*models.UserSchema, *models.ExpenseServiceError)
	GetUserByContext(ctx context.Context) (*models.UserSchema, *models.ExpenseServiceError)
	CreateUser(user *models.UserSchema) *models.ExpenseServiceError
	UpdateUser(user *models.UserSchema) *models.ExpenseServiceError
	DeleteUser(userId *uuid.UUID) *models.ExpenseServiceError

	ValidateIfUserExists(userId *uuid.UUID) *models.ExpenseServiceError

	GetActivationTokenByUserId(userId *uuid.UUID) (*models.ActivationTokenSchema, *models.ExpenseServiceError)
	GetActivationTokenByToken(token string) (*models.ActivationTokenSchema, *models.ExpenseServiceError)
	CreateActivationToken(userId *uuid.UUID) (*models.ActivationTokenSchema, *models.ExpenseServiceError)
	DeleteActivationToken(userId *uuid.UUID) *models.ExpenseServiceError
	ActivateUser(userId *uuid.UUID) *models.ExpenseServiceError
	ConfirmActivationToken(userId *uuid.UUID) *models.ExpenseServiceError

	ValidateIfUserIsActivated(userId *uuid.UUID) *models.ExpenseServiceError

	FindUsersLikeUsername(username string) ([]*models.UserSchema, *models.ExpenseServiceError)

	ValidateEmailExistence(email string) *models.ExpenseServiceError
	ValidateUsernameExistence(username string) *models.ExpenseServiceError
}

type UserRepository struct {
	DatabaseMgr managers.DatabaseMgr
}

func (ur *UserRepository) CreateUser(user *models.UserSchema) *models.ExpenseServiceError {
	_, err := ur.DatabaseMgr.ExecuteStatement("INSERT INTO \"user\" (id, username, email, password, activated) VALUES ($1, $2, $3, $4, $5)", user.UserID, user.Username, user.Email, user.Password, user.Activated)
	if err != nil {
		// Check if duplicate key was violated, if so return error EXPENSE_USER_EXISTS
		if pqErr := err.(*pq.Error); pqErr.Code.Name() == "unique_violation" {
			log.Printf("Error while creating user: %v", err)
			return expense_errors.EXPENSE_USER_EXISTS
		}

		log.Println("Error while creating user: ", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

func (ur *UserRepository) UpdateUser(user *models.UserSchema) *models.ExpenseServiceError {
	result, err := ur.DatabaseMgr.ExecuteStatement("UPDATE \"user\" SET username = $1, email = $2, password = $3, activated = $4 WHERE id = $5", user.Username, user.Email, user.Password, user.Activated, user.UserID)
	if err != nil {

		// Check if unique constraint was violated, if so return error EXPENSE_CONFLICT
		if pqErr := err.(*pq.Error); pqErr.Code.Name() == "unique_violation" {
			log.Printf("Error while updating user: %v", err)
			return expense_errors.EXPENSE_CONFLICT
		}

		if pqErr := err.(*pq.Error); pqErr.Code.Name() == "not_null_violation" {
			return expense_errors.EXPENSE_BAD_REQUEST
		}

		log.Printf("Error while updating user: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		log.Printf("No rows affected while updating user: %v", user)
		return expense_errors.EXPENSE_USER_NOT_FOUND
	}

	return nil
}

func (ur *UserRepository) DeleteUser(userId *uuid.UUID) *models.ExpenseServiceError {
	result, err := ur.DatabaseMgr.ExecuteStatement("DELETE FROM \"user\" WHERE id = $1", userId)
	if err != nil {

		// Check if foreign key constraint was violated, if so return error EXPENSE_CONFLICT because user is still referenced in other tables
		if pqErr := err.(*pq.Error); pqErr.Code.Name() == "foreign_key_violation" {
			log.Printf("Error while deleting user: %v", err)
			return expense_errors.EXPENSE_CONFLICT
		}

		log.Printf("Error while deleting user: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		log.Printf("No rows affected while deleting user with id: %v", userId)
		return expense_errors.EXPENSE_USER_NOT_FOUND
	}

	return nil
}

func (ur *UserRepository) GetUserById(userId *uuid.UUID) (*models.UserSchema, *models.ExpenseServiceError) {
	user := &models.UserSchema{}
	row := ur.DatabaseMgr.ExecuteQueryRow("SELECT id, username, email, activated FROM \"user\" WHERE id = $1", userId)
	if err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Activated); err != nil {
		// Check if no rows were returned, if so return error EXPENSE_USER_NOT_FOUND
		if err == sql.ErrNoRows {
			return nil, expense_errors.EXPENSE_USER_NOT_FOUND
		}

		log.Printf("Error while getting user by id: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return user, nil
}

func (ur *UserRepository) GetUserBySchema(request *models.UserSchema) (*models.UserSchema, *models.ExpenseServiceError) {
	user := &models.UserSchema{}
	row := ur.DatabaseMgr.ExecuteQueryRow("SELECT id, username, email, password, activated FROM \"user\" WHERE username = $1 OR email = $2", request.Username, request.Email)
	if err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.Activated); err != nil {
		// Check if no rows were returned, if so return error EXPENSE_USER_NOT_FOUND
		if err == sql.ErrNoRows {
			return nil, expense_errors.EXPENSE_USER_NOT_FOUND
		}

		log.Printf("Error while getting user by id: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return user, nil
}

func (ur *UserRepository) GetUserByContext(ctx context.Context) (*models.UserSchema, *models.ExpenseServiceError) {
	userId, ok := ctx.Value(models.ExpenseContextKeyUserID).(*uuid.UUID)
	if !ok {
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}
	user, err := ur.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) ValidateIfUserExists(userId *uuid.UUID) *models.ExpenseServiceError {
	rows, err := ur.DatabaseMgr.ExecuteQuery("SELECT id FROM \"user\" WHERE id = $1", userId)

	if err != nil {
		log.Printf("Error while validating if user exists: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if !rows.Next() {
		return expense_errors.EXPENSE_USER_NOT_FOUND
	}

	return expense_errors.EXPENSE_USER_EXISTS
}

func (ur *UserRepository) GetActivationTokenByUserId(userId *uuid.UUID) (*models.ActivationTokenSchema, *models.ExpenseServiceError) {
	token := &models.ActivationTokenSchema{}
	row := ur.DatabaseMgr.ExecuteQueryRow("SELECT id_user, token, created_at, confirmed_at, expires_at FROM activation_token WHERE id_user = $1", userId)
	if err := row.Scan(&token.UserID, &token.Token, &token.CreatedAt, &token.ConfirmedAt, &token.ExpiresAt); err != nil {
		// Check if no rows were returned, if so return error EXPENSE_USER_NOT_FOUND
		if err == sql.ErrNoRows {
			return nil, expense_errors.EXPENSE_USER_NOT_FOUND
		}

		log.Printf("Error while getting activation token by user id: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return token, nil
}

func (ur *UserRepository) CreateActivationToken(userId *uuid.UUID) (*models.ActivationTokenSchema, *models.ExpenseServiceError) {
	tokenString := utils.GenerateRandomString(6)
	creationDate := time.Now()
	expirationDate := creationDate.Add(time.Hour * 1)

	token := &models.ActivationTokenSchema{
		UserID:    userId,
		Token:     &tokenString,
		CreatedAt: &creationDate,
		ExpiresAt: &expirationDate,
	}

	for _, err := ur.DatabaseMgr.ExecuteStatement("INSERT INTO activation_token (id_user, token, created_at, expires_at) VALUES ($1, $2, $3, $4) RETURNING token", token.UserID, token.Token, token.CreatedAt, token.ExpiresAt); err != nil; {
		if pqErr := err.(*pq.Error); pqErr.Code.Name() == "unique_violation" {
			// If the token already exists, generate a new one and try again
			tokenString = utils.GenerateRandomString(6)
			continue
		}

		log.Println("Error while creating activation token: ", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return token, nil
}

func (ur *UserRepository) DeleteActivationToken(userId *uuid.UUID) *models.ExpenseServiceError {
	result, err := ur.DatabaseMgr.ExecuteStatement("DELETE FROM activation_token WHERE id_user = $1 AND confirmed_at IS NULL", userId)
	if err != nil {
		log.Printf("Error while deleting activation token: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return expense_errors.EXPENSE_MAIL_ALREADY_VERIFIED
	}

	return nil
}

func (ur *UserRepository) ActivateUser(userId *uuid.UUID) *models.ExpenseServiceError {
	result, err := ur.DatabaseMgr.ExecuteStatement("UPDATE \"user\" SET activated = true WHERE id = $1", userId)
	if err != nil {
		log.Printf("Error while activating user: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return expense_errors.EXPENSE_MAIL_ALREADY_VERIFIED
	}

	return nil
}

func (ur *UserRepository) ConfirmActivationToken(userId *uuid.UUID) *models.ExpenseServiceError {
	_, err := ur.DatabaseMgr.ExecuteStatement("UPDATE activation_token SET confirmed_at = $1, expires_at = $2, token = $3 WHERE id_user = $4", time.Now(), nil, nil, userId)
	if err != nil {
		_, err := ur.DatabaseMgr.ExecuteStatement("UPDATE \"user\" SET activated = false WHERE id = $1", userId)
		if err != nil {
			log.Printf("Error while deactivating user: %v", err)
		}

		log.Printf("Error while updating activation token: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

func (ur *UserRepository) ValidateIfUserIsActivated(userId *uuid.UUID) *models.ExpenseServiceError {
	row := ur.DatabaseMgr.ExecuteQueryRow("SELECT activated FROM \"user\" WHERE id = $1", userId)

	var activated bool
	if err := row.Scan(&activated); err != nil {
		if err == sql.ErrNoRows {
			return expense_errors.EXPENSE_USER_NOT_FOUND
		}

		log.Printf("Error while getting user by id: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if !activated {
		return expense_errors.EXPENSE_USER_NOT_ACTIVATED
	}

	return nil
}

func (ur *UserRepository) FindUsersLikeUsername(username string) ([]*models.UserSchema, *models.ExpenseServiceError) {
	rows, err := ur.DatabaseMgr.ExecuteQuery("SELECT id, username, email, activated FROM \"user\" WHERE username LIKE $1", "%"+username+"%")
	if err != nil {
		log.Printf("Error while getting user by username: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	users := make([]*models.UserSchema, 0)
	for rows.Next() {
		user := &models.UserSchema{}
		if err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Activated); err != nil {
			log.Printf("Error while scanning user: %v", err)
			return nil, expense_errors.EXPENSE_INTERNAL_ERROR
		}

		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) GetActivationTokenByToken(token string) (*models.ActivationTokenSchema, *models.ExpenseServiceError) {
	tokenSchema := &models.ActivationTokenSchema{}

	row := ur.DatabaseMgr.ExecuteQueryRow("SELECT id_user, token, created_at, confirmed_at, expires_at FROM activation_token WHERE token = $1", token)
	if err := row.Scan(&tokenSchema.UserID, &tokenSchema.Token, &tokenSchema.CreatedAt, &tokenSchema.ConfirmedAt, &tokenSchema.ExpiresAt); err != nil {
		// Check if no rows were returned, if so then the token has expired
		if err == sql.ErrNoRows {
			return nil, expense_errors.EXPENSE_NOT_FOUND
		}

		log.Printf("Error while getting activation token by user id: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return tokenSchema, nil
}

func (ur *UserRepository) ValidateEmailExistence(email string) *models.ExpenseServiceError {
	exists, err := ur.DatabaseMgr.CheckIfExists("SELECT COUNT(*) FROM \"user\" WHERE email = $1", email)
	if err != nil {
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if exists {
		return expense_errors.EXPENSE_EMAIL_EXISTS
	}

	return nil

}

func (ur *UserRepository) ValidateUsernameExistence(username string) *models.ExpenseServiceError {
	queryString := "SELECT COUNT(*) FROM \"user\" WHERE username = $1"
	exists, err := ur.DatabaseMgr.CheckIfExists(queryString, username)
	if err != nil {
		return expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	if exists {
		return expense_errors.EXPENSE_USERNAME_EXISTS
	}

	return nil
}
