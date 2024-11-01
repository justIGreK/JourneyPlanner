package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRepo struct {
	ChatColl *mongo.Collection
}

func NewChatRepository(db *mongo.Client) *ChatRepo {
	return &ChatRepo{
		ChatColl: db.Database(dbname).Collection(chatCollection),
	}
}

func (r *ChatRepo) InsertMessage(ctx context.Context, msg models.Message) error {
	_, err := r.ChatColl.InsertOne(ctx, msg)
	if err != nil {
		return fmt.Errorf("InsertMessage error: %v", err)
	}
	return nil
}

func (r *ChatRepo) FindMessagesByChatID(ctx context.Context, groupID string) ([]models.Message, error) {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return nil, fmt.Errorf("InvalidID: %v", err)
	}
	var messages []models.Message
	cursor, err := r.ChatColl.Find(ctx, bson.M{"group_id": oid[0]})
	if err != nil {
		return nil, fmt.Errorf("FindMessagesByChatID error: %v", err)
	}
	err = cursor.All(ctx, &messages)
	if err != nil {
		return nil, fmt.Errorf("error loading messages: %v", err)
	}
	return messages, nil
}
