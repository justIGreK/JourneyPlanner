package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var logs *zap.SugaredLogger

func SetLogger(l *zap.Logger) {
	logs = l.Sugar()
}

type MongoGroupRepo struct {
	GroupColl *mongo.Collection
}

func NewMongoGroupRepo(db *mongo.Client) *MongoGroupRepo {
	return &MongoGroupRepo{GroupColl: db.Database(dbname).Collection(groupCollection)}
}

func (r *MongoGroupRepo) CreateGroup(ctx context.Context, group models.Group) (string, error) {
	result, err := r.GroupColl.InsertOne(ctx, group)
	if err != nil {
		return "", fmt.Errorf("CreateGroup: %v", err)
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *MongoGroupRepo) GetGroupList(ctx context.Context, userLogin string) ([]models.Group, error) {
	var groupList []models.Group
	filter := bson.M{
		"$and": []bson.M{
			{"members": userLogin},
			{"isActive": true},
		},
	}
	cursor, err := r.GroupColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("GetGroupList error: %v", err)
	}
	err = cursor.All(ctx, &groupList)
	if err != nil {
		return nil, fmt.Errorf("cursorAll: %v", err)
	}
	return groupList, nil
}

func (r *MongoGroupRepo) GetGroup(ctx context.Context, groupID string, userLogin ...string) (*models.Group, error) {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return nil, fmt.Errorf("InvalidID: %v", err)
	}
	var group models.Group

	filters := []bson.M{
		{"_id": oid[0]},
		{"isActive": true},
	}

	if len(userLogin) > 0 {
		filters = append(filters, bson.M{"members": userLogin[0]})
	}

	filter := bson.M{"$and": filters}

	err = r.GroupColl.FindOne(ctx, filter).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("GetGroupById error: %v", err)
	}
	return &group, nil
}

func (r *MongoGroupRepo) ChangeGroupLeader(ctx context.Context, groupID, userLogin string) error {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	filter := bson.M{
		"$and": []bson.M{
			{"_id": oid[0]},
			{"isActive": true},
		},
	}
	update := bson.M{"$set": bson.M{"leader_login": userLogin}}
	_, err = r.GroupColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("ChangeGroupLeader error: %v", err)
	}

	return nil
}

func (r *MongoGroupRepo) DeleteGroup(ctx context.Context, groupID string) error {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	filter := bson.M{
		"$and": []bson.M{
			{"_id": oid[0]},
			{"isActive": true},
		},
	}
	update := bson.M{"$set": bson.M{"isActive": false}}
	_, err = r.GroupColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("DeleteGroup error: %v", err)
	}
	return nil
}

func (r *MongoGroupRepo) LeaveGroup(ctx context.Context, groupID, userLogin string) error {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	filter := bson.M{
		"$and": []bson.M{
			{"_id": oid[0]},
			{"isActive": true},
		},
	}
	update := bson.M{
		"$pull": bson.M{
			"members": userLogin,
		},
	}
	_, err = r.GroupColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("LeaveGroup error: %v", err)
	}

	return nil
}

func (r *MongoGroupRepo) JoinGroup(ctx context.Context, groupID, userLogin string) error {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	filter := bson.M{
		"$and": []bson.M{
			{"_id": oid[0]},
			{"isActive": true},
		},
	}
	update := bson.M{
		"$push": bson.M{
			"members": userLogin,
		},
	}
	_, err = r.GroupColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Join group error: %v", err)
	}

	return nil
}
