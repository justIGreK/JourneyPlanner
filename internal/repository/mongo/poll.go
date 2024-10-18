package mongorepo

import (
	"go.mongodb.org/mongo-driver/mongo"
)


type MongoPollRepo struct {
	PollColl *mongo.Collection
}

func NewMongoPollRepo(db *mongo.Client) *MongoPollRepo{
	return &MongoPollRepo{PollColl: db.Database(dbname).Collection(pollCollection)}
}
