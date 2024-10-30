package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"

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
	return err
}

func (r *ChatRepo) FindMessagesByChatID(ctx context.Context, groupID string) ([]models.Message, error) {
	var messages []models.Message
	cursor, err := r.ChatColl.Find(ctx, bson.M{"group_id": groupID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &messages)
	if err!=nil{
		logs.Errorf("error loading messages:%v", err)
		return nil, err
	}
	return messages, nil
}
