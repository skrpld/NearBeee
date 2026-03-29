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

	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	ctx        context.Context
	postgresDB *postgres.PostgresDB
}

func NewPostgresRepository(postgresDB *postgres.PostgresDB) *PostgresRepository {
	return &PostgresRepository{
		postgresDB: postgresDB,
	}
}

const (
	usersTableName  = "users"
	topicsTableName = "topics"
)

func (r *PostgresRepository) CreateUser(email, passwordHash, refreshToken string, refreshTokenExpiryTime time.Time) (*entities.User, error) {
	var user entities.User

	query := fmt.Sprintf(`INSERT INTO %s (email, password_hash, refresh_token, refresh_token_expiry_time) 
			VALUES ($1, $2, $3, $4) RETURNING *`, usersTableName)

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

func (r *PostgresRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE email = $1`, usersTableName)

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

func (r *PostgresRepository) GetUserById(userId uuid.UUID) (*entities.User, error) {
	var user entities.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, usersTableName)

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

func (r *PostgresRepository) UpdateRefreshTokenByUserId(userId uuid.UUID, refreshToken string, refreshTokenExpiryTime time.Time) error {
	query := fmt.Sprintf(`UPDATE %s SET refresh_token = $1, refresh_token_expiry_time = $2
          WHERE user_id = $3`, usersTableName)

	_, err := r.postgresDB.Exec(query, refreshToken, refreshTokenExpiryTime, userId)
	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return errors.ErrInvalidToken
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) CreateTopic(userId uuid.UUID, title, content, idempotencyKey string, latitude, longitude float64) (*entities.Topic, error) {
	var topic entities.Topic

	query := fmt.Sprintf(`INSERT INTO %s (user_id, title, content, idempotency_key, latitude, longitude)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`, topicsTableName)

	err := r.postgresDB.QueryRow(query, userId, title, content, idempotencyKey, latitude, longitude).
		Scan(&topic.TopicId, &topic.UserId, &topic.Title, &topic.Content, &topic.IdempotencyKey, &topic.Latitude, &topic.Longitude, &topic.CreatedAt, &topic.UpdatedAt)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23505" { // 23505 - unique_violation
			return nil, errors.ErrIdempotencyKeyAlreadyExists
		}
		return nil, err
	}

	return &topic, nil
}

func (r *PostgresRepository) GetTopicsByUserId(userId uuid.UUID, count int64) ([]*entities.Topic, error) {
	var topics []*entities.Topic

	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1 
                 ORDER BY created_at DESC LIMIT $2`, topicsTableName)
	rows, err := r.postgresDB.Query(query, userId, parsePostgresLimit(count))
	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrExpiredToken
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var topic entities.Topic

		err = rows.Scan(&topic.TopicId, &topic.UserId,
			&topic.Title, &topic.Content,
			&topic.IdempotencyKey, &topic.Latitude,
			&topic.Longitude, &topic.CreatedAt, &topic.UpdatedAt)
		if err != nil {
			return nil, err
		}

		topics = append(topics, &topic)
	}

	return topics, nil
}

func (r *PostgresRepository) GetTopicsByLocation(latitude, longitude, radius float64, count int64) ([]*entities.Topic, error) {
	var topics []*entities.Topic

	query := fmt.Sprintf(`SELECT * FROM %s
         WHERE calculate_distance($1, $2, latitude, longitude) <= $3
         ORDER BY calculate_distance($1, $2, latitude, longitude) 
         LIMIT $4`, topicsTableName)

	rows, err := r.postgresDB.Query(query, latitude, longitude, radius, parsePostgresLimit(count))
	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrInvalidCoords
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var topic entities.Topic

		err = rows.Scan(&topic.TopicId, &topic.UserId,
			&topic.Title, &topic.Content,
			&topic.IdempotencyKey, &topic.Latitude,
			&topic.Longitude, &topic.CreatedAt, &topic.UpdatedAt)
		if err != nil {
			return nil, err
		}

		topics = append(topics, &topic)
	}

	return topics, nil
}

func (r *PostgresRepository) GetTopicByTopicId(topicId uuid.UUID) (*entities.Topic, error) {
	var topic entities.Topic

	query := fmt.Sprintf(`SELECT * FROM %s WHERE topic_id = $1`, topicsTableName)
	err := r.postgresDB.QueryRow(query, topicId).
		Scan(&topic.TopicId, &topic.UserId,
			&topic.Title, &topic.Content,
			&topic.IdempotencyKey, &topic.Latitude,
			&topic.Longitude, &topic.CreatedAt, &topic.UpdatedAt)

	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrInvalidTopicId
		}
		return nil, err
	}

	return &topic, nil
}

func (r *PostgresRepository) UpdateTopicById(title, content string, topicId, userId uuid.UUID) (*entities.Topic, error) {
	var topic entities.Topic

	query := fmt.Sprintf(`UPDATE %s SET title = $1, content = $2 
          WHERE topic_id = $3 AND user_id = $4 RETURNING *`, topicsTableName)
	//TODO: по хорошему добавить проверку на доступ к посту (и месаги) а не просто инвалид пост ид
	err := r.postgresDB.QueryRow(query, title, content, topicId, userId).
		Scan(&topic.TopicId, &topic.UserId, &topic.Title,
			&topic.Content, &topic.IdempotencyKey,
			&topic.Latitude, &topic.Longitude,
			&topic.CreatedAt, &topic.UpdatedAt)
	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrInvalidTopicId
		}
		return nil, err
	}

	return &topic, nil
}

func (r *PostgresRepository) DeleteTopicById(topicId, userId uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE topic_id = $1 AND user_id = $2`, topicsTableName)

	result, err := r.postgresDB.Exec(query, topicId, userId)
	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return errors.ErrInvalidTopicId
		}
		return err
	}

	countRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if countRows != 1 {
		return errors.ErrInvalidTopicId
	}

	return nil
}

func parsePostgresLimit(limit int64) any {
	if limit < 1 {
		return nil
	}
	return limit
}
