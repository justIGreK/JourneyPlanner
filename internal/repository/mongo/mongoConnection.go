package mongorepo

import (
	"context"

	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	dbname          = "journeydb"
	pollCollection  = "polls"
	taskCollection  = "tasks"
	userCollection  = "users"
	groupCollection = "groups"
	inviteCollection = "invites"
	blacklistCollection = "blacklist"
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
