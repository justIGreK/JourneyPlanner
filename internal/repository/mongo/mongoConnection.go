package mongorepo

import (
	"context"
	"errors"

	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	dbname              = "journeydb"
	pollCollection      = "polls"
	taskCollection      = "tasks"
	userCollection      = "users"
	groupCollection     = "groups"
	inviteCollection    = "invites"
	blacklistCollection = "blacklist"
	chatCollection      = "messages"
)

func CreateMongoClient(ctx context.Context) *mongo.Client {
	dbURI := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		logs.Fatal("Failed to create MongoDB client: ", zap.Error(err))
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		logs.Fatal("MongoDB is not connected: ", zap.Error(err))
	}
	return client
}

func convertToObjectIDs(ids ...string) ([]primitive.ObjectID, error) {
	objectIDs := make([]primitive.ObjectID, 0, len(ids))

	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New("InvalidID: " + id)
		}
		objectIDs = append(objectIDs, oid)
	}

	return objectIDs, nil
}
