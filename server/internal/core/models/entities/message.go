package entities

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	MessageId string    `json:"message_id"`
	PostId    uuid.UUID `json:"post_id"`
	UserId    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
