package repository

import (
	"context"
	"database/sql"
	stderr "errors"
	"fmt"
	"time"

	"github.com/skrpld/NearBeee/internal/core/database/postgres"
	"github.com/skrpld/NearBeee/internal/core/models/entities"
	"github.com/skrpld/NearBeee/pkg/errors"

	_ "embed"

	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var (
	//go:embed sql/auth/user_create.sql
	qUserCreate string
	//go:embed sql/auth/user_get_by_email.sql
	qUserGetByEmail string
	//go:embed sql/auth/user_get_by_id.sql
	qUserGetById string
	//go:embed sql/auth/user_update_refresh_token_by_user_id.sql
	qUserUpdateRefreshTokenByUserId string
)

type AuthRepository struct {
	ctx        context.Context
	postgresDB *postgres.PostgresDB
}

func NewAuthRepository(postgresDB *postgres.PostgresDB) *AuthRepository {
	return &AuthRepository{
		postgresDB: postgresDB,
	}
}

const (
	usersTableName = "users"
)

func (r *AuthRepository) CreateUser(email, passwordHash, refreshToken string, refreshTokenExpiryTime time.Time) (*entities.User, error) {
	var user entities.User

	query := fmt.Sprintf(qUserCreate, usersTableName)

	err := r.postgresDB.QueryRow(query, email, passwordHash, refreshToken, refreshTokenExpiryTime).
		Scan(&user.UserId, &user.Email, &user.PasswordHash, &user.RefreshToken, &user.RefreshTokenExpiryTime)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23505" { // 23505 - unique_violation
			return nil, errors.ErrUserAlreadyExists
		}
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User

	query := fmt.Sprintf(qUserGetByEmail, usersTableName)

	err := r.postgresDB.QueryRow(query, email).
		Scan(&user.UserId, &user.Email, &user.PasswordHash, &user.RefreshToken, &user.RefreshTokenExpiryTime)

	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrInvalidEmail
		}
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) GetUserById(userId uuid.UUID) (*entities.User, error) {
	var user entities.User

	query := fmt.Sprintf(qUserGetById, usersTableName)

	err := r.postgresDB.QueryRow(query, userId).
		Scan(&user.UserId, &user.Email, &user.PasswordHash, &user.RefreshToken, &user.RefreshTokenExpiryTime)

	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrInvalidEmail
		}
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) UpdateRefreshTokenByUserId(userId uuid.UUID, refreshToken string, refreshTokenExpiryTime time.Time) error {
	query := fmt.Sprintf(qUserUpdateRefreshTokenByUserId, usersTableName)

	_, err := r.postgresDB.Exec(query, refreshToken, refreshTokenExpiryTime, userId)
	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return errors.ErrInvalidToken
		}
		return err
	}
	return nil
}
