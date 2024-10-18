package mongorepo

import (
	logger "JourneyPlanner/pkg/log"
	"context"

	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	dbname         = "journeydb"
	pollCollection = "polls"
	taskCollection = "tasks"
	userCollection = "users"
)

func CreateMongoClient(ctx context.Context) *mongo.Client {
	logger := logger.GetLogger()
	dbURI := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		logger.Fatal("Failed to create MongoDB client: ", zap.Error(err))
	}

	return client
}
