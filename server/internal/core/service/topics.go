package service

import (
	"github.com/skrpld/NearBeee/internal/core/models/dto"
	"github.com/skrpld/NearBeee/internal/core/models/entities"
	"github.com/skrpld/NearBeee/pkg/errors"

	"github.com/google/uuid"
)

type TopicsRepository interface {
	CreateTopic(userId uuid.UUID, title, content, idempotencyKey string, latitude, longitude float64) (*entities.Topic, error)
	GetTopicsByUserId(userId uuid.UUID, count int64) ([]*entities.Topic, error)
	GetTopicsByLocation(latitude, longitude, radius float64, count int64) ([]*entities.Topic, error)
	GetTopicById(topicId uuid.UUID) (*entities.Topic, error)
	UpdateTopicById(title, content string, topicId, userId uuid.UUID) (*entities.Topic, error)
	DeleteTopicById(topicId, userId uuid.UUID) error
}

type TopicsService struct {
	repo TopicsRepository
}

func NewTopicsService(repo TopicsRepository) *TopicsService {
	return &TopicsService{repo: repo}
}

func (s *TopicsService) CreateTopic(rows *dto.CreateTopicRequest) (*dto.CreateTopicResponse, error) {
	topic, err := s.repo.CreateTopic(rows.UserId, rows.Title, rows.Content, rows.IdempotencyKey, rows.Latitude, rows.Longitude)
	if err != nil {
		return nil, err
	}

	response := dto.CreateTopicResponse{
		TopicId: topic.TopicId.String(),
	}

	return &response, nil
}

func (s *TopicsService) GetTopicsByUserId(rows *dto.GetTopicsByUserIdRequest) (*dto.GetTopicsByUserIdResponse, error) {
	topics, err := s.repo.GetTopicsByUserId(rows.UserId, rows.Count)
	if err != nil {
		return nil, err
	}

	response := dto.GetTopicsByUserIdResponse{
		Topics: topics,
	}

	return &response, nil
}

func (s *TopicsService) GetTopicsByLocation(rows *dto.GetTopicsByLocationRequest) (*dto.GetTopicsByLocationResponse, error) {
	topics, err := s.repo.GetTopicsByLocation(rows.Latitude, rows.Longitude, rows.Radius, rows.Count)
	if err != nil {
		return nil, err
	}

	response := dto.GetTopicsByLocationResponse{
		Topics: topics,
	}

	return &response, nil
}

func (s *TopicsService) GetTopicById(rows *dto.GetTopicByTopicIdRequest) (*dto.GetTopicByTopicIdResponse, error) {
	topicId, err := uuid.Parse(rows.TopicId)
	if err != nil {
		return nil, errors.ErrInvalidTopicId
	}

	topic, err := s.repo.GetTopicById(topicId)
	if err != nil {
		return nil, err
	}

	response := dto.GetTopicByTopicIdResponse{
		Topic: topic,
	}

	return &response, nil
}

func (s *TopicsService) UpdateTopicById(rows *dto.UpdateTopicByIdRequest) (*dto.UpdateTopicByIdResponse, error) {
	topicId, err := uuid.Parse(rows.TopicId)
	if err != nil {
		return nil, errors.ErrInvalidTopicId
	}

	topic, err := s.repo.UpdateTopicById(rows.Title, rows.Content, topicId, rows.UserId)
	if err != nil {
		return nil, err
	}

	response := dto.UpdateTopicByIdResponse{
		Topic: topic,
	}

	return &response, nil
}

func (s *TopicsService) DeleteTopicById(rows *dto.DeleteTopicByIdRequest) (*dto.DeleteTopicResponse, error) {
	topicId, err := uuid.Parse(rows.TopicId)
	if err != nil {
		return nil, errors.ErrInvalidTopicId
	}

	err = s.repo.DeleteTopicById(topicId, rows.UserId)
	if err != nil {
		return nil, err
	}

	response := dto.DeleteTopicResponse{
		TopicId: rows.TopicId,
	}

	return &response, nil
}
