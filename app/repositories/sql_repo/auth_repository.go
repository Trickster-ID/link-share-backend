package sql_repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"linkshare/app/constants"
	"linkshare/app/global/db"
	"linkshare/app/global/helper"
	"linkshare/app/global/model"
	"linkshare/app/models"
	"linkshare/generated"
	"net/http"
	"time"
)

type IAuthRepository interface {
	Create(sqlTx pgx.Tx, request *generated.RegisterRequest, ctx context.Context) *model.ErrorLog
	GetUserByUsernameOrEmail(username, email string, ctx context.Context) (*models.Users, *model.ErrorLog)
}

type authRepository struct {
	db db.PgxIface
}

func NewAuthRepository(db db.PgxIface) IAuthRepository {
	return &authRepository{
		db: db,
	}
}

// Create inserts a new user into the database using the provided transaction and registration data.
// Returns an ErrorLog if the operation fails, or nil on success.
func (r *authRepository) Create(sqlTx pgx.Tx, request *generated.RegisterRequest, ctx context.Context) *model.ErrorLog {
	query := `INSERT INTO users (username, email, password_hash, role_id, uri) VALUES ($1, $2, $3, $4, $5)`
	_, err := sqlTx.Exec(ctx, query, request.Username, request.Email, request.Password, constants.UserRoleID, helper.GenerateURLPath(5))
	if err != nil {
		var asdf *pgconn.PgError
		if errors.As(err, &asdf) {
			if asdf.Code == "23505" {
				return helper.WriteLog(nil, http.StatusBadRequest, "username or email already taken")
			}
		}
		return helper.WriteLog(err, http.StatusInternalServerError, "error while creating user")
	}
	return nil
}

// GetUserByUsernameOrEmail fetches a user by either username or email
// It will prioritize username if both are provided
// Returns a user object if found, or an appropriate error log
func (r *authRepository) GetUserByUsernameOrEmail(username, email string, ctx context.Context) (*models.Users, *model.ErrorLog) {
	user := &models.Users{}

	// Validate inputs
	if username == "" && email == "" {
		err := errors.New("missing credentials")
		errLog := helper.WriteLog(err, 400, "please provide either username or email")
		return nil, errLog
	}

	// Use parameterized query to prevent SQL injection
	query := `
	SELECT id, username, email, password_hash, role_id, created_at, updated_at
	FROM users
	WHERE `

	var param string
	if username != "" {
		query += "username = $1"
		param = username
	} else {
		query += "email = $1"
		param = email
	}

	err := r.db.QueryRow(ctx, query, param).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.RoleID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// More specific error message based on the credential used
			var message string
			if username != "" {
				message = fmt.Sprintf("user with username '%s' not found", username)
			} else {
				message = fmt.Sprintf("user with email '%s' not found", email)
			}
			errLog := helper.WriteLog(err, 404, message)
			return nil, errLog
		}
		// Log the actual error for better debugging
		errLog := helper.WriteLog(err, 500, fmt.Sprintf("database error while fetching user: %v", err))
		return nil, errLog
	}
	return user, nil
}

func (r *authRepository) AssignToken(userId int64, token string, expired *time.Time, ctx context.Context) error {

	return nil
}
