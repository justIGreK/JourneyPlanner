package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoInviteRepo struct {
	InviteColl *mongo.Collection
}

func NewMongoInviteRepo(db *mongo.Client) *MongoInviteRepo {
	return &MongoInviteRepo{InviteColl: db.Database(dbname).Collection(inviteCollection)}
}

func (r *MongoInviteRepo) AddInvitation(ctx context.Context, inviteCreator, invitedUser, groupName, token string) error {
	invite := models.Invitation{
		Sender:    inviteCreator,
		Receiver:  invitedUser,
		GroupName: groupName,
		Token:     token,
		IsUsed: false,
	}
	_, err := r.InviteColl.InsertOne(ctx, invite)
	return err
}

func (r *MongoInviteRepo) GetInvites(ctx context.Context, userLogin string) ([]models.Invitation, error) {
	filter := bson.M{
		"$and":[]bson.M{
			{"isUsed": false},
			{"receiver":userLogin},
		},
	}
	cursor, err := r.InviteColl.Find(ctx, filter)
	if err != nil{
		logs.Errorf("getinvites Find() error: %v", err)
	}
	var invites []models.Invitation
	cursor.All(ctx, &invites)
	if err != nil{
		logs.Errorf("getinvites All() error: %v", err)
		return nil, err
	}
	return invites, nil
}

func (r *MongoInviteRepo) DeleteInvite(ctx context.Context, token string) error {
	filter := bson.M{
		"$and": []bson.M{
			{"token": token},
			{"isUsed": false},
		},
	}
	update := bson.M{"$set": bson.M{"isUsed": true}}
	_, err := r.InviteColl.UpdateOne(ctx, filter, update)
	if err != nil {
		logs.Error("DeleteGroup error", err)
		return err
	}
	return nil
}