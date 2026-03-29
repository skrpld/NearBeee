package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/skrpld/NearBeee/internal/core/models/dto"
	"github.com/skrpld/NearBeee/internal/core/models/entities"
	"github.com/skrpld/NearBeee/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MessagesRepository interface {
	CreateMessage(ctx context.Context, topicId, userId uuid.UUID, content string) (*entities.Message, error)
	GetMessageByMessageId(ctx context.Context, messageId bson.ObjectID) (*entities.Message, error)
	GetMessageByUserId(ctx context.Context, userId uuid.UUID, count int64) ([]*entities.Message, error)
	GetMessagesByTopicId(ctx context.Context, topicId uuid.UUID, count int64) ([]*entities.Message, error)
	UpdateMessageById(ctx context.Context, messageId bson.ObjectID, userId uuid.UUID, content string) (*entities.Message, error)
	DeleteMessageById(ctx context.Context, messageId bson.ObjectID, userId uuid.UUID) error
}

type MessagesService struct {
	repo MessagesRepository
}

func NewMessagesService(repo MessagesRepository) *MessagesService {
	return &MessagesService{repo: repo}
}

func (s *MessagesService) CreateMessage(ctx context.Context, rows *dto.CreateMessageRequest) (*dto.CreateMessageResponse, error) {
	topicId, err := uuid.Parse(rows.TopicId)
	if err != nil {
		return nil, errors.ErrInvalidTopicId
	}

	message, err := s.repo.CreateMessage(ctx, topicId, rows.UserId, rows.Content)
	if err != nil {
		return nil, err
	}

	response := dto.CreateMessageResponse{
		Message: message,
	}

	return &response, nil
}

func (s *MessagesService) GetMessageByMessageId(ctx context.Context, rows *dto.GetMessageByMessageIdRequest) (*dto.GetMessageByMessageIdResponse, error) {
	objectId, err := bson.ObjectIDFromHex(rows.MessageId)
	if err != nil {
		return nil, errors.ErrInvalidMsgId
	}

	message, err := s.repo.GetMessageByMessageId(ctx, objectId)
	if err != nil {
		return nil, err
	}

	response := dto.GetMessageByMessageIdResponse{
		Message: message,
	}

	return &response, nil
}

func (s *MessagesService) GetMessageByUserId(ctx context.Context, rows *dto.GetMessageByUserIdRequest) (*dto.GetMessageByUserIdResponse, error) {
	messages, err := s.repo.GetMessageByUserId(ctx, rows.UserId, rows.Count)
	if err != nil {
		return nil, err
	}

	response := dto.GetMessageByUserIdResponse{
		Messages: messages,
	}

	return &response, nil
}

func (s *MessagesService) GetMessagesByTopicId(ctx context.Context, rows *dto.GetMessagesByTopicIdRequest) (*dto.GetMessagesByTopicIdResponse, error) {
	topicId, err := uuid.Parse(rows.TopicId)
	if err != nil {
		return nil, errors.ErrInvalidTopicId
	}

	messages, err := s.repo.GetMessagesByTopicId(ctx, topicId, rows.Count)
	if err != nil {
		return nil, err
	}

	response := dto.GetMessagesByTopicIdResponse{
		Messages: messages,
	}

	return &response, nil
}

func (s *MessagesService) UpdateMessageById(ctx context.Context, rows *dto.UpdateMessageByIdRequest) (*dto.UpdateMessageByIdResponse, error) {
	objectId, err := bson.ObjectIDFromHex(rows.MessageId)
	if err != nil {
		return nil, errors.ErrInvalidMsgId
	}

	message, err := s.repo.UpdateMessageById(ctx, objectId, rows.UserId, rows.Content)
	if err != nil {
		return nil, err
	}

	response := dto.UpdateMessageByIdResponse{
		Message: message,
	}

	return &response, nil
}

func (s *MessagesService) DeleteMessageById(ctx context.Context, rows *dto.DeleteMessageByIdRequest) (*dto.DeleteMessageByIdResponse, error) {
	objectId, err := bson.ObjectIDFromHex(rows.MessageId)
	if err != nil {
		return nil, errors.ErrInvalidMsgId
	}

	err = s.repo.DeleteMessageById(ctx, objectId, rows.UserId)
	if err != nil {
		return nil, err
	}

	response := dto.DeleteMessageByIdResponse{
		Success: true,
	}

	return &response, nil
}
