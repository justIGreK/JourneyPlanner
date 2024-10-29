package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoBlacklistRepo struct {
	BlacklistColl *mongo.Collection
}

func NewMongoBlacklistRepo(db *mongo.Client) *MongoBlacklistRepo {
	return &MongoBlacklistRepo{BlacklistColl: db.Database(dbname).Collection(blacklistCollection)}
}

func (r *MongoBlacklistRepo) CreateBlacklist(ctx context.Context, groupOID primitive.ObjectID) error {
	blacklist := models.BlackList{
		GroupID:   groupOID,
		Blacklist: []string{},
	}
	_, err := r.BlacklistColl.InsertOne(ctx, blacklist)
	return err
}

func (r *MongoBlacklistRepo) BanUser(ctx context.Context, groupOID primitive.ObjectID, userLogin string) error {
	filter := bson.M{
		"group_id": groupOID,
	}
	update := bson.M{"$push": bson.M{"blacklist": userLogin}}

	_, err := r.BlacklistColl.UpdateOne(ctx, filter, update)
	if err != nil {
		logs.Error("Blacklist error", err)
		return err
	}
	return nil
}

func (r *MongoBlacklistRepo) UnbanUser(ctx context.Context, groupOID primitive.ObjectID, userLogin string) error {
	filter := bson.M{
		"group_id": groupOID,
	}
	update := bson.M{"$pull": bson.M{"blacklist": userLogin}}

	_, err := r.BlacklistColl.UpdateOne(ctx, filter, update)
	if err != nil {
		logs.Error("Blacklist error", err)
		return err
	}
	return nil
}

func (r *MongoBlacklistRepo) GetBlacklist(ctx context.Context, groupOID primitive.ObjectID) (*models.BlackList, error) {
	var blacklist models.BlackList
	filter := bson.M{
		"group_id": groupOID,
	}
	err := r.BlacklistColl.FindOne(ctx, filter).Decode(&blacklist)
	if err != nil {
		logs.Error("Get Blacklist error", err)
		return nil, err
	}
	return &blacklist, nil
}
