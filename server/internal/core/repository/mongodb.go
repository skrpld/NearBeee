package repository

import (
	"context"
	stderr "errors"
	"time"

	"github.com/google/uuid"
	"github.com/skrpld/NearBeee/internal/core/database/mongodb"
	"github.com/skrpld/NearBeee/internal/core/models/dao"
	"github.com/skrpld/NearBeee/internal/core/models/entities"
	"github.com/skrpld/NearBeee/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongodbRepository struct {
	mongoDB *mongodb.MongoDB
}

func NewMongodbRepository(mongoDB *mongodb.MongoDB) *MongodbRepository {
	return &MongodbRepository{
		mongoDB: mongoDB,
	}
}

const msgCollectionName = "messages"

func (r *MongodbRepository) CreateMessage(ctx context.Context, topicId, userId uuid.UUID, content string) (*entities.Message, error) {
	newMsg := &dao.Message{
		TopicId:   topicId,
		UserId:    userId,
		Content:   content,
		CreatedAt: currentTimeUTC(),
		UpdatedAt: currentTimeUTC(),
	}

	result, err := r.mongoDB.Collection(msgCollectionName).InsertOne(ctx, newMsg)
	if err != nil {
		return nil, err
	}

	newMsg.MessageId = result.InsertedID.(bson.ObjectID)

	return newMsg.ToEntity(), nil
}

func (r *MongodbRepository) GetMessageByMessageId(ctx context.Context, msgId bson.ObjectID) (*entities.Message, error) {
	var msg dao.Message

	err := r.mongoDB.Collection(msgCollectionName).
		FindOne(ctx, bson.M{"_id": msgId}).
		Decode(&msg)
	if err != nil {
		if stderr.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.ErrMsgNotFound
		}
		return nil, err
	}

	return msg.ToEntity(), nil
}

func (r *MongodbRepository) GetMessageByUserId(ctx context.Context, userId uuid.UUID, count int64) ([]*entities.Message, error) {
	opts := options.Find().
		SetLimit(parseMongoLimit(count)).
		SetSort(bson.M{"created_at": -1})

	result, err := r.mongoDB.Collection(msgCollectionName).
		Find(ctx, bson.M{"user_id": userId}, opts)
	if err != nil {
		return nil, err
	}

	defer result.Close(ctx)

	var msgs []dao.Message

	if err = result.All(ctx, &msgs); err != nil {
		return nil, err
	}

	response := make([]*entities.Message, 0, len(msgs))
	for _, msg := range msgs {
		response = append(response, msg.ToEntity())
	}

	return response, nil
}

func (r *MongodbRepository) GetMessagesByTopicId(ctx context.Context, topicId uuid.UUID, count int64) ([]*entities.Message, error) {
	opts := options.Find().
		SetLimit(parseMongoLimit(count)).
		SetSort(bson.M{"created_at": -1})

	result, err := r.mongoDB.Collection(msgCollectionName).
		Find(ctx, bson.M{"topic_id": topicId}, opts)
	if err != nil {
		return nil, err
	}

	defer result.Close(ctx)

	var msgs []dao.Message

	if err = result.All(ctx, &msgs); err != nil {
		return nil, err
	}

	response := make([]*entities.Message, 0, len(msgs))
	for _, msg := range msgs {
		response = append(response, msg.ToEntity())
	}

	return response, nil
}

func (r *MongodbRepository) UpdateMessageById(ctx context.Context, messageId bson.ObjectID, userId uuid.UUID, content string) (*entities.Message, error) {
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	update := bson.M{
		"$set": bson.M{
			"content":    content,
			"updated_at": currentTimeUTC(),
		},
	}

	var msg dao.Message

	err := r.mongoDB.Collection(msgCollectionName).
		FindOneAndUpdate(ctx, bson.M{"_id": messageId, "user_id": userId}, update, opts).
		Decode(&msg)
	if err != nil {
		if stderr.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.ErrMsgNotFound
		}
		return nil, err
	}

	return msg.ToEntity(), nil
}

func (r *MongodbRepository) DeleteMessageById(ctx context.Context, messageId bson.ObjectID, userId uuid.UUID) error {
	result, err := r.mongoDB.Collection(msgCollectionName).DeleteOne(ctx, bson.M{"_id": messageId, "user_id": userId})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.ErrMsgNotFound
	}

	return nil
}

func parseMongoLimit(limit int64) int64 {
	if limit < 1 {
		return 0
	}
	return limit
}

func currentTimeUTC() time.Time {
	return time.Now().UTC()
}
