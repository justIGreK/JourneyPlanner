package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	return err
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
		logs.Errorf("getinvites Find() error: %v", err)
	}
	var invites []models.Invitation
	cursor.All(ctx, &invites)
	if err != nil {
		logs.Errorf("getinvites All() error: %v", err)
		return nil, err
	}
	return invites, nil
}

func (r *MongoInviteRepo) IsAlreadyInvited(ctx context.Context, groupOID primitive.ObjectID, userLogin string) bool{
	var invite models.Invitation
	filter := bson.M{
		"$and": []bson.M{
			{"isUsed": false},
			{"receiver": userLogin},
			{"group_id": groupOID},
		},
	}
	err := r.InviteColl.FindOne(ctx, filter).Decode(&invite)
	if err != nil{
		if err == mongo.ErrNoDocuments{
			return true
		}
		logs.Error(err)
	}
	return false
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
		logs.Error("DeleteGroup error", err)
		return err
	}
	return nil
}

func (r *MongoInviteRepo) DeleteInviteByID(ctx context.Context, inviteID primitive.ObjectID, userLogin string) (int64, error) {
	filter := bson.M{
		"$and": []bson.M{
			{"_id": inviteID},
			{"receiver": userLogin},
			{"isUsed": false},
		},
	}
	update := bson.M{"$set": bson.M{"isUsed": true}}
	result, err := r.InviteColl.UpdateOne(ctx, filter, update)
	if err != nil {
		logs.Error("DeleteInviteByID error", err)
		return 0, err
	}
	return result.ModifiedCount, nil
}
