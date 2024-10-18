package mongorepo

import "go.mongodb.org/mongo-driver/mongo"

type MongoTaskRepo struct {
	TaskColl *mongo.Collection
}

func NewMongoTaskRepo(db *mongo.Client) *MongoTaskRepo{
	return &MongoTaskRepo{TaskColl: db.Database(dbname).Collection(taskCollection)}
}
