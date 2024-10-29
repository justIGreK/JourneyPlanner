package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoPollRepo struct {
	PollColl *mongo.Collection
}

func NewMongoPollRepo(db *mongo.Client) *MongoPollRepo {
	return &MongoPollRepo{PollColl: db.Database(dbname).Collection(pollCollection)}
}

func (r *MongoPollRepo) CreatePoll(ctx context.Context, poll models.Poll) error {
	_, err := r.PollColl.InsertOne(ctx, poll)
	return err
}

func (r *MongoPollRepo) GetPollList(ctx context.Context, groupOID primitive.ObjectID) ([]models.Poll, []models.Poll, error) {
	now := time.Now().UTC()
	filter := bson.M{
		"group_id": groupOID,
		"$or": []bson.M{
			{"endtime": bson.M{"$lt": now}},
			{"isEarlyClosed": true},
		},
	}
	openPolls, closedPolls := []models.Poll{}, []models.Poll{}
	cursor, err := r.PollColl.Find(ctx, filter)
	if err != nil {
		logs.Error("GetPollList error", err)
		return nil, nil, err
	}
	err = cursor.All(ctx, &closedPolls)
	if err != nil {
		logs.Error("GetPollList error, cursor.All()", err)
		return nil, nil, err
	}
	filter = bson.M{
		"$and": []bson.M{
			{"group_id": groupOID},
			{"isEarlyClosed": false},
			{"endtime": bson.M{"$gt": now}},
		},
	}
	cursor, err = r.PollColl.Find(ctx, filter)
	if err != nil {
		logs.Error("GetPollList error", err)
		return nil, nil, err
	}
	err = cursor.All(ctx, &openPolls)
	if err != nil {
		logs.Error("GetPollList error, cursor.All()", err)
		return nil, nil, err
	}
	
	return openPolls, closedPolls, nil
}

func (r *MongoPollRepo) GetPollById(ctx context.Context, pollOID primitive.ObjectID) (*models.Poll, error) {
	var poll models.Poll
	filter := bson.M{
		"_id": pollOID,
	}
	err := r.PollColl.FindOne(ctx, filter).Decode(&poll)
	if err != nil {
		logs.Error("GetPollById error", err)
		return nil, err
	}
	return &poll, nil
}

func (r *MongoPollRepo) ClosePoll(ctx context.Context, pollOID primitive.ObjectID) error {
	filter := bson.M{
		"$and": []bson.M{
			{"_id": pollOID},
			{"isEarlyClosed": false},
		},
	}
	update := bson.M{"$set": bson.M{"isEarlyClosed": true}}
	_, err := r.PollColl.UpdateOne(ctx, filter, update)
	if err != nil {
		logs.Error("ClosePoll error", err)
		return err
	}
	return nil
}

func (r *MongoPollRepo) DeletePoll(ctx context.Context, pollOID primitive.ObjectID) error {
	filter := bson.M{"_id": pollOID}
	_, err := r.PollColl.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

const (
	firstOption  = "firstOption"
	secondOption = "secondOption"
)

func (r *MongoPollRepo) AddVote(ctx context.Context, pollOID primitive.ObjectID, voteOption, userLogin string) error {
	now := time.Now().UTC()
	filter := bson.M{
		"$and": []bson.M{
			{"_id": pollOID},
			{"endtime": bson.M{"$gt": now}},
			{"isEarlyClosed": false},
		},
	}
	var update bson.M
	if voteOption == firstOption {
		update = bson.M{"$push": bson.M{
			"votes1": userLogin,
		}}
	} else if voteOption == secondOption {
		update = bson.M{"$push": bson.M{
			"votes2": userLogin,
		}}
	} else {
		return errors.New("unexpected error. Invalid voteOption")
	}
	_, err := r.PollColl.UpdateOne(ctx, filter, update)
	if err != nil {
		logs.Error("ClosePoll error", err)
		return err
	}
	return nil
}
func (r *MongoPollRepo) RemoveVote(ctx context.Context, pollOID primitive.ObjectID, userLogin string) error {
	now := time.Now().UTC()
	filter := bson.M{
		"$and": []bson.M{
			{"_id": pollOID},
			{"endtime": bson.M{"$gt": now}},
			{"isEarlyClosed": false},
		},
	}

	update := bson.M{"$pull": bson.M{
		"votes1": userLogin,
		"votes2": userLogin,
	}}
	_, err := r.PollColl.UpdateOne(ctx, filter, update)
	if err != nil {
		logs.Error("ClosePoll error", err)
		return err
	}
	return nil
}
