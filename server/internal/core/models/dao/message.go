package dao

import (
	"time"

	"github.com/google/uuid"
	"github.com/skrpld/NearBeee/internal/core/models/entities"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Message struct {
	MessageId bson.ObjectID `bson:"_id,omitempty"`
	PostId    uuid.UUID     `bson:"post_id"`
	UserId    uuid.UUID     `bson:"user_id"`
	Content   string        `bson:"content"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

func (m *Message) ToEntity() *entities.Message {
	return &entities.Message{
		MessageId: m.MessageId.Hex(),
		PostId:    m.PostId,
		UserId:    m.UserId,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
