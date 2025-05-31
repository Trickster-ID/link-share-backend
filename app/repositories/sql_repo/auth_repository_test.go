package sql_repo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"linkshare/app/configuration"
	"linkshare/app/models"
	"testing"
	"time"
)

func TestGetUserByUsernameOrEmail(t *testing.T) {
	configuration.InitialConfigForUnitTest()
	// Create mock for pgx.Conn
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create pgxmock: %v", err)
	}
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		_ = mock.Close(ctx)
	}(mock, context.Background())

	repo := NewAuthRepository(mock)

	// Create a context
	ctx := context.Background()

	t.Run("+: return data", func(t *testing.T) {
		// Expected result
		now := time.Now()
		expectedUser := &models.Users{
			Id:           1,
			Username:     "testuser",
			Email:        "test@example.com",
			PasswordHash: "hashedpassword",
			RoleID:       1,
			CreatedAt:    &now,
			UpdatedAt:    &now,
		}

		// Mock database response
		rows := pgxmock.NewRows([]string{"id", "username", "email", "password_hash", "role_id", "created_at", "updated_at"}).
			AddRow(expectedUser.Id, expectedUser.Username, expectedUser.Email, expectedUser.PasswordHash, expectedUser.RoleID, expectedUser.CreatedAt, expectedUser.UpdatedAt)

		mock.ExpectQuery(`select id, username, email, password_hash, role_id, created_at, updated_at from users where username = \$1`).
			WithArgs("testuser").
			WillReturnRows(rows)

		// Call the method
		user, errLog := repo.GetUserByUsernameOrEmail("testuser", "", ctx)

		// Assertions
		assert.Nil(t, errLog)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.Id, user.Id)
		assert.Equal(t, expectedUser.Username, user.Username)
		assert.Equal(t, expectedUser.Email, user.Email)

		// Ensure that expectations were met
		err = mock.ExpectationsWereMet()
		assert.Nil(t, err)
	})

	t.Run("-: empty user pass", func(t *testing.T) {
		// Call the method
		user, errLog := repo.GetUserByUsernameOrEmail("", "", ctx)

		// Assertions
		assert.Nil(t, user)
		assert.NotNil(t, errLog)
		assert.Equal(t, 404, errLog.StatusCode)

		// Ensure no queries were made
		err = mock.ExpectationsWereMet()
		assert.Nil(t, err)
	})

	t.Run("-: user notfound", func(t *testing.T) {
		// Mock database response for no rows
		mock.ExpectQuery(`select id, username, email, password_hash, role_id, created_at, updated_at from users where username = \$1`).
			WithArgs("nonexistent").
			WillReturnError(pgx.ErrNoRows)

		// Call the method
		user, errLog := repo.GetUserByUsernameOrEmail("nonexistent", "", ctx)

		// Assertions
		assert.Nil(t, user)
		assert.NotNil(t, errLog)
		assert.Equal(t, 404, errLog.StatusCode)

		// Ensure that expectations were met
		err = mock.ExpectationsWereMet()
		assert.Nil(t, err)
	})

	t.Run("-: error database", func(t *testing.T) {
		// Mock database response with an error
		mock.ExpectQuery(`select id, username, email, password_hash, role_id, created_at, updated_at from users where username = \$1`).
			WithArgs("testuser").
			WillReturnError(errors.New("db error"))

		// Call the method
		user, errLog := repo.GetUserByUsernameOrEmail("testuser", "", ctx)

		// Assertions
		assert.Nil(t, user)
		assert.NotNil(t, errLog)
		assert.Equal(t, 500, errLog.StatusCode)

		// Ensure that expectations were met
		err = mock.ExpectationsWereMet()
		assert.Nil(t, err)
	})
}
