package repository

import (
	"context"
	"database/sql"
	stderr "errors"
	"fmt"

	_ "embed"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/skrpld/NearBeee/internal/core/database/postgres"
	"github.com/skrpld/NearBeee/internal/core/models/entities"
	"github.com/skrpld/NearBeee/pkg/errors"
)

var (
	//go:embed sql/topics/topic_create.sql
	qTopicCreate string
	//go:embed sql/topics/topic_get_by_user_id.sql
	qTopicGetByUserId string
	//go:embed sql/topics/topic_get_by_location.sql
	qTopicGetByLocation string
	//go:embed sql/topics/topic_get_by_id.sql
	qTopicGetById string
	//go:embed sql/topics/topic_update_by_id.sql
	qTopicUpdateById string
	//go:embed sql/topics/topic_delete_by_id.sql
	qTopicDeleteById string
)

type TopicsRepository struct {
	ctx        context.Context
	postgresDB *postgres.PostgresDB
}

func NewTopicsRepository(postgresDB *postgres.PostgresDB) *TopicsRepository {
	return &TopicsRepository{
		postgresDB: postgresDB,
	}
}

const (
	topicsTableName = "topics"
)

func (r *TopicsRepository) CreateTopic(userId uuid.UUID, title, content, idempotencyKey string, latitude, longitude float64) (*entities.Topic, error) {
	var topic entities.Topic

	query := fmt.Sprintf(qTopicCreate, topicsTableName)

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

func (r *TopicsRepository) GetTopicsByUserId(userId uuid.UUID, count int64) ([]*entities.Topic, error) {
	var topics []*entities.Topic

	query := fmt.Sprintf(qTopicGetByUserId, topicsTableName)
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

func (r *TopicsRepository) GetTopicsByLocation(latitude, longitude, radius float64, count int64) ([]*entities.Topic, error) {
	var topics []*entities.Topic

	query := fmt.Sprintf(qTopicGetByLocation, topicsTableName)

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

func (r *TopicsRepository) GetTopicById(topicId uuid.UUID) (*entities.Topic, error) {
	var topic entities.Topic

	query := fmt.Sprintf(qTopicGetById, topicsTableName)
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

func (r *TopicsRepository) UpdateTopicById(title, content string, topicId, userId uuid.UUID) (*entities.Topic, error) {
	var topic entities.Topic

	query := fmt.Sprintf(qTopicUpdateById, topicsTableName)
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

func (r *TopicsRepository) DeleteTopicById(topicId, userId uuid.UUID) error {
	query := fmt.Sprintf(qTopicDeleteById, topicsTableName)

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
