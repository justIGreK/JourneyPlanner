package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoInviteRepo struct {
	InviteColl *mongo.Collection
}

func NewMongoInviteRepo(db *mongo.Client) *MongoInviteRepo {
	return &MongoInviteRepo{InviteColl: db.Database(dbname).Collection(inviteCollection)}
}

func (r *MongoInviteRepo) AddInvitation(ctx context.Context, invite models.Invitation) error {
	_, err := r.InviteColl.InsertOne(ctx, invite)
	if err != nil {
		return fmt.Errorf("AddInvitation error: %v", err)
	}
	return nil
}

func (r *MongoInviteRepo) GetInvites(ctx context.Context, userLogin string) ([]models.Invitation, error) {
	filter := bson.M{
		"$and": []bson.M{
			{"isUsed": false},
			{"receiver": userLogin},
		},
	}
	cursor, err := r.InviteColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("getInvites error: %v", err)
	}
	var invites []models.Invitation
	err = cursor.All(ctx, &invites)
	if err != nil {
		return nil, fmt.Errorf("getinvites All(): %v", err)
	}
	return invites, nil
}

func (r *MongoInviteRepo) IsAlreadyInvited(ctx context.Context, groupID, userLogin string) (bool, error) {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return false, fmt.Errorf("InvalidID: %v", err)
	}
	var invite models.Invitation
	filter := bson.M{
		"$and": []bson.M{
			{"isUsed": false},
			{"receiver": userLogin},
			{"group_id": oid[0]},
		},
	}
	err = r.InviteColl.FindOne(ctx, filter).Decode(&invite)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return true, nil
		}
		return false, fmt.Errorf("isAlreadyInvites error: %v", err)
	}
	return false, nil
}

func (r *MongoInviteRepo) DeleteInviteByToken(ctx context.Context, token string) error {
	filter := bson.M{
		"$and": []bson.M{
			{"token": token},
			{"isUsed": false},
		},
	}
	update := bson.M{"$set": bson.M{"isUsed": true}}
	_, err := r.InviteColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("DeleteGroup error: %v", err)
	}
	return nil
}

func (r *MongoInviteRepo) DeleteInviteByID(ctx context.Context, inviteID, userLogin string) (int64, error) {
	oid, err := convertToObjectIDs(inviteID)
	if err != nil {
		return 0, fmt.Errorf("InvalidID: %v", err)
	}
	filter := bson.M{
		"$and": []bson.M{
			{"_id": oid[0]},
			{"receiver": userLogin},
			{"isUsed": false},
		},
	}
	update := bson.M{"$set": bson.M{"isUsed": true}}
	result, err := r.InviteColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf("DeleteInviteByID error: %v", err)
	}
	return result.ModifiedCount, nil
}
