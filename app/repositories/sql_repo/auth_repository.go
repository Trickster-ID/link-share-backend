package sql_repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"linkshare/app/global/db"
	"linkshare/app/global/helper"
	"linkshare/app/global/model"
	"linkshare/app/models"
	"time"
)

type IAuthRepository interface {
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

func (r *authRepository) GetUserByUsernameOrEmail(username, email string, ctx context.Context) (*models.Users, *model.ErrorLog) {
	user := &models.Users{}
	whereClause := ""
	valueOfWhereCluase := ""
	if username != "" {
		whereClause = "username = $1"
		valueOfWhereCluase = username
	} else if email != "" {
		whereClause = "email = $1"
		valueOfWhereCluase = email
	} else {
		err := errors.New("username is empty")
		errLog := helper.WriteLog(err, 404, "please enter valid username")
		return nil, errLog
	}
	err := r.db.QueryRow(ctx, fmt.Sprintf(`
	select 
    id, username, email, password_hash, role_id, created_at, updated_at 
	from users 
	where %s
	`, whereClause), valueOfWhereCluase).Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errLog := helper.WriteLog(err, 404, "please enter valid username or email or password")
			return nil, errLog
		}
		errLog := helper.WriteLog(err, 500, "error while getting user")
		return nil, errLog
	}
	return user, nil
}

func (r *authRepository) AssignToken(userId int64, token string, expired *time.Time, ctx context.Context) error {

	return nil
}
