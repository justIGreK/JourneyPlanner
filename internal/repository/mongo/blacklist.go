package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoBlacklistRepo struct {
	BlacklistColl *mongo.Collection
}

func NewMongoBlacklistRepo(db *mongo.Client) *MongoBlacklistRepo {
	return &MongoBlacklistRepo{BlacklistColl: db.Database(dbname).Collection(blacklistCollection)}
}

func (r *MongoBlacklistRepo) CreateBlacklist(ctx context.Context, groupID string) error {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return errors.New("InvalidID")
	}
	blacklist := models.BlackList{
		GroupID:   oid[0],
		Blacklist: []string{},
	}
	_, err = r.BlacklistColl.InsertOne(ctx, blacklist)
	return err
}

func (r *MongoBlacklistRepo) BanUser(ctx context.Context, groupID, userLogin string) error {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return errors.New("InvalidID")
	}
	filter := bson.M{
		"group_id": oid[0],
	}
	update := bson.M{"$push": bson.M{"blacklist": userLogin}}

	_, err = r.BlacklistColl.UpdateOne(ctx, filter, update)
	if err != nil {
		logs.Error("Blacklist error", err)
		return err
	}
	return nil
}

func (r *MongoBlacklistRepo) UnbanUser(ctx context.Context, groupID, userLogin string) error {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return errors.New("InvalidID")
	}
	filter := bson.M{
		"group_id": oid[0],
	}
	update := bson.M{"$pull": bson.M{"blacklist": userLogin}}

	_, err = r.BlacklistColl.UpdateOne(ctx, filter, update)
	if err != nil {
		logs.Error("Blacklist error", err)
		return err
	}
	return nil
}

func (r *MongoBlacklistRepo) GetBlacklist(ctx context.Context, groupID string) (*models.BlackList, error) {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return nil, errors.New("InvalidID")
	}
	var blacklist models.BlackList
	filter := bson.M{
		"group_id": oid[0],
	}
	err = r.BlacklistColl.FindOne(ctx, filter).Decode(&blacklist)
	if err != nil {
		logs.Error("Get Blacklist error", err)
		return nil, err
	}
	return &blacklist, nil
}
