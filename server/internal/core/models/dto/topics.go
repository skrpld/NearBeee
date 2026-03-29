package dto

import (
	"github.com/google/uuid"
	"github.com/skrpld/NearBeee/internal/core/models/entities"
)

type CreateTopicRequest struct {
	UserId         uuid.UUID `json:"-"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	IdempotencyKey string    `json:"idempotency_key"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
}

type CreateTopicResponse struct {
	TopicId string `json:"topic_id"` //TODO: возвращать сам topic
}

type GetTopicsByUserIdRequest struct {
	UserId uuid.UUID `json:"-"`
	Count  int64     `json:"count"`
}

type GetTopicsByUserIdResponse struct {
	Topics []*entities.Topic `json:"topics"`
}

type GetTopicsByLocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Count     int64   `json:"count"`
	Radius    float64 `json:"radius"`
}

type GetTopicsByLocationResponse struct {
	Topics []*entities.Topic `json:"topics"`
}

type GetTopicByTopicIdRequest struct {
	TopicId string `json:"-"`
}

type GetTopicByTopicIdResponse struct {
	Topic *entities.Topic `json:"topic"`
}

type UpdateTopicByIdRequest struct {
	TopicId string    `json:"-"`
	UserId  uuid.UUID `json:"-"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}

type UpdateTopicByIdResponse struct {
	Topic *entities.Topic `json:"topic"`
}

type DeleteTopicByIdRequest struct {
	TopicId string    `json:"-"`
	UserId  uuid.UUID `json:"-"`
}

type DeleteTopicResponse struct {
	TopicId string `json:"topic_id"` //TODO: возвращать success
}
