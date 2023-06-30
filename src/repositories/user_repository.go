package repositories

import (
	"context"
	"errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"time"
)

type UserRepo interface {
	GetUserById(ctx context.Context, userId *uuid.UUID) (*models.UserSchema, *models.ExpenseServiceError)
	GetUserBySchema(ctx context.Context, user *models.UserSchema) (*models.UserSchema, *models.ExpenseServiceError)
	GetUserByContext(ctx context.Context) (*models.UserSchema, *models.ExpenseServiceError)

	CreateUser(ctx context.Context, user *models.UserSchema) *models.ExpenseServiceError
	UpdateUser(ctx context.Context, user *models.UserSchema) *models.ExpenseServiceError
	DeleteUser(ctx context.Context, userId *uuid.UUID) *models.ExpenseServiceError

	UpdatePassword(ctx context.Context, userId *uuid.UUID, newPassword string) *models.ExpenseServiceError

	ValidateIfUserExists(ctx context.Context, userId *uuid.UUID) *models.ExpenseServiceError
	ValidateIfUserIsActivated(ctx context.Context, userId *uuid.UUID) *models.ExpenseServiceError
	ValidateEmailExistence(ctx context.Context, email string) *models.ExpenseServiceError
	ValidateUsernameExistence(ctx context.Context, username string) *models.ExpenseServiceError

	ActivateUser(ctx context.Context, userId *uuid.UUID) *models.ExpenseServiceError

	CreateTokenByUserIdAndType(ctx context.Context, userId *uuid.UUID, tokenType string) (*models.TokenSchema, *models.ExpenseServiceError)
	DeleteTokenByUserIdAndType(ctx context.Context, userId *uuid.UUID, tokenType string) (int64, *models.ExpenseServiceError)
	GetTokenByUserIdAndType(ctx context.Context, userId *uuid.UUID, tokenType string) (*models.TokenSchema, *models.ExpenseServiceError)
	GetTokenByTokenAndType(ctx context.Context, token, tokenType string) (*models.TokenSchema, *models.ExpenseServiceError)
	ConfirmTokenByType(ctx context.Context, userId *uuid.UUID, tokenType string) *models.ExpenseServiceError

	FindUsersLikeUsername(ctx context.Context, username string) ([]*models.UserSchema, *models.ExpenseServiceError)
}

type UserRepository struct {
	DatabaseMgr managers.DatabaseMgr
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *models.UserSchema) *models.ExpenseServiceError {
	_, err := ur.DatabaseMgr.ExecuteStatement(ctx, "INSERT INTO \"user\" (id, username, email, password, activated) VALUES ($1, $2, $3, $4, $5)", user.UserID, user.Username, user.Email, user.Password, user.Activated)
	if err != nil {
		// Check if duplicate key was violated, if so return error EXPENSE_USER_EXISTS
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr); pgxErr.Code == "unique_violation" {
			log.Printf("Error while creating user: %v", err)
			return expense_errors.EXPENSE_USER_EXISTS
		}

		log.Println("Error while creating user: ", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, user *models.UserSchema) *models.ExpenseServiceError {
	result, err := ur.DatabaseMgr.ExecuteStatement(ctx, "UPDATE \"user\" SET username = $1, email = $2, password = $3, activated = $4 WHERE id = $5", user.Username, user.Email, user.Password, user.Activated, user.UserID)

	if err != nil {
		// Check if unique constraint was violated, if so return error EXPENSE_CONFLICT
		var pgxErr *pgconn.PgError
		errors.As(err, &pgxErr)

		if pgxErr.Code == "unique_violation" {
			log.Printf("Error while updating user: %v", err)
			return expense_errors.EXPENSE_CONFLICT
		} else if pgxErr.Code == "not_null_violation" {
			log.Printf("Error while updating user: %v", err)
			return expense_errors.EXPENSE_BAD_REQUEST
		}

		log.Printf("Error while updating user: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		log.Printf("No rows affected while updating user: %v", user)
		return expense_errors.EXPENSE_USER_NOT_FOUND
	}

	return nil
}

func (ur *UserRepository) DeleteUser(ctx context.Context, userId *uuid.UUID) *models.ExpenseServiceError {
	result, err := ur.DatabaseMgr.ExecuteStatement(ctx, "DELETE FROM \"user\" WHERE id = $1", userId)

	if err != nil {
		// Check if foreign key constraint was violated, if so return error EXPENSE_CONFLICT because user is still referenced in other tables
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr); pgxErr.Code == "foreign_key_violation" {
			log.Printf("Error while deleting user: %v", err)
			return expense_errors.EXPENSE_CONFLICT
		}

		log.Printf("Error while deleting user: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		log.Printf("No rows affected while deleting user with id: %v", userId)
		return expense_errors.EXPENSE_USER_NOT_FOUND
	}

	return nil
}

func (ur *UserRepository) GetUserById(ctx context.Context, userId *uuid.UUID) (*models.UserSchema, *models.ExpenseServiceError) {
	user := &models.UserSchema{}
	row := ur.DatabaseMgr.ExecuteQueryRow(ctx, "SELECT id, username, email, activated FROM \"user\" WHERE id = $1", userId)
	if err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Activated); err != nil {
		// Check if no rows were returned, if so return error EXPENSE_USER_NOT_FOUND
		if err == pgx.ErrNoRows {
			return nil, expense_errors.EXPENSE_USER_NOT_FOUND
		}

		log.Printf("Error while getting user by id: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return user, nil
}

func (ur *UserRepository) GetUserBySchema(ctx context.Context, request *models.UserSchema) (*models.UserSchema, *models.ExpenseServiceError) {
	user := &models.UserSchema{}
	row := ur.DatabaseMgr.ExecuteQueryRow(ctx, "SELECT id, username, email, password, activated FROM \"user\" WHERE username = $1 OR email = $2", request.Username, request.Email)
	if err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.Activated); err != nil {
		// Check if no rows were returned, if so return error EXPENSE_USER_NOT_FOUND
		log.Printf("Error while getting user by schema: %v", err)
		if err == pgx.ErrNoRows {
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

	user, err := ur.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) UpdatePassword(ctx context.Context, userId *uuid.UUID, newPassword string) *models.ExpenseServiceError {
	_, err := ur.DatabaseMgr.ExecuteStatement(ctx, "UPDATE \"user\" SET password = $1 WHERE id = $2", newPassword, userId)
	if err != nil {
		log.Printf("Error while updating password: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

func (ur *UserRepository) ValidateIfUserExists(ctx context.Context, userId *uuid.UUID) *models.ExpenseServiceError {
	rows, err := ur.DatabaseMgr.ExecuteQuery(ctx, "SELECT id FROM \"user\" WHERE id = $1", userId)
	if err != nil {
		log.Printf("Error while validating if user exists: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	defer rows.Close()

	if !rows.Next() {
		return expense_errors.EXPENSE_USER_NOT_FOUND
	}

	return nil
}

func (ur *UserRepository) GetTokenByUserIdAndType(ctx context.Context, userId *uuid.UUID, tokenType string) (*models.TokenSchema, *models.ExpenseServiceError) {
	token := &models.TokenSchema{}

	row := ur.DatabaseMgr.ExecuteQueryRow(ctx, "SELECT id_user, token, created_at, confirmed_at, expires_at FROM token WHERE id_user = $1 AND type = $2", userId, tokenType)
	if err := row.Scan(&token.UserID, &token.Token, &token.CreatedAt, &token.ConfirmedAt, &token.ExpiresAt); err != nil {
		// Check if no rows were returned, if so return error EXPENSE_USER_NOT_FOUND
		if err == pgx.ErrNoRows {
			return nil, expense_errors.EXPENSE_USER_NOT_FOUND
		}

		log.Printf("Error while getting activation token by user id: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return token, nil
}

func (ur *UserRepository) CreateTokenByUserIdAndType(ctx context.Context, userId *uuid.UUID, tokenType string) (*models.TokenSchema, *models.ExpenseServiceError) {
	creationDate := time.Now()
	expirationDate := creationDate.Add(time.Hour * 1)

	token := &models.TokenSchema{
		UserID:    userId,
		Token:     utils.GenerateRandomString(6),
		Type:      tokenType,
		CreatedAt: &creationDate,
		ExpiresAt: &expirationDate,
	}

	tokenId := uuid.New()

	for _, err := ur.DatabaseMgr.ExecuteStatement(ctx, "INSERT INTO token(id, id_user, token, type, created_at, expires_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING token", tokenId, token.UserID, token.Token, token.Type, token.CreatedAt, token.ExpiresAt); err != nil; {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr); pgxErr.Code == "unique_violation" {
			// If the token already exists, generate a new one and try again
			token.Token = utils.GenerateRandomString(6)
			continue
		}

		log.Println("Error while creating activation token: ", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return token, nil
}

func (ur *UserRepository) DeleteTokenByUserIdAndType(ctx context.Context, userId *uuid.UUID, tokenType string) (int64, *models.ExpenseServiceError) {
	result, err := ur.DatabaseMgr.ExecuteStatement(ctx, "DELETE FROM token WHERE id_user = $1 AND type = $2 AND confirmed_at IS NULL", userId, tokenType)
	if err != nil {
		log.Printf("Error while deleting %v token: %v", tokenType, err)
		return 0, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	rowsAffected := result.RowsAffected()

	return rowsAffected, nil
}

func (ur *UserRepository) ActivateUser(ctx context.Context, userId *uuid.UUID) *models.ExpenseServiceError {
	result, err := ur.DatabaseMgr.ExecuteStatement(ctx, "UPDATE \"user\" SET activated = true WHERE id = $1", userId)
	if err != nil {
		log.Printf("Error while activating user: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return expense_errors.EXPENSE_MAIL_ALREADY_VERIFIED
	}

	return nil
}

func (ur *UserRepository) ConfirmTokenByType(ctx context.Context, userId *uuid.UUID, tokenType string) *models.ExpenseServiceError {
	_, err := ur.DatabaseMgr.ExecuteStatement(ctx, "UPDATE token SET confirmed_at = $1, expires_at = $2, token = $3 WHERE id_user = $4 AND type = $5", time.Now(), nil, nil, userId, tokenType)
	if err != nil {
		log.Printf("Error while updating activation token: %v", err)
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return nil
}

func (ur *UserRepository) ValidateIfUserIsActivated(ctx context.Context, userId *uuid.UUID) *models.ExpenseServiceError {
	row := ur.DatabaseMgr.ExecuteQueryRow(ctx, "SELECT activated FROM \"user\" WHERE id = $1", userId)

	var activated bool
	if err := row.Scan(&activated); err != nil {
		if err == pgx.ErrNoRows {
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

func (ur *UserRepository) FindUsersLikeUsername(ctx context.Context, username string) ([]*models.UserSchema, *models.ExpenseServiceError) {
	rows, err := ur.DatabaseMgr.ExecuteQuery(ctx, "SELECT id, username, email, activated FROM \"user\" WHERE username LIKE $1", "%"+username+"%")
	if err != nil {
		log.Printf("Error while getting user by username: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}
	defer rows.Close()

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

func (ur *UserRepository) GetTokenByTokenAndType(ctx context.Context, token, tokenType string) (*models.TokenSchema, *models.ExpenseServiceError) {
	tokenSchema := &models.TokenSchema{}

	row := ur.DatabaseMgr.ExecuteQueryRow(ctx, "SELECT id_user, token, created_at, confirmed_at, expires_at FROM token WHERE token = $1 AND type = $2", token, tokenType)
	if err := row.Scan(&tokenSchema.UserID, &tokenSchema.Token, &tokenSchema.CreatedAt, &tokenSchema.ConfirmedAt, &tokenSchema.ExpiresAt); err != nil {
		// Check if no rows were returned, if so then the token has expired
		if err == pgx.ErrNoRows {
			return nil, expense_errors.EXPENSE_NOT_FOUND
		}

		log.Printf("Error while getting activation token by user id: %v", err)
		return nil, expense_errors.EXPENSE_INTERNAL_ERROR
	}

	return tokenSchema, nil
}

func (ur *UserRepository) ValidateEmailExistence(ctx context.Context, email string) *models.ExpenseServiceError {
	exists, err := ur.DatabaseMgr.CheckIfExists(ctx, "SELECT COUNT(*) FROM \"user\" WHERE email = $1", email)
	if err != nil {
		return expense_errors.EXPENSE_INTERNAL_ERROR
	}

	if exists {
		return expense_errors.EXPENSE_EMAIL_EXISTS
	}

	return nil

}

func (ur *UserRepository) ValidateUsernameExistence(ctx context.Context, username string) *models.ExpenseServiceError {
	queryString := "SELECT COUNT(*) FROM \"user\" WHERE username = $1"
	exists, err := ur.DatabaseMgr.CheckIfExists(ctx, queryString, username)
	if err != nil {
		return expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	if exists {
		return expense_errors.EXPENSE_USERNAME_EXISTS
	}

	return nil
}
